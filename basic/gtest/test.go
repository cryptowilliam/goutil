package gtest

import "testing"

type (
	Case struct {
		inputs  []interface{}
		expects []interface{}
	}

	CaseReadOnly struct {
		Inputs  []interface{}
		Expects []interface{}
	}

	CaseList struct {
		items []Case
	}
)

func NewCaseList() *CaseList {
	return &CaseList{}
}

func (cl *CaseList) New() *Case {
	cl.items = append(cl.items, Case{})
	return &cl.items[len(cl.items)-1]
}

func (cl *CaseList) Get() []CaseReadOnly {
	var r []CaseReadOnly
	for _, v := range cl.items {
		r = append(r, CaseReadOnly{Inputs: v.inputs, Expects: v.expects})
	}
	return r
}

func (c *Case) Input(in interface{}) *Case {
	c.inputs = append(c.inputs, in)
	return c
}

func (c *Case) Expect(expect interface{}) *Case {
	c.expects = append(c.expects, expect)
	return c
}

/**
testing.T 报告出错并终止
FailNow / Fatal / Fatalf

testing.T 报告出错并继续
Fail / Error / Errorf
*/

// print error and end testing
func Assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// print message and end testing
func PrintlnExit(t *testing.T, format string, args ...interface{}) {
	t.Fatalf(format, args...)
}
