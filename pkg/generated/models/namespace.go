package models

// Namespace

import "encoding/json"

// Namespace
type Namespace struct {
	NamespaceCidr *SubnetType    `json:"namespace_cidr"`
	Perms2        *PermType2     `json:"perms2"`
	ParentUUID    string         `json:"parent_uuid"`
	ParentType    string         `json:"parent_type"`
	IDPerms       *IdPermsType   `json:"id_perms"`
	UUID          string         `json:"uuid"`
	FQName        []string       `json:"fq_name"`
	DisplayName   string         `json:"display_name"`
	Annotations   *KeyValuePairs `json:"annotations"`
}

// String returns json representation of the object
func (model *Namespace) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakeNamespace makes Namespace
func MakeNamespace() *Namespace {
	return &Namespace{
		//TODO(nati): Apply default
		UUID:          "",
		FQName:        []string{},
		DisplayName:   "",
		Annotations:   MakeKeyValuePairs(),
		NamespaceCidr: MakeSubnetType(),
		Perms2:        MakePermType2(),
		ParentUUID:    "",
		ParentType:    "",
		IDPerms:       MakeIdPermsType(),
	}
}

// InterfaceToNamespace makes Namespace from interface
func InterfaceToNamespace(iData interface{}) *Namespace {
	data := iData.(map[string]interface{})
	return &Namespace{
		NamespaceCidr: InterfaceToSubnetType(data["namespace_cidr"]),

		//{"description":"All networks in this namespace belong to this list of Prefixes. Not implemented.","type":"object","properties":{"ip_prefix":{"type":"string"},"ip_prefix_len":{"type":"integer"}}}
		Perms2: InterfaceToPermType2(data["perms2"]),

		//{"type":"object","properties":{"global_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7},"share":{"type":"array","item":{"type":"object","properties":{"tenant":{"type":"string"},"tenant_access":{"type":"integer","minimum":0,"maximum":7}}}}}}
		ParentUUID: data["parent_uuid"].(string),

		//{"type":"string"}
		ParentType: data["parent_type"].(string),

		//{"type":"string"}
		IDPerms: InterfaceToIdPermsType(data["id_perms"]),

		//{"type":"object","properties":{"created":{"type":"string"},"creator":{"type":"string"},"description":{"type":"string"},"enable":{"type":"boolean"},"last_modified":{"type":"string"},"permissions":{"type":"object","properties":{"group":{"type":"string"},"group_access":{"type":"integer","minimum":0,"maximum":7},"other_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7}}},"user_visible":{"type":"boolean"}}}
		UUID: data["uuid"].(string),

		//{"type":"string"}
		FQName: data["fq_name"].([]string),

		//{"type":"array","item":{"type":"string"}}
		DisplayName: data["display_name"].(string),

		//{"type":"string"}
		Annotations: InterfaceToKeyValuePairs(data["annotations"]),

		//{"type":"object","properties":{"key_value_pair":{"type":"array","item":{"type":"object","properties":{"key":{"type":"string"},"value":{"type":"string"}}}}}}

	}
}

// InterfaceToNamespaceSlice makes a slice of Namespace from interface
func InterfaceToNamespaceSlice(data interface{}) []*Namespace {
	list := data.([]interface{})
	result := MakeNamespaceSlice()
	for _, item := range list {
		result = append(result, InterfaceToNamespace(item))
	}
	return result
}

// MakeNamespaceSlice() makes a slice of Namespace
func MakeNamespaceSlice() []*Namespace {
	return []*Namespace{}
}