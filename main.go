package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#fff")).
Background(lipgloss.Color("#7D54F2"))

var backgroundStyle = lipgloss.NewStyle().
Foreground(lipgloss.Color("#fff")).
Background(lipgloss.Color("#333")) 


type tile int

const (
	empty tile = iota
	wall
)

type model struct {
	x, y int
	grid [][]tile
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var s string
	height := len(m.grid)
	if height == 0 {
		return "Empty world"
	}
	width := len(m.grid[0])

	for y := 0; y < height; y++ {
		line := ""
		for x := 0; x < width; x++ {
			if x == m.x && y == m.y {
				line += baseStyle.Render("我")
			} else if m.grid[y][x] == wall {
				line += "牆"
			} else {
				line += backgroundStyle.Render("　")
			}
		}
		s += line + "\n"
	}

	s += "\n(use h, j, k, l to move, X is a wall, ctrl+c to quit)"
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	height := len(m.grid)
	if height == 0 {
		return m, nil
	}
	width := len(m.grid[0])

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h": // left
			if m.x > 0 && m.grid[m.y][m.x-1] != wall {
				m.x--
			}
		case "j": // down
			if m.y < height-1 && m.grid[m.y+1][m.x] != wall {
				m.y++
			}
		case "k": // up
			if m.y > 0 && m.grid[m.y-1][m.x] != wall {
				m.y--
			}
		case "l": // right
			if m.x < width-1 && m.grid[m.y][m.x+1] != wall {
				m.x++
			}
		}
	}

	return m, nil
}

func main() {
	grid := [][]tile{
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, wall,  wall,  wall,  empty, empty, empty, empty, empty},
		{empty, empty, wall,  empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, wall,  empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
	}

	p := tea.NewProgram(model{x: 0, y: 0, grid: grid})

	if _, pErr := p.Run(); pErr != nil {
		log.Panic(pErr)
	}
}