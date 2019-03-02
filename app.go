package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/components"
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
	app                    *tview.Application
	requestFormPrmt        *tview.Form
	responseTextPrmt       *tview.TextView
	responseHeaderTextPrmt *tview.TextView
	messageInfoTextPrmt    *tview.TextView
	shortcutInfoTextPrmt   *tview.TextView
	pages                  *tview.Pages
)

var (
	output models.Output
	event  *models.Event

	appPathFileName = ""
	bufferLog       = ""
	responseData    = ""

	makeRequestData    = models.NewMakeRequestData()
	mapPrmtToShortcuts = make(map[tview.Primitive]string)

	// List of views of the application
	expertModeView  *views.RequestExpertModeView
	settingsView    *views.SettingsView
	saveRequestView *views.SaveRequestView

	// List of components of the application
	treeAPICpnt *components.TreeAPICpnt
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
		mapPrmtToShortcuts[responseTextPrmt] = utils.ResultShortcutsText
		mapPrmtToShortcuts[expertModeView.TitlePrmt] = utils.ExpertModeShortcutsText
		mapPrmtToShortcuts[settingsView.TitlePrmt] = utils.SettingsShortcutsText
		mapPrmtToShortcuts[saveRequestView.TitlePrmt] = utils.SaveRequestShortcutsText

		json.Unmarshal([]byte(getDataFromTheDisk()), &output)

		refreshingTreeAPICpn()
		refreshingConfig()
		refreshingContext()

		displaySettingsViewPage()
	}

	app = tview.NewApplication()
	root := drawMainComponents(app)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		setMessageInfoTextPrmt("Shortcut: " + event.Name() + " - " + time.Now().Format(time.RFC850))

		switch event.Key() {
		case tcell.KeyEsc:
			focusPrimitive(messageInfoTextPrmt)
		case tcell.KeyCtrlJ:
			focusPrimitive(treeAPICpnt.Tree)
		case tcell.KeyCtrlD:
			displayRequestResponseViewPage()
		case tcell.KeyCtrlR:
			focusPrimitive(requestFormPrmt)
		case tcell.KeyCtrlH:
			displayRequestExpertModeViewPage()
		case tcell.KeyCtrlE:
			executeRequest()
		case tcell.KeyCtrlO:
			displaySettingsViewPage()
		case tcell.KeyCtrlC:
			if responseTextPrmt.GetFocusable().HasFocus() {
				utils.WriteToClipboard(responseData)
			}
			if prmt := app.GetFocus(); prmt != nil {
				input, er := app.GetFocus().(*tview.InputField)
				if er == true {
					utils.WriteToClipboard(input.GetText())
				}
			}
			// Disable "Ctrl+C" exit application default shortcut
			return nil
		case tcell.KeyCtrlA:
			if responseTextPrmt.GetFocusable().HasFocus() {
				utils.WriteToClipboard(bufferLog)
			}
		case tcell.KeyCtrlQ:
			app.Stop()
		}
		return event
	})

	initializeData()

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}

func drawMainComponents(app *tview.Application) tview.Primitive {
	messageInfoTextPrmt = tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true)
	messageInfoTextPrmt.Box.SetBackgroundColor(utils.BackColorPrmt)

	shortcutInfoTextPrmt = tview.NewTextView().
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true)
	shortcutInfoTextPrmt.Box.SetBackgroundColor(utils.BackColor)
	shortcutInfoTextPrmt.SetText(utils.MainShortcutsText)

	grid := tview.NewGrid().
		SetRows(1, 0, 2).
		SetColumns(0, 10, -4).
		SetBorders(false).
		AddItem(messageInfoTextPrmt, 0, 0, 1, 3, 0, 0, false).
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
	treeAPICpnt = components.NewTreeAPICpnt(app, event)
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
	pages.AddPage("RequestResponseViewPage", drawResponsePanel(), true, false)
	pages.AddPage("RequestExpertModeViewPage", makeRequestExportModeView(), true, false)
	pages.AddPage("SettingsViewPage", makeSettingsView(), true, true)
	pages.AddPage("SaveRequestViewPage", makeSaveRequestView(), true, false)

	flex.AddItem(drawMakeRequestPanel(), 9, 0, false)
	flex.AddItem(pages, 0, 1, false)

	return flex
}

func drawResponsePanel() tview.Primitive {
	responseTextPrmt = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	responseTextPrmt.SetBackgroundColor(utils.BackColorPrmt)

	responseHeaderTextPrmt = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).SetWrap(true)
	responseHeaderTextPrmt.SetBackgroundColor(utils.BackColorPrmt)

	titleAndMenuFlexPrmt := utils.MakeTitlePrmt("Execute Request")
	titleAndMenuFlexPrmt.AddItem(responseHeaderTextPrmt, 0, 1, false)

	flex := tview.NewFlex()
	flex.AddItem(titleAndMenuFlexPrmt, 0, 1, false)
	flex.AddItem(tview.NewBox().SetBorder(false), 2, 0, false)
	flex.AddItem(responseTextPrmt, 0, 2, false)

	frame := tview.NewFrame(flex).SetBorders(2, 2, 0, 0, 0, 0)
	frame.SetBorder(false)

	return frame
}

func makeRequestExportModeView() tview.Primitive {
	// Create the view
	expertModeView = views.NewRequestExpertModeView(app, event)

	// Build components
	expertModeView.InitView()

	return expertModeView.ParentPrmt
}

func makeSettingsView() tview.Primitive {
	// Create the view
	settingsView = views.NewSettingsView(app, event)

	// Build components
	settingsView.InitView()

	return settingsView.ParentPrmt
}

func makeSaveRequestView() tview.Primitive {
	// Create the view
	saveRequestView = views.NewSaveRequestView(app, event)

	// Build components
	saveRequestView.InitView()

	return saveRequestView.ParentPrmt
}

func displayRequestInfo(client *httpclient.HTTPClient, data string) {

	format := func(key string, value string) string {
		if key == "" {
			return "[white]" + value
		}
		return "[blue]" + key + "[white] " + value
	}

	var sb strings.Builder

	// Add request header info
	sb.WriteString(format(client.Request.Method, client.Request.URL+" [blue]HTTP[white]/"+client.Request.HTTP))
	sb.WriteString("\r\n")
	for k, v := range client.HeadersRequest {
		if k != "Content-Type" {
			sb.WriteString(format(k, v))
			sb.WriteString("\r\n")
		}
	}

	sb.WriteString("\r\n")
	sb.WriteString(format("Host", client.Request.Host))

	// Add request body
	sb.WriteString("\r\n")
	sb.WriteString(format("Content-Type", client.HeadersRequest["Content-Type"]))
	sb.WriteString("\r\n\r\n")
	sb.WriteString(format("", client.Request.Body))

	// Add response header info
	sb.WriteString("\r\n\r\n")
	sb.WriteString(format("HTTP[white]/"+client.Response.HTTP, client.Response.Status))
	sb.WriteString("\r\n")
	sb.WriteString(format("Content-Length", client.Response.Contentlength))
	sb.WriteString("\r\n")
	sb.WriteString(format("Content-Type", client.Response.ContentType))
	sb.WriteString("\r\n")
	sb.WriteString(format("Date", client.Response.Date))
	sb.WriteString("\r\n")
	sb.WriteString(format("Referrer-Policy", client.Response.Referrerpolicy))
	sb.WriteString("\r\n")
	sb.WriteString(format("Connection", client.Response.Connection))

	responseHeaderTextPrmt.SetText(sb.String()).SetTextAlign(tview.AlignLeft)

	// Log if error status
	status := client.Response.StatusCode
	if (status < "200") || (status > "205") {
		logger("Status: "+client.Response.Status, "error")
	}

	setBodyTextPrmt(utils.FormatLog(data, "data"))
}

func displayRequestResponseViewPage() {
	pages.SwitchToPage("RequestResponseViewPage")
	focusPrimitive(responseTextPrmt)
}

func displaySettingsViewPage() {
	pages.SwitchToPage("SettingsViewPage")
	focusPrimitive(settingsView.TitlePrmt)
}

func displaySaveRequestViewPage() {
	pages.SwitchToPage("SaveRequestViewPage")
	focusPrimitive(saveRequestView.TitlePrmt)
}

func displayRequestExpertModeViewPage() {
	pages.SwitchToPage("RequestExpertModeViewPage")
	focusPrimitive(expertModeView.TitlePrmt)
}

func executeRequest() {
	displayRequestResponseViewPage()
	bufferLog = ""
	setMessageInfoTextPrmt("")

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

	HTTPClient, error := httpclient.Call(method, URL, contentType, body, httpHeaderValues, logger)
	if error != nil {
		logger(fmt.Sprint(error), "error")
	} else {
		responseData = fmt.Sprintf("%+s", HTTPClient.Body)
		displayRequestInfo(HTTPClient, responseData)
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
	return utils.GetByteFromPathFileName(appPathFileName, logger)
}

func saveCurrentRequest() {
	// Read data from the disk
	var data models.Output
	json.Unmarshal([]byte(getDataFromTheDisk()), &data)

	// Save only the current request
	data.AddOrReplace(makeRequestData)

	// Write on the disk
	json, _ := json.Marshal(data)
	_ = ioutil.WriteFile(appPathFileName, json, 0644)

	output = data
}

func getRequestURLPrmtText() string {
	prmt := utils.GetInputFieldForm(requestFormPrmt, requestURLPrmtLabel)
	return prmt.GetText()
}

func newInputView(placeHolder string) *tview.InputField {
	inputField := tview.NewInputField().SetPlaceholder(placeHolder)
	return inputField
}

func newTextView(text string) *tview.TextView {
	textView := tview.NewTextView().
		SetChangedFunc(func() {
		})
	textView.SetText(text)
	textView.SetWrap(true)
	textView.SetWordWrap(false)
	textView.SetDynamicColors(true)
	textView.SetRegions(true)

	return textView
}

func focusPrimitive(prmt tview.Primitive) {
	app.SetFocus(prmt)
	shortcuts := mapPrmtToShortcuts[prmt]
	if prmt == nil || shortcuts == "" {
		shortcutInfoTextPrmt.SetText(utils.MainShortcutsText)
	} else {
		shortcutInfoTextPrmt.SetText(shortcuts)
	}
	if prmt != responseTextPrmt {
		responseTextPrmt.SetBackgroundColor(utils.BackColorPrmt)
	}
}

func focusTextViewWithColor(prmt tview.Primitive, BackColor tcell.Color, backFocusColor tcell.Color) {
	focusPrimitive(prmt)
	if app.GetFocus() == prmt {
		prmt.(*tview.TextView).SetBackgroundColor(backFocusColor)
	} else {
		prmt.(*tview.TextView).SetBackgroundColor(BackColor)
	}
}

func setMessageInfoTextPrmt(message string) {
	if message != "" {
		messageInfoTextPrmt.SetText(utils.FormatLog(message, "info"))
	}
}

func setBodyTextPrmt(data string) {
	responseTextPrmt.Clear()
	responseTextPrmt.SetTextAlign(tview.AlignLeft)
	bufferLog = bufferLog + data + "\r\n"
	responseTextPrmt.SetText(bufferLog)
}

func logger(message string, mode string) {
	setBodyTextPrmt(utils.FormatLog(message, mode))
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
	// Read data from the disk
	var data models.Output
	json.Unmarshal([]byte(getDataFromTheDisk()), &data)

	// Save the config
	data.Config = value

	// Write on the disk
	json, _ := json.Marshal(data)
	_ = ioutil.WriteFile(appPathFileName, json, 0644)

	output = data

	refreshingConfig()
	refreshingTreeAPICpn()
}

func updateContext(value models.Context) {
	// Read data from the disk
	var data models.Output
	json.Unmarshal([]byte(getDataFromTheDisk()), &data)

	// Save the config
	data.Context = value

	// Write on the disk
	json, _ := json.Marshal(data)
	_ = ioutil.WriteFile(appPathFileName, json, 0644)

	output = data

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
