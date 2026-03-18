package game

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_History_Cleared_On_Level_Change(t *testing.T) {
	grid2 := [][]Tile{
		{e, e, e},
		{e, e, e},
		{e, e, e},
	}
	grid1 := [][]Tile{
		{e, doorTile{baseTile: baseTile{kind: doorKind}, targetLevel: Level{Title: "Level 2", Grid: grid2, StartX: 1, StartY: 1}}},
	}

	m := model{
		x:    0,
		y:    0,
		grid: grid1,
	}

	// Move right onto door
	msgL := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	newModel, _ := m.Update(msgL)
	res := newModel.(model)

	// Verify we moved to level 2
	if res.title != "Level 2" {
		t.Errorf("expected level title 'Level 2', got '%s'", res.title)
	}
	if res.x != 1 || res.y != 1 {
		t.Errorf("expected teleport to (1,1), got (%d,%d)", res.x, res.y)
	}

	// Before the fix, history will have 1 state (from level 1)
	// After the fix, history should be empty
	if len(res.history) != 0 {
		t.Errorf("expected history to be cleared, but it has %d states", len(res.history))
	}
	if res.moves != 0 {
		t.Errorf("expected moves to be reset, but it is %d", res.moves)
	}

	// Try to undo - it should not go back to level 1
	msgU := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")}
	mUndo, _ := res.Update(msgU)
	resUndo := mUndo.(model)

	if resUndo.title != "Level 2" {
		t.Errorf("undo took player back to previous level! title is '%s'", resUndo.title)
	}
}
