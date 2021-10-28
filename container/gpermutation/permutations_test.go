package gpermutation

import (
	"strings"
	"testing"
)

func TestKindSelectCombinations_ListAll(t *testing.T) {
	ksc := Permutations{}
	ksc.AddKind("a")
	ksc.AddKind("b", "c", "d")
	ksc.AddKind("e", "f")
	ksc.AddKind("g", "h")

	rst := ksc.ListAll()
	var rststring []string
	for i := range rst {
		item := ""
		for j := range rst[i].Items {
			item += (rst[i].Items[j]).(string)
			if j != len(rst[i].Items)-1 {
				item += ", "
			}
		}
		rststring = append(rststring, item)
	}

	rstcorrect := []string{
		"a, b, e, g",
		"a, b, e, h",
		"a, b, f, g",
		"a, b, f, h",
		"a, c, e, g",
		"a, c, e, h",
		"a, c, f, g",
		"a, c, f, h",
		"a, d, e, g",
		"a, d, e, h",
		"a, d, f, g",
		"a, d, f, h",
	}

	if strings.Join(rststring, "+") != strings.Join(rstcorrect, "+") {
		t.Errorf("Correct result is %s, but returns %s", strings.Join(rstcorrect, ", "), strings.Join(rststring, ", "))
	}
}
