package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type tileKind int

const (
	emptyKind tileKind = iota
	wallKind
	doorKind
	waterKind
	boxKind
	lockKind
)

type MoveResult struct {
	CanMove bool
	Sound   string
}

type Tile interface {
	Kind() tileKind
	MoveInto(m *model, nx, ny, dx, dy int) MoveResult
	Count() int
	Clone() Tile
}

type baseTile struct {
	kind  tileKind
	sound string
}

func (t baseTile) Kind() tileKind {
	return t.kind
}

func (t baseTile) Count() int {
	return 0
}

type emptyTile struct {
	baseTile
}

func (t emptyTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.x, m.y = nx, ny
	return MoveResult{CanMove: true, Sound: t.sound}
}

func (t emptyTile) Clone() Tile {
	return t
}

type wallTile struct {
	baseTile
}

func (t wallTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	return MoveResult{CanMove: false}
}

func (t wallTile) Clone() Tile {
	return t
}

type waterTile struct {
	baseTile
}

func (t waterTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.x, m.y = nx, ny
	return MoveResult{CanMove: true, Sound: t.sound}
}

func (t waterTile) Clone() Tile {
	return t
}

type doorTile struct {
	baseTile
	targetLevel Level
}

func (t doorTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	m.grid = t.targetLevel.Grid
	m.x = t.targetLevel.StartX
	m.y = t.targetLevel.StartY
	m.title = t.targetLevel.Title
	m.description = t.targetLevel.Description
	return MoveResult{CanMove: true, Sound: t.sound}
}

func (t doorTile) Clone() Tile {
	return t
}

type boxTile struct {
	baseTile
	count int
}

func (t *boxTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	height := len(m.grid)
	width := len(m.grid[0])
	nnx, nny := nx+dx, ny+dy
	if nnx >= 0 && nnx < width && nny >= 0 && nny < height {
		behindBoxTile := m.grid[nny][nnx]
		if behindBoxTile.Kind() == emptyKind {
			// Push the box
			t.count++
			m.grid[nny][nnx] = m.grid[ny][nx]
			m.grid[ny][nx] = empty
			m.x, m.y = nx, ny
			return MoveResult{CanMove: true, Sound: t.sound}
		} else if behindBoxTile.Kind() == boxKind {
			// Merge the box
			t.count++
			if targetBox, ok := behindBoxTile.(*boxTile); ok {
				targetBox.count += t.count
				m.grid[ny][nx] = empty
				m.x, m.y = nx, ny
				return MoveResult{CanMove: true, Sound: t.sound}
			}
		} else if behindBoxTile.Kind() == lockKind {
			// Push into lock
			if targetLock, ok := behindBoxTile.(*lockTile); ok {
				if t.count == targetLock.requiredCount {
					m.grid[nny][nnx] = targetLock.targetDoor
					m.grid[ny][nx] = empty
					m.x, m.y = nx, ny
					return MoveResult{CanMove: true, Sound: "unlock"}
				}
			}
		}
	}
	return MoveResult{CanMove: false}
}

func (t *boxTile) Clone() Tile {
	return &boxTile{
		baseTile: t.baseTile,
		count:    t.count,
	}
}

type lockTile struct {
	baseTile
	requiredCount int
	targetDoor    Tile
}

func (t *lockTile) MoveInto(m *model, nx, ny, dx, dy int) MoveResult {
	return MoveResult{CanMove: false}
}

func (t *lockTile) Clone() Tile {
	return &lockTile{
		baseTile:      t.baseTile,
		requiredCount: t.requiredCount,
		targetDoor:    t.targetDoor.Clone(),
	}
}

func (t *lockTile) Count() int {
	return t.requiredCount
}

func (t *lockTile) DisplayChar() string {
	if t.requiredCount <= 9 {
		fullWidthNums := []string{"０", "１", "２", "３", "４", "５", "６", "７", "８", "９"}
		return fullWidthNums[t.requiredCount]
	}
	if t.requiredCount <= 99 {
		return fmt.Sprintf("%02d", t.requiredCount)
	}
	return "鎖"
}

var (
	empty = emptyTile{baseTile{kind: emptyKind, sound: "walk"}}
	wall  = wallTile{baseTile{kind: wallKind}}
	water = waterTile{baseTile{kind: waterKind, sound: "splash"}}
)

func newBox() *boxTile {
	return &boxTile{baseTile: baseTile{kind: boxKind, sound: "push"}}
}

func newLock(requiredCount int, targetDoor Tile) *lockTile {
	return &lockTile{
		baseTile:      baseTile{kind: lockKind},
		requiredCount: requiredCount,
		targetDoor:    targetDoor,
	}
}

func newDoor(targetLevel Level) *doorTile {
	return &doorTile{
		baseTile:    baseTile{kind: doorKind, sound: "door"},
		targetLevel: targetLevel,
	}
}

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
	lockKind: {
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fff")).
			Background(lipgloss.Color("#000000")),
		char: "鎖",
	},
}

func (t *boxTile) DisplayChar() string {
	if t.count == 0 {
		return tiles[boxKind].char
	}
	if t.count <= 9 {
		fullWidthNums := []string{"０", "１", "２", "３", "４", "５", "６", "７", "８", "９"}
		return fullWidthNums[t.count]
	}
	if t.count <= 99 {
		return fmt.Sprintf("%02d", t.count)
	}
	return "＋"
}

func (t *boxTile) Count() int {
	return t.count
}

var playerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fff")).
	Background(lipgloss.Color("#7D54F2"))
