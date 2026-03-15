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

type Tile interface {
	Kind() tileKind
	MoveInto(m *model, nx, ny, dx, dy int) MoveResult
}

type baseTile struct {
	kind  tileKind
	sound string
}

func (t baseTile) Kind() tileKind {
	return t.kind
}

type emptyTile struct {
	baseTile
}

func (t emptyTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.x, m.y = nx, ny
	return MoveResult{CanMove: true, Sound: t.sound}
}

type wallTile struct {
	baseTile
}

func (t wallTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	return MoveResult{CanMove: false}
}

type waterTile struct {
	baseTile
}

func (t waterTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.x, m.y = nx, ny
	return MoveResult{CanMove: true, Sound: t.sound}
}

type doorTile struct {
	baseTile
	targetGrid [][]Tile
	targetX    int
	targetY    int
}

func (t doorTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.grid = t.targetGrid
	m.x = t.targetX
	m.y = t.targetY
	return MoveResult{CanMove: true, Sound: t.sound}
}

type boxTile struct {
	baseTile
}

func (t boxTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	height := len(m.grid)
	width := len(m.grid[0])
	nnx, nny := nx+dx, ny+dy
	if nnx >= 0 && nnx < width && nny >= 0 && nny < height {
		behindBoxTile := m.grid[nny][nnx]
		if behindBoxTile.Kind() == emptyKind {
			// Push the box
			m.grid[nny][nnx] = m.grid[ny][nx]
			m.grid[ny][nx] = empty
			m.x, m.y = nx, ny
			return MoveResult{CanMove: true, Sound: t.sound}
		}
	}
	return MoveResult{CanMove: false}
}

var (
	empty = emptyTile{baseTile{kind: emptyKind, sound: "walk"}}
	wall  = wallTile{baseTile{kind: wallKind}}
	water = waterTile{baseTile{kind: waterKind, sound: "splash"}}
	box   = boxTile{baseTile{kind: boxKind, sound: "push"}}
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
