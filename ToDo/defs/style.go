package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Default_Style lipgloss.Style
var Header_Style lipgloss.Style
var Remark_Style lipgloss.Style
var Defocused_Default_Style lipgloss.Style
var Defocused_Remark_Style lipgloss.Style
var Amber_Style lipgloss.Style
var Light_Amber_Style lipgloss.Style
var Green_Style lipgloss.Style
var Light_Green_Style lipgloss.Style
var Print_Width int = 30
var Complete_Column_Width int = 3

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

	Defocused_Default_Style = Default_Style.Copy().Background(lipgloss.Color("#616362"))

	Remark_Style = lipgloss.NewStyle().
		Width(Print_Width).
		Italic(true).
		Foreground(lipgloss.Color("#c1becc")).
		Background(lipgloss.Color("#a68df2"))

	Defocused_Remark_Style = Remark_Style.Copy().Background(lipgloss.Color("#747575"))

	Amber_Style = lipgloss.NewStyle().
		Width(Complete_Column_Width).
		Background(lipgloss.Color("#e69525")).
		Foreground(lipgloss.Color("#ffffff"))

	Light_Amber_Style = Amber_Style.Copy().Background(lipgloss.Color("#cfa66d"))

	Green_Style = lipgloss.NewStyle().
		Width(Complete_Column_Width).
		Background(lipgloss.Color("#05961a")).
		Foreground(lipgloss.Color("#ffffff"))

	Light_Green_Style = Green_Style.Copy().Background(lipgloss.Color("#5cbf69"))
}

func GetTaskStyles(indentation int, taskComplete bool) (taskStyle lipgloss.Style, taskStatusStyle lipgloss.Style, descStyle lipgloss.Style, descStatusStyle lipgloss.Style) {
	if taskComplete {
		taskStyle = Defocused_Default_Style.Copy()
		descStyle = Defocused_Remark_Style.Copy()
	} else {
		taskStyle = Default_Style.Copy()
		descStyle = Remark_Style.Copy()
	}
	taskStyle.MarginLeft(2 * indentation).Width(taskStyle.GetWidth() - 2*indentation)
	descStyle.MarginLeft(2 * indentation).Width(descStyle.GetWidth() - 2*indentation)

	if taskComplete {
		taskStatusStyle = Green_Style
		descStatusStyle = Light_Green_Style
	} else {
		taskStatusStyle = Amber_Style
		descStatusStyle = Light_Amber_Style
	}

	return taskStyle, taskStatusStyle, descStyle, descStatusStyle
}
