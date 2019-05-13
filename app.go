package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/components/tree"
	"github.com/joakim-ribier/gttp/httpclient"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/joakim-ribier/gttp/views"
	"github.com/rivo/tview"
)

const (
	requestURLPrmtLabel       = "Request URL"
	requestMethodPrmtLabel    = "Request Method"
	requestExContextPrmtLabel = "Execution Context"
)

var (
	app                  *tview.Application
	requestFormPrmt      *tview.Form
	logEventTextPrmt     *tview.TextView
	shortcutInfoTextPrmt *tview.TextView
	pages                *tview.Pages
)

var (
	output models.Output
	event  *models.Event

	appPathFileName = ""
	responseData    = ""

	makeRequestData           = models.NewMakeRequestData()
	mapFocusPrmtToShortutText = make(map[tview.Primitive]string)
	focusPrmts                = []*tview.TextView{}

	// List of views of the application
	expertModeView      *views.RequestExpertModeView
	settingsView        *views.SettingsView
	saveRequestView     *views.SaveRequestView
	requestResponseView *views.RequestResponseView

	// List of components of the application
	treeAPICpnt *components.TreeCpnt
)

// App main method
func App() {
	if len(os.Args) != 2 {
		fmt.Println(`Please provide a data json file {0}.`)
		fmt.Println("\nSee https://github.com/joakim-ribier/gttp for details.")
		return
	}

	event = models.NewEvent(getMDR, updateMDR, getConfig, updateConfig, getOutput, updateContext)
	appPathFileName = os.Args[1]

	initializeData := func() {
		mapFocusPrmtToShortutText[requestResponseView.ResponsePrmt] = utils.ResultShortcutsText
		mapFocusPrmtToShortutText[expertModeView.TitlePrmt] = utils.ExpertModeShortcutsText
		mapFocusPrmtToShortutText[settingsView.TitlePrmt] = utils.SettingsShortcutsText
		mapFocusPrmtToShortutText[saveRequestView.TitlePrmt] = utils.SaveRequestShortcutsText

		unmarshal()

		refreshingTreeAPICpn()
		refreshingConfig()
		refreshingContext()

		displaySettingsViewPage()
	}

	app = tview.NewApplication()
	root := drawMainComponents(app)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logEventText("Shortcut: "+event.Name()+" - "+time.Now().Format(time.RFC850), "info")

		switch event.Key() {
		case tcell.KeyCtrlA:
			if requestResponseView.ResponsePrmt.GetFocusable().HasFocus() {
				utils.WriteToClipboard(requestResponseView.LogBuffer, logEventText)
			}
		case tcell.KeyCtrlC:
			if requestResponseView.ResponsePrmt.GetFocusable().HasFocus() {
				utils.WriteToClipboard(responseData, logEventText)
			}
			if prmt := app.GetFocus(); prmt != nil {
				if input, er := app.GetFocus().(*tview.InputField); er {
					utils.WriteToClipboard(input.GetText(), logEventText)
				}
			}
			// Disable "Ctrl+C" exit application default shortcut
			return nil
		case tcell.KeyCtrlD:
			displayRequestResponseViewPage(requestResponseView.ResponsePrmt)
		case tcell.KeyCtrlE:
			executeRequest()
		case tcell.KeyCtrlF:
			focusPrimitive(requestFormPrmt, nil)
		case tcell.KeyCtrlH:
			displayRequestExpertModeViewPage()
		case tcell.KeyCtrlJ:
			focusPrimitive(treeAPICpnt.RootPrmt, nil)
		case tcell.KeyCtrlO:
			displaySettingsViewPage()
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlR:
			displayRequestResponseViewPage(requestResponseView.RequestPrmt)
		case tcell.KeyEsc:
			focusPrimitive(logEventTextPrmt, nil)
		}
		return event
	})

	initializeData()

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}

func drawMainComponents(app *tview.Application) tview.Primitive {
	logEventTextPrmt = tview.NewTextView()
	logEventTextPrmt.SetBackgroundColor(utils.BackColorPrmt)
	logEventTextPrmt.
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)

	shortcutInfoTextPrmt = tview.NewTextView()
	shortcutInfoTextPrmt.SetBackgroundColor(utils.BackColor)
	shortcutInfoTextPrmt.
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true).
		SetText(utils.MainShortcutsText)

	grid := tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(0, 10, -4).
		SetBorders(false).
		AddItem(logEventTextPrmt, 0, 0, 1, 3, 0, 0, false).
		AddItem(drawLeftPanel(), 1, 0, 1, 2, 0, 0, false).
		AddItem(drawRightPanel(), 1, 2, 1, 1, 0, 0, false).
		AddItem(shortcutInfoTextPrmt, 2, 0, 1, 3, 0, 0, false)

	frame := tview.NewFrame(grid).SetBorders(0, 0, 0, 0, 0, 0)
	frame.SetBorder(true).SetTitle(" " + utils.Subtitle + " ")

	return frame
}

func drawLeftPanel() tview.Primitive {
	titlePrmt := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText(utils.TitleAPIText)

	// Add tree APIs component
	treeAPICpnt = components.NewTreeCpnt(app, event)
	tree := treeAPICpnt.Make(func(it models.MakeRequestData) {
		refreshMDRView(it)
	}, func(page string) {
		pages.SwitchToPage(page)
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(titlePrmt, 8, 0, false)
	flex.AddItem(tree, 0, 1, false)
	return flex
}

func drawRightPanel() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	pages = tview.NewPages()
	pages.SetBorder(false).SetBorderPadding(1, 1, 1, 0)

	pages.AddPage("RequestResponseViewPage", makeRequestResponseView(), true, false)
	pages.AddPage("RequestExpertModeViewPage", makeRequestExportModeView(), true, false)
	pages.AddPage("SettingsViewPage", makeSettingsView(), true, true)
	pages.AddPage("SaveRequestViewPage", makeSaveRequestView(), true, false)

	flex.AddItem(drawMakeRequestPanel(), 9, 0, false)
	flex.AddItem(pages, 0, 1, false)

	return flex
}

func displayRequestResponseViewPage(focusOn *tview.TextView) {
	pages.SwitchToPage("RequestResponseViewPage")
	focusPrimitive(focusOn, focusOn.Box)
}

func displaySettingsViewPage() {
	pages.SwitchToPage("SettingsViewPage")
	focusPrimitive(settingsView.TitlePrmt, nil)
}

func displaySaveRequestViewPage() {
	pages.SwitchToPage("SaveRequestViewPage")
	focusPrimitive(saveRequestView.TitlePrmt, nil)
}

func displayRequestExpertModeViewPage() {
	pages.SwitchToPage("RequestExpertModeViewPage")
	focusPrimitive(expertModeView.TitlePrmt, nil)
}

func executeRequest() {
	displayRequestResponseViewPage(requestResponseView.ResponsePrmt)
	requestResponseView.ResetLogBuffer()
	logEventText("", "info")

	// Get current context to replace all variables
	_, currentContext := utils.GetDropDownFieldForm(requestFormPrmt, requestExContextPrmtLabel).GetCurrentOption()
	currentContextValues := getOutput().Context.GetAllKeyValue(currentContext)

	URL := types.URL(getRequestURLPrmtText()).
		ReplaceContext(makeRequestData.MapRequestHeaderKeyValue).
		ReplaceContext(currentContextValues)

	method := makeRequestData.Method
	contentType := makeRequestData.ContentType
	body := []byte(makeRequestData.Body)
	httpHeaderValues := makeRequestData.GetHTTPHeaderValues().ReplaceContext(currentContextValues)

	HTTPClient, error := httpclient.Call(method, URL, contentType, body, httpHeaderValues, requestResponseView.Logger)
	if error != nil {
		requestResponseView.Logger(fmt.Sprint(error), "error")
	} else {
		responseData = fmt.Sprintf("%+s", HTTPClient.Body)
		requestResponseView.Display(HTTPClient, responseData)
	}
}

func drawMakeRequestPanel() tview.Primitive {
	methodValues := utils.MethodValues

	flex := tview.NewFlex()

	requestFormPrmt = tview.NewForm()
	requestFormPrmt.SetBorder(false)

	setDropDownExContextDefaultValue := func() {
		envs := getOutput().Context.GetEnvsName()

		prmt := utils.GetDropDownFieldForm(requestFormPrmt, requestExContextPrmtLabel)
		prmt.SetOptions(envs, nil)

		index := envs.GetIndex("default")
		prmt.SetCurrentOption(index)
	}

	// New Field - "Ex. Context"
	requestFormPrmt.AddDropDown(requestExContextPrmtLabel, nil, 0, nil)

	// New Field - "Request Method"
	requestFormPrmt.AddDropDown(requestMethodPrmtLabel, methodValues, 0, func(option string, index int) {
		makeRequestData.Method = types.Method(option)
	})

	// New Field - "Request URL"
	requestFormPrmt.AddInputField(requestURLPrmtLabel, makeRequestData.URL.String(), 0, nil, func(text string) {
		makeRequestData.URL = types.URL(text)
	})

	// New Field - "Execute"
	requestFormPrmt.AddButton("Execute", func() {
		executeRequest()
	})

	// New Field - "Expert mode"
	requestFormPrmt.AddButton("Expert mode", func() {
		pages.SwitchToPage("RequestExpertModeViewPage")
	})

	// New Field - "Save Request"
	requestFormPrmt.AddButton("Save request", func() {
		saveCurrentRequest()
		refreshingTreeAPICpn()
		refreshingConfig()
		refreshMDRView(getMDR())
		displaySaveRequestViewPage()
	})

	utils.AddInputFieldEventForm(requestFormPrmt, requestURLPrmtLabel)

	flex.AddItem(requestFormPrmt, 0, 1, false)

	event.AddListenerMRD["refreshRequestPanelView"] = func(makeRequestData models.MakeRequestData) {
		utils.GetInputFieldForm(requestFormPrmt, requestURLPrmtLabel).SetText(makeRequestData.URL.String())

		methodSelectedIndex := methodValues.GetIndex(makeRequestData.Method.String())
		utils.GetDropDownFieldForm(requestFormPrmt, requestMethodPrmtLabel).SetCurrentOption(methodSelectedIndex)
	}

	event.AddContextListener["refreshRequestPanelView"] = func(context models.Context) {
		setDropDownExContextDefaultValue()
	}

	return flex
}

func getDataFromTheDisk() []byte {
	return utils.GetByteFromPathFileName(appPathFileName, logEventText)
}

func saveCurrentRequest() {
	// 1. Read disk data
	unmarshal()
	// 2. Update output
	output.AddOrReplace(makeRequestData)
	// 3. Update disk data
	marshal()
}

func getRequestURLPrmtText() string {
	prmt := utils.GetInputFieldForm(requestFormPrmt, requestURLPrmtLabel)
	return prmt.GetText()
}

func focusPrimitive(prmt tview.Primitive, box *tview.Box) {
	app.SetFocus(prmt)

	// Set border false to all focus prmt
	for v := range focusPrmts {
		focusPrmts[v].SetBorder(false)
	}
	if box != nil {
		box.SetBorder(true)
	}

	// Display the right shortcuts text
	if text, exists := mapFocusPrmtToShortutText[prmt]; exists {
		shortcutInfoTextPrmt.SetText(text)
	} else {
		shortcutInfoTextPrmt.SetText(utils.MainShortcutsText)
	}
}

func logEventText(message string, status string) {
	if message != "" {
		logEventTextPrmt.SetText(utils.FormatLog(message, status))
	}
}

func updateMDR(value models.MakeRequestData) {
	makeRequestData = value
}

func refreshMDRView(makeRequestData models.MakeRequestData) {
	updateMDR(makeRequestData)
	for _, value := range event.AddListenerMRD {
		value(makeRequestData)
	}
}

func getMDR() models.MakeRequestData {
	return makeRequestData
}

func getConfig() models.Config {
	return output.Config
}

func updateConfig(value models.Config) {
	// 1. Read disk data
	unmarshal()
	// 2. Update output
	output.Config = value
	// 3. Update disk data
	marshal()
	// 4. Refresh views
	refreshingConfig()
	refreshingTreeAPICpn()
}

func updateContext(value models.Context) {
	// 1. Read disk data
	unmarshal()
	// 2. Update output
	output.Context = value
	// 3. Update disk data
	marshal()
	// 4. Refresh views
	refreshingContext()
}

func refreshingConfig() {
	for _, value := range event.AddListenerConfig {
		value(output.Config)
	}
}

func refreshingContext() {
	for _, value := range event.AddContextListener {
		value(output.Context)
	}
}

func refreshingTreeAPICpn() {
	treeAPICpnt.Refresh()
}

func getOutput() models.Output {
	return output
}

// ## -- Marshal & unmarshal json

func unmarshal() {
	var data models.Output
	if error := json.Unmarshal([]byte(getDataFromTheDisk()), &data); error != nil {
		logEventText("Error to decode '"+appPathFileName+"' json data file.", "error")
	} else {
		output = data
	}
}

func marshal() {
	if json, error := json.Marshal(output); error != nil {
		logEventText("Error to encode 'output' model.", "error")
	} else {
		if error := ioutil.WriteFile(appPathFileName, json, 0644); error != nil {
			logEventText("Error to encode '"+appPathFileName+"' json data file.", "error")
		}
	}
}

// -- ##

// ## -- Make all views

func makeRequestExportModeView() tview.Primitive {
	expertModeView = views.NewRequestExpertModeView(app, event)
	expertModeView.InitView()

	return expertModeView.ParentPrmt
}

func makeSettingsView() tview.Primitive {
	settingsView = views.NewSettingsView(app, event)
	settingsView.InitView()

	return settingsView.ParentPrmt
}

func makeSaveRequestView() tview.Primitive {
	saveRequestView = views.NewSaveRequestView(app, event)
	saveRequestView.InitView()

	return saveRequestView.ParentPrmt
}

func makeRequestResponseView() tview.Primitive {
	requestResponseView = views.NewRequestResponseView(app, event)
	requestResponseView.InitView()

	focusPrmts = append(focusPrmts, requestResponseView.ResponsePrmt)
	focusPrmts = append(focusPrmts, requestResponseView.RequestPrmt)

	return requestResponseView.ParentPrmt
}

// -- ##
