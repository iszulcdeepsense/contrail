package models

// VPNGroup

import "encoding/json"

// VPNGroup
type VPNGroup struct {
	ProvisioningStartTime     string         `json:"provisioning_start_time"`
	ParentUUID                string         `json:"parent_uuid"`
	ProvisioningProgress      int            `json:"provisioning_progress"`
	Annotations               *KeyValuePairs `json:"annotations"`
	ParentType                string         `json:"parent_type"`
	FQName                    []string       `json:"fq_name"`
	ProvisioningState         string         `json:"provisioning_state"`
	ProvisioningLog           string         `json:"provisioning_log"`
	ProvisioningProgressStage string         `json:"provisioning_progress_stage"`
	Type                      string         `json:"type"`
	DisplayName               string         `json:"display_name"`
	IDPerms                   *IdPermsType   `json:"id_perms"`
	Perms2                    *PermType2     `json:"perms2"`
	UUID                      string         `json:"uuid"`

	LocationRefs []*VPNGroupLocationRef `json:"location_refs"`
}

// VPNGroupLocationRef references each other
type VPNGroupLocationRef struct {
	UUID string   `json:"uuid"`
	To   []string `json:"to"` //FQDN

}

// String returns json representation of the object
func (model *VPNGroup) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakeVPNGroup makes VPNGroup
func MakeVPNGroup() *VPNGroup {
	return &VPNGroup{
		//TODO(nati): Apply default
		Type:                      "",
		DisplayName:               "",
		FQName:                    []string{},
		ProvisioningState:         "",
		ProvisioningLog:           "",
		ProvisioningProgressStage: "",
		Perms2:                MakePermType2(),
		UUID:                  "",
		IDPerms:               MakeIdPermsType(),
		ParentUUID:            "",
		ProvisioningProgress:  0,
		ProvisioningStartTime: "",
		Annotations:           MakeKeyValuePairs(),
		ParentType:            "",
	}
}

// InterfaceToVPNGroup makes VPNGroup from interface
func InterfaceToVPNGroup(iData interface{}) *VPNGroup {
	data := iData.(map[string]interface{})
	return &VPNGroup{
		IDPerms: InterfaceToIdPermsType(data["id_perms"]),

		//{"type":"object","properties":{"created":{"type":"string"},"creator":{"type":"string"},"description":{"type":"string"},"enable":{"type":"boolean"},"last_modified":{"type":"string"},"permissions":{"type":"object","properties":{"group":{"type":"string"},"group_access":{"type":"integer","minimum":0,"maximum":7},"other_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7}}},"user_visible":{"type":"boolean"}}}
		Perms2: InterfaceToPermType2(data["perms2"]),

		//{"type":"object","properties":{"global_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7},"share":{"type":"array","item":{"type":"object","properties":{"tenant":{"type":"string"},"tenant_access":{"type":"integer","minimum":0,"maximum":7}}}}}}
		UUID: data["uuid"].(string),

		//{"type":"string"}
		ProvisioningStartTime: data["provisioning_start_time"].(string),

		//{"title":"Time provisioning started","default":"","type":"string","permission":["create","update"]}
		ParentUUID: data["parent_uuid"].(string),

		//{"type":"string"}
		ProvisioningProgress: data["provisioning_progress"].(int),

		//{"title":"Provisioning Progress","default":0,"type":"integer","permission":["create","update"]}
		Annotations: InterfaceToKeyValuePairs(data["annotations"]),

		//{"type":"object","properties":{"key_value_pair":{"type":"array","item":{"type":"object","properties":{"key":{"type":"string"},"value":{"type":"string"}}}}}}
		ParentType: data["parent_type"].(string),

		//{"type":"string"}
		FQName: data["fq_name"].([]string),

		//{"type":"array","item":{"type":"string"}}
		ProvisioningState: data["provisioning_state"].(string),

		//{"title":"Provisioning Status","default":"CREATED","type":"string","permission":["create","update"],"enum":["CREATED","IN_CREATE_PROGRESS","UPDATED","IN_UPDATE_PROGRESS","DELETED","IN_DELETE_PROGRESS","ERROR"]}
		ProvisioningLog: data["provisioning_log"].(string),

		//{"title":"Provisioning Log","default":"","type":"string","permission":["create","update"]}
		ProvisioningProgressStage: data["provisioning_progress_stage"].(string),

		//{"title":"Provisioning Progress Stage","default":"","type":"string","permission":["create","update"]}
		Type: data["type"].(string),

		//{"title":"VPN Type","description":"Type of VPN","default":"ipsec","type":"string"}
		DisplayName: data["display_name"].(string),

		//{"type":"string"}

	}
}

// InterfaceToVPNGroupSlice makes a slice of VPNGroup from interface
func InterfaceToVPNGroupSlice(data interface{}) []*VPNGroup {
	list := data.([]interface{})
	result := MakeVPNGroupSlice()
	for _, item := range list {
		result = append(result, InterfaceToVPNGroup(item))
	}
	return result
}

// MakeVPNGroupSlice() makes a slice of VPNGroup
func MakeVPNGroupSlice() []*VPNGroup {
	return []*VPNGroup{}
}