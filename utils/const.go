package utils

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/core"
)

// Represents App global information
const (
	GitHubTViewURL = "https://godoc.org/github.com/rivo/tview"
	Title          = "GTTP"
	Subtitle       = Title + " - Go HTTP Client for Terminal UIs"
	TitleShortcuts = "Press [blue::ub]Escape[white::-] or Ctrl+[blue::ub]Q[white::-] to exit"
	GitHubLink     = "~//github.com/joakim-ribier/gttp"
	TitleAPIText   = `
   /\/|   ____ _____ _____ ____    /\/|
  |/\/   / ___|_   _|_   _|  _ \  |/\/
        | |  _  | |   | | | |_) |
        | |_| | | |   | | |  __/
         \____| |_|   |_| |_|
                            o
                             o
`
)

// Represents App color list
const (
	BackColor          = tcell.ColorDefault
	BackColorPrmt      = tcell.ColorGray
	BackFocusColorPrmt = tcell.ColorSilver
	TitleBackColor     = tcell.ColorBlue
	TitleTextColor     = tcell.ColorWhite
)

// Represents App shortcuts list
const (
	ShortcutQ = "Ctrl+[blue::ub]Q[white::-] Exit"
	ShortcutD = "Ctrl+[blue::ub]D[white::-] Result View"
	ShortcutR = "Ctrl+[blue::ub]R[white::-] Request"
	ShortcutE = "Ctrl+[blue::ub]E[white::-] Execute"
	ShortcutJ = "Ctrl+[blue::ub]J[white::-] Select API"

	ShortcutH        = "Ctrl+[blue::ub]H[white::-] Expert Mode"
	ShortcutHSubMenu = ShortcutH + " >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"

	ShortcutO        = "Ctrl+[blue::ub]O[white::-] Settings"
	ShortcutOSubMenu = ShortcutO + " >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"

	ShortcutDC = "Ctrl+[blue::ub]C[white::-] Copy Result"
	ShortcutDA = "Ctrl+[blue::ub]A[white::-] Copy All"

	ShortcutSRSubMenu = " Save Request >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"

	ShortcutPressEscape = "Press Escape"
	ShortcutSeparator   = " | "
)

// Represents data shortcuts to display to the user
var (
	MainShortcutsText        = strings.Join([]string{ShortcutJ, ShortcutE, ShortcutR, ShortcutH, ShortcutD, ShortcutO, ShortcutQ}, ShortcutSeparator)
	ResultShortcutsText      = strings.Join([]string{ShortcutDC, ShortcutDA, ShortcutPressEscape}, ShortcutSeparator)
	ExpertModeShortcutsText  = strings.Join([]string{ShortcutHSubMenu, ShortcutPressEscape}, ShortcutSeparator)
	SettingsShortcutsText    = strings.Join([]string{ShortcutOSubMenu, ShortcutPressEscape}, ShortcutSeparator)
	SaveRequestShortcutsText = strings.Join([]string{ShortcutSRSubMenu, ShortcutPressEscape}, ShortcutSeparator)
)

// Represents data to make a new request
var (
	MethodValues      = core.StringSlice{"GET", "POST", "PUT"}
	ContentTypeValues = core.StringSlice{
		"application/javascript",
		"application/json",
		"application/x-www-form-urlencoded",
		"application/xml",
		"application/zip",
		"application/pdf",
		"application/sql",
		"application/graphql",
		"application/ld+json",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.oasis.opendocument.text",
		"audio/mpeg",
		"audio/ogg",
		"multipart/form-data",
		"text/css",
		"text/html",
		"text/xml",
		"text/csv",
		"text/plain",
		"image/png",
		"image/jpeg",
		"image/gif",
	}
)
