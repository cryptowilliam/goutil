package gstring

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"sort"
	"strings"
	"unicode"
)

func HeadToLowerASCII(s string, n int) string {
	if n <= 0 {
		return s
	}
	if len(s) <= n {
		return strings.ToLower(s)
	}

	head := SubstrByLenAscii(s, n)
	tail := ""
	if n < len(s) {
		tail = s[n:]
	}
	return strings.ToLower(head) + tail
}

func EqualLower(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

func EqualUpper(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

// 以UTF8为单元的Len函数
func LenUTF8(s string) int {
	return strings.Count(s, "") - 1
}

// 以Byte为单元的Len函数
func LenAscii(s string) int {
	return len(s)
}

// []rune的StartWith函数
func StartWithRune(s, sub []rune) bool {
	for i := range sub {
		if s[i] != sub[i] {
			return false
		}
	}
	return true
}

// 以UTF8为单元的Index函数
func IndexUTF8(s, substr string) int {
	s_r := []rune(s)
	substr_r := []rune(substr)

	for i := range s_r {
		if StartWithRune(s_r[i:], substr_r) {
			return i
		}
	}

	return -1
}

func LastIndexUTF8(s, substr string) int {
	s_r := []rune(s)
	substr_r := []rune(substr)

	for i := len(s_r) - 1; i >= 0; i-- {
		if StartWithRune(s_r[i:], substr_r) {
			return i
		}
	}

	return -1
}

// 以Byte为单元的Index函数
func IndexAscii(s, substr string) int {
	return strings.Index(s, substr)
}

/*
   WARNING
   strings.Trim(src, " ") may be invalid,
   because space character is not only the utf char " ".
   You can using strings.TrimSpace instead of it and check the source code.
*/

// Find index of sep string in src string, and start searching from given index - fromIndex
// Sample: IndexAfter("01234567893", "3", 5) = 10
func IndexAfter(s string, substr string, fromIndex int) int {
	if fromIndex <= 0 {
		return strings.Index(s, substr)
	}
	if fromIndex > (len(s) - 1) {
		return -1
	}
	rs := []rune(s)
	ss := string(rs[fromIndex:len(s)])
	result := strings.Index(ss, substr)
	if result < 0 {
		return result
	}
	return result + fromIndex
}

func IndexAfterUTF8(s string, substr string, fromIndex int) int {
	if fromIndex <= 0 {
		return IndexUTF8(s, substr)
	}
	if fromIndex > (LenUTF8(s) - 1) {
		return -1
	}
	rs := []rune(s)
	ss := string(rs[fromIndex:LenUTF8(s)])
	result := IndexUTF8(ss, substr)
	//fmt.Println("substr", substr, result)
	if result < 0 {
		return result
	}
	//fmt.Println("hia")
	return result + fromIndex
}

// 以字节为单元的Substr函数
// 如果包含中文等UTF8字符，则不可以使用此函数
func SubstrAscii(s string, begin, end int) (string, error) {
	if begin < 0 || begin > (len(s)-1) {
		return "", gerrors.New("begin index is error")
	}
	if end < 0 || end > (len(s)-1) {
		return "", gerrors.New("end index is error")
	}
	if begin > end {
		return "", gerrors.New("begin is bigger than end")
	}
	rs := []rune(s)
	// NOTICE: when rune[a:b], char of index a is included, but b is NOT!
	return string(rs[begin : end+1]), nil
}

// 以字节为单元的Substr函数
// 如果begine、length过大过小，都不会报错
func TrySubstrLenAscii(s string, begin, length int) string {
	if begin < 0 {
		begin = 0
	}
	if len(s) == 0 || (begin > len(s)-1) || length <= 0 {
		return ""
	}

	// end index
	end := begin + length - 1
	if end > len(s)-1 {
		end = len(s) - 1
	}

	rs := []rune(s)
	// NOTICE: when rune[a:b], char of index a is included, but b is NOT!
	//fmt.Println(len(src), begin, end+1)
	//fmt.Println("[", src, "]")
	return string(rs[begin : end+1])
}

// 以UTF8为单元的Substr函数
func TrySubstrLenUTF8(s string, begin, length int) string {
	if begin < 0 {
		begin = 0
	}
	if len(s) == 0 || (begin > LenUTF8(s)-1) || length <= 0 {
		return ""
	}

	// end index
	end := begin + length - 1
	if end > LenUTF8(s)-1 {
		end = LenUTF8(s) - 1
	}

	rs := []rune(s)
	// NOTICE: when rune[a:b], char of index a is included, but b is NOT!
	//fmt.Println(len(src), begin, end+1)
	//fmt.Println("[", src, "]")
	return string(rs[begin : end+1])
}

// Get substring of source string
// sample:
// SubstrAscii("ab)cd(efgh)ij)kl", "(", ")", true, true, true, true) == "efgh"
func SubstrBetween(s string, begin string, end string,
	searchBeginFromLeft bool, searchEndFromLeft bool,
	includBegin bool, includEnd bool) (string, error) {
	if len(s) == 0 || len(begin) == 0 || len(end) == 0 {
		return "", gerrors.New("SubstrAscii has invalid input")
	}

	var a, b int
	if searchBeginFromLeft {
		a = strings.Index(s, begin)
	} else {
		a = strings.LastIndex(s, begin)
	}
	if a < 0 {
		return "", gerrors.New("can't find begin string")
	}
	if searchEndFromLeft {
		b = IndexAfter(s, end, a+len(begin))
	} else {
		b = strings.LastIndex(s, end)
	}
	if b < 0 {
		return "", gerrors.New("can't find end string")
	}
	if !includBegin {
		a += len(begin)
	}
	if includEnd {
		b += len(end)
	}
	if a > b {
		return "", gerrors.New("begin index is bigger than end")
	}

	return s[a:b], nil
}

func SubstrBetweenUTF8(s string, begin string, end string,
	searchBeginFromLeft bool, searchEndFromLeft bool,
	includBegin bool, includEnd bool) (string, error) {
	if LenUTF8(s) == 0 || LenUTF8(begin) == 0 || LenUTF8(end) == 0 {
		return "", gerrors.New("SubstrByteUTF8 has invalid input")
	}

	var a, b int
	if searchBeginFromLeft {
		a = IndexUTF8(s, begin)
	} else {
		a = LastIndexUTF8(s, begin)
	}
	if a < 0 {
		return "", gerrors.New("cann't find begin string")
	}

	if searchEndFromLeft {
		b = IndexAfterUTF8(s, end, a+LenUTF8(begin))
	} else {
		b = LastIndexUTF8(s, end)
	}
	if b < 0 {
		return "", gerrors.New("cann't find end string")
	}
	if !includBegin {
		a += LenUTF8(begin)
	}
	if includEnd {
		b += LenUTF8(end)
	}
	if a > b {
		return "", gerrors.Errorf("begin index [%d] is bigger than end [%d]", a, b)
	}

	return string([]rune(s)[a:b]), nil
}

// 取前面一段字符串，最大长度为length
func SubstrByLenAscii(s string, length int) string {
	if len(s) == 0 || length <= 0 {
		return ""
	}
	if length >= len(s) {
		return s
	}
	return s[0:length] // 将src中从下标0到(length-1)下的元素创建为一个新的切片
}

func SubstrByLenUTF8(s string, length int) string {
	if len(s) == 0 || length <= 0 {
		return ""
	}
	if length >= LenUTF8(s) {
		return s
	}
	return string([]rune(s)[0:length]) // 将src中从下标0到(length-1)下的元素创建为一个新的切片
}

func LastSubstrByLenAscii(s string, length int) string {
	if len(s) == 0 || length <= 0 {
		return ""
	}
	if length >= len(s) {
		return s
	}
	return s[len(s)-length:]
}

func LastSubstrByLenUTF8(s string, length int) string {
	if len(s) == 0 || length <= 0 {
		return ""
	}
	if length >= LenUTF8(s) {
		return s
	}
	return string([]rune(s)[LenUTF8(s)-length:])
}

func RemoveIndex(s string, index int) string {
	if len(s) == 0 || index < 0 || index >= len(s) {
		return s
	}
	if index == 0 {
		return s[1:]
	}
	if index == len(s)-1 {
		return s[:len(s)-1]
	}
	return s[:index] + s[index+1:]
}

func RemoveHead(s string, length int) string {
	if len(s) == 0 || length >= len(s) {
		return ""
	}
	if length <= 0 {
		return s
	}
	return s[length:]
}

func RemoveTail(s string, length int) string {
	if len(s) == 0 || length >= len(s) {
		return ""
	}
	if length <= 0 {
		return s
	}
	return s[:len(s)-length] // NOTICE: end pos char is not included in return string
}

func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

func ReplaceReverse(s, old, new string, n int) string {
	s = Reverse(s)
	old = Reverse(old)
	new = Reverse(new)
	s = strings.Replace(s, old, new, n)
	return Reverse(s)
}

func ReplaceWithTags(s string, begin string, endAfterBegin string, replaceStr string, replaceTime int) (string, error) {
	if len(s) == 0 || len(begin) == 0 || len(endAfterBegin) == 0 {
		return "", gerrors.New("ReplaceWithTags has invalid input")
	}
	if replaceTime <= 0 {
		replaceTime = strings.Count(s, begin)
	}

	for i := 0; i < replaceTime; i++ {
		ss, e := SubstrBetween(s, begin, endAfterBegin, true, true, true, true)
		if e != nil {
			return s, nil
		}
		s = strings.Replace(s, ss, replaceStr, 1)
	}

	return s, nil
}

// Get count of digit character in string
func CountDigit(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsDigit(c) {
			count++
		}
	}
	return count
}

func StartWith(s, toFind string) bool {
	if len(s) == 0 || len(toFind) == 0 {
		return false
	}

	pos := strings.Index(s, toFind)
	return pos == 0
}

func EndWith(s, toFind string) bool {
	if len(s) == 0 || len(toFind) == 0 {
		return false
	}

	pos := strings.LastIndex(s, toFind)

	// pos < 0: can't find toFind
	// if don't add pos >= 0, there will be a bug if EndWith("astring", "*astring")
	return pos >= 0 && pos == len(s)-len(toFind)
}

func SplitByLen(s string, length int) []string {
	var r []string
	for len(s) >= length {
		splitLen := len(s)
		if length < splitLen {
			splitLen = length
		}
		r = append(r, s[:splitLen])
		s = s[splitLen:]
	}
	if len(s) > 0 {
		r = append(r, s)
	}
	return r
}

// SplitChunks("1234567", 3, false) == "1 234 567"
func SplitChunksAscii(s string, chunkSize int, fromLeft bool) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) || chunkSize <= 0 {
		return []string{s}
	}

	ss := []string{}

	smallSliceLen := len(s) % chunkSize
	completeSliceCount := (len(s) - smallSliceLen) / chunkSize
	if fromLeft {
		begin := 0
		end := begin + chunkSize
		for i := 0; i < completeSliceCount; i++ {
			ss = append(ss, s[begin:end])
			begin += chunkSize
			end += chunkSize
		}
		if smallSliceLen > 0 {
			begin = chunkSize * completeSliceCount
			end = begin + smallSliceLen
			substr := s[begin:end]
			ss = append(ss, substr)
		}
	} else {
		if smallSliceLen > 0 {
			substr := s[0:smallSliceLen]
			ss = append(ss, substr)
		}
		begin := smallSliceLen
		end := begin + chunkSize
		for i := 0; i < completeSliceCount; i++ {
			ss = append(ss, s[begin:end])
			begin += chunkSize
			end += chunkSize
		}
	}
	return ss
}

func IsASCII(s string) bool {
	for _, v := range []rune(s) {
		if v > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func TrimUnASCII(s string) string {
	res := []rune{}
	for _, v := range []rune(s) {
		if v <= unicode.MaxASCII {
			res = append(res, v)
		}
	}
	return string(res)
}

// trim s until s doesn't start with cutset
func TrimLeftAll(s, cutset string) string {
	if s == "" || cutset == "" {
		return s
	}

	for {
		if strings.Index(s, cutset) == 0 {
			s = strings.TrimLeft(s, cutset)
		} else {
			return s
		}
	}
}

// trim s until s doesn't end with cutset
func TrimRightAll(s, cutset string) string {
	if s == "" || cutset == "" {
		return s
	}

	for {
		if EndWith(s, cutset) {
			s = strings.TrimRight(s, cutset)
		} else {
			return s
		}
	}
}

func OnlyFirstLetterUpperCase(s string) string {
	if s == "" {
		return s
	}
	if IsASCII(s[0:1]) {
		return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
	} else {
		return strings.ToLower(s)
	}
}

type sortString struct {
	runes []rune
}

func (rs sortString) Len() int {
	return len(rs.runes)
}

func (rs sortString) Swap(i, j int) {
	rs.runes[i], rs.runes[j] = rs.runes[j], rs.runes[i]
}

func (rs sortString) Less(i, j int) bool {
	return rs.runes[i] < rs.runes[j]
}

// sort string by rune hex value
// example:
// 722abBCcA -> 277ABCabc
func SortByHex(s string) string {
	ss := sortString{
		runes: []rune(s),
	}
	sort.Sort(ss)
	return string(ss.runes)
}

func SplitRunes(s string) []rune {
	var res []rune
	for i := 0; i < len(s); i++ {
		res = append(res, rune(s[i]))
	}
	return res
}

func SplitToLines(s string) []string {
	// You have to use "\n", don't use '\n'
	// Splitting on `\n`, searches for an actual \ followed by n in the text, not the newline byte.
	return strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n")
}
