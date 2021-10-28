package gcsv

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"testing"
	"time"
)

func TestOpenCsv(t *testing.T) {
	var csvInput = []byte(`
name,age,createat,address
jacek,26,2012-04-01T15:00:00Z,beijing

john,,0001-01-01T00:00:00Z,beijing`,
	)

	type User struct {
		Name      string    `csv:"name"`
		Age       int       `csv:"age,omitempty"`
		CreatedAt time.Time `csv:"createat"`
		Address   string    `csv:"address"`
	}

	var users []User

	if err := gocsv.UnmarshalBytes(csvInput, &users); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(users)
}
