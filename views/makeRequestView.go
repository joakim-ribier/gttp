package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joakim-ribier/gttp/components"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// MakeRequestView represents the creation of a new request
type MakeRequestView struct {
	App    *tview.Application
	AppCtx *models.AppCtx

	Labels map[string]string

	RootPrmt tview.Primitive
	FormPrmt *tview.Form

	// Actions
	Save   func(callback func())
	Remove func(callback func())
}

// NewMakeRequestView returns the view for the make request view
func NewMakeRequestView(
	app *tview.Application,
	ctx *models.AppCtx,
	save func(callback func()),
	remove func(callback func())) *MakeRequestView {

	labels := make(map[string]string)
	labels["execution_context"] = "Execution Context"
	labels["request_method"] = "Request Method"
	labels["request_url"] = "Request URL"
	labels["execute"] = "[::ub]E[-:-:-]xecute"
	labels["save_request"] = "[::ub]S[-:-:-]ave"
	labels["expert_mode"] = "([::ub]H[-:-:-])Expert mode"
	labels["delete_request"] = "[::ub]D[-:-:-]elete"
	labels["new_request"] = "[::ub]N[-:-:-]ew"
	labels["project"] = "Project"
	labels["alias"] = "Alias"
	labels["cancel"] = "Cancel"
	labels["save"] = "Save"

	return &MakeRequestView{
		App:    app,
		AppCtx: ctx,
		Labels: labels,
		Save:   save,
		Remove: remove,
	}
}

// InitView build all components to display correctly the view
func (view *MakeRequestView) InitView(
	executeRequest func(),
	displayExpertMode func(),
	newRequest func()) {

	methodValues := utils.MethodValues

	flex := tview.NewFlex()
	flex.SetBorder(false)
	flex.SetBorderPadding(0, 0, 0, 0)

	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBorderPadding(0, 0, 0, 0)

	setDropDownExContextDefaultValue := func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.setDropDownExContextDefaultValue")

		envs := view.AppCtx.GetOutput().Context.GetEnvsName()

		prmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["execution_context"])
		prmt.SetOptions(envs, nil)

		index := envs.GetIndex("default")
		prmt.SetCurrentOption(index)
	}

	// New Field - "Ex. Context"
	formPrmt.AddDropDown(view.Labels["execution_context"], nil, 0, nil)

	// New Field - "Request Method"
	formPrmt.AddDropDown(view.Labels["request_method"], methodValues, 0, func(option string, index int) {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddDropDown@" + view.Labels["request_method"])

		makeRequestData := view.AppCtx.GetMDR()
		makeRequestData.Method = types.Method(option)

		view.AppCtx.UpdateMDR(makeRequestData)
	})

	// New Field - "Request URL"
	formPrmt.AddInputField(view.Labels["request_url"], view.AppCtx.GetMDR().URL.String(), 0, nil, func(text string) {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddInputField@" + view.Labels["request_url"])

		makeRequestData := view.AppCtx.GetMDR()
		makeRequestData.URL = types.URL(text)

		view.AppCtx.UpdateMDR(makeRequestData)
	})
	utils.AddInputFieldEventForm(formPrmt, view.Labels["request_url"])

	// New Field - "Execute"
	formPrmt.AddButton(view.Labels["execute"], func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["execute"])

		executeRequest()
	})

	// New Field - "Expert mode"
	formPrmt.AddButton(view.Labels["expert_mode"], func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["expert_mode"])

		displayExpertMode()
	})

	formPrmt.AddButton("", func() {
	})

	// New Field - "Save Request"
	formPrmt.AddButton(view.Labels["save_request"], func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["save_request"])
		view.DisplaySaveView()
	})

	formPrmt.AddButton("", func() {
	})

	formPrmt.GetButton(0)

	// New Field - "New Request"
	formPrmt.AddButton(view.Labels["new_request"], func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["new_request"])

		newRequest()
	})

	// New Field - "Delete Request"
	formPrmt.AddButton(view.Labels["delete_request"], func() {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddButton@" + view.Labels["delete_request"])
		view.DisplayRemoveView()
	})

	flex.AddItem(formPrmt, 0, 1, false)

	view.AppCtx.AddListenerMRD["refreshRequestPanelView"] = func(makeRequestData models.MakeRequestData) {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddListenerMRD")

		utils.GetInputFieldForm(formPrmt, view.Labels["request_url"]).SetText(makeRequestData.URL.String())

		methodSelectedIndex := methodValues.GetIndex(makeRequestData.Method.String())
		utils.GetDropDownFieldForm(formPrmt, view.Labels["request_method"]).SetCurrentOption(methodSelectedIndex)
	}

	view.AppCtx.AddContextListener["refreshRequestPanelView"] = func(context models.Context) {
		view.AppCtx.PrintTrace("MakeRequestView.InitView{...}.AddContextListener")

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

// DisplaySaveView displays request saving/updating view
func (view *MakeRequestView) DisplaySaveView() {
	textViewError := tview.NewTextView()
	textViewError.SetTextColor(tcell.ColorRed)
	textViewError.SetText("")

	form := tview.NewForm()

	// New field - "Project name"
	form.AddInputField(view.Labels["project"], view.AppCtx.GetMDR().ProjectName, 0, nil, nil)
	utils.AddInputFieldEventForm(form, view.Labels["project"])

	// New field - "Request Alias"
	form.AddInputField(view.Labels["alias"], view.AppCtx.GetMDR().Alias, 0, nil, nil)
	utils.AddInputFieldEventForm(form, view.Labels["alias"])

	// New Field - "Cancel"
	form.AddButton(view.Labels["cancel"], func() {
		view.AppCtx.CloseModal()
	})

	// New Field - "Save"
	form.AddButton(view.Labels["save"], func() {
		if mrd := view.AppCtx.GetMDR(); mrd.URL != "" {
			mrd.ProjectName = utils.GetInputFieldForm(form, view.Labels["project"]).GetText()
			mrd.Alias = utils.GetInputFieldForm(form, view.Labels["alias"]).GetText()

			view.AppCtx.UpdateMDR(mrd)

			view.Save(func() {
				view.AppCtx.CloseModal()
			})
		} else {
			textViewError.SetText(" empty 'request url' field")
		}
	})

	flexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	flexPrmt.SetBorder(true).SetTitle(" Save / Update ")
	flexPrmt.AddItem(form, 0, 1, true)
	flexPrmt.AddItem(textViewError, 1, 0, false)

	view.AppCtx.DisplayModal(components.BuildModal(flexPrmt, 45, 10))

	view.App.SetFocus(form)
}

// DisplaySaveView displays request saving/updating view
func (view *MakeRequestView) DisplayRemoveView() {
	modal := components.BuildYesNoModal(
		"Do you confirm the deletion?",
		" Remove request ", func() {
			view.AppCtx.CloseModal()
		}, func() {
			view.Remove(func() {
				view.AppCtx.CloseModal()
			})
		}, func(form tview.Primitive) {
			view.App.SetFocus(form)
		})
	view.AppCtx.DisplayModal(modal)
}
