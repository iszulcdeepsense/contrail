package models

// ServiceApplianceSet

import "encoding/json"

// ServiceApplianceSet
type ServiceApplianceSet struct {
	IDPerms                       *IdPermsType   `json:"id_perms"`
	ServiceApplianceHaMode        string         `json:"service_appliance_ha_mode"`
	ServiceApplianceDriver        string         `json:"service_appliance_driver"`
	ParentType                    string         `json:"parent_type"`
	DisplayName                   string         `json:"display_name"`
	Annotations                   *KeyValuePairs `json:"annotations"`
	Perms2                        *PermType2     `json:"perms2"`
	UUID                          string         `json:"uuid"`
	ServiceApplianceSetProperties *KeyValuePairs `json:"service_appliance_set_properties"`
	ParentUUID                    string         `json:"parent_uuid"`
	FQName                        []string       `json:"fq_name"`

	ServiceAppliances []*ServiceAppliance `json:"service_appliances"`
}

// String returns json representation of the object
func (model *ServiceApplianceSet) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakeServiceApplianceSet makes ServiceApplianceSet
func MakeServiceApplianceSet() *ServiceApplianceSet {
	return &ServiceApplianceSet{
		//TODO(nati): Apply default
		Annotations: MakeKeyValuePairs(),
		Perms2:      MakePermType2(),
		UUID:        "",
		ServiceApplianceSetProperties: MakeKeyValuePairs(),
		ParentUUID:                    "",
		FQName:                        []string{},
		DisplayName:                   "",
		ServiceApplianceHaMode:        "",
		ServiceApplianceDriver:        "",
		ParentType:                    "",
		IDPerms:                       MakeIdPermsType(),
	}
}

// InterfaceToServiceApplianceSet makes ServiceApplianceSet from interface
func InterfaceToServiceApplianceSet(iData interface{}) *ServiceApplianceSet {
	data := iData.(map[string]interface{})
	return &ServiceApplianceSet{
		ServiceApplianceDriver: data["service_appliance_driver"].(string),

		//{"description":"Name of the provider driver for this service appliance set.","type":"string"}
		ParentType: data["parent_type"].(string),

		//{"type":"string"}
		IDPerms: InterfaceToIdPermsType(data["id_perms"]),

		//{"type":"object","properties":{"created":{"type":"string"},"creator":{"type":"string"},"description":{"type":"string"},"enable":{"type":"boolean"},"last_modified":{"type":"string"},"permissions":{"type":"object","properties":{"group":{"type":"string"},"group_access":{"type":"integer","minimum":0,"maximum":7},"other_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7}}},"user_visible":{"type":"boolean"}}}
		ServiceApplianceHaMode: data["service_appliance_ha_mode"].(string),

		//{"description":"High availability mode for the service appliance set, active-active or active-backup.","type":"string"}
		ParentUUID: data["parent_uuid"].(string),

		//{"type":"string"}
		FQName: data["fq_name"].([]string),

		//{"type":"array","item":{"type":"string"}}
		DisplayName: data["display_name"].(string),

		//{"type":"string"}
		Annotations: InterfaceToKeyValuePairs(data["annotations"]),

		//{"type":"object","properties":{"key_value_pair":{"type":"array","item":{"type":"object","properties":{"key":{"type":"string"},"value":{"type":"string"}}}}}}
		Perms2: InterfaceToPermType2(data["perms2"]),

		//{"type":"object","properties":{"global_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7},"share":{"type":"array","item":{"type":"object","properties":{"tenant":{"type":"string"},"tenant_access":{"type":"integer","minimum":0,"maximum":7}}}}}}
		UUID: data["uuid"].(string),

		//{"type":"string"}
		ServiceApplianceSetProperties: InterfaceToKeyValuePairs(data["service_appliance_set_properties"]),

		//{"description":"List of Key:Value pairs that are used by the provider driver and opaque to system.","type":"object","properties":{"key_value_pair":{"type":"array","item":{"type":"object","properties":{"key":{"type":"string"},"value":{"type":"string"}}}}}}

	}
}

// InterfaceToServiceApplianceSetSlice makes a slice of ServiceApplianceSet from interface
func InterfaceToServiceApplianceSetSlice(data interface{}) []*ServiceApplianceSet {
	list := data.([]interface{})
	result := MakeServiceApplianceSetSlice()
	for _, item := range list {
		result = append(result, InterfaceToServiceApplianceSet(item))
	}
	return result
}

// MakeServiceApplianceSetSlice() makes a slice of ServiceApplianceSet
func MakeServiceApplianceSetSlice() []*ServiceApplianceSet {
	return []*ServiceApplianceSet{}
}