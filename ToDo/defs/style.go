package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Default_Style lipgloss.Style
var Header_Style lipgloss.Style
var Remark_Style lipgloss.Style
var Print_Width int = 30

func init() {
	Header_Style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center).
		Width(Print_Width)

	Default_Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#8662f5")).
		Width(Print_Width)

	Remark_Style = lipgloss.NewStyle().
		Width(Print_Width).
		Italic(true).
		Foreground(lipgloss.Color("#c1becc")).
		Background(lipgloss.Color("#a68df2"))
}