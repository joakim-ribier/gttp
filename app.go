package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	components "github.com/joakim-ribier/gttp/components/tree"
	"github.com/joakim-ribier/gttp/httpclient"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/joakim-ribier/gttp/views"
	"github.com/rivo/tview"
)

var (
	logFile      *os.File
	logOn        = os.Getenv("GTTP_LOG") == "ON"
	logRequestOn = os.Getenv("GTTP_LOG_REQUEST") == "ON"
	logLevels    = os.Getenv("GTTP_LOG_LEVEL")

	app                  *tview.Application
	logEventTextPrmt     *tview.TextView
	shortcutInfoTextPrmt *tview.TextView
	pages                *tview.Pages
)

var (
	output models.Output
	event  *models.Event

	appPathFileName = ""
	responseData    = ""

	makeRequestData           = models.EmptyMakeRequestData()
	mapFocusPrmtToShortutText = make(map[tview.Primitive]string)
	focusPrmts                = []*tview.TextView{}

	// List of views of the application
	expertModeView      *views.RequestExpertModeView
	settingsView        *views.SettingsView
	saveRequestView     *views.SaveRequestView
	requestResponseView *views.RequestResponseView
	requestView         *views.MakeRequestView

	// List of components of the application
	treeAPICpnt *components.TreeCpnt
)

// GetDataFromDisk reads @filename and stores the data on @models.Output.
func GetDataFromDisk(filename string, log func(string, string)) models.Output {
	var value models.Output

	bytes := utils.ReadFile(filename, log)
	if error := json.Unmarshal(bytes, &value); error != nil {
		log("Error to decode '"+filename+"' json data file.", "error")
	}

	return value
}

// SaveDataOnDisk writes @value on the @filename file disk.
func SaveDataOnDisk(filename string, value models.Output, log func(string, string)) {
	if json, error := json.Marshal(value); error != nil {
		log("Encoding 'output' model error...", "error")
	} else {
		if error := ioutil.WriteFile(filename, json, 0644); error != nil {
			log("Writing data to '"+filename+"' file error...", "error")
		}
	}
}

// PrintOut writes message to a log file
func PrintOut(level string, value string) {
	getLevel := func() string {
		switch level {
		case "trace":
			return "[TRACE]"
		case "debug":
			return "[DEBUG]"
		case "error":
			return "[ERROR]"
		case "info":
			return "[INFO] "
		default:
			return ""
		}
	}

	if logOn && strings.Contains(logLevels, level) {
		logFile.WriteString(
			strings.ToUpper(getLevel()) + " " +
				time.Now().Format("Jan 02 15:04:05.000") + " " +
				value + "\n")
	}
}

func App() {

	getFilenameFromArgs := func(args []string) string {
		if len(args) > 1 {
			return args[1]
		} else {
			return "gttp-tmp.json"
		}
	}

	// log file
	os.Remove("application.log")
	if logOn {
		logFile, _ = os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		defer logFile.Close()
	}

	event = models.NewEvent(getMDR, updateMDR, deleteMDR, getConfig, updateConfig, getOutput, updateContext, PrintOut)
	appPathFileName = getFilenameFromArgs(os.Args)

	initializeData := func() {
		mapFocusPrmtToShortutText[requestResponseView.ResponsePrmt] = utils.ResultShortcutsText
		mapFocusPrmtToShortutText[expertModeView.TitlePrmt] = utils.ExpertModeShortcutsText
		mapFocusPrmtToShortutText[settingsView.TitlePrmt] = utils.SettingsShortcutsText
		mapFocusPrmtToShortutText[saveRequestView.TitlePrmt] = utils.SaveRequestShortcutsText

		output = GetDataFromDisk(appPathFileName, log)

		refreshingTreeAPICpn()
		refreshingConfig()
		refreshingContext()
	}

	app = tview.NewApplication()
	root := drawMainComponents(app)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log("Shortcut: "+event.Name()+" - "+time.Now().Format(time.RFC850), "info")

		switch event.Key() {
		case tcell.KeyCtrlA:
			if requestResponseView.ResponsePrmt.HasFocus() {
				utils.WriteToClipboard(requestResponseView.LogBuffer, log)
			}
		case tcell.KeyCtrlC:
			if requestResponseView.ResponsePrmt.HasFocus() {
				utils.WriteToClipboard(responseData, log)
			}
			if prmt := app.GetFocus(); prmt != nil {
				if input, er := app.GetFocus().(*tview.InputField); er {
					utils.WriteToClipboard(input.GetText(), log)
				}
			}
			// Disable "Ctrl+C" exit application default shortcut
			return nil
		case tcell.KeyCtrlD:
			displayRequestResponseViewPage(requestResponseView.ResponsePrmt)
		case tcell.KeyCtrlE:
			executeRequest()
		case tcell.KeyCtrlF:
			focusPrimitive(requestView.FormPrmt, nil)
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

	drawLeftPanel := func() tview.Primitive {
		treeAPICpnt = components.NewTreeCpnt(app, event)
		tree := treeAPICpnt.Make(func(it models.MakeRequestData) {
			refreshMDRView(it)
		}, func(page string) {
			pages.SwitchToPage(page)
		})

		flex := utils.MakeTitlePrmt(utils.TreePrmtTitle)
		flex.SetBorder(false)
		flex.SetBorderPadding(1, 0, 1, 1)
		flex.SetBackgroundColor(utils.BackColor)

		flex.AddItem(tree, 0, 1, false)

		return flex
	}

	drawRightPanel := func() tview.Primitive {

		makeRequestExportModeView := func() tview.Primitive {
			expertModeView = views.NewRequestExpertModeView(app, event)
			expertModeView.InitView()

			return expertModeView.ParentPrmt
		}

		makeSettingsView := func() tview.Primitive {
			settingsView = views.NewSettingsView(app, event)
			settingsView.InitView()

			return settingsView.ParentPrmt
		}

		makeSaveRequestView := func() tview.Primitive {
			saveRequestView = views.NewSaveRequestView(app, event)
			saveRequestView.InitView()

			return saveRequestView.ParentPrmt
		}

		makeRequestResponseView := func() tview.Primitive {
			requestResponseView = views.NewRequestResponseView(app, event)
			requestResponseView.InitView()

			focusPrmts = append(focusPrmts, requestResponseView.ResponsePrmt)
			focusPrmts = append(focusPrmts, requestResponseView.RequestPrmt)

			return requestResponseView.ParentPrmt
		}

		makeRequestView := func() tview.Primitive {
			saveRequestAction := func() {
				saveRequest(getMDR())
				refreshingTreeAPICpn()
				refreshingConfig()
				refreshMDRView(getMDR())
				displaySaveRequestViewPage()
			}

			removeRequestAction := func() {
				removeRequest(getMDR())
				updateMDR(models.MakeRequestData{})
				refreshingTreeAPICpn()
				refreshingConfig()
				refreshMDRView(getMDR())
			}

			newRequestAction := func() {
				updateMDR(models.MakeRequestData{})
				refreshMDRView(getMDR())
			}

			requestView = views.NewMakeRequestView(app, event)
			requestView.InitView(
				executeRequest,
				displayRequestExpertModeViewPage,
				saveRequestAction,
				removeRequestAction,
				newRequestAction)

			return requestView.RootPrmt
		}

		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(false)
		flex.SetBorderPadding(1, 0, 0, 0)

		pages = tview.NewPages()
		pages.SetBorder(false).SetBorderPadding(0, 1, 0, 0)

		pages.AddPage("RequestResponseViewPage", makeRequestResponseView(), true, false)
		pages.AddPage("RequestExpertModeViewPage", makeRequestExportModeView(), true, false)
		pages.AddPage("SettingsViewPage", makeSettingsView(), true, true)
		pages.AddPage("SaveRequestViewPage", makeSaveRequestView(), true, false)

		flex.AddItem(makeRequestView(), 9, 0, false)
		flex.AddItem(pages, 0, 1, false)

		return flex
	}

	logEventTextPrmt = tview.NewTextView()
	logEventTextPrmt.SetBackgroundColor(utils.BackGrayColor)
	logEventTextPrmt.SetTextAlign(tview.AlignLeft).SetDynamicColors(true)

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
	log("", "info")

	// Get current context to replace all variables
	_, currentContext := requestView.GetContext()
	currentContextValues := getOutput().Context.GetAllKeyValue(currentContext)

	URL := types.URL(requestView.GetURL()).
		ReplaceContext(makeRequestData.MapRequestHeaderKeyValue).
		ReplaceContext(currentContextValues)

	method := makeRequestData.Method
	contentType := makeRequestData.ContentType
	body := []byte(makeRequestData.Body)
	httpHeaderValues := makeRequestData.GetHTTPHeaderValues().ReplaceContext(currentContextValues)

	prefix := "[" + strconv.Itoa(rand.Intn(100)) + "] "

	HTTPClient, error := httpclient.Call(method, URL, contentType, body, httpHeaderValues, requestResponseView.Logger)
	if error != nil {
		event.PrintInfo(prefix + makeRequestData.ToLog(URL))
		event.PrintError(prefix + fmt.Sprint(error))

		requestResponseView.Logger(fmt.Sprint(error), "error")
	} else {
		event.PrintInfo(prefix + makeRequestData.ToLog(URL))

		responseData = fmt.Sprintf("%+s", HTTPClient.Body)
		if os.Getenv("GTTP_LOG_REQUEST") == "ON" {
			event.PrintInfo(prefix + responseData)
		}

		requestResponseView.Display(HTTPClient, responseData)
	}
}

func saveRequest(value models.MakeRequestData) {
	output = GetDataFromDisk(appPathFileName, log)

	// update
	output.AddOrReplace(value)

	SaveDataOnDisk(appPathFileName, output, log)
}

func removeRequest(value models.MakeRequestData) {
	output = GetDataFromDisk(appPathFileName, log)

	// update
	output.Remove(value)

	SaveDataOnDisk(appPathFileName, output, log)
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

func updateMDR(value models.MakeRequestData) {
	makeRequestData = value
}

func deleteMDR(value models.MakeRequestData) {
	removeRequest(value)
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
	output = GetDataFromDisk(appPathFileName, log)

	// update
	output.Config = value

	SaveDataOnDisk(appPathFileName, output, log)

	// refresh views
	refreshingConfig()
	refreshingTreeAPICpn()
}

func updateContext(value models.Context) {
	output = GetDataFromDisk(appPathFileName, log)

	// update
	output.Context = value

	SaveDataOnDisk(appPathFileName, output, log)

	refreshingContext()
}

func refreshingConfig() {
	for key, value := range event.AddListenerConfig {
		event.PrintTrace("App.refreshingConfig." + key)
		value(output.Config)
	}
}

func refreshingContext() {
	for key, value := range event.AddContextListener {
		event.PrintTrace("App.refreshingContext." + key)
		value(output.Context)
	}
}

func refreshingTreeAPICpn() {
	treeAPICpnt.Refresh()
}

func getOutput() models.Output {
	return output
}

// Log displays UI message to user.
func log(message string, status string) {
	if message != "" {
		logEventTextPrmt.SetText(utils.FormatLog(message, status))
	}
}
