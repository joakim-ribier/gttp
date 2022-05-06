package views

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	components "github.com/joakim-ribier/gttp/components/tree"
	"github.com/joakim-ribier/gttp/models"
	"github.com/joakim-ribier/gttp/utils"
	"github.com/rivo/tview"
)

// SettingsView represents the application settings view
type SettingsView struct {
	App   *tview.Application
	Event *models.Event

	Labels map[string]string

	TitlePrmt  tview.Primitive
	ParentPrmt tview.Primitive
}

// NewSettingsView returns the settings view of the app
func NewSettingsView(app *tview.Application, ev *models.Event) *SettingsView {
	var legendSB strings.Builder
	legendSB.WriteString("[" + utils.GreenColorName + "]Update the display format of the API(s) tree.\r\n\r\n")
	legendSB.WriteString("Change the patterns order to update the view:\r\n\r\n")
	legendSB.WriteString("* {m}       => method: \"GET\"\r\n")
	legendSB.WriteString("* {url}     => url: \"http//...\"\r\n")
	legendSB.WriteString("* {u}       => short url: \"/stats\"\r\n")
	legendSB.WriteString("* {a}|{u}   => alias or short url\r\n")
	legendSB.WriteString("* {a}|{url} => alias or url\r\n\r\n")
	legendSB.WriteString("* {color}   => default color\r\n\r\n")

	legendSB.WriteString("Example:\r\n\r\n")
	legendSB.WriteString(tview.Escape("{color}[::b] {m} [-:black:-] [white]{a}|{u}") + "\r\n\r\n\r\n")
	legendSB.WriteString(tview.Escape("For more details ([::b] or [-:black:-]):") + "\r\n\r\n")
	legendSB.WriteString("@see [" + utils.BlueColorName + "]" + utils.GitHubTViewURL)

	var gttpPageSB strings.Builder
	gttpPageSB.WriteString("[" + utils.GreenColorName + "]Go Rich Http Client\r\n\r\n")
	gttpPageSB.WriteString("@see [" + utils.BlueColorName + "]https://github.com/joakim-ribier/gttp")

	var executePageSB strings.Builder
	executePageSB.WriteString("[" + utils.GreenColorName + "]Execute http request\r\n\r\n")
	executePageSB.WriteString("Choose a request (" + string(9658) + " " + utils.SelectAPIShortcut + ") and press (" + utils.ExecuteShortcut + ") to execute it.\r\n\n")
	executePageSB.WriteString("* Select the execution context, depend on environment settings (" + utils.SettingsShortcut + ").")

	labels := make(map[string]string)
	labels["title"] = "Application Settings"

	labels["menu_tree_format_title"] = "Change tree API(s) format"
	labels["menu_tree_format_desc"] = "Update the display format of the API(s) tree"
	labels["menu_tree_overview_title"] = "Example of tree formatting"

	labels["menu_env_title"] = "Environment"
	labels["menu_env_desc"] = "Add variables for specific env"

	labels["menu_man_title"] = "man " + strings.ToUpper(utils.Title)
	labels["menu_man_desc"] = "Documentation..."
	labels["menu_man_gttp_title"] = "About " + strings.ToUpper(utils.Title)
	labels["menu_man_gttp_page"] = gttpPageSB.String()
	labels["menu_man_execute_title"] = "Execute request"
	labels["menu_man_execute_page"] = executePageSB.String()

	labels["add"] = "Add"
	labels["description"] = legendSB.String()
	labels["env"] = "Env."
	labels["envs"] = "Env."
	labels["new_env"] = "New env."
	labels["overview"] = "Overview"
	labels["patterns"] = "Pattern"
	labels["remove"] = "Remove"
	labels["save"] = "Save"
	labels["variables"] = "Variables"
	labels["value"] = "Value"
	labels["variable"] = "Variable"

	return &SettingsView{
		App:    app,
		Event:  ev,
		Labels: labels,
	}
}

// InitView build all components to display correctly the view
func (view *SettingsView) InitView() {
	mapMenuToFocusPrmt := make(map[string]tview.Primitive)

	// Pages for each menu content
	pages := tview.NewPages()
	pages.SetBackgroundColor(utils.BackGrayColor)
	pages.AddPage("EnvPage", view.makeEnvPage(mapMenuToFocusPrmt), true, false)
	pages.AddPage("APITreeFormatPage", view.makeAPITreeFormatPage(mapMenuToFocusPrmt), true, false)
	pages.AddPage("ManPage", view.makeManPage(mapMenuToFocusPrmt), true, false)

	// Menu
	menu := view.makeMenu(pages, mapMenuToFocusPrmt)

	// Title
	titleAndMenuFlexPrmt := utils.MakeTitlePrmt(view.Labels["title"])
	titleAndMenuFlexPrmt.AddItem(menu, 0, 1, false)

	flexPrmt := tview.NewFlex()
	flexPrmt.AddItem(titleAndMenuFlexPrmt, 0, 1, false)
	flexPrmt.AddItem(tview.NewBox().SetBorder(false), 2, 0, false)
	flexPrmt.AddItem(pages, 0, 2, false)

	frame := tview.NewFrame(flexPrmt).SetBorders(0, 0, 0, 0, 0, 0)

	// Display the "man page" menu
	menu.SetCurrentItem(menu.GetItemCount() - 1)
	pages.SwitchToPage("ManPage")
	view.App.SetFocus(menu)
	view.App.SetFocus(mapMenuToFocusPrmt["menu_man"])

	// Don't forget!
	view.TitlePrmt = menu
	view.ParentPrmt = frame
}

func (view *SettingsView) makeMenu(pages *tview.Pages, mapMenuToFocusPrmt map[string]tview.Primitive) *tview.List {
	menu := tview.NewList().
		AddItem(view.Labels["menu_env_title"], view.Labels["menu_env_desc"], 'e', func() {
			pages.SwitchToPage("EnvPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_env"])
		}).
		AddItem(view.Labels["menu_tree_format_title"], view.Labels["menu_tree_format_desc"], 't', func() {
			pages.SwitchToPage("APITreeFormatPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_tree_format"])
		}).
		AddItem(view.Labels["menu_man_title"], view.Labels["menu_man_desc"], 'z', func() {
			pages.SwitchToPage("ManPage")
			view.App.SetFocus(mapMenuToFocusPrmt["menu_man"])
		})

	menu.
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(utils.BackGrayColor)

	return menu
}

func (view *SettingsView) makeManPage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	makeGTTPPage := func() *tview.TextView {
		prmt := tview.NewTextView().SetDynamicColors(true)
		prmt.SetText(view.Labels["menu_man_gttp_page"])
		prmt.SetBackgroundColor(utils.BackGrayColor)
		return prmt
	}

	makeExecuteRequestPage := func() *tview.TextView {
		prmt := tview.NewTextView().SetDynamicColors(true)
		prmt.SetText(view.Labels["menu_man_execute_page"])
		prmt.SetBackgroundColor(utils.BackGrayColor)
		return prmt
	}

	pages := tview.NewPages()
	pages.SetBackgroundColor(utils.BackGrayColor)
	pages.AddPage("GTTPPage", makeGTTPPage(), true, true)
	pages.AddPage("ExecuteRequestPage", makeExecuteRequestPage(), true, false)

	menu := tview.NewList().
		AddItem(view.Labels["menu_man_gttp_title"], "", 'a', func() {
			pages.SwitchToPage("GTTPPage")
		}).
		AddItem(view.Labels["menu_man_execute_title"], "", 'b', func() {
			pages.SwitchToPage("ExecuteRequestPage")
		})

	menu.
		SetBorderPadding(1, 1, 1, 1).
		SetBackgroundColor(utils.BackGrayColor)

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(menu, 0, 1, false)
	flex.AddItem(pages, 0, 3, false)

	mapMenuToFocusPrmt["menu_man"] = menu

	return flex
}

func (view *SettingsView) makeAPITreeFormatPage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	overview := func(treeAPICpnt *components.TreeCpnt, form *tview.Form) {
		view.Event.PrintTrace("SettingsView.makeAPITreeFormatPage{...}.overview")
		prmt := utils.GetInputFieldForm(form, view.Labels["patterns"])

		values := []models.MakeRequestData{
			models.SimpleMakeRequestData("GET", "https://api.weather.com/day/{city}", "", "Get the city weather for today (alias)"),
			models.SimpleMakeRequestData("GET", "https://api.bank.com/balance/{account_number}", "", "Get my account balance (alias)"),
			models.SimpleMakeRequestData("GET", "https://api.github.com/zen", "Github", "~/zen (alias)"),
			models.SimpleMakeRequestData("POST", "https://api.github.com/new/commit", "Github", ""),
			models.SimpleMakeRequestData("GET", "https://api.jira.com/list/ticket/{project}", "Jira", "List all tickets of project (alias)"),
			models.SimpleMakeRequestData("PUT", "https://api.jira.com/update/ticket", "Jira", "Update a ticket (alias)"),
			models.SimpleMakeRequestData("DELETE", "https://api.jira.com/delete/ticket", "Jira", ""),
		}

		treeAPICpnt.RefreshWithPattern(prmt.GetText(), view.Event.GetOutput().UpdateMakeRequestData(values))
	}

	// Add tree APIs component

	treeAPICpnt := components.NewTreeCpnt(view.App, view.Event)
	tree := treeAPICpnt.Make(nil, nil)
	tree.SetBackgroundColor(utils.BackGrayColor)
	treeAPICpnt.UpdateTitle(view.Labels["menu_tree_overview_title"])

	// Description prmt
	descPrmt := tview.NewTextView().SetDynamicColors(true)
	descPrmt.SetText(view.Labels["description"])
	descPrmt.SetBackgroundColor(utils.BackGrayColor)

	// Form to update pattern prmt
	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	mapMenuToFocusPrmt["menu_tree_format"] = formPrmt

	// New field - "Patterns"
	formPrmt.AddInputField(view.Labels["patterns"], "", 0, nil, nil)

	// Add generic events to inputField
	utils.AddInputFieldEventForm(formPrmt, view.Labels["patterns"])

	// New field - "see the overview"
	formPrmt.AddButton(view.Labels["overview"], func() {
		overview(treeAPICpnt, formPrmt)
	})

	// New field - "Save"
	formPrmt.AddButton(view.Labels["save"], func() {
		prmt := utils.GetInputFieldForm(formPrmt, view.Labels["patterns"])

		config := view.Event.GetConfig()
		config.Pattern = prmt.GetText()

		view.Event.UpdateConfig(config)
	})

	formFlexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	formFlexPrmt.AddItem(descPrmt, 22, 0, false)
	formFlexPrmt.AddItem(formPrmt, 0, 1, false)

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(formFlexPrmt, 0, 1, false)
	flex.AddItem(tree, 0, 1, false)

	view.Event.AddListenerConfig["makeAPITreeFormatPage"] = func(data models.Config) {
		view.Event.PrintTrace("SettingsView.makeAPITreeFormatPage{...}.listener")

		prmt := utils.GetInputFieldForm(formPrmt, view.Labels["patterns"])
		prmt.SetText(view.Event.GetConfig().Pattern)

		overview(treeAPICpnt, formPrmt)
	}

	return flex
}

func (view *SettingsView) makeEnvPage(mapMenuToFocusPrmt map[string]tview.Primitive) *tview.Flex {
	// Overview displays the selected environment
	overview := func(table *tview.Table, env string, data map[string]string) {
		table.Clear()

		// Add Env value
		table.SetCell(0, 0, tview.NewTableCell(view.Labels["env"]).SetTextColor(tcell.ColorYellow))
		table.SetCell(0, 1, tview.NewTableCell(env))

		// Add break line
		table.SetCell(1, 0, tview.NewTableCell(""))

		// Add all variables
		var i = 2
		for key, value := range data {
			table.SetCell(i, 0, tview.NewTableCell(key).SetTextColor(tcell.ColorYellow))
			table.SetCell(i, 1, tview.NewTableCell(value))
			i = i + 1
		}
	}

	// Add Overview table
	table := tview.NewTable().SetBorders(false)
	table.SetBackgroundColor(utils.BackGrayColor)

	// Form
	formPrmt := tview.NewForm()
	formPrmt.SetBorder(false)
	formPrmt.SetBackgroundColor(utils.BackGrayColor)

	mapMenuToFocusPrmt["menu_env"] = formPrmt

	selectVariableDropDownPrmtOption := func(variable string) {
		_, env := utils.GetDropDownFieldForm(formPrmt, view.Labels["envs"]).GetCurrentOption()

		contextVariable := view.Event.GetOutput().Context.FindVariableByEnv(env, variable)

		utils.GetInputFieldForm(formPrmt, view.Labels["variable"]).SetText(contextVariable.Variable)
		utils.GetInputFieldForm(formPrmt, view.Labels["value"]).SetText(contextVariable.Value)
		utils.GetInputFieldForm(formPrmt, view.Labels["new_env"]).SetText("")
	}

	selectVariablesDropDownPrmtOption := func(env string, variable string) {
		variables := view.Event.GetOutput().Context.GetAllVariables(env)

		variablesDropDownPrmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["variables"])
		variablesDropDownPrmt.SetOptions(variables, func(option string, index int) {
			selectVariableDropDownPrmtOption(option)
		})

		displayDefault := func() {
			variablesDropDownPrmt.SetCurrentOption(0)
			_, variable := variablesDropDownPrmt.GetCurrentOption()
			selectVariableDropDownPrmtOption(variable)
		}

		if env != "" && variable != "" {
			index := variables.GetIndex(variable)
			if index != -1 {
				variablesDropDownPrmt.SetCurrentOption(index)
				selectVariableDropDownPrmtOption(variable)
			} else {
				displayDefault()
			}
		} else {
			displayDefault()
		}

		overview(table, env, view.Event.GetOutput().Context.GetAllKeyValue(env))
	}

	refreshContext := func(env string, variable string) {
		view.Event.PrintTrace("SettingsView.makeAPITreeFormatPage{...}.refreshContext")

		envs := view.Event.GetOutput().Context.GetEnvsName()

		envsDropDownPrmt := utils.GetDropDownFieldForm(formPrmt, view.Labels["envs"])
		envsDropDownPrmt.SetOptions(envs, func(option string, index int) {
			selectVariablesDropDownPrmtOption(option, "")
		})

		displayDefault := func() {
			index := envs.GetIndex("default")
			envsDropDownPrmt.SetCurrentOption(index)
			selectVariablesDropDownPrmtOption("default", "")
		}

		if env != "" {
			index := envs.GetIndex(env)
			if index != -1 {
				envsDropDownPrmt.SetCurrentOption(index)
				selectVariablesDropDownPrmtOption(env, variable)
			} else {
				displayDefault()
			}
		} else {
			displayDefault()
		}
	}

	// New field - "Env(s)"
	formPrmt.AddDropDown(view.Labels["envs"], nil, 0, nil)

	// New field - "Variables"
	formPrmt.AddDropDown(view.Labels["variables"], nil, 0, nil)

	// New field - "Variable"
	formPrmt.AddInputField(view.Labels["variable"], "", 0, nil, nil)
	// New field - "Value"
	formPrmt.AddInputField(view.Labels["value"], "", 0, nil, nil)
	// New field - "New Env"
	formPrmt.AddInputField(view.Labels["new_env"], "", 0, nil, nil)

	// Add generic events to inputField
	utils.AddInputFieldEventForm(formPrmt, view.Labels["variable"])
	utils.AddInputFieldEventForm(formPrmt, view.Labels["value"])
	utils.AddInputFieldEventForm(formPrmt, view.Labels["new_env"])

	// New field - "Add"
	formPrmt.AddButton(view.Labels["add"], func() {
		_, env := utils.GetDropDownFieldForm(formPrmt, view.Labels["envs"]).GetCurrentOption()
		variable := utils.GetInputFieldForm(formPrmt, view.Labels["variable"]).GetText()
		value := utils.GetInputFieldForm(formPrmt, view.Labels["value"]).GetText()
		newEnv := utils.GetInputFieldForm(formPrmt, view.Labels["new_env"]).GetText()
		if newEnv != "" {
			env = newEnv
		}
		context := view.Event.GetOutput().Context
		context.Add(env, variable, value)
		view.Event.UpdateContext(context)

		refreshContext(env, variable)
	})

	// New field - "Remove"
	formPrmt.AddButton(view.Labels["remove"], func() {
		_, env := utils.GetDropDownFieldForm(formPrmt, view.Labels["envs"]).GetCurrentOption()
		_, variable := utils.GetDropDownFieldForm(formPrmt, view.Labels["variables"]).GetCurrentOption()
		if variable != "" {
			context := view.Event.GetOutput().Context
			context.Remove(env, variable)
			view.Event.UpdateContext(context)

			refreshContext(env, "")
		}
	})

	flex := tview.NewFlex()
	flex.SetBorderPadding(1, 1, 1, 1)
	flex.AddItem(formPrmt, 0, 1, false)
	flex.AddItem(table, 0, 2, false)

	view.Event.AddContextListener["makeEnvPage"] = func(data models.Context) {
		refreshContext("", "")
	}

	return flex
}
