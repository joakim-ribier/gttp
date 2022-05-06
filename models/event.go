package models

// Event contains all events to manage application data
type Event struct {
	AddListenerMRD     map[string]func(data MakeRequestData)
	AddListenerConfig  map[string]func(data Config)
	AddContextListener map[string]func(data Context)

	GetMDR    func() MakeRequestData
	UpdateMDR func(data MakeRequestData)
	DeleteMDR func(data MakeRequestData)

	GetConfig    func() Config
	UpdateConfig func(data Config)

	GetOutput func() Output

	UpdateContext func(data Context)

	printOut func(level string, value string)
}

// NewEvent makes a new event struct
func NewEvent(
	getMDR func() MakeRequestData,
	upMDR func(data MakeRequestData),
	deleteMDR func(data MakeRequestData),
	getConfig func() Config,
	updateConfig func(data Config),
	getOutput func() Output,
	updateContext func(data Context),
	printOut func(level string, value string)) *Event {

	return &Event{
		AddListenerMRD:     make(map[string]func(data MakeRequestData)),
		AddListenerConfig:  make(map[string]func(data Config)),
		AddContextListener: make(map[string]func(data Context)),
		DeleteMDR:          deleteMDR,
		UpdateMDR:          upMDR,
		GetMDR:             getMDR,
		GetConfig:          getConfig,
		UpdateConfig:       updateConfig,
		GetOutput:          getOutput,
		UpdateContext:      updateContext,
		printOut:           printOut,
	}
}

// PrintInfo prints "info" log to file
func (evt Event) PrintInfo(value string) {
	evt.printOut("info", value)
}

// PrintError prints "error" log to file
func (evt Event) PrintError(value string) {
	evt.printOut("error", value)
}

// PrintDebug prints "debug" log to file
func (evt Event) PrintDebug(value string) {
	evt.printOut("debug", value)
}

// PrintTrace prints "trace" log to file
func (evt Event) PrintTrace(value string) {
	evt.printOut("trace", value)
}
