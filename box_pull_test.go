package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestBoxMoveIntoWithPull(t *testing.T) {
	// Setup a boxTile
	box := b()
	box.count = 5
	m := &model{} // model is not used much in the initial check but required by signature

	// nx, ny, dx, dy are not used in the initial check for pull
	res := box.MoveInto(m, 0, 0, 1, 0, true)

	if res.CanMove {
		t.Error("Expected CanMove to be false when pull is true in boxTile.MoveInto")
	}
}

func TestModelPullLogic(t *testing.T) {
	// Setup a grid: [Empty, Player, Box]
	// Player is at (1, 0). Box is at (0, 0).
	// We want to move right (to (2, 0)) and pull the box.
	box := b()
	box.count = 5
	grid := [][]Tile{
		{box, e, e},
	}
	m := model{
		x:    1,
		y:    0,
		grid: grid,
	}

	// Simulate shift+l (Right, Pull=true)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'L'}}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(model)

	// Player should have moved to (2, 0)
	if updatedModel.x != 2 || updatedModel.y != 0 {
		t.Errorf("Player should be at (2, 0), got (%d, %d)", updatedModel.x, updatedModel.y)
	}

	// Box should have moved from (0, 0) to (1, 0) (player's previous position)
	if updatedModel.grid[0][1].Kind() != boxKind {
		t.Errorf("Box should have moved to (1, 0), got %v", updatedModel.grid[0][1].Kind())
	}
	if b, ok := updatedModel.grid[0][1].(*boxTile); ok {
		if b.count != 4 {
			t.Errorf("Box count should be 4, got %d", b.count)
		}
	} else {
		t.Error("Tile at (1, 0) is not a boxTile")
	}

	// Old box position (0, 0) should be empty
	if updatedModel.grid[0][0].Kind() != emptyKind {
		t.Errorf("Old box position (0, 0) should be empty, got %v", updatedModel.grid[0][0].Kind())
	}
}
