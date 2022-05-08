package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// BuildModal builds simple modal.
func BuildModal(p tview.Primitive, width, height int) tview.Primitive {
	cpnt := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(p, height, 1, false).
		AddItem(nil, 0, 1, false)

	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(cpnt, width, 1, false).
		AddItem(nil, 0, 1, false)
}

// BuildYesNoModal builds a simple "no / yes" modal.
func BuildYesNoModal(
	title string,
	text string,
	no func(),
	yes func(),
	focus func(tview.Primitive)) tview.Primitive {

	titleTView := tview.NewTextView()
	titleTView.SetTextColor(tcell.ColorWhite)
	titleTView.SetText(title)

	form := tview.NewForm()
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlN:
			no()
			return nil
		case tcell.KeyCtrlY:
			yes()
		}
		return event
	})

	form.SetButtonsAlign(tview.AlignCenter)
	form.AddButton("[::ub]N[-:-:-]o", func() {
		no()
	})

	form.AddButton("[::ub]Y[-:-:-]es", func() {
		yes()
	})

	flexPrmt := tview.NewFlex().SetDirection(tview.FlexRow)
	flexPrmt.SetBorder(true).SetTitle(text)
	flexPrmt.AddItem(tview.NewBox(), 1, 0, false)
	flexPrmt.AddItem(titleTView, 1, 0, false)
	flexPrmt.AddItem(form, 0, 1, true)

	focus(form)

	return BuildModal(flexPrmt, 45, 7)
}
