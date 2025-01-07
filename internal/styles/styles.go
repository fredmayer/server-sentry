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

// StatusBar - Красивый ответ
func StatusBar(prefix string, value string, suffix string) string {
	// Status Bar.

	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 1).
		MarginRight(1)

	statusNugget := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)
	suffixStyle := statusNugget.
		Background(lipgloss.Color("#A550DF")).
		Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusStyle.Render(prefix),
		statusBarStyle.Render(value),
		suffixStyle.Render(suffix),
	)
}
