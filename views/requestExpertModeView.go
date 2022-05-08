package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// RequestExpertModeView represents the request expert mode view
type RequestExpertModeView struct {
	App    *tview.Application
	AppCtx *models.AppCtx

	Labels map[string]string

	TitlePrmt  tview.Primitive
	ParentPrmt tview.Primitive
}

// NewRequestExpertModeView returns the view for the request expert mode
func NewRequestExpertModeView(app *tview.Application, ev *models.AppCtx) *RequestExpertModeView {
	labels := make(map[string]string)
	labels["menu_content_type_title"] = "Define specific \"Content-Type\""
	labels["menu_content_type_desc"] = "application/json,text/plain,multipart/f..."
	labels["menu_header_title"] = "Add request Header"
	labels["menu_header_desc"] = "& {param} url or {value} ex. context"
	labels["menu_body_title"] = "Add request Body"
	labels["menu_body_desc"] = ""
	labels["menu_preview_title"] = "Display request"
	labels["menu_preview_desc"] = ""

	labels["title"] = "Request Expert Mode"
	labels["requestPreview"] = "Request Preview"
	labels["headers"] = "Headers/Params"
	labels["headersPreview"] = "Headers/Params Preview"
	labels["key"] = "Key"
	labels["value"] = "Value"
	labels["body"] = "Body"
	labels["bodyPreview"] = "Body Preview"
	labels["contentType"] = "Content-Type"
	labels["contentTypePreview"] = "\"Content-Type\" list Preview"
	labels["add"] = "Add"
	labels["remove"] = "Remove"
	labels["projectName"] = "Project Name"
	labels["alias"] = "Alias"
	labels["method"] = "Method"
	labels["url"] = "URL"

	return &RequestExpertModeView{
		App:    app,
		AppCtx: ev,
		Labels: labels,
	}
}

// InitView build all components to display correctly the view
func (view *RequestExpertModeView) InitView() {
	mapMenuToFocusPrmt := make(map[string]tview.Primitive)

	// Make pages for each menu content
	pages := tview.NewPages()
	pages.SetBackgroundColor(utils.BackGrayColor)
	pages.AddPage("AddContentTypePage", view.makeAddContentTypePage(mapMenuToFocusPrmt), true, false)
	pages.AddPage("AddHeaderPage", view.makeAddHeaderPage(mapMenuToFocusPrmt), true, false)
	pages.AddPage("AddBodyPage", view.makeAddBodyPage(mapMenuToFocusPrmt), true, false)
	pages.AddPage("PreviewPage", view.makePreviewPage(), true, false)

	// Make menu
	menu := view.makeMenu(pages, mapMenuToFocusPrmt)

	// Make title
	titleAndMenuFlexPrmt := utils.MakeTitlePrmt(view.Labels["title"])
	titleAndMenuFlexPrmt.AddItem(menu, 0, 1, false)

	flexPrmt := tview.NewFlex()
	flexPrmt.AddItem(titleAndMenuFlexPrmt, 0, 1, false)
	flexPrmt.AddItem(tview.NewBox().SetBorder(false), 2, 0, false)
	flexPrmt.AddItem(pages, 0, 2, false)

	frame := tview.NewFrame(flexPrmt).SetBorders(0, 0, 0, 0, 0, 0)

	titleAndMenuFlexPrmt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Name() {
		case "Ctrl+Down":
			view.App.SetFocus(menu)
		}
		return event
	})

	// Display the "Preview" page menu by default
	menu.SetCurrentItem(menu.GetItemCount() - 1)
	pages.SwitchToPage("PreviewPage")
	view.App.SetFocus(menu)
	view.App.SetFocus(mapMenuToFocusPrmt["menu_preview"])

	// Don't forget!
	view.TitlePrmt = menu
	view.ParentPrmt = frame
}

func (view *RequestExpertModeView) makeMenu(pages *tview.Pages, mapMenuToFocusPrmt map[string]tview.Primitive) *tview.List {
	menu := tview.NewList().
		AddItem(view.Labels["menu_body_title"], view.Labels["menu_body_desc"], 'b', func() {
			pages.SwitchToPage("AddBodyPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_body"])
		}).
		AddItem(view.Labels["menu_content_type_title"], view.Labels["menu_content_type_desc"], 'c', func() {
			pages.SwitchToPage("AddContentTypePage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_content_type"])
		}).
		AddItem(view.Labels["menu_header_title"], view.Labels["menu_header_desc"], 'h', func() {
			pages.SwitchToPage("AddHeaderPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_header"])
		}).
		AddItem(view.Labels["menu_preview_title"], view.Labels["menu_preview_desc"], 'p', func() {
			pages.SwitchToPage("PreviewPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_preview"])
		})

	menu.
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(utils.BackGrayColor)

	return menu
}

func (view *RequestExpertModeView) makeAddContentTypePage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	contentTypeValues := utils.ContentTypeValues

	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	// Add "Content-Type" field
	formPrmt.AddDropDown(view.Labels["contentType"], contentTypeValues, 1, func(option string, index int) {
		makeRequestData := view.AppCtx.GetMDR()
		makeRequestData.ContentType = option

		// Update request
		view.updateMDR(makeRequestData)
	})

	// Add listener to refresh primitive when the MakeRequestData is changing...
	view.AppCtx.AddListenerMRD["requestExpertModeViewContentTypePage"] = func(makeRequestData models.MakeRequestData) {
		methodSelectedIndex := contentTypeValues.GetIndex(makeRequestData.ContentType)
		utils.GetDropDownFieldForm(formPrmt, view.Labels["contentType"]).SetCurrentOption(methodSelectedIndex)
	}

	// Make table prmt
	makeTablePrmt := func(values []string) *tview.Flex {
		table := tview.NewTable().SetBorders(false)
		table.SetBackgroundColor(utils.BackGrayColor)

		titlePrmt := tview.NewTextView()
		titlePrmt.SetText(view.Labels["contentTypePreview"])
		titlePrmt.SetTextColor(tcell.ColorGreen)
		titlePrmt.
			SetTextAlign(tview.AlignCenter).
			SetBackgroundColor(utils.BackGrayColor)

		flexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
		flexPrmt.AddItem(titlePrmt, 1, 0, false)
		flexPrmt.AddItem(tview.NewBox().SetBackgroundColor(utils.BackGrayColor), 1, 0, false)
		flexPrmt.AddItem(table, 0, 1, false)

		// Fill table with values
		var i = 1
		for _, value := range values {
			table.SetCell(i, 0, tview.NewTableCell(" ; "+value))
			i = i + 1
		}

		return flexPrmt
	}

	// Map menu with form
	mapMenuToFocusPrmt["menu_content_type"] = formPrmt

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(formPrmt, 0, 1, false)
	flex.AddItem(makeTablePrmt(utils.ContentTypeValues), 0, 2, false)

	return flex
}

func (view *RequestExpertModeView) makeAddHeaderPage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	// Display headers preview
	displayPreview := func(textView *tview.TextView) {
		sortedHeaderKeys := view.AppCtx.GetMDR().MapRequestHeaderKeyValue.ToSortedKeys()
		header := view.AppCtx.GetMDR().MapRequestHeaderKeyValue

		var sb strings.Builder
		for _, key := range sortedHeaderKeys {
			sb.WriteString("[" + utils.BlueColorName + "]" + key + "[white] " + header[key])
			sb.WriteString("\r\n\r\n")
		}
		textView.SetText(sb.String())
	}

	// Make preview prmt
	previewTitlePrmt := tview.NewTextView()
	previewTitlePrmt.SetText(view.Labels["headersPreview"])
	previewTitlePrmt.SetTextColor(tcell.ColorGreen)
	previewTitlePrmt.
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(utils.BackGrayColor)

	previewPrmt := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	previewPrmt.SetBackgroundColor(utils.BackGrayColor)

	previewFlexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	previewFlexPrmt.AddItem(previewTitlePrmt, 1, 0, false)
	previewFlexPrmt.AddItem(tview.NewBox().SetBackgroundColor(utils.BackGrayColor), 1, 0, false)
	previewFlexPrmt.AddItem(previewPrmt, 0, 1, false)

	// Make header form
	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	selectedEventDropDown := func(key string) {
		makeRequestData := view.AppCtx.GetMDR()
		value := makeRequestData.MapRequestHeaderKeyValue[key]

		if item := utils.GetInputFieldForm(formPrmt, view.Labels["key"]); item != nil {
			item.SetText(key)
		}

		if item := utils.GetInputFieldForm(formPrmt, view.Labels["value"]); item != nil {
			item.SetText(value)
		}
	}

	saveAndRefreshView := func(makeRequestData models.MakeRequestData) {
		// update object
		view.updateMDR(makeRequestData)

		headerKeys := makeRequestData.MapRequestHeaderKeyValue.ToSortedKeys()

		dropDrownPrmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["headers"])
		dropDrownPrmt.SetOptions(headerKeys, func(option string, index int) {
			selectedEventDropDown(option)
		})
		// Very important, fill the component with values before to SetCurrentOption
		dropDrownPrmt.SetCurrentOption(0)

		displayPreview(previewPrmt)
	}

	// Add "Headers" field
	formPrmt.AddDropDown(view.Labels["headers"], nil, 0, func(option string, index int) {
		selectedEventDropDown(option)
	})

	// Add "Key" field
	formPrmt.AddInputField(view.Labels["key"], "", 0, nil, nil)
	utils.AddInputFieldEventForm(formPrmt, view.Labels["key"])

	// Add "Value" field
	formPrmt.AddInputField(view.Labels["value"], "", 0, nil, nil)
	utils.AddInputFieldEventForm(formPrmt, view.Labels["value"])

	// Add "Add" button
	formPrmt.AddButton(view.Labels["add"], func() {
		makeRequestData := view.AppCtx.GetMDR()

		keyFieldPrmt := utils.GetInputFieldForm(formPrmt, view.Labels["key"])
		valueFieldPrmt := utils.GetInputFieldForm(formPrmt, view.Labels["value"])

		key := keyFieldPrmt.GetText()
		value := valueFieldPrmt.GetText()

		// add new value
		makeRequestData.MapRequestHeaderKeyValue[key] = value

		saveAndRefreshView(makeRequestData)
	})

	// Add "Remove" button
	formPrmt.AddButton(view.Labels["remove"], func() {
		utils.GetInputFieldForm(formPrmt, view.Labels["key"]).SetText("")
		utils.GetInputFieldForm(formPrmt, view.Labels["value"]).SetText("")

		dropDrownPrmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["headers"])

		makeRequestData := view.AppCtx.GetMDR()
		// delete value
		_, value := dropDrownPrmt.GetCurrentOption()
		delete(makeRequestData.MapRequestHeaderKeyValue, value)

		saveAndRefreshView(makeRequestData)
	})

	// Add listener to refresh primitive when the MakeRequestData is changing...
	view.AppCtx.AddListenerMRD["requestExpertModeViewHeaderPage"] = func(makeRequestData models.MakeRequestData) {
		utils.GetInputFieldForm(formPrmt, view.Labels["key"]).SetText("")
		utils.GetInputFieldForm(formPrmt, view.Labels["value"]).SetText("")

		saveAndRefreshView(view.AppCtx.GetMDR())
	}

	// Map menu with form
	mapMenuToFocusPrmt["menu_header"] = formPrmt

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(formPrmt, 0, 1, false)
	flex.AddItem(previewFlexPrmt, 0, 2, false)

	return flex
}

func (view *RequestExpertModeView) makeAddBodyPage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	// Make preview prmt
	previewTitlePrmt := tview.NewTextView()
	previewTitlePrmt.SetText(view.Labels["bodyPreview"])
	previewTitlePrmt.SetTextColor(tcell.ColorGreen)
	previewTitlePrmt.
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(utils.BackGrayColor)

	previewPrmt := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	previewPrmt.SetBackgroundColor(utils.BackGrayColor)

	previewFlexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	previewFlexPrmt.AddItem(previewTitlePrmt, 1, 0, false)
	previewFlexPrmt.AddItem(tview.NewBox().SetBackgroundColor(utils.BackGrayColor), 1, 0, false)
	previewFlexPrmt.AddItem(previewPrmt, 0, 1, false)

	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	// Add "Body" field
	formPrmt.AddInputField(view.Labels["body"], "", 0, nil, func(value string) {
		previewPrmt.SetText(value)
	})

	// Add generic events to inputField
	utils.AddInputFieldEventForm(formPrmt, view.Labels["body"])

	// Add "Add" button
	formPrmt.AddButton(view.Labels["add"], func() {
		makeRequestData := view.AppCtx.GetMDR()
		makeRequestData.Body = utils.GetInputFieldForm(formPrmt, view.Labels["body"]).GetText()

		// Update request
		view.updateMDR(makeRequestData)
		previewPrmt.SetText(makeRequestData.Body)
	})

	// Add listener to refresh primitive when the MakeRequestData is changing...
	view.AppCtx.AddListenerMRD["requestExpertModeViewBodyPage"] = func(makeRequestData models.MakeRequestData) {
		utils.GetInputFieldForm(formPrmt, view.Labels["body"]).SetText(makeRequestData.Body)
		previewPrmt.SetText(makeRequestData.Body)
	}

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(formPrmt, 0, 1, false)
	flex.AddItem(previewFlexPrmt, 0, 2, false)

	// Map menu with form
	mapMenuToFocusPrmt["menu_body"] = formPrmt

	return flex
}

func (view *RequestExpertModeView) makePreviewPage() *tview.Flex {
	titlePrmt := tview.NewTextView()
	titlePrmt.SetText(view.Labels["requestPreview"])
	titlePrmt.SetTextColor(tcell.ColorGreen)
	titlePrmt.
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(utils.BackGrayColor)

	previewPrmt := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		ScrollTo(0, 0)
	previewPrmt.SetBackgroundColor(utils.BackGrayColor)
	previewPrmt.Box.SetBorderPadding(1, 1, 1, 1)

	// Add listener to refresh primitive when the MakeRequestData is changing...
	view.AppCtx.AddListenerMRD["requestExpertModeViewPreviewPage"] = func(makeRequestData models.MakeRequestData) {
		view.displayPreview(previewPrmt, makeRequestData)
	}

	view.displayPreview(previewPrmt, view.AppCtx.GetMDR())

	flexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	flexPrmt.SetBorderPadding(1, 0, 0, 0)
	flexPrmt.AddItem(titlePrmt, 1, 0, false)
	flexPrmt.AddItem(tview.NewBox().SetBackgroundColor(utils.BackGrayColor), 1, 0, false)
	flexPrmt.AddItem(previewPrmt, 0, 1, false)

	return flexPrmt
}

func (view *RequestExpertModeView) displayPreview(textView *tview.TextView, makeRequestData models.MakeRequestData) {
	textView.SetText("")
	var sb strings.Builder

	sb.WriteString("[yellow]" + view.Labels["projectName"] + "[white]: " + makeRequestData.ProjectName)
	sb.WriteString("\r\n")
	sb.WriteString("[yellow]" + view.Labels["alias"] + "[white]: " + makeRequestData.Alias)
	sb.WriteString("\r\n\r\n")

	sb.WriteString("[yellow]" + view.Labels["method"] + "[white]: " + makeRequestData.Method.String())
	sb.WriteString("\r\n")
	sb.WriteString("[yellow]" + view.Labels["url"] + "[white]: " + makeRequestData.URL.ReplaceContext(makeRequestData.MapRequestHeaderKeyValue).String())
	sb.WriteString("\r\n\r\n")

	sb.WriteString("[yellow]" + view.Labels["contentType"] + "[white]: " + makeRequestData.ContentType)
	sb.WriteString("\r\n")
	sb.WriteString("[yellow]" + view.Labels["headers"] + ":\r\n")
	for k, v := range makeRequestData.MapRequestHeaderKeyValue {
		sb.WriteString("[" + utils.BlueColorName + "]" + k + "[white] " + v)
		sb.WriteString("\r\n")
	}
	sb.WriteString("\r\n")

	sb.WriteString("[yellow]" + view.Labels["body"] + ":")
	if makeRequestData.Body != "" {
		sb.WriteString("\r\n")
		sb.WriteString(makeRequestData.Body)
	}
	textView.SetText(sb.String())
}

func (view *RequestExpertModeView) updateMDR(data models.MakeRequestData) {
	view.AppCtx.UpdateMDR(data)
	if update, is := view.AppCtx.AddListenerMRD["requestExpertModeViewPreviewPage"]; is {
		update(data)
	}
}
