package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GetDropDownFieldForm get a dropdown field by label from the form
func GetDropDownFieldForm(form *tview.Form, itemLabel string) *tview.DropDown {
	return form.GetFormItemByLabel(itemLabel).(*tview.DropDown)
}

// GetInputFieldForm get an inputfield field by label from the form
func GetInputFieldForm(form *tview.Form, itemLabel string) *tview.InputField {
	if item := form.GetFormItemByLabel(itemLabel); item != nil {
		return item.(*tview.InputField)
	} else {
		return nil
	}
}

// AddInputFieldEventForm add generic event to inputfield
func AddInputFieldEventForm(form *tview.Form, itemLabel string) {
	prmt := GetInputFieldForm(form, itemLabel)
	prmt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlP {
			text, _ := ReadFromClipboard()
			prmt.SetText(prmt.GetText() + text)
		}
		return event
	})
}

// MakeSeparator create a simple tview.NewBox
func MakeSeparator(color tcell.Color) *tview.Box {
	return tview.NewBox().
		SetBorder(false).
		SetBackgroundColor(color)
}

// MakeTitlePrmt builds a title widget
func MakeTitlePrmt(title string) *tview.Flex {
	titleBackColor := BackBlueColor
	titleTextColor := tcell.ColorWhite

	// Title
	titlePrmt := tview.NewTextView()
	titlePrmt.SetBackgroundColor(titleBackColor)
	titlePrmt.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText("[::b]" + title).
		SetTextColor(titleTextColor)

	flexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	flexPrmt.SetBackgroundColor(BackGrayColor)
	flexPrmt.AddItem(MakeSeparator(titleBackColor), 1, 0, false)
	flexPrmt.AddItem(titlePrmt, 1, 0, false)
	flexPrmt.AddItem(MakeSeparator(titleBackColor), 1, 0, false)
	flexPrmt.AddItem(MakeSeparator(BackColor), 1, 0, false)

	return flexPrmt
}
