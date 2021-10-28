package ghtml

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

/**
xpath samples:

find all A elements: "//a"
find all A elements that have href attribute: "//a[@href]"
find all A elements with href attribute and only return href value: "//a/@href"
find the third A element: "//a[3]"
*/

type (
	// html node
	Node struct {
		// use struct member instead of "type Node html.Node", because need to hide methods of html.Node
		raw *html.Node
	}

	// nodes list with it's own methods
	Nodes struct {
		items []*Node
	}

	// callback function of Each()
	NodeEachCallback func(i int, n *Node)
)

// build node from HTML string
func NewFromString(s string) (*Node, error) {
	node, err := htmlquery.Parse(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	return &Node{raw: node}, nil
}

// get attribute value by key
func (n *Node) Attr(attrKey string) string {
	if n == nil {
		return ""
	}
	for _, v := range n.raw.Attr {
		if v.Key == attrKey {
			return v.Val
		}
	}
	return ""
}

// get text of node
func (n *Node) Text() string {
	if n == nil {
		return ""
	}
	return htmlquery.InnerText(n.raw)
}

// query with xpath expr
func (n *Node) Query(xpathExpr string) *Nodes {
	if n == nil {
		return nil
	}
	ns := htmlquery.Find(n.raw, xpathExpr)
	res := &Nodes{}
	for _, v := range ns {
		if v == nil {
			continue
		}
		res.items = append(res.items, &Node{raw: v})
	}
	return res
}

// query nodes by element + attribute key + attribute full value
func (n *Node) QueryByElementAndFullAttrVal(element, attrKey, fullAttrVal string) *Nodes {
	if n == nil {
		return nil
	}
	ns := htmlquery.Find(n.raw, "//"+element)
	res := &Nodes{}
	for _, v := range ns {
		if v == nil {
			continue
		}
		n := &Node{raw: v}
		if n.Attr(attrKey) == fullAttrVal {
			res.items = append(res.items, n)
		}
	}
	return res
}

// query nodes by element + attribute key + attribute part value
func (n *Node) QueryByElementAndPartAttrVal(element, attrKey, partAttrVal string) *Nodes {
	if n == nil {
		return nil
	}
	ns := htmlquery.Find(n.raw, "//"+element)
	res := &Nodes{}
	for _, v := range ns {
		if v == nil {
			continue
		}
		n := &Node{raw: v}
		if strings.Contains(n.Attr(attrKey), partAttrVal) {
			res.items = append(res.items, n)
		}
	}
	return res
}

// get father node
func (n *Node) Parent() *Node {
	return &Node{raw: n.raw.Parent}
}

// get by index, return nil if index not exist
func (ns *Nodes) Get(i int) *Node {
	if ns == nil || len(ns.items) == 0 || i > ns.Len()-1 {
		return nil
	}
	return ns.items[i]
}

// get first node, return nil if index not exist
func (ns *Nodes) First() *Node {
	if ns == nil || len(ns.items) == 0 {
		return nil
	}
	return ns.items[0]
}

// get last node, return nil if index not exist
func (ns *Nodes) Last() *Node {
	if ns == nil || len(ns.items) == 0 {
		return nil
	}
	return ns.items[ns.Len()-1]
}

// get length
func (ns *Nodes) Len() int {
	if ns == nil {
		return 0
	}
	return len(ns.items)
}

// loop through all members
func (ns *Nodes) Each(cb NodeEachCallback) {
	for i := 0; i < ns.Len(); i++ {
		cb(i, ns.Get(i))
	}
}

// get attribute values of all nodes by attribute key
func (ns *Nodes) Attrs(attrKey string) []string {
	var res []string
	for i := 0; i < ns.Len(); i++ {
		res = append(res, ns.Get(i).Attr(attrKey))
	}
	return res
}

// get texts of all nodes
func (ns *Nodes) Texts() []string {
	var res []string
	for i := 0; i < ns.Len(); i++ {
		res = append(res, ns.Get(i).Text())
	}
	return res
}
