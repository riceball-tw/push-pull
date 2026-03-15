package main

import (
	"github.com/charmbracelet/lipgloss"
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

type tile struct {
	kind       tileKind
	targetGrid [][]tile
	targetX    int
	targetY    int
	sound      string
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
