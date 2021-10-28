// author - 2019: https://github.com/mxschmitt/golang-combinations
// author 2019 -: me
// Package combinations provides a method to generate all combinations out of a given string array.

/*
Background
The algorithm iterates over each number from 1 to 2^length(input), separating it by binary components and utilizes the true/false interpretation of binary 1's and 0's to extract all unique ordered combinations of the input slice.

E.g. a binary number 0011 means selecting the first and second index from the slice and ignoring the third and fourth. For input {"A", "B", "C", "D"} this signifies the combination {"A", "B"}.

For input slice {"A", "B", "C", "D"} there are 2^4 - 1 = 15 binary combinations, so mapping each bit position to a slice index and selecting the entry for binary 1 and discarding for binary 0 gives the full subset as:

1	=	0001	=>	---A	=>	{"A"}
2	=	0010	=>	--B-	=>	{"B"}
3	=	0011	=>	--BA	=>	{"A", "B"}
4	=	0100	=>	-C--	=>	{"C"}
5	=	0101	=>	-C-A	=>	{"A", "C"}
6	=	0110	=>	-CB-	=>	{"B", "C"}
7	=	0111	=>	-CBA	=>	{"A", "B", "C"}
8	=	1000	=>	D---	=>	{"D"}
9	=	1001	=>	D--A	=>	{"A", "D"}
10	=	1010	=>	D-B-	=>	{"B", "D"}
11	=	1011	=>	D-BA	=>	{"A", "B", "D"}
12	=	1100	=>	DC--	=>	{"C", "D"}
13	=	1101	=>	DC-A	=>	{"A", "C", "D"}
14	=	1110	=>	DCB-	=>	{"B", "C", "D"}
15	=	1111	=>	DCBA	=>	{"A", "B", "C", "D"}
*/

package gcombination

import (
	"github.com/cryptowilliam/goutil/container/gbit"
)

// All returns all combinations for a given interface array.
// This is essentially a powerset of the given set except that the empty set is disregarded.
func All(set []interface{}) (subsets [][]interface{}) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []interface{}

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

// All returns all combinations for a given interface array and fixed set length.
// This is essentially a powerset of the given set except that the empty set is disregarded.
func AllWithLen(set []interface{}, minLen, maxLen int) (subsets [][]interface{}) {
	if maxLen < minLen || minLen > len(set) || maxLen <= 0 {
		return nil
	}
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		bit1Count := int(gbit.Count1BitsHamming32(uint32(subsetBits)))
		if bit1Count < minLen || bit1Count > maxLen {
			continue
		}

		var subset []interface{}
		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}

func strSet2ItfSet(src []string) []interface{} {
	var r []interface{}
	for _, v := range src {
		r = append(r, v)
	}
	return r
}

func itfSs2StrSs(src [][]interface{}) [][]string {
	var r [][]string
	for _, set := range src {
		var set2 []string
		for _, val := range set {
			set2 = append(set2, val.(string))
		}
		r = append(r, set2)
	}
	return r
}

func intSet2ItfSet(src []int) []interface{} {
	var r []interface{}
	for _, v := range src {
		r = append(r, v)
	}
	return r
}

func itfSs2IntSs(src [][]interface{}) [][]int {
	var r [][]int
	for _, set := range src {
		var set2 []int
		for _, val := range set {
			set2 = append(set2, val.(int))
		}
		r = append(r, set2)
	}
	return r
}

func AllString(set []string) (subsets [][]string) { return itfSs2StrSs(All(strSet2ItfSet(set))) }
func AllInt(set []int) (subsets [][]int)          { return itfSs2IntSs(All(intSet2ItfSet(set))) }
func AllStringWithLen(set []string, minLen, maxLen int) (subsets [][]string) {
	return itfSs2StrSs(AllWithLen(strSet2ItfSet(set), minLen, maxLen))
}
func AllIntWithLen(set []int, minLen, maxLen int) (subsets [][]int) {
	return itfSs2IntSs(AllWithLen(intSet2ItfSet(set), minLen, maxLen))
}
