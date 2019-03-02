package components

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/models"
	"github.com/rivo/tview"
)

const treeReferencePattern = "#@#"

// TreeAPICpnt represents the tree primitive which list the API(s)
type TreeAPICpnt struct {
	App   *tview.Application
	Event *models.Event

	labels map[string]string

	treeNode *tview.TreeNode
	Tree     *tview.TreeView
}

// NewTreeAPICpnt returns a new TreeAPIList struct
func NewTreeAPICpnt(app *tview.Application, ev *models.Event) *TreeAPICpnt {
	labels := make(map[string]string)
	labels["title"] = "API(s)"
	labels["root"] = "root#@#root"

	return &TreeAPICpnt{
		App:    app,
		Event:  ev,
		labels: labels,
	}
}

// Make makes the tree component
func (cpnt *TreeAPICpnt) Make(refreshMDRView func(it models.MakeRequestData), switchToPage func(page string)) *tview.TreeView {
	cpnt.treeNode = tview.NewTreeNode(cpnt.labels["title"]).
		SetReference(cpnt.labels["root"]).
		SetColor(tcell.ColorBlue)

	cpnt.Tree = tview.NewTreeView().
		SetRoot(cpnt.treeNode).
		SetCurrentNode(cpnt.treeNode)

	cpnt.Tree.SetChangedFunc(func(node *tview.TreeNode) {
		value := node.GetReference().(string)

		s := strings.Split(value, treeReferencePattern)
		method, url := s[0], s[1]

		if it, error := cpnt.Event.GetOutput().Find(method, url); error == nil {
			refreshMDRView(it)
			switchToPage("RequestExpertModeViewPage")
		}
	})

	return cpnt.Tree
}

// RefreshWithPattern refreshes the tree data with specific pattern
func (cpnt *TreeAPICpnt) RefreshWithPattern(pattern string) {
	cpnt.treeNode.ClearChildren()
	for key, slice := range cpnt.Event.GetOutput().SortDataByProject() {
		node := tview.NewTreeNode(key).SetSelectable(false)
		cpnt.treeNode.AddChild(node)
		for _, data := range slice {
			if data.URL != "" {
				value := data.TreeFormat(pattern)
				childNode := tview.NewTreeNode(value).
					SetReference(data.Method.String() + treeReferencePattern + data.URL.String()).
					SetSelectable(true)
				node.AddChild(childNode)
			}
		}
	}
}

// Refresh refreshes the tree data
func (cpnt *TreeAPICpnt) Refresh() {
	cpnt.RefreshWithPattern(cpnt.Event.GetConfig().Pattern)
}
