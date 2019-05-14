package views

import (
	"strings"

	"github.com/joakim-ribier/gttp/httpclient"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// RequestResponseView represents the response of the request view
type RequestResponseView struct {
	App   *tview.Application
	Event *models.Event

	Labels    map[string]string
	LogBuffer string

	TitlePrmt    *tview.Flex
	ParentPrmt   tview.Primitive
	RequestPrmt  *tview.TextView
	ResponsePrmt *tview.TextView
}

// NewRequestResponseView returns the view for the request response view
func NewRequestResponseView(app *tview.Application, ev *models.Event) *RequestResponseView {
	labels := make(map[string]string)
	labels["title"] = "Execute Request"
	labels["http"] = "HTTP"
	labels["contentType"] = "Content-Type"
	labels["host"] = "Host"
	labels["contentLength"] = "Content-Length"
	labels["date"] = "Date"
	labels["referrerPolicy"] = "Referrer-Policy"
	labels["connection"] = "Connection"
	labels["status"] = "Status"

	return &RequestResponseView{
		App:       app,
		Event:     ev,
		Labels:    labels,
		LogBuffer: "",
	}
}

// InitView builds all components to display correctly the view
func (view *RequestResponseView) InitView() {
	view.ResponsePrmt = tview.NewTextView()
	view.ResponsePrmt.SetBackgroundColor(utils.BackGrayColor)
	view.ResponsePrmt.SetDynamicColors(true).SetScrollable(true)

	view.RequestPrmt = tview.NewTextView()
	view.RequestPrmt.SetBackgroundColor(utils.BackGrayColor).SetBorderPadding(0, 0, 0, 0)
	view.RequestPrmt.SetDynamicColors(true).SetScrollable(true).SetWrap(true)

	view.TitlePrmt = utils.MakeTitlePrmt(view.Labels["title"])
	view.TitlePrmt.AddItem(view.RequestPrmt, 0, 1, false)

	flex := tview.NewFlex()
	flex.AddItem(view.TitlePrmt, 0, 1, false)
	flex.AddItem(tview.NewBox().SetBorder(false), 2, 0, false)
	flex.AddItem(view.ResponsePrmt, 0, 2, false)

	view.ParentPrmt = tview.NewFrame(flex).SetBorders(0, 0, 0, 0, 0, 0)
}

// Display displays request & response data
func (view *RequestResponseView) Display(client *httpclient.HTTPClient, data string) {
	var sb strings.Builder
	format := func(key string, value string) string {
		if key == "" {
			return "[white]" + value
		}
		return "[" + utils.BlueColorName + "]" + key + "[white] " + value
	}

	// Request header
	sb.WriteString(format(client.Request.Method, client.Request.URL+" ["+utils.BlueColorName+"]"+view.Labels["http"]+"[white]/"+client.Request.HTTP))
	sb.WriteString("\r\n")
	for k, v := range client.HeadersRequest {
		if k != "Content-Type" {
			sb.WriteString(format(k, v))
			sb.WriteString("\r\n")
		}
	}

	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["host"], client.Request.Host))

	// Content-Type
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["contentType"], client.HeadersRequest["Content-Type"]))
	sb.WriteString("\r\n\r\n")

	// Body
	body := strings.Replace(client.Request.Body, " ", "", -1)
	body = strings.Replace(body, "\n", "", -1)
	body = strings.Replace(body, "\r", "", -1)
	sb.WriteString(format("", tview.Escape(body)))

	// Response header
	sb.WriteString("\r\n\r\n")
	sb.WriteString(format(view.Labels["http"]+"[white]/"+client.Response.HTTP, client.Response.Status))
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["contentLength"], client.Response.Contentlength))
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["contentType"], client.Response.ContentType))
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["date"], client.Response.Date))
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["referrerPolicy"], client.Response.Referrerpolicy))
	sb.WriteString("\r\n")
	sb.WriteString(format(view.Labels["connection"], client.Response.Connection))

	// Set request prmt text
	view.RequestPrmt.SetText(sb.String()).SetTextAlign(tview.AlignLeft)

	// Log if error status
	status := client.Response.StatusCode
	if (status < "200") || (status > "205") {
		view.Logger(view.Labels["status"]+": "+client.Response.Status, "error")
	}

	// Set the body response prmt text
	view.setResponsePrmtText(utils.FormatLog(data, "data"))
}

// Logger logs to the response prmt
func (view *RequestResponseView) Logger(message string, mode string) {
	view.setResponsePrmtText(utils.FormatLog(message, mode))
}

// ResetLogBuffer resets the log buffer
func (view *RequestResponseView) ResetLogBuffer() {
	view.LogBuffer = ""
}

func (view *RequestResponseView) setResponsePrmtText(data string) {
	view.ResponsePrmt.
		Clear().
		SetTextAlign(tview.AlignLeft)

	view.LogBuffer = view.LogBuffer + data + "\r\n"

	view.ResponsePrmt.SetText(view.LogBuffer)
}
