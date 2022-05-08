package utils

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/joakim-ribier/gttp/core"
)

// Represents App color list
const (
	BackColorName     = "#2B2B2B"
	BackBlueColorName = "#214283"
	BackGrayColorName = "#424445"
	BlueColorName     = "dodgerblue"
	GreenColorName    = "#629755"
)

var (
	BackColor     = tcell.GetColor(BackColorName)
	BackBlueColor = tcell.GetColor(BackBlueColorName)
	BackGrayColor = tcell.GetColor(BackGrayColorName)
)

// Represents App global information
const (
	TreePrmtTitle  = "Requests"
	GitHubTViewURL = "https://godoc.org/github.com/rivo/tview"
	Title          = "GTTP"
	Subtitle       = Title + " - Go HTTP Client for Terminal UIs"
	TitleShortcuts = "Press [" + BlueColorName + "::ub]Escape[white::-] or Ctrl+[" + BlueColorName + "::ub]Q[white::-] to exit"
	GitHubLink     = "~//github.com/joakim-ribier/gttp"
)

// Represents App shortcuts list
const (
	JShortcut = "Ctrl+[" + BlueColorName + "::ub]J[white::-]"
	EShortcut = "Ctrl+[" + BlueColorName + "::ub]E[white::-]"

	ExecuteShortcut   = EShortcut + " Execute"
	SelectAPIShortcut = JShortcut + " Select API"
	SettingsShortcut  = "Ctrl+[" + BlueColorName + "::ub]O[white::-] Settings"

	ShortcutD  = "Ctrl+[" + BlueColorName + "::ub]W[white::-] Response View"
	ShortcutF  = "Ctrl+[" + BlueColorName + "::ub]F[white::-] Make Request"
	ShortcutH  = "Ctrl+[" + BlueColorName + "::ub]H[white::-] Expert Mode"
	ShortcutQ  = "Ctrl+[" + BlueColorName + "::ub]Q[white::-] Exit"
	ShortcutR  = "Ctrl+[" + BlueColorName + "::ub]R[white::-] Request Header View"
	ShortcutDC = "Ctrl+[" + BlueColorName + "::ub]C[white::-] Copy Response"
	ShortcutDA = "Ctrl+[" + BlueColorName + "::ub]A[white::-] Copy All (log)"

	ShortcutHSubMenu        = ShortcutH + " >> Ctrl+[" + BlueColorName + "::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"
	SettingsShortcutSubMenu = SettingsShortcut + " >> Ctrl+[" + BlueColorName + "::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"
	ShortcutSRSubMenu       = " Save Request >> Ctrl+[" + BlueColorName + "::ub]Down[white::-] Left Menu >> Select Letter (or press down/up)"

	ShortcutPressEscape = "Press Escape"
	ShortcutSeparator   = " | "
)

// Represents data shortcuts to display to the user
var (
	MainShortcutsText       = strings.Join([]string{SelectAPIShortcut, ShortcutF, ShortcutH, ShortcutD, SettingsShortcut, ShortcutQ}, ShortcutSeparator)
	ResultShortcutsText     = strings.Join([]string{ShortcutR, ShortcutDC, ShortcutDA, ShortcutPressEscape}, ShortcutSeparator)
	ExpertModeShortcutsText = strings.Join([]string{ShortcutHSubMenu, ShortcutPressEscape}, ShortcutSeparator)
	SettingsShortcutsText   = strings.Join([]string{SettingsShortcutSubMenu, ShortcutPressEscape}, ShortcutSeparator)
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
