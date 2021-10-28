package gstring

import (
	"math/rand"
	"strings"
	"time"
)

func ArrayEqual(a, b []string, orderSensitive bool) bool {
	return false
}

func FindMaxLength(elements []string) []string {
	if len(elements) <= 1 {
		return elements
	}

	maxlen := 0
	for _, str := range elements {
		if len(str) > maxlen {
			maxlen = len(str)
		}
	}

	rst := make([]string, 0)
	for _, val := range elements {
		if maxlen == len(val) {
			rst = append(rst, val)
		}
	}
	return rst
}

// Delete empty elements - "".
func RemoveEmpty(elements []string) []string {
	var res []string
	for _, v := range elements {
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	return res
}

// Delete elements composed entirely of spaces, like "  ".
func RemoveEntirelySpaces(elements []string) []string {
	var res []string
	for _, v := range elements {
		if strings.Replace(v, " ", "", -1) == "" {
			continue
		}
		res = append(res, v)
	}
	return res
}

// Delete empty elements - "".
func RemoveNull(elements []string) []string {
	var res []string
	for _, v := range elements {
		if v == "" {
			continue
		}
		res = append(res, v)
	}
	return res
}

// Delete elements start with 'toFind'.
func RemoveStartWith(elements []string, toFind string) []string {
	var res []string
	for _, v := range elements {
		if StartWith(v, toFind) {
			continue
		}
		res = append(res, v)
	}
	return res
}

func RemoveDuplicate(elements []string) []string {
	if len(elements) <= 1 {
		return elements
	}

	// another way to initialize map
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func ToLower(elements []string) []string {
	if len(elements) == 0 {
		return []string{}
	}

	rst := []string{}
	for _, v := range elements {
		rst = append(rst, strings.ToLower(v))
	}
	return rst
}

func ToUpper(elements []string) []string {
	if len(elements) == 0 {
		return []string{}
	}

	rst := []string{}
	for _, v := range elements {
		rst = append(rst, strings.ToUpper(v))
	}
	return rst
}

func RemoveByValue(elements []string, toRemove string) []string {
	if len(elements) == 0 {
		return nil
	}

	result := make([]string, 0)
	for _, val := range elements {
		if toRemove != val {
			result = append(result, val)
		}
	}
	return result
}

func RemoveByValues(elements, toRemove []string) []string {
	if len(elements) == 0 {
		return nil
	}
	if len(toRemove) == 0 {
		return elements
	}

	result := make([]string, 0)
	for _, val := range elements {
		if CountByValue(toRemove, val) <= 0 {
			result = append(result, val)
		}
	}
	return result
}

func CountByValue(elements []string, toFind string) int {
	var result = 0
	for _, val := range elements {
		if toFind == val {
			result++
		}
	}
	return result
}

func Contains(elements []string, toFind string) bool {
	for _, val := range elements {
		if toFind == val {
			return true
		}
	}
	return false
}

func ContainsNotCaseSensitive(elements []string, toFind string) bool {
	for _, val := range elements {
		if strings.ToUpper(toFind) == strings.ToUpper(val) {
			return true
		}
	}
	return false
}

// Random sort
func Shuffle(elements []string) []string {
	final := make([]string, len(elements))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(elements))

	for i, v := range perm {
		final[v] = elements[i]
	}

	return final
}

func Merge(elements1, elements2 []string) []string {
	if elements1 == nil {
		return elements2
	}
	if elements2 == nil {
		return elements1
	}
	result := elements1
	for i := range elements2 {
		result = append(result, elements2[i])
	}
	return result
}

// 插入到前面
func Prepend(slice []string, elems ...string) []string {
	slice = append(elems, slice...)
	return slice
}

func Equal(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}

func IndexInArray(ss []string, tofind string) int {
	for i := 0; i < len(ss); i++ {
		if ss[i] == tofind {
			return i
		}
	}
	return -1
}
