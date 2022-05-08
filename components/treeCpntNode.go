package components

import (
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/rivo/tview"
)

// TreeCpntNode contains data for tree cpnt
type TreeCpntNode struct {
	textView *tview.TextView
	label    string
	method   types.Method
	url      types.URL
}

// NewTreeCpntNode creates new TreeCpntNode struct
func NewTreeCpntNode(prmt *tview.TextView, label string, method types.Method, url types.URL) TreeCpntNode {
	return TreeCpntNode{
		textView: prmt,
		label:    label,
		method:   method,
		url:      url,
	}
}
