package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type model struct {
	x, y    int
	grid    [][]Tile
	sound   string
	moves   int
	history []state
}

type state struct {
	x, y int
	grid [][]Tile
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
		x:    m.x,
		y:    m.y,
		grid: newGrid,
	}
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

	s += "\n(use h, j, k, l to move, u to undo, q to quit)"
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
		}

		if dx != 0 || dy != 0 {
			nx, ny := m.x+dx, m.y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				currentState := m.snapshot()
				targetTile := m.grid[ny][nx]
				res := targetTile.MoveInto(&m, nx, ny, dx, dy)
				if res.CanMove {
					m.history = append(m.history, currentState)
					m.moves++
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
		x:    l.StartX,
		y:    l.StartY,
		grid: l.Grid,
	}
}

func main() {
	p := tea.NewProgram(NewModel(level2))

	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	if _, pErr := p.Run(); pErr != nil {
		log.Panic(pErr)
	}
}