package db

import (
	"context"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/Juniper/contrail/pkg/errutil"
	"github.com/Juniper/contrail/pkg/format"
	"github.com/Juniper/contrail/pkg/models"
	"github.com/Juniper/contrail/pkg/services"
	"github.com/Juniper/contrail/pkg/types/ipam"
)

// CreateIpamSubnet creates IPAM subnet
func (db *Service) CreateIpamSubnet(
	ctx context.Context, request *ipam.CreateIpamSubnetRequest,
) (subnetUUID string, err error) {
	if request.IpamSubnet == nil {
		return "", errors.Errorf("can't create ipamSubnet for nil subnet")
	}

	subnetUUID = request.IpamSubnet.GetSubnetUUID()
	if subnetUUID == "" {
		subnetUUID = uuid.NewV4().String()
	}

	// TODO: check allocation pool
	// TODO: check and reserve gw
	// TODO: check and reserve service addr
	// TODO: check and reserve dns nameservers
	// TODO: check allocation units

	_, err = db.deleteIPPools(ctx, &ipPool{
		key: subnetUUID,
	})
	if err != nil {
		return "", err
	}

	ipPools, err := prepareIPPools(request.IpamSubnet, subnetUUID)
	if err != nil {
		return "", err
	}

	for _, ipPool := range ipPools {
		err = db.createIPPool(ctx, ipPool)
		if err != nil {
			return "", err
		}
	}

	return subnetUUID, err
}

// CheckIfIpamSubnetExists checks if subnet with provided subnet UUID already exists
func (db *Service) CheckIfIpamSubnetExists(ctx context.Context, subnetUUID string) (bool, error) {
	if subnetUUID == "" {
		return false, nil
	}

	res, err := db.getIPPools(ctx, &ipPool{
		key: subnetUUID,
	})

	return len(res) != 0, err
}

// DeleteIpamSubnet deletes IPAM subnet
func (db *Service) DeleteIpamSubnet(ctx context.Context, request *ipam.DeleteIpamSubnetRequest) (err error) {
	if request.SubnetUUID == "" {
		return errors.Errorf("empty subnet uuid in DeleteIpamSubnet")
	}

	deletedCount, err := db.deleteIPPools(ctx, &ipPool{
		key: request.SubnetUUID,
	})

	if err != nil {
		return err
	}

	if deletedCount == 0 {
		return errors.Errorf("ipam subnet with uuid %s doesn't exist", request.SubnetUUID)
	}

	return nil
}

// AllocateIP allocates ip
func (db *Service) AllocateIP(
	ctx context.Context, request *ipam.AllocateIPRequest,
) (address string, subnetUUID string, err error) {
	// TODO: Implement:
	//		- IPAM based ip allocation from flat subnet where instance ip is directly referring
	//		  IPAM for internal ip address
	//		- IPAM based ip allocation from flat subnet where instance ip is referring to vrouter which has
	// 		  allocation pools on vrouter->ipam link.
	//		- Virtual network based ip allocation from flat-subnet ipam
	//		- Virtual network based ip allocation from user-defined connected ipam

	// TODO: Handle allocation methods:
	//		- user-defined-subnet-preferred
	//		- flat-subnet-only
	//		- flat-subnet-preferred

	// This is a simple Virtual network based ip allocation from user-defined subnet

	// TODO: virtual network can be absent in the request
	virtualNetwork := request.VirtualNetwork
	if virtualNetwork != nil && virtualNetwork.HasNetworkBasedAllocationMethod() {
		return db.performNetworkBasedIPAllocation(ctx, request)
	}

	// TODO: we don't really allocate an IP for other methods
	return request.IPAddress, request.SubnetUUID, nil
}

// DeallocateIP deallocates ip
func (db *Service) DeallocateIP(ctx context.Context, request *ipam.DeallocateIPRequest) (err error) {
	// TODO: Implement other allocation methods
	if !request.VirtualNetwork.HasNetworkBasedAllocationMethod() {
		return nil
	}

	subnets, err := db.getRefSubnets(ctx, request.VirtualNetwork)
	if err != nil {
		return err
	}

	for _, subnet := range subnets {
		hit, err := subnet.Contains(net.ParseIP(request.IPAddress))
		if err != nil {
			return err
		}
		if !hit {
			continue
		}
		return db.deallocateIP(ctx, subnet.SubnetUUID, net.ParseIP(request.IPAddress))
	}

	return errors.Errorf("could not deallocate address %s from any of available subnets in virtual network %v",
		request.IPAddress, request.VirtualNetwork.GetUUID())
}

func (db *Service) getRefSubnets(ctx context.Context, vn *models.VirtualNetwork) ([]*models.IpamSubnetType, error) {
	var subnets []*models.IpamSubnetType
	if vn.GetAddressAllocationMethod() == models.UserDefinedSubnetOnly {
		subnets = vn.GetIpamSubnets().GetSubnets()
	} else if vn.GetAddressAllocationMethod() == models.FlatSubnetOnly {

		for _, ipamRef := range vn.GetNetworkIpamRefs() {
			ipamResponse, err := db.GetNetworkIpam(ctx, &services.GetNetworkIpamRequest{
				ID: ipamRef.GetUUID(),
			})
			if err != nil {
				return nil, errutil.ErrorBadRequestf("getting referenced network IPAM with UUID %s failed: %v",
					ipamRef.GetUUID(), err)
			}
			subnets = append(subnets, ipamResponse.GetNetworkIpam().GetIpamSubnets().GetSubnets()...)
		}
	}
	return subnets, nil
}

// IsIPAllocated checks if ip is allocated
func (db *Service) IsIPAllocated(
	ctx context.Context, request *ipam.IsIPAllocatedRequest,
) (isAllocated bool, err error) {
	// TODO: Implement other allocation methods
	if !request.VirtualNetwork.HasNetworkBasedAllocationMethod() {
		return false, nil
	}

	subnets, err := db.getRefSubnets(ctx, request.VirtualNetwork)
	if err != nil {
		return false, err
	}

	for _, subnet := range subnets {
		hit, err := subnet.Contains(net.ParseIP(request.IPAddress))
		if err != nil {
			return false, err
		}
		if !hit {
			continue
		}
		ip := net.ParseIP(request.IPAddress)
		reqPool := ipPool{subnet.SubnetUUID, ip, cidr.Inc(ip)}
		res, err := db.getIPPools(ctx, &reqPool)
		if err != nil {
			return false, err
		}
		return len(res) == 0, nil
	}
	return false, errors.Errorf(
		"provided ip %v, doesn't belong to any subnet in virtual network %v",
		request.IPAddress, request.VirtualNetwork.GetUUID())
}

// performNetworkBasedIPAllocation performs virtual network based ip allocation in a user-defined subnet
func (db *Service) performNetworkBasedIPAllocation(
	ctx context.Context, request *ipam.AllocateIPRequest,
) (address string, subnetUUID string, err error) {
	virtualNetwork := request.VirtualNetwork
	subnetUUIDs := virtualNetwork.GetIpamSubnets().UUIDs()
	if request.SubnetUUID != "" {
		if !format.ContainsString(subnetUUIDs, request.SubnetUUID) {
			return "", "", errors.Errorf("could not find subnet %s in in virtual network %v", request.SubnetUUID,
				virtualNetwork.GetUUID())
		}
		subnetUUIDs = []string{request.SubnetUUID}
	}

	for _, subnetUUID := range subnetUUIDs {
		addr, err := db.allocateIPForSubnetUUID(ctx, subnetUUID, request.IPAddress)
		if errutil.IsNotFound(err) {
			continue
		}
		if err != nil {
			return "", "", err
		}

		return addr, subnetUUID, nil
	}

	return "", "", errors.Errorf("could not allocate address %s in any available subnets in virtual network %v",
		request.IPAddress, virtualNetwork.GetUUID())
}

func (db *Service) allocateIPForSubnetUUID(
	ctx context.Context, subnetUUID string, ipRequested string,
) (address string, err error) {
	if ipRequested != "" {
		err = db.setIP(ctx, subnetUUID, net.ParseIP(ipRequested))
		if err != nil {
			return "", err
		}
		return ipRequested, nil
	}

	ip, err := db.allocateIP(ctx, subnetUUID)
	if err != nil {
		return "", err
	}
	return ip.String(), nil

}

func prepareIPPools(ipamSubnet *models.IpamSubnetType, subnetUUID string) ([]*ipPool, error) {
	var ipPools []*ipPool
	for _, pool := range ipamSubnet.GetAllocationPools() {
		ipPools = append(ipPools, &ipPool{
			key:   subnetUUID,
			start: net.ParseIP(pool.Start),
			end:   net.ParseIP(pool.End),
		})
	}

	// If there are no allocation pools just allocate the whole subnet
	if len(ipPools) == 0 {
		net, err := ipamSubnet.GetSubnet().Net()
		if err != nil {
			return nil, err
		}
		ipPool := &ipPool{
			key: subnetUUID,
		}
		ipPool.start, ipPool.end = cidr.AddressRange(net)
		ipPool.start = cidr.Inc(ipPool.start)
		ipPools = append(ipPools, ipPool)
	}

	return ipPools, nil
}
