package models

import "github.com/rivo/tview"

// Event contains all events to manage application data
type AppCtx struct {
	GetRootPrmt func() *tview.Pages

	AddListenerMRD     map[string]func(data MakeRequestData)
	AddListenerConfig  map[string]func(data Config)
	AddContextListener map[string]func(data Context)

	GetMDR    func() MakeRequestData
	UpdateMDR func(data MakeRequestData)

	GetConfig    func() Config
	UpdateConfig func(data Config)

	GetOutput func() Output

	UpdateContext func(data Context)

	RefreshViews func(views string)
	SwitchView   func(view string)

	printOut func(level string, value string)
}

// NewEvent makes a new event struct
func NewAppCtx(
	getRootPrmt func() *tview.Pages,
	getMDR func() MakeRequestData,
	upMDR func(data MakeRequestData),
	getConfig func() Config,
	updateConfig func(data Config),
	getOutput func() Output,
	updateContext func(data Context),
	printOut func(level string, value string),
	refreshViews func(views string),
	switchView func(views string)) *AppCtx {

	return &AppCtx{
		GetRootPrmt:        getRootPrmt,
		AddListenerMRD:     make(map[string]func(data MakeRequestData)),
		AddListenerConfig:  make(map[string]func(data Config)),
		AddContextListener: make(map[string]func(data Context)),
		UpdateMDR:          upMDR,
		GetMDR:             getMDR,
		GetConfig:          getConfig,
		UpdateConfig:       updateConfig,
		GetOutput:          getOutput,
		UpdateContext:      updateContext,
		printOut:           printOut,
		RefreshViews:       refreshViews,
		SwitchView:         switchView,
	}
}

// PrintInfo prints "info" log to file
func (ctx AppCtx) PrintInfo(value string) {
	ctx.printOut("info", value)
}

// PrintError prints "error" log to file
func (ctx AppCtx) PrintError(value string) {
	ctx.printOut("error", value)
}

// PrintDebug prints "debug" log to file
func (ctx AppCtx) PrintDebug(value string) {
	ctx.printOut("debug", value)
}

// PrintTrace prints "trace" log to file
func (ctx AppCtx) PrintTrace(value string) {
	ctx.printOut("trace", value)
}

// DisplayModal displays a model to the screen.
func (ctx AppCtx) DisplayModal(prmt tview.Primitive) {
	ctx.CloseModal()

	ctx.GetRootPrmt().AddPage("modal", prmt, true, true)
}

// CloseModal closes modal from the screen.
func (ctx AppCtx) CloseModal() {
	ctx.GetRootPrmt().RemovePage("modal")
}
