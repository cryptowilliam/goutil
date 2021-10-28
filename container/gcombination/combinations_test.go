package gcombination

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAllString(t *testing.T) {
	tt := []struct {
		name string
		in   []string
		out  [][]string
	}{
		{
			name: "Empty slice",
			in:   []string{},
			out:  nil,
		},
		{
			name: "Single item",
			in:   []string{"A"},
			out: [][]string{
				{"A"},
			},
		},
		{
			name: "Two items",
			in:   []string{"A", "B"},
			out: [][]string{
				{"A"},
				{"B"},
				{"A", "B"},
			},
		},
		{
			name: "Three items",
			in:   []string{"A", "B", "C"},
			out: [][]string{
				{"A"},
				{"B"},
				{"A", "B"},
				{"C"},
				{"A", "C"},
				{"B", "C"},
				{"A", "B", "C"},
			},
		},
		{
			name: "Four items",
			in:   []string{"A", "B", "C", "D"},
			out: [][]string{
				{"A"},
				{"B"},
				{"A", "B"},
				{"C"},
				{"A", "C"},
				{"B", "C"},
				{"A", "B", "C"},
				{"D"},
				{"A", "D"},
				{"B", "D"},
				{"A", "B", "D"},
				{"C", "D"},
				{"A", "C", "D"},
				{"B", "C", "D"},
				{"A", "B", "C", "D"},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := AllString(tc.in)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("error: \nreturn:\t%v\nwant:\t%v", out, tc.out)
			}
		})
	}
}

func ExampleAllString() {
	combinations := AllString([]string{"A", "B", "C"})
	fmt.Println(combinations)
	// Output:
	// [[A] [B] [A B] [C] [A C] [B C] [A B C]]
}

func ExampleAllStringWithLen() {
	ss := AllStringWithLen([]string{"A", "B", "C"}, 1, 2)
	fmt.Println(ss)
	// Output:
	// [[A] [B] [A B] [C] [A C] [B C]]
}

func ExampleAllStringWithLen2() {
	ss := AllStringWithLen([]string{"A", "B", "C"}, 2, 2)
	fmt.Println(ss)
	// Output:
	// [[A B] [A C] [B C]]
}
