package gstring

import (
	"fmt"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	checker := NewChecker().Allow('a', 'z').AllowRune('-')
	fmt.Println(checker.Check("apple-amazon"))
	fmt.Println(checker.Check("helloWorld"))
	fmt.Println(checker.Check("1apple"))
}
