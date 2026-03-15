package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type tileKind int

const (
	emptyKind tileKind = iota
	wallKind
	doorKind
	waterKind
	boxKind
)

type MoveResult struct {
	CanMove bool
	Sound   string
}

func (t tile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	height := len(m.grid)
	width := len(m.grid[0])

	switch t.kind {
	case emptyKind, waterKind:
		m.x, m.y = nx, ny
		return MoveResult{CanMove: true, Sound: t.sound}
	case doorKind:
		m.grid = t.targetGrid
		m.x = t.targetX
		m.y = t.targetY
		return MoveResult{CanMove: true, Sound: t.sound}
	case boxKind:
		nnx, nny := nx+dx, ny+dy
		if nnx >= 0 && nnx < width && nny >= 0 && nny < height {
			behindBoxTile := m.grid[nny][nnx]
			if behindBoxTile.kind == emptyKind {
				// Push the box
				m.grid[nny][nnx] = m.grid[ny][nx]
				m.grid[ny][nx] = empty
				m.x, m.y = nx, ny
				return MoveResult{CanMove: true, Sound: t.sound}
			}
		}
		return MoveResult{CanMove: false}
	case wallKind:
		return MoveResult{CanMove: false}
	default:
		return MoveResult{CanMove: false}
	}
}

type tile struct {
	kind       tileKind
	targetGrid [][]tile
	targetX    int
	targetY    int
	sound      string
}

var (
	empty = tile{kind: emptyKind, sound: "walk"}
	wall  = tile{kind: wallKind}
	water = tile{kind: waterKind, sound: "splash"}
	box   = tile{kind: boxKind, sound: "push"}
)

type tileInfo struct {
	style lipgloss.Style
	char  string
}

var tiles = map[tileKind]tileInfo{
	emptyKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#333")),
		char: "　",
	},
	wallKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#481309")),
		char: "牆",
	},
	doorKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#553311")),
		char: "門",
	},
	waterKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#0077be")),
		char: "水",
	},
	boxKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#B87333")),
		char: "箱",
	},
}

var playerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fff")).
	Background(lipgloss.Color("#7D54F2"))

type model struct {
	x, y  int
	grid  [][]tile
	sound string
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
				info := tiles[m.grid[y][x].kind]
				line += info.style.Render(info.char)
			}
		}
		s += line + "\n"
	}

	if m.sound != "" {
		s += "\n* Playing sound: " + m.sound + " *"
	}

	s += "\n(use h, j, k, l to move, q to quit)"
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
				targetTile := m.grid[ny][nx]
				res := targetTile.MoveInto(&m, nx, ny, dx, dy)
				if res.CanMove && res.Sound != "" {
					cmd = playSound(res.Sound)
				}
			}
		}
	}

	return m, cmd
}

func main() {
	grid2 := [][]tile{
		{wall, wall, wall, wall, wall},
		{wall, empty, empty, empty, wall},
		{wall, empty, wall, empty, wall},
		{wall, empty, empty, empty, wall},
		{wall, wall, wall, wall, wall},
	}

	grid1 := [][]tile{
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, wall, wall, wall, empty, empty, empty, empty, empty},
		{empty, water, wall, empty, empty, empty, empty, empty, empty, empty},
		{empty, water, wall, box, empty, empty, empty, box, empty, empty},
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
		{empty, empty, empty, empty, empty, empty, empty, empty, empty, empty},
	}

	// Add a door from grid1 to grid2
	grid1[0][5] = tile{
		kind:       doorKind,
		targetGrid: grid2,
		targetX:    2,
		targetY:    2,
		sound:      "creak",
	}

	// Add a door from grid2 to grid1
	grid2[1][1] = tile{
		kind:       doorKind,
		targetGrid: grid1,
		targetX:    0,
		targetY:    0,
		sound:      "creak",
	}

	p := tea.NewProgram(model{x: 0, y: 0, grid: grid1})

	sr := beep.SampleRate(44100)
	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	if _, pErr := p.Run(); pErr != nil {
		log.Panic(pErr)
	}
}