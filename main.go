package main

import (
	"fmt"
	"strings"

	"github.com/joakim-ribier/gttp/utils"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const logo = `
  ____           _____ _____ ____
 / ___|         |_   _|_   _|  _ \
| |  _   _____    | |   | | | |_) |
| |_| | |_____|   | |   | | |  __/
 \____|           |_|   |_| |_|
        o   ^__^
         o  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
`

func main() {
	theme()

	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			app.Stop()
			App()
		}
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(makePrmt(), true).Run(); err != nil {
		panic(err)
	}
}

func makePrmt() tview.Primitive {
	// Add logo prmt
	logoWidth, logoHeight := computeStringSize(logo)
	logoBox := tview.NewTextView().SetTextColor(tcell.ColorGreen)
	fmt.Fprint(logoBox, logo)

	// Add GitHub link + shortcuts prmts
	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(utils.GitHubLink, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(utils.TitleShortcuts, true, tview.AlignCenter, tcell.ColorWhite)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	return flex
}

func computeStringSize(value string) (int, int) {
	lines := strings.Split(value, "\n")
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	return width, len(lines)
}

// override default tview theme
func theme() {
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    utils.BackColor,
		ContrastBackgroundColor:     tcell.GetColor(utils.BlueColorName),
		MoreContrastBackgroundColor: tcell.ColorGreen,
		BorderColor:                 tcell.ColorWhite,
		TitleColor:                  tcell.ColorWhite,
		GraphicsColor:               tcell.ColorWhite,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorYellow,
		TertiaryTextColor:           tcell.GetColor(utils.GreenColorName),
		InverseTextColor:            tcell.ColorBlue,
		ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
	}
}
