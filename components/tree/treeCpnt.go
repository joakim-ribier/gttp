package components

import (
	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/models"
	"github.com/rivo/tview"
)

// TreeCpnt represents the tree primitive which list the API(s)
type TreeCpnt struct {
	App   *tview.Application
	Event *models.Event

	labels map[string]string

	RootPrmt  *tview.Flex
	treeIndex int
	nodes     map[int]TreeCpntNode

	refreshMDRView func(it models.MakeRequestData)
	switchToPage   func(page string)
}

// NewTreeCpnt returns a new TreeCpnt struct
func NewTreeCpnt(app *tview.Application, ev *models.Event) *TreeCpnt {
	labels := make(map[string]string)
	labels["title"] = "API(s)"

	return &TreeCpnt{
		App:       app,
		Event:     ev,
		labels:    labels,
		treeIndex: -1,
		nodes:     make(map[int]TreeCpntNode),
	}
}

// Make makes the tree (home made) component
func (cpnt *TreeCpnt) Make(refreshMDRView func(it models.MakeRequestData), switchToPage func(page string)) *tview.Flex {
	cpnt.RootPrmt = tview.NewFlex().SetDirection(tview.FlexRow)
	cpnt.RootPrmt.SetBorder(false)
	cpnt.RootPrmt.SetBorderPadding(0, 0, 0, 0)

	cpnt.refreshMDRView = refreshMDRView
	cpnt.switchToPage = switchToPage

	titleTextView := tview.NewTextView().SetText(cpnt.labels["title"])
	cpnt.RootPrmt.AddItem(titleTextView, 2, 0, false)

	cpnt.RootPrmt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			cpnt.pressKeyDown()
		case tcell.KeyUp:
			cpnt.pressKeyUp()
		}
		return event
	})

	return cpnt.RootPrmt
}

// removeAll removes all children (prmt)
func (cpnt *TreeCpnt) removeAll() {
	for _, node := range cpnt.nodes {
		cpnt.RootPrmt.RemoveItem(node.textView)
	}
}

func (cpnt *TreeCpnt) pressKeyDown() {
	if cpnt.treeIndex >= (len(cpnt.nodes) - 1) {
		cpnt.treeIndex = 0
	} else {
		cpnt.treeIndex = cpnt.treeIndex + 1
	}
	cpnt.selectNode(cpnt.treeIndex)
}

func (cpnt *TreeCpnt) pressKeyUp() {
	if cpnt.treeIndex <= 0 {
		cpnt.treeIndex = len(cpnt.nodes) - 1
	} else {
		cpnt.treeIndex = cpnt.treeIndex - 1
	}
	cpnt.selectNode(cpnt.treeIndex)
}

func (cpnt *TreeCpnt) selectNode(index int) {
	node := cpnt.refreshNodeText(index)
	it, error := cpnt.Event.GetOutput().Find(node.method.String(), node.url.String())
	if error == nil {
		cpnt.refreshMDRView(it)
		cpnt.switchToPage("RequestExpertModeViewPage")
	}
}

func (cpnt *TreeCpnt) refreshNodeText(index int) TreeCpntNode {
	// initialize all nodes
	for _, node := range cpnt.nodes {
		node.textView.SetText(node.label)
	}

	// selected node
	node := cpnt.nodes[index]
	node.textView.SetText(string(9658) + " " + node.label)

	return node
}

// RefreshWithPattern refreshes the tree data with specific pattern
func (cpnt *TreeCpnt) RefreshWithPattern(pattern string) {
	cpnt.removeAll()

	addSetInputCaptureCallback := func(prmt *tview.TextView) {
		prmt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyDown:
				cpnt.pressKeyDown()
			case tcell.KeyUp:
				cpnt.pressKeyUp()
			}
			return event
		})
	}

	index := -1
	for key, slice := range cpnt.Event.GetOutput().SortDataByProject() {
		// Add 'project' new node
		parentNodeLabel := cpnt.formatParentNodeLabel(key)
		textView := tview.NewTextView().SetDynamicColors(true).SetText(parentNodeLabel)
		addSetInputCaptureCallback(textView)

		cpnt.RootPrmt.AddItem(textView, 1, 0, true)

		index++
		cpnt.nodes[index] = NewTreeCpntNode(textView, parentNodeLabel, "", "")
		for _, data := range slice {
			// Add 'request' new child node
			value := data.TreeFormat(pattern)
			childNodePrmt := tview.NewTextView().SetDynamicColors(true).SetText(value)

			addSetInputCaptureCallback(childNodePrmt)

			cpnt.RootPrmt.AddItem(childNodePrmt, 1, 0, true)

			index++
			cpnt.nodes[index] = NewTreeCpntNode(childNodePrmt, value, data.Method, data.URL)
		}
	}
}

func (cpnt *TreeCpnt) formatParentNodeLabel(value string) string {
	return "[black:white] " + value + " "
}

// Refresh refreshes the tree data
func (cpnt *TreeCpnt) Refresh() {
	cpnt.RefreshWithPattern(cpnt.Event.GetConfig().Pattern)
}
