package ui

import "github.com/charmbracelet/lipgloss"

// Centralized styling configuration
type Styles struct {
	width int
}

func NewStyles(width int) *Styles {
	return &Styles{width: width}
}

func (s *Styles) Main() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(s.width).
		Padding(0, 1, 0, 1).
		MarginBottom(1)
}

func (s *Styles) Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#33B5E5")).
		Width(s.width)
}

func (s *Styles) Config() lipgloss.Style {
	return lipgloss.NewStyle().MarginBottom(1)
}

func (s *Styles) Stats() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (s *Styles) Error() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (s *Styles) Result() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(s.width - 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true)
}

func (s *Styles) SectionTitle(text, color string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Underline(true).
		Render(text)
}
