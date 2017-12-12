package models

// TagType

import "encoding/json"

// TagType
type TagType struct {
	TagTypeID   U16BitHexInt   `json:"tag_type_id"`
	IDPerms     *IdPermsType   `json:"id_perms"`
	UUID        string         `json:"uuid"`
	ParentUUID  string         `json:"parent_uuid"`
	ParentType  string         `json:"parent_type"`
	FQName      []string       `json:"fq_name"`
	DisplayName string         `json:"display_name"`
	Annotations *KeyValuePairs `json:"annotations"`
	Perms2      *PermType2     `json:"perms2"`
}

// String returns json representation of the object
func (model *TagType) String() string {
	b, _ := json.Marshal(model)
	return string(b)
}

// MakeTagType makes TagType
func MakeTagType() *TagType {
	return &TagType{
		//TODO(nati): Apply default
		DisplayName: "",
		Annotations: MakeKeyValuePairs(),
		Perms2:      MakePermType2(),
		UUID:        "",
		ParentUUID:  "",
		ParentType:  "",
		FQName:      []string{},
		IDPerms:     MakeIdPermsType(),
		TagTypeID:   MakeU16BitHexInt(),
	}
}

// InterfaceToTagType makes TagType from interface
func InterfaceToTagType(iData interface{}) *TagType {
	data := iData.(map[string]interface{})
	return &TagType{
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
		ParentUUID: data["parent_uuid"].(string),

		//{"type":"string"}
		ParentType: data["parent_type"].(string),

		//{"type":"string"}
		TagTypeID: InterfaceToU16BitHexInt(data["tag_type_id"]),

		//{"description":"Internal Tag type ID                  coded on 16 bits where the first 255 IDs are reserved                  and pre-defined. Users (principally cloud admin) can define                  arbitrary types but its automatically shared to all project as                  it is a global resource.","type":"string"}
		IDPerms: InterfaceToIdPermsType(data["id_perms"]),

		//{"type":"object","properties":{"created":{"type":"string"},"creator":{"type":"string"},"description":{"type":"string"},"enable":{"type":"boolean"},"last_modified":{"type":"string"},"permissions":{"type":"object","properties":{"group":{"type":"string"},"group_access":{"type":"integer","minimum":0,"maximum":7},"other_access":{"type":"integer","minimum":0,"maximum":7},"owner":{"type":"string"},"owner_access":{"type":"integer","minimum":0,"maximum":7}}},"user_visible":{"type":"boolean"}}}

	}
}

// InterfaceToTagTypeSlice makes a slice of TagType from interface
func InterfaceToTagTypeSlice(data interface{}) []*TagType {
	list := data.([]interface{})
	result := MakeTagTypeSlice()
	for _, item := range list {
		result = append(result, InterfaceToTagType(item))
	}
	return result
}

// MakeTagTypeSlice() makes a slice of TagType
func MakeTagTypeSlice() []*TagType {
	return []*TagType{}
}