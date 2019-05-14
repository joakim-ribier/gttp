package views

import (
	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// SaveRequestView represents the settings view of the saving request
type SaveRequestView struct {
	App   *tview.Application
	Event *models.Event

	Labels map[string]string

	TitlePrmt  tview.Primitive
	ParentPrmt tview.Primitive
}

// NewSaveRequestView returns the new view
func NewSaveRequestView(app *tview.Application, ev *models.Event) *SaveRequestView {
	labels := make(map[string]string)

	labels["title"] = "Save Request"
	labels["project"] = "Project Name"
	labels["alias"] = "Request Alias"
	labels["save"] = "Add"

	labels["menu_1_title"] = "Set Project & Alias"
	labels["menu_1_desc"] = "Update project and request alias name"

	return &SaveRequestView{
		App:    app,
		Event:  ev,
		Labels: labels,
	}
}

// InitView build all components to display correctly the view
func (view *SaveRequestView) InitView() {
	mapMenuToFocusPrmt := make(map[string]tview.Primitive)

	// Pages for each menu content
	pages := tview.NewPages()
	pages.SetBackgroundColor(utils.BackGrayColor)
	pages.AddPage("SaveRequestProjectAliasName", view.makeForm(mapMenuToFocusPrmt), true, false)

	// Menu
	menu := view.makeMenu(pages, mapMenuToFocusPrmt)

	// Title
	titleAndMenuFlexPrmt := utils.MakeTitlePrmt(view.Labels["title"])
	titleAndMenuFlexPrmt.AddItem(menu, 0, 1, false)

	flexPrmt := tview.NewFlex()
	flexPrmt.AddItem(titleAndMenuFlexPrmt, 0, 1, false)
	flexPrmt.AddItem(tview.NewBox().SetBorder(false), 2, 0, false)
	flexPrmt.AddItem(pages, 0, 2, false)

	frame := tview.NewFrame(flexPrmt).SetBorders(0, 0, 0, 0, 0, 0)

	titleAndMenuFlexPrmt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Name() {
		case "Ctrl+Down":
			view.App.SetFocus(menu)
		}
		return event
	})

	// Display the "man page" menu
	menu.SetCurrentItem(menu.GetItemCount() - 1)
	pages.SwitchToPage("SaveRequestProjectAliasName")
	view.App.SetFocus(menu)
	view.App.SetFocus(mapMenuToFocusPrmt["menu_1"])

	// Don't forget!
	view.TitlePrmt = titleAndMenuFlexPrmt
	view.ParentPrmt = frame
}

func (view *SaveRequestView) makeMenu(pages *tview.Pages, mapMenuToFocusPrmt map[string]tview.Primitive) *tview.List {
	menu := tview.NewList().
		AddItem(view.Labels["menu_1_title"], view.Labels["menu_1_desc"], 'a', func() {
			pages.SwitchToPage("SaveRequestProjectAliasName")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_1"])
		})

	menu.
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(utils.BackGrayColor)

	return menu
}

func (view *SaveRequestView) makeForm(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Form {
	// Left side Form
	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	// New field - "Project name"
	formPrmt.AddInputField(view.Labels["project"], "", 0, nil, nil)
	utils.AddInputFieldEventForm(formPrmt, view.Labels["project"])

	// New field - "Request Alias"
	formPrmt.AddInputField(view.Labels["alias"], "", 0, nil, nil)
	utils.AddInputFieldEventForm(formPrmt, view.Labels["alias"])

	// New field - "Save"
	formPrmt.AddButton(view.Labels["save"], func() {
		makeRequestData := view.Event.GetMDR()
		makeRequestData.ProjectName = utils.GetInputFieldForm(formPrmt, view.Labels["project"]).GetText()
		makeRequestData.Alias = utils.GetInputFieldForm(formPrmt, view.Labels["alias"]).GetText()

		view.Event.UpdateMDR(makeRequestData)
	})

	// Add listener to refresh primitive when the MakeRequestData is changing...
	view.Event.AddListenerMRD["makeForm"] = func(makeRequestData models.MakeRequestData) {
		projectPrmt := utils.GetInputFieldForm(formPrmt, view.Labels["project"])
		aliasPrmt := utils.GetInputFieldForm(formPrmt, view.Labels["alias"])

		projectPrmt.SetText(makeRequestData.ProjectName)
		aliasPrmt.SetText(makeRequestData.Alias)
	}

	mapMenuToFocusPrmt["menu_1"] = formPrmt
	return formPrmt
}
