package main

import (
	"testing"
)

func TestLockTile(t *testing.T) {
	// Setup:
	// Player at (0, 0)
	// Box at (1, 0) with count 9
	// Lock at (2, 0) requiring 10
	// Door is the target of the lock.
	
	targetDoor := doorTile{baseTile: baseTile{kind: doorKind}}
	lock := &lockTile{
		baseTile:      baseTile{kind: lockKind},
		requiredCount: 10,
		targetDoor:    targetDoor,
	}
	box := b()
	box.count = 10 // Exact count required

	grid := [][]Tile{
		{e, box, lock, e},
	}
	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	// Move right into box
	res := m.grid[0][1].MoveInto(&m, 1, 0, 1, 0, false)

	if !res.CanMove {
		t.Fatal("Should be able to push 10-count box into 10-count lock")
	}

	if m.x != 1 || m.y != 0 {
		t.Errorf("Player should be at (1, 0), got (%d, %d)", m.x, m.y)
	}

	if m.grid[0][1].Kind() != emptyKind {
		t.Errorf("Tile at (1, 0) should be empty after push, got %v", m.grid[0][1].Kind())
	}

	if m.grid[0][2].Kind() != doorKind {
		t.Errorf("Tile at (2, 0) should be door after unlock, got %v", m.grid[0][2].Kind())
	}
}

func TestLockTileInsufficientCount(t *testing.T) {
	// Setup:
	// Player at (0, 0)
	// Box at (1, 0) with count 1
	// Lock at (2, 0) requiring 10
	
	targetDoor := doorTile{baseTile: baseTile{kind: doorKind}}
	lock := &lockTile{
		baseTile:      baseTile{kind: lockKind},
		requiredCount: 10,
		targetDoor:    targetDoor,
	}
	box := b()
	box.count = 9 // Will NOT become 10 when pushed, still < 10

	grid := [][]Tile{
		{e, box, lock, e},
	}
	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	// Move right into box
	res := m.grid[0][1].MoveInto(&m, 1, 0, 1, 0, false)

	if res.CanMove {
		t.Fatal("Should NOT be able to push box into lock with insufficient count")
	}

	if m.x != 0 || m.y != 0 {
		t.Errorf("Player should still be at (0, 0), got (%d, %d)", m.x, m.y)
	}

	if m.grid[0][1].Kind() != boxKind {
		t.Errorf("Tile at (1, 0) should still be box, got %v", m.grid[0][1].Kind())
	}

	if m.grid[0][2].Kind() != lockKind {
		t.Errorf("Tile at (2, 0) should still be lock, got %v", m.grid[0][2].Kind())
	}
}
