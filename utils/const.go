package utils

import (
	"strings"

	"github.com/gdamore/tcell"
	"github.com/joakim-ribier/gttp/core"
)

// Represents App global information
const (
	GitHubTViewURL = "https://godoc.org/github.com/rivo/tview"
	Title          = "G-TTP"
	Subtitle       = Title + " - Go HTTP Client for Terminal UIs"
	TitleShortcuts = "Press [blue::ub]Escape[white::-] or Ctrl+[blue::ub]Q[white::-] to exit"
	GitHubLink     = "~//github.com/joakim-ribier/gttp"
	TitleAPIText   = `
	 ____           _____ _____ ____  
	/ ___|         |_   _|_   _|  _ \ 
   | |  _   _____    | |   | | | |_) |
   | |_| | |_____|   | |   | | |  __/ 
	\____|           |_|   |_| |_|    
									 
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
	ShortcutD  = "Ctrl+[blue::ub]D[white::-] Response View"
	ShortcutE  = "Ctrl+[blue::ub]E[white::-] Execute"
	ShortcutF  = "Ctrl+[blue::ub]F[white::-] Make Request"
	ShortcutH  = "Ctrl+[blue::ub]H[white::-] Expert Mode"
	ShortcutJ  = "Ctrl+[blue::ub]J[white::-] Select API"
	ShortcutO  = "Ctrl+[blue::ub]O[white::-] Settings"
	ShortcutQ  = "Ctrl+[blue::ub]Q[white::-] Exit"
	ShortcutR  = "Ctrl+[blue::ub]R[white::-] Request Header View"
	ShortcutDC = "Ctrl+[blue::ub]C[white::-] Copy Response"
	ShortcutDA = "Ctrl+[blue::ub]A[white::-] Copy All (log)"

	ShortcutHSubMenu  = ShortcutH + " >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"
	ShortcutOSubMenu  = ShortcutO + " >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"
	ShortcutSRSubMenu = " Save Request >> Ctrl+[blue::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"

	ShortcutPressEscape = "Press Escape"
	ShortcutSeparator   = " | "
)

// Represents data shortcuts to display to the user

var (
	MainShortcutsText        = strings.Join([]string{ShortcutJ, ShortcutE, ShortcutF, ShortcutH, ShortcutD, ShortcutO, ShortcutQ}, ShortcutSeparator)
	ResultShortcutsText      = strings.Join([]string{ShortcutR, ShortcutDC, ShortcutDA, ShortcutPressEscape}, ShortcutSeparator)
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
