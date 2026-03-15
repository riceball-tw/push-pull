package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#fff")).
Background(lipgloss.Color("#7D54F2"))


type model struct {
	x int
	y int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	// Simple rendering of the character "@" at (x, y)
	// We'll use newlines for Y and spaces for X
	var s string
	for i := 0; i < m.y; i++ {
		s += "\n"
	}
	spaces := ""
	for i := 0; i < m.x; i++ {
		spaces += " "
	}
	s += spaces
	s += baseStyle.Render("@")
	s += "\n\n(use h, j, k, l to move, ctrl+c to quit)"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h": // left
			if m.x > 0 {
				m.x--
			}
		case "j": // down
			m.y++
		case "k": // up
			if m.y > 0 {
				m.y--
			}
		case "l": // right
			m.x++
		}
	}

	return m, nil
}

func main() {
	p := tea.NewProgram(model{x: 0, y: 0})

	if _, pErr := p.Run(); pErr != nil {
		log.Panic(pErr)
	}
}