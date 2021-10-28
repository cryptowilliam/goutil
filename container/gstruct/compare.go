package gstruct

import "github.com/r3labs/diff"

// compare map / structures ...
func Equal(a, b interface{}) (bool, error) {
	changelog, err := diff.Diff(a, b)
	if err != nil {
		return false, err
	}
	return len(([]diff.Change)(changelog)) == 0, nil
}
