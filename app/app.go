package app

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/joakim-ribier/gttp/actions"
	"github.com/joakim-ribier/gttp/components"
	"github.com/joakim-ribier/gttp/controllers"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/services"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/joakim-ribier/gttp/views"
	"github.com/rivo/tview"
)

var (
	logFile   *os.File
	logOn     = os.Getenv("GTTP_LOG") == "ON"
	logLevels = os.Getenv("GTTP_LOG_LEVEL")

	app      *tview.Application
	ctx      *models.AppCtx
	rootPrmt *tview.Pages

	logEventTextPrmt     *tview.TextView
	shortcutInfoTextPrmt *tview.TextView
	pages                *tview.Pages
)

var (
	output models.Output

	responseData = ""

	makeRequestData           = models.EmptyMakeRequestData()
	mapFocusPrmtToShortutText = make(map[tview.Primitive]string)
	focusPrmts                = []*tview.TextView{}

	// List of controllers
	makeRequestController *controllers.MakeRequestController

	// List of services
	appDataService *services.ApplicationDataService

	// List of views of the application
	expertModeView      *views.RequestExpertModeView
	settingsView        *views.SettingsView
	requestResponseView *views.RequestResponseView

	// List of components of the application
	treeAPICpnt *components.TreeCpnt
)

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

	appDataService = services.NewApplicationDataService(getFilenameFromArgs(os.Args), log)

	ctx = models.NewAppCtx(
		getRootPrmt,
		getMDR,
		updateMDR,
		getConfig,
		updateConfig,
		getOutput,
		updateContext,
		PrintOut,
		refresh,
		switchPage)

	app = tview.NewApplication()
	rootPrmt = drawMainComponents(app)

	// Fixme: To delete
	mapFocusPrmtToShortutText[requestResponseView.ResponsePrmt] = utils.ResultShortcutsText
	mapFocusPrmtToShortutText[expertModeView.TitlePrmt] = utils.ExpertModeShortcutsText
	mapFocusPrmtToShortutText[settingsView.TitlePrmt] = utils.SettingsShortcutsText

	refresh("all")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		log("Shortcut: "+event.Name()+" - "+time.Now().Format(time.RFC850), "info")

		// disable all shortcuts (except for 'app.Stop()') if it's the root modal page which has focus
		if page, _ := rootPrmt.GetFrontPage(); page == "modal" && event.Key() != tcell.KeyCtrlQ {
			return event
		}

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
			makeRequestController.Remove()
		case tcell.KeyCtrlE:
			executeRequest()
		case tcell.KeyCtrlF:
			focusPrimitive(makeRequestController.View.FormPrmt, nil)
		case tcell.KeyCtrlH:
			switchPage("ExpertRequestView")
		case tcell.KeyCtrlJ:
			focusPrimitive(treeAPICpnt.RootPrmt, nil)
		case tcell.KeyCtrlN:
			makeRequestController.New()
		case tcell.KeyCtrlO:
			switchPage("SettingsView")
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlR:
			displayRequestResponseViewPage(requestResponseView.RequestPrmt)
		case tcell.KeyCtrlS:
			makeRequestController.Save()
		case tcell.KeyCtrlW:
			displayRequestResponseViewPage(requestResponseView.ResponsePrmt)
		case tcell.KeyEsc:
			focusPrimitive(logEventTextPrmt, nil)
		}
		return event
	})

	if err := app.SetRoot(rootPrmt, true).Run(); err != nil {
		panic(err)
	}
}

func drawMainComponents(app *tview.Application) *tview.Pages {

	drawLeftPanel := func() tview.Primitive {
		treeAPICpnt = components.NewTreeCpnt(app, ctx)
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
			expertModeView = views.NewRequestExpertModeView(app, ctx)
			expertModeView.InitView()

			return expertModeView.ParentPrmt
		}

		makeSettingsView := func() tview.Primitive {
			settingsView = views.NewSettingsView(app, ctx)
			settingsView.InitView()

			return settingsView.ParentPrmt
		}

		// build request response view
		requestResponseView = views.NewRequestResponseView(app, ctx)
		requestResponseView.InitView()

		focusPrmts = append(focusPrmts, requestResponseView.ResponsePrmt)
		focusPrmts = append(focusPrmts, requestResponseView.RequestPrmt)

		// build "make/execute request" controller
		makeRequestController = controllers.NewMakeRequestController(
			app,
			appDataService,
			ctx,
			actions.NewMakeRequestAction(requestResponseView.Display, requestResponseView.Logger))

		flex := tview.NewFlex().SetDirection(tview.FlexRow)
		flex.SetBorder(false)
		flex.SetBorderPadding(1, 0, 0, 0)

		pages = tview.NewPages()
		pages.SetBorder(false).SetBorderPadding(0, 1, 0, 0)

		pages.AddPage("RequestResponseViewPage", requestResponseView.ParentPrmt, true, false)
		pages.AddPage("RequestExpertModeViewPage", makeRequestExportModeView(), true, false)
		pages.AddPage("SettingsViewPage", makeSettingsView(), true, true)

		flex.AddItem(makeRequestController.Draw(), 9, 0, false)
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

	return tview.NewPages().AddPage("root", grid, true, true)
}

func displayRequestResponseViewPage(focusOn *tview.TextView) {
	pages.SwitchToPage("RequestResponseViewPage")
	focusPrimitive(focusOn, focusOn.Box)
}

func executeRequest() {
	displayRequestResponseViewPage(requestResponseView.ResponsePrmt)
	requestResponseView.ResetLogBuffer()
	log("", "info")

	makeRequestController.Execute()
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

func refreshMDRView(makeRequestData models.MakeRequestData) {
	updateMDR(makeRequestData)
	for _, value := range ctx.AddListenerMRD {
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
	output = appDataService.Load()

	// update
	output.Config = value

	appDataService.Save(output)

	// refresh views
	refreshingConfig()
	refreshingTreeAPICpn()
}

func updateContext(value models.Context) {
	output = appDataService.Load()

	// update
	output.Context = value

	appDataService.Save(output)

	refreshingContext()
}

func refreshingConfig() {
	for key, value := range ctx.AddListenerConfig {
		ctx.PrintTrace("App.refreshingConfig." + key)
		value(output.Config)
	}
}

func refreshingContext() {
	for key, value := range ctx.AddContextListener {
		ctx.PrintTrace("App.refreshingContext." + key)
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

func switchPage(page string) {
	switch page {
	case "ExpertRequestView":
		pages.SwitchToPage("RequestExpertModeViewPage")
		focusPrimitive(expertModeView.TitlePrmt, nil)
	case "SettingsView":
		pages.SwitchToPage("SettingsViewPage")
		focusPrimitive(settingsView.TitlePrmt, nil)
	}
}

func refresh(value string) {
	output = appDataService.Load()

	refreshMRDAllViews := func() {
		for _, value := range ctx.AddListenerMRD {
			value(makeRequestData)
		}
	}

	if value == "all" {
		refreshingTreeAPICpn()
		refreshingConfig()
		refreshingContext()
		refreshMRDAllViews()
	} else {
		if strings.Contains(value, "tree") {
			refreshingTreeAPICpn()
		}
		if strings.Contains(value, "config") {
			refreshingConfig()
		}
		if strings.Contains(value, "ctx") {
			refreshingContext()
		}
		if strings.Contains(value, "request") {
			refreshMRDAllViews()
		}
	}
}

func getRootPrmt() *tview.Pages {
	return rootPrmt
}
