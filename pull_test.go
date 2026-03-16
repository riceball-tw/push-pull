package main

import (
	"testing"
)

func TestPullBox(t *testing.T) {
	// Setup a grid
	box := &boxTile{baseTile: baseTile{kind: boxKind, sound: "push"}, count: 5}
	grid := [][]Tile{
		{e, e, e, e},
		{e, e, box, e},
		{e, e, e, e},
	}
	m := model{
		x:    1,
		y:    1,
		grid: grid,
	}

	// Player is at (1, 1). Box is at (2, 1).
	// We want to pull left (dx=-1).
	dx, dy := -1, 0
	nx, ny := m.x+dx, m.y+dy // (0, 1)

	// In model.Update:
	currentState := m.snapshot()
	targetTile := m.grid[ny][nx] // emptyTile at (0, 1)
	res := targetTile.MoveInto(&m, nx, ny, dx, dy, true) // pull=true

	if !res.CanMove {
		t.Fatal("Should be able to move into empty tile")
	}

	// m.x, m.y are updated by MoveInto
	if m.x != 0 || m.y != 1 {
		t.Errorf("Player should be at (0, 1), got (%d, %d)", m.x, m.y)
	}

	// Manually execute the pull logic from main.go
	pull := true
	if pull && res.CanMove {
		bx, by := currentState.x-dx, currentState.y-dy // (1-(-1), 1-0) = (2, 1)
		if bx >= 0 && bx < 4 && by >= 0 && by < 3 {
			behindTile := m.grid[by][bx]
			if behindTile.Kind() == boxKind {
				if boxObj, ok := behindTile.(*boxTile); ok {
					m.grid[currentState.y][currentState.x] = boxObj
					m.grid[by][bx] = e
					boxObj.count--
				}
			}
		}
	}

	// Check if box moved to (1, 1)
	if m.grid[1][1].Kind() != boxKind {
		t.Errorf("Box should have moved to (1, 1), got %v", m.grid[1][1].Kind())
	}
	if box.count != 4 {
		t.Errorf("Box count should be 4, got %d", box.count)
	}
	// Check if old box position (2, 1) is empty
	if m.grid[1][2].Kind() != emptyKind {
		t.Errorf("Old box position (2, 1) should be empty, got %v", m.grid[1][2].Kind())
	}
}

func TestPullIntoWall(t *testing.T) {
	// Setup a grid: [Wall, Player, Box]
	box := &boxTile{baseTile: baseTile{kind: boxKind}, count: 5}
	grid := [][]Tile{
		{w, e, box, e},
	}
	m := model{
		x:    1,
		y:    0,
		grid: grid,
	}

	dx, dy := -1, 0
	nx, ny := m.x+dx, m.y+dy // (0, 0)

	m.snapshot() // currentState
	targetTile := m.grid[ny][nx] // wallTile at (0, 0)
	res := targetTile.MoveInto(&m, nx, ny, dx, dy, true) // pull=true

	if res.CanMove {
		t.Fatal("Should NOT be able to move into wall")
	}

	// Player should still be at (1, 0)
	if m.x != 1 || m.y != 0 {
		t.Errorf("Player should still be at (1, 0), got (%d, %d)", m.x, m.y)
	}

	// Pull logic should not execute because res.CanMove is false.
	// But let's verify if we manually followed main.go's logic
	if res.CanMove {
		// (this won't happen)
	}

	// Box should still be at (2, 0)
	if m.grid[0][2].Kind() != boxKind {
		t.Errorf("Box should still be at (2, 0), got %v", m.grid[0][2].Kind())
	}
}
