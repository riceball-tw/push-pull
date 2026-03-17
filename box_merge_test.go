package main

import (
	"testing"
)

func TestBoxMerging(t *testing.T) {
	// Setup:
	// Player at (0, 0)
	// Box1 at (1, 0)
	// Box2 at (2, 0)
	// Expected: Player moves to (1, 0), Box1 merges into Box2 at (2, 0).
	// Resulting Box at (2, 0) should have count = Box1.count + Box2.count + 1 (since pushing adds 1).
    // Let's check how count currently works.

	box1 := b()
	box2 := b()
	box1.count = 2
	box2.count = 3
	
	grid := [][]Tile{
		{e, box1, box2, e},
	}
	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	// Move right into box1
	res := m.grid[0][1].MoveInto(&m, 1, 0, 1, 0, false)

	if !res.CanMove {
		t.Fatal("Should be able to push box1 into box2")
	}

	if m.x != 1 || m.y != 0 {
		t.Errorf("Player should be at (1, 0), got (%d, %d)", m.x, m.y)
	}

	if m.grid[0][1].Kind() != emptyKind {
		t.Errorf("Tile at (1, 0) should be empty after push, got %v", m.grid[0][1].Kind())
	}

	mergedBox, ok := m.grid[0][2].(*boxTile)
	if !ok {
		t.Fatalf("Tile at (2, 0) should be a box, got %T", m.grid[0][2])
	}

	// Current logic adds 1 to count when pushing.
	// If they merge:
	// box1 has 2.
	// box2 has 3.
	// Pushing box1: it gets +1 (now 3).
	// Then it merges with box2: 3 + 3 = 6.
	// Or maybe it should be (box1.count + 1) + box2.count.
	expectedCount := (2 + 1) + 3
	if mergedBox.count != expectedCount {
		t.Errorf("Merged box count should be %d, got %d", expectedCount, mergedBox.count)
	}
}
