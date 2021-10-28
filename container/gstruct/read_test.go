package gstruct

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDemo(t *testing.T) {
	type Server struct {
		Name        string `json:"name,omitempty"`
		ID          int
		Enabled     bool
		users       []string // not exported
		http.Server          // embedded
	}

	server := &Server{
		Name:    "gopher",
		ID:      123456,
		Enabled: true,
	}
	// Convert a struct to a map[string]interface{}
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	m := ToMap(server)
	fmt.Println(m)

	// Convert the values of a struct to a []interface{}
	// => ["gopher", 123456, true]
	v := MemberValues(server)
	fmt.Println(v)

	// Convert the names of a struct to a []string
	// (see "Names methods" for more info about fields)
	n := MemberNames(server)
	fmt.Println(n)

	// Convert the values of a struct to a []*Field
	// (see "Field methods" for more info about fields)
	f := Members(server)
	fmt.Println(f[0])

	// Return the struct name => "Server"
	n2 := StructName(server)
	fmt.Println(n2)

	// Check if any field of a struct is initialized or not.
	h := HasZero(server)
	fmt.Println(h)

	// Check if all fields of a struct is initialized or not.
	z := IsZero(server)
	fmt.Println(z)

	// Check if server is a struct or a pointer to struct
	i := IsStruct(server)
	fmt.Println(i)
}
