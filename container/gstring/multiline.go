package gstring

import "strings"

func GetFirstLines(src string, count int) string {
	if count <= 0 {
		return ""
	}
	sa := strings.Split(src, "\n")
	if count >= len(sa) {
		return src
	}
	sa = sa[:count-1]
	return strings.Join(sa, "\n")
}

func GetLastLines(src string, count int) string {
	if count <= 0 {
		return ""
	}
	sa := strings.Split(src, "\n")
	if count >= len(sa) {
		return src
	}
	sa = sa[len(sa)-count:]
	return strings.Join(sa, "\n")
}

func RemoveFirstLines(src string, count int) string {
	if count <= 0 {
		return src
	}
	sa := strings.Split(src, "\n")
	if len(sa) <= count {
		return ""
	}
	sa = sa[count:]
	return strings.Join(sa, "\n")
}

func RemoveLastLines(src string, count int) string {
	if count <= 0 {
		return src
	}
	sa := strings.Split(src, "\n")
	if count >= len(sa) {
		return ""
	}
	sa = sa[:len(sa)-count]
	return strings.Join(sa, "\n")
}
