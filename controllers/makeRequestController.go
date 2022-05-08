package controllers

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/joakim-ribier/gttp/actions"
	"github.com/joakim-ribier/gttp/httpclient"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/models/types"
	"github.com/joakim-ribier/gttp/services"
	"github.com/joakim-ribier/gttp/views"
	"github.com/rivo/tview"
)

var (
	logRequestOn = os.Getenv("GTTP_LOG_REQUEST") == "ON"
)

type MakeRequestController struct {
	View   *views.MakeRequestView
	Action *actions.MakeRequestAction

	// services
	AppDataService *services.ApplicationDataService

	// models
	AppCtx *models.AppCtx

	// lib
	App *tview.Application
}

func NewMakeRequestController(
	app *tview.Application,
	appDataService *services.ApplicationDataService,
	ctx *models.AppCtx,
	action *actions.MakeRequestAction) *MakeRequestController {

	return &MakeRequestController{
		App:            app,
		AppCtx:         ctx,
		View:           nil,
		AppDataService: appDataService,
		Action:         action,
	}
}

// Draw contructs and initializes the view.
func (c *MakeRequestController) Draw() tview.Primitive {
	if c.View == nil {
		c.View = views.NewMakeRequestView(c.App, c.AppCtx, c.saveC, c.deleteC)
	}
	c.View.InitView(c.Execute, c.ExpertMode, c.New)
	return c.View.RootPrmt
}

// saveC reloads file and saves/updates the current request.
func (c *MakeRequestController) saveC(callback func()) {
	// reload the data from file to save only the current updated request
	output := c.AppDataService.Load()

	output.AddOrReplace(c.AppCtx.GetMDR())

	c.AppDataService.Save(output)
	c.AppCtx.RefreshViews("all")

	callback()
}

// deleteC reloads file and deletes the current request.
func (c *MakeRequestController) deleteC(callback func()) {
	// reload the data from file to remove only the current request
	output := c.AppDataService.Load()

	output.Remove(c.AppCtx.GetMDR())
	c.AppCtx.UpdateMDR(models.EmptyMakeRequestData())

	c.AppDataService.Save(output)
	c.AppCtx.RefreshViews("all")

	callback()
}

// Execute calls the request and display the response.
func (c *MakeRequestController) Execute() {
	makeRequestData := c.AppCtx.GetMDR()
	prefix := "[" + strconv.Itoa(rand.Intn(100)) + "] "

	// Get current context to replace all variables
	_, currentContext := c.View.GetContext()
	currentContextValues := c.AppCtx.GetOutput().Context.GetAllKeyValue(currentContext)

	URL := types.URL(c.View.GetURL()).
		ReplaceContext(makeRequestData.MapRequestHeaderKeyValue).
		ReplaceContext(currentContextValues)

	method := makeRequestData.Method
	contentType := makeRequestData.ContentType
	body := []byte(makeRequestData.Body)
	httpHeaderValues := makeRequestData.GetHTTPHeaderValues().ReplaceContext(currentContextValues)

	HTTPClient, error := httpclient.Call(method, URL, contentType, body, httpHeaderValues, c.Action.DisplayErrorRequest)
	if error != nil {
		c.AppCtx.PrintInfo(prefix + makeRequestData.ToLog(URL))
		c.AppCtx.PrintError(prefix + fmt.Sprint(error))

		c.Action.DisplayErrorRequest(fmt.Sprint(error), "error")
	} else {
		c.AppCtx.PrintInfo(prefix + makeRequestData.ToLog(URL))

		response := fmt.Sprintf("%+s", HTTPClient.Body)
		if logRequestOn {
			c.AppCtx.PrintInfo(prefix + response)
		}

		c.Action.DisplayResponse(HTTPClient, response)
	}
}

// Save displays saving/updating request view.
func (c *MakeRequestController) Save() {
	c.View.DisplaySaveView()
}

// Remove displays removing request view.
func (c *MakeRequestController) Remove() {
	c.View.DisplayRemoveView()
}

// New cleans and refreshs the current request view
func (c *MakeRequestController) New() {
	c.AppCtx.UpdateMDR(models.EmptyMakeRequestData())
	c.AppCtx.RefreshViews("request")
}

// ExpertMode displays expert request mode view
func (c *MakeRequestController) ExpertMode() {
	c.AppCtx.SwitchView("ExpertRequestView")
}
