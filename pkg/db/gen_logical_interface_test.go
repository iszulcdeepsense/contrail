// nolint
package db

import (
	"context"
	"github.com/satori/go.uuid"
	"testing"
	"time"

	"github.com/Juniper/contrail/pkg/common"
	"github.com/Juniper/contrail/pkg/models"
	"github.com/pkg/errors"
)

//For skip import error.
var _ = errors.New("")

func TestLogicalInterface(t *testing.T) {
	t.Parallel()
	db := &DB{
		DB:      testDB,
		Dialect: NewDialect("mysql"),
	}
	db.initQueryBuilders()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	model := models.MakeLogicalInterface()
	model.UUID = uuid.NewV4().String()
	model.FQName = []string{"default", "default-domain", model.UUID}
	model.Perms2.Owner = "admin"
	var err error

	// Create referred objects

	var VirtualMachineInterfaceCreateRef []*models.LogicalInterfaceVirtualMachineInterfaceRef
	var VirtualMachineInterfaceRefModel *models.VirtualMachineInterface

	VirtualMachineInterfaceRefUUID := uuid.NewV4().String()
	VirtualMachineInterfaceRefUUID1 := uuid.NewV4().String()
	VirtualMachineInterfaceRefUUID2 := uuid.NewV4().String()

	VirtualMachineInterfaceRefModel = models.MakeVirtualMachineInterface()
	VirtualMachineInterfaceRefModel.UUID = VirtualMachineInterfaceRefUUID
	VirtualMachineInterfaceRefModel.FQName = []string{"test", VirtualMachineInterfaceRefUUID}
	_, err = db.CreateVirtualMachineInterface(ctx, &models.CreateVirtualMachineInterfaceRequest{
		VirtualMachineInterface: VirtualMachineInterfaceRefModel,
	})
	VirtualMachineInterfaceRefModel.UUID = VirtualMachineInterfaceRefUUID1
	VirtualMachineInterfaceRefModel.FQName = []string{"test", VirtualMachineInterfaceRefUUID1}
	_, err = db.CreateVirtualMachineInterface(ctx, &models.CreateVirtualMachineInterfaceRequest{
		VirtualMachineInterface: VirtualMachineInterfaceRefModel,
	})
	VirtualMachineInterfaceRefModel.UUID = VirtualMachineInterfaceRefUUID2
	VirtualMachineInterfaceRefModel.FQName = []string{"test", VirtualMachineInterfaceRefUUID2}
	_, err = db.CreateVirtualMachineInterface(ctx, &models.CreateVirtualMachineInterfaceRequest{
		VirtualMachineInterface: VirtualMachineInterfaceRefModel,
	})
	if err != nil {
		t.Fatal("ref create failed", err)
	}
	VirtualMachineInterfaceCreateRef = append(VirtualMachineInterfaceCreateRef,
		&models.LogicalInterfaceVirtualMachineInterfaceRef{UUID: VirtualMachineInterfaceRefUUID, To: []string{"test", VirtualMachineInterfaceRefUUID}})
	VirtualMachineInterfaceCreateRef = append(VirtualMachineInterfaceCreateRef,
		&models.LogicalInterfaceVirtualMachineInterfaceRef{UUID: VirtualMachineInterfaceRefUUID2, To: []string{"test", VirtualMachineInterfaceRefUUID2}})
	model.VirtualMachineInterfaceRefs = VirtualMachineInterfaceCreateRef

	//create project to which resource is shared
	projectModel := models.MakeProject()

	projectModel.UUID = uuid.NewV4().String()
	projectModel.FQName = []string{"default-domain-test", projectModel.UUID}
	projectModel.Perms2.Owner = "admin"

	var createShare []*models.ShareType
	createShare = append(createShare, &models.ShareType{Tenant: "default-domain-test:" + projectModel.UUID, TenantAccess: 7})
	model.Perms2.Share = createShare

	_, err = db.CreateProject(ctx, &models.CreateProjectRequest{
		Project: projectModel,
	})
	if err != nil {
		t.Fatal("project create failed", err)
	}

	_, err = db.CreateLogicalInterface(ctx,
		&models.CreateLogicalInterfaceRequest{
			LogicalInterface: model,
		})

	if err != nil {
		t.Fatal("create failed", err)
	}

	response, err := db.ListLogicalInterface(ctx, &models.ListLogicalInterfaceRequest{
		Spec: &models.ListSpec{Limit: 1,
			Filters: []*models.Filter{
				&models.Filter{
					Key:    "uuid",
					Values: []string{model.UUID},
				},
			},
		}})
	if err != nil {
		t.Fatal("list failed", err)
	}
	if len(response.LogicalInterfaces) != 1 {
		t.Fatal("expected one element", err)
	}

	ctxDemo := context.WithValue(ctx, "auth", common.NewAuthContext("default", "demo", "demo", []string{}))
	_, err = db.DeleteLogicalInterface(ctxDemo,
		&models.DeleteLogicalInterfaceRequest{
			ID: model.UUID},
	)
	if err == nil {
		t.Fatal("auth failed")
	}

	_, err = db.CreateLogicalInterface(ctx,
		&models.CreateLogicalInterfaceRequest{
			LogicalInterface: model})
	if err == nil {
		t.Fatal("Raise Error On Duplicate Create failed", err)
	}

	_, err = db.DeleteLogicalInterface(ctx,
		&models.DeleteLogicalInterfaceRequest{
			ID: model.UUID})
	if err != nil {
		t.Fatal("delete failed", err)
	}

	_, err = db.GetLogicalInterface(ctx, &models.GetLogicalInterfaceRequest{
		ID: model.UUID})
	if err == nil {
		t.Fatal("expected not found error")
	}

	//Delete the project created for sharing
	_, err = db.DeleteProject(ctx, &models.DeleteProjectRequest{
		ID: projectModel.UUID})
	if err != nil {
		t.Fatal("delete project failed", err)
	}
	return
}