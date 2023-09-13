package style

import (
	"github.com/charmbracelet/lipgloss"
)

var Header_Style lipgloss.Style
var Footer_Style lipgloss.Style
var Warning_Header_Style lipgloss.Style
var Warning_Task_Style lipgloss.Style
var Warning_Desc_Style lipgloss.Style

var Default_Style lipgloss.Style
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
		Background(lipgloss.Color("#5f2485")).
		Foreground(lipgloss.Color("#ffffff")).
		Width(Print_Width + Complete_Column_Width).
		Align(lipgloss.Center)

	Footer_Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#050000")).
		Foreground(lipgloss.Color("#ffffff")).
		Width(Print_Width + Complete_Column_Width).
		Align(lipgloss.Center)

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

	Warning_Header_Style = lipgloss.NewStyle().
		Width(Print_Width).
		Background(lipgloss.Color("#e60b0b")).
		Foreground(lipgloss.Color("#ffffff")).
		Align(lipgloss.Center)

	Warning_Task_Style = lipgloss.NewStyle().
		Width(Print_Width).
		Background(lipgloss.Color("#e60b0b")).
		Foreground(lipgloss.Color("#ffffff"))

	Warning_Desc_Style = lipgloss.NewStyle().
		Width(Print_Width).
		Background(lipgloss.Color("#bf3232")).
		Foreground(lipgloss.Color("#ffffff"))

}

func GetTaskStyles(indentation int, taskComplete bool, deleting bool) (taskStyle lipgloss.Style, taskStatusStyle lipgloss.Style, descStyle lipgloss.Style, descStatusStyle lipgloss.Style) {
	if deleting {
		taskStyle = Warning_Task_Style.Copy()
		descStyle = Warning_Desc_Style.Copy()
	} else {
		if taskComplete {
			taskStyle = Defocused_Default_Style.Copy()
			descStyle = Defocused_Remark_Style.Copy()
		} else {
			taskStyle = Default_Style.Copy()
			descStyle = Remark_Style.Copy()
		}
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

func GetJustifiedString(words []string) (justifiedString string) {
	maxWidth := Print_Width + Complete_Column_Width
	lineStartIndex := 0
	lineLength := 0
	var strs []string
	i := 0
	for i < len(words) {
		if lineLength+len(words[i]) < maxWidth {
			// word and space leaves space to fill
			// or is end of line
			lineLength += len(words[i]) + 1
			i++
		} else if lineLength+len(words[i]) == maxWidth {
			// addition of word takes us to end
			// no justification required
			line := ""
			for j := lineStartIndex; j < i; j++ {
				line += words[j] + " "
			}
			strs = append(strs, (line + words[i]))
			lineStartIndex = i + 1
			lineLength = 0
			i++
		} else {
			// addition of word takes us past end
			// justify line buffer
			numWords := i - lineStartIndex
			numSpaces := numWords - 1
			if numSpaces == 0 {
				line := words[i-1]
				for len(line) < maxWidth {
					line += " "
				}
				strs = append(strs, (line))
			} else {
				spacesToAdd := (maxWidth - (lineLength - numWords))
				j := 0
				line := ""
				spaces := ""
				for j < spacesToAdd/numSpaces {
					spaces += " "
					j++
				}
				j = 0
				for j < numWords-1 {
					line += words[j+lineStartIndex] + spaces
					if j < spacesToAdd%numSpaces {
						line += " "
					}
					j++
				}
				strs = append(strs, (line + words[j+lineStartIndex]))
			}
			lineStartIndex = i
			lineLength = 0
		}
	}
	if lineStartIndex != len(words) {
		// we have words in buffer
		numWords := i - lineStartIndex
		numSpaces := numWords - 1
		if numSpaces == 0 {
			line := words[i-1]
			for len(line) < maxWidth {
				line += " "
			}
			strs = append(strs, (line))
		} else {
			spacesToAdd := (maxWidth - (lineLength - numWords))
			j := 0
			line := ""
			spaces := ""
			for j < spacesToAdd/numSpaces {
				spaces += " "
				j++
			}
			j = 0
			for j < numWords-1 {
				line += words[j+lineStartIndex] + spaces
				if j < spacesToAdd%numSpaces {
					line += " "
				}
				j++
			}
			strs = append(strs, (line + words[j+lineStartIndex]))
		}
	}
	returnString := ""
	for index := range strs {
		if index == len(strs)-1 {
			returnString += strs[index]
		} else {
			returnString += strs[index] + "\n"
		}
	}

	return returnString
}
