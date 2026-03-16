package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type model struct {
	x, y        int
	grid        [][]Tile
	sound       string
	moves       int
	history     []state
	title       string
	description string
}

type state struct {
	x, y        int
	grid        [][]Tile
	title       string
	description string
}

func (m model) snapshot() state {
	newGrid := make([][]Tile, len(m.grid))
	for y := range m.grid {
		newGrid[y] = make([]Tile, len(m.grid[y]))
		for x := range m.grid[y] {
			newGrid[y][x] = m.grid[y][x].Clone()
		}
	}
	return state{
		x:           m.x,
		y:           m.y,
		grid:        newGrid,
		title:       m.title,
		description: m.description,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var s string
	s += lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#7D54F2")).Padding(0, 1).Render(m.title) + "\n"
	s += lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#EEE")).Render(m.description) + "\n\n"

	height := len(m.grid)
	if height == 0 {
		return "Empty world"
	}
	width := len(m.grid[0])

	for y := 0; y < height; y++ {
		line := ""
		for x := 0; x < width; x++ {
			if x == m.x && y == m.y {
				line += playerStyle.Render("我")
			} else {
				tile := m.grid[y][x]
				info := tiles[tile.Kind()]
				char := info.char
				if bt, ok := tile.(*boxTile); ok {
					char = bt.DisplayChar()
				} else if lt, ok := tile.(*lockTile); ok {
					char = lt.DisplayChar()
				}
				line += info.style.Render(char)
			}
		}
		s += line + "\n"
	}

	if m.sound != "" {
		s += "\n* Playing sound: " + m.sound + " *"
	}
	s += fmt.Sprintf("\nMoves: %d", m.moves)
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	height := len(m.grid)
	if height == 0 {
		return m, nil
	}
	width := len(m.grid[0])

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case soundMsg:
		m.sound = string(msg)
		return m, nil
	case tea.KeyMsg:
		m.sound = ""
		dx, dy := 0, 0
		pull := false
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "u":
			if len(m.history) > 0 {
				prev := m.history[len(m.history)-1]
				m.history = m.history[:len(m.history)-1]
				m.x = prev.x
				m.y = prev.y
				m.grid = prev.grid
				m.title = prev.title
				m.description = prev.description
				m.moves--
			}
			return m, nil
		case "h": // left
			dx = -1
		case "j": // down
			dy = 1
		case "k": // up
			dy = -1
		case "l": // right
			dx = 1
		case "shift+h", "H":
			dx, pull = -1, true
		case "shift+j", "J":
			dy, pull = 1, true
		case "shift+k", "K":
			dy, pull = -1, true
		case "shift+l", "L":
			dx, pull = 1, true
		}

		if dx != 0 || dy != 0 {
			nx, ny := m.x+dx, m.y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				currentState := m.snapshot()
				targetTile := m.grid[ny][nx]
				res := targetTile.MoveInto(&m, nx, ny, dx, dy, pull)
				if res.CanMove {
					if targetTile.Kind() != doorKind {
						m.history = append(m.history, currentState)
						m.moves++
					}
					
					if pull {
						// Pull logic: if we moved, check if there was a box behind us
						bx, by := currentState.x-dx, currentState.y-dy
						if bx >= 0 && bx < width && by >= 0 && by < height {
							behindTile := m.grid[by][bx]
							if behindTile.Kind() == boxKind {
								if box, ok := behindTile.(*boxTile); ok {
									// Move box to player's previous position
									m.grid[currentState.y][currentState.x] = box
									m.grid[by][bx] = e
									box.count--
									if box.count < 0 {
										box.count = 0
									}
								}
							}
						}
					}

					if res.Sound != "" {
						cmd = playSound(res.Sound)
					}
				}
			}
		}
	}

	return m, cmd
}

func NewModel(l Level) model {
	return model{
		x:           l.StartX,
		y:           l.StartY,
		grid:        l.Grid,
		title:       l.Title,
		description: l.Description,
	}
}

func main() {
	p := tea.NewProgram(NewModel(level1))

	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	if _, pErr := p.Run(); pErr != nil {
		log.Panic(pErr)
	}
}