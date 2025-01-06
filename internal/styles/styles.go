package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func H2(str string) {
	style := lipgloss.NewStyle().Blink(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		// BorderBackground(lipgloss.Color("63")).
		BorderLeft(true).
		BorderBottom(true).
		BorderTop(true).
		BorderRight(true).
		PaddingLeft(1).
		PaddingRight(1).MarginTop(2)

	fmt.Println(style.Render(str))
}

func ReturnWithX(str string) string {
	cross := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6347")).
		PaddingRight(1)
	text := lipgloss.NewStyle().Bold(true)

	return fmt.Sprintf("%s %s", cross.Render("×"), text.Render(str))
}

func ReturnWithOk(str string) string {
	cross := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#07b804")).
		PaddingRight(1)
	text := lipgloss.NewStyle().Bold(true)

	return fmt.Sprintf("\n%s %s", cross.Render("✓"), text.Render(str))
}
