package views

import (
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// MakeRequestView represents the creation of a new request
type MakeRequestView struct {
	App   *tview.Application
	Event *models.Event

	Labels map[string]string

	RootPrmt tview.Primitive
	FormPrmt *tview.Form
}

// NewMakeRequestView returns the view for the make request view
func NewMakeRequestView(app *tview.Application, ev *models.Event) *MakeRequestView {
	labels := make(map[string]string)
	labels["execution_context"] = "Execution Context"
	labels["request_method"] = "Request Method"
	labels["request_url"] = "Request URL"
	labels["execute"] = "[white::ub]E[white::-]xecute"
	labels["save_request"] = "Save"
	labels["expert_mode"] = "Expert mode"
	labels["delete_request"] = "Delete"
	labels["new_request"] = "New"

	return &MakeRequestView{
		App:    app,
		Event:  ev,
		Labels: labels,
	}
}

// InitView build all components to display correctly the view
func (view *MakeRequestView) InitView(
	executeRequest func(),
	displayExpertMode func(),
	saveRequest func(),
	removeRequest func(),
	newRequest func()) {

	methodValues := utils.MethodValues

	flex := tview.NewFlex()
	flex.SetBorder(false)
	flex.SetBorderPadding(0, 0, 0, 0)

	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBorderPadding(0, 0, 0, 0)

	setDropDownExContextDefaultValue := func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.setDropDownExContextDefaultValue")

		envs := view.Event.GetOutput().Context.GetEnvsName()

		prmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["execution_context"])
		prmt.SetOptions(envs, nil)

		index := envs.GetIndex("default")
		prmt.SetCurrentOption(index)
	}

	// New Field - "Ex. Context"
	formPrmt.AddDropDown(view.Labels["execution_context"], nil, 0, nil)

	// New Field - "Request Method"
	formPrmt.AddDropDown(view.Labels["request_method"], methodValues, 0, func(option string, index int) {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddDropDown@" + view.Labels["request_method"])

		makeRequestData := view.Event.GetMDR()
		makeRequestData.Method = types.Method(option)

		view.Event.UpdateMDR(makeRequestData)
	})

	// New Field - "Request URL"
	formPrmt.AddInputField(view.Labels["request_url"], view.Event.GetMDR().URL.String(), 0, nil, func(text string) {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddInputField@" + view.Labels["request_url"])

		makeRequestData := view.Event.GetMDR()
		makeRequestData.URL = types.URL(text)

		view.Event.UpdateMDR(makeRequestData)
	})
	utils.AddInputFieldEventForm(formPrmt, view.Labels["request_url"])

	// New Field - "Execute"
	formPrmt.AddButton(view.Labels["execute"], func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["execute"])

		executeRequest()
	})

	// New Field - "Expert mode"
	formPrmt.AddButton(view.Labels["expert_mode"], func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["expert_mode"])

		displayExpertMode()
	})

	formPrmt.AddButton("", func() {
	})

	// New Field - "Save Request"
	formPrmt.AddButton(view.Labels["save_request"], func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["save_request"])

		saveRequest()
	})

	formPrmt.AddButton("", func() {
	})

	// New Field - "New Request"
	formPrmt.AddButton(view.Labels["new_request"], func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["new_request"])

		newRequest()
	})

	// New Field - "Delete Request"
	formPrmt.AddButton(view.Labels["delete_request"], func() {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["delete_request"])

		removeRequest()
	})

	flex.AddItem(formPrmt, 0, 1, false)

	view.Event.AddListenerMRD["refreshRequestPanelView"] = func(makeRequestData models.MakeRequestData) {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddListenerMRD")

		utils.GetInputFieldForm(formPrmt, view.Labels["request_url"]).SetText(makeRequestData.URL.String())

		methodSelectedIndex := methodValues.GetIndex(makeRequestData.Method.String())
		utils.GetDropDownFieldForm(formPrmt, view.Labels["request_method"]).SetCurrentOption(methodSelectedIndex)
	}

	view.Event.AddContextListener["refreshRequestPanelView"] = func(context models.Context) {
		view.Event.PrintTrace("MakeRequestView.InitView{...}.AddContextListener")

		setDropDownExContextDefaultValue()
	}

	// Define root view
	view.RootPrmt = flex
	view.FormPrmt = formPrmt
}

// GetURL gets value form the URL input text prmt
func (view *MakeRequestView) GetURL() string {
	prmt := utils.GetInputFieldForm(view.FormPrmt, view.Labels["request_url"])
	return prmt.GetText()
}

// GetContext gets value form the context dropdown prmt
func (view *MakeRequestView) GetContext() (int, string) {
	return utils.GetDropDownFieldForm(view.FormPrmt, view.Labels["execution_context"]).GetCurrentOption()
}
