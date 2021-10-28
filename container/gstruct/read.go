package gstruct

import (
	"github.com/fatih/structs"
)

// Convert a struct to a map[string]interface{}
// => {"Name":"gopher", "ID":123456, "Enabled":true}
func ToMap(x interface{}) map[string]interface{} {
	return structs.Map(x)
}

// Convert the names of a struct to a []string
// (see "Names methods" for more info about fields)
func MemberNames(x interface{}) []string {
	return structs.Names(x)
}

func MemberValues(x interface{}) []interface{} {
	return structs.Values(x)
}

// Convert the values of a struct to a []*Field
// (see "Field methods" for more info about fields)
func Members(x interface{}) []*structs.Field {
	return structs.Fields(x)
}

// Return the struct name => "Server"
func StructName(x interface{}) string {
	return structs.Name(x)
}

// Check if any field of a struct is initialized or not.
func HasZero(x interface{}) bool {
	return structs.HasZero(x)
}

// Check if all fields of a struct is initialized or not.
func IsZero(x interface{}) bool {
	return structs.IsZero(x)
}

// Check if server is a struct or a pointer to struct
func IsStruct(x interface{}) bool {
	return structs.IsStruct(x)
}
