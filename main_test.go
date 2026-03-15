package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_Character_Movement(t *testing.T) {
	grid := [][]Tile{
		{empty, empty, empty},
		{empty, wall,  empty},
		{empty, empty, empty},
	}

	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	tests := []struct {
		name     string
		msg      tea.Msg
		wantX    int
		wantY    int
		wantQuit bool
	}{
		{
			name:  "move right",
			msg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")},
			wantX: 1,
			wantY: 0,
		},
		{
			name:  "move down",
			msg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			wantX: 1,
			wantY: 0, // (1,1) is wall
		},
		{
			name:  "move right again",
			msg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")},
			wantX: 2,
			wantY: 0,
		},
		{
			name:  "move down at col 2",
			msg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			wantX: 2,
			wantY: 1,
		},
	}

	// Sequential test
	currModel := tea.Model(m)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd tea.Cmd
			currModel, cmd = currModel.Update(tt.msg)
			m = currModel.(model)
			if m.x != tt.wantX || m.y != tt.wantY {
				t.Errorf("%s: got (x:%d, y:%d), want (x:%d, y:%d)", tt.name, m.x, m.y, tt.wantX, tt.wantY)
			}
			if tt.wantQuit && (cmd == nil || cmd() != tea.Quit()) {
				t.Errorf("%s: expected tea.Quit", tt.name)
			}
		})
	}
}

func Test_Individual_Character_Movement(t *testing.T) {
	grid := [][]Tile{
		{empty, empty, empty},
		{empty, wall,  empty},
		{empty, empty, empty},
	}

	tests := []struct {
		name  string
		startX int
		startY int
		key   string
		wantX int
		wantY int
	}{
		{"move right", 0, 0, "l", 1, 0},
		{"move left", 1, 0, "h", 0, 0},
		{"move down", 0, 0, "j", 0, 1},
		{"move up", 0, 1, "k", 0, 0},
		{"boundary left", 0, 0, "h", 0, 0},
		{"boundary right", 2, 0, "l", 2, 0},
		{"boundary up", 0, 0, "k", 0, 0},
		{"boundary down", 0, 2, "j", 0, 2},
		{"wall right", 0, 1, "l", 0, 1},
		{"wall left", 2, 1, "h", 2, 1},
		{"wall down", 1, 0, "j", 1, 0},
		{"wall up", 1, 2, "k", 1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{x: tt.startX, y: tt.startY, grid: grid}
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if len(tt.key) > 1 { // handle special keys if needed
				// but here we only use single runes
			}
			newModel, _ := m.Update(msg)
			res := newModel.(model)
			if res.x != tt.wantX || res.y != tt.wantY {
				t.Errorf("%s: moved from (%d,%d) to (%d,%d) with %s, want (%d,%d)", 
					tt.name, tt.startX, tt.startY, res.x, res.y, tt.key, tt.wantX, tt.wantY)
			}
		})
	}
}

func Test_Quit(t *testing.T) {
	m := model{x: 0, y: 0, grid: [][]Tile{{empty}}}
	
	// Test 'q'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	_, cmd := m.Update(msg)
	if cmd == nil || cmd() != tea.Quit() {
		t.Error("expected tea.Quit for 'q'")
	}

	// Test 'ctrl+c'
	msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd = m.Update(msg)
	if cmd == nil || cmd() != tea.Quit() {
		t.Error("expected tea.Quit for 'ctrl+c'")
	}
}

func Test_Door_Teleport(t *testing.T) {
	grid2 := [][]Tile{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}
	grid1 := [][]Tile{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}

	// Place a door in grid1 at (1,1) leading to grid2 at (2,2)
	grid1[1][1] = doorTile{
		baseTile:   baseTile{kind: doorKind},
		targetGrid: grid2,
		targetX:    2,
		targetY:    2,
	}

	m := model{
		x:    0,
		y:    1,
		grid: grid1,
	}

	// Move right from (0,1) to (1,1), which is a door
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	newModel, _ := m.Update(msg)
	res := newModel.(model)

	if res.x != 2 || res.y != 2 {
		t.Errorf("expected teleport to (2,2), got (%d,%d)", res.x, res.y)
	}

	// Verify grid changed (grid2 should have different pointer/identity if we want to be strict,
	// but here we check if a tile change in grid2 is reflected in res.grid)
	grid2[0][0] = wall
	if res.grid[0][0].Kind() != wallKind {
		t.Errorf("expected grid to be grid2")
	}
}

func Test_Sound_Triggers(t *testing.T) {
	grid := [][]Tile{
		{empty, water},
		{empty, empty},
	}
	// Note: water has sound: "splash"

	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	// Move right onto water
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	newModel, cmd := m.Update(msg)
	res := newModel.(model)

	if res.x != 1 || res.y != 0 {
		t.Errorf("expected movement to (1,0), got (%d,%d)", res.x, res.y)
	}

	if cmd == nil {
		t.Fatal("expected a command for sound")
	}

	// In tests, speaker might not be initialized, but we can still check the message
	// returned by the command.
	soundM := cmd()
	if sm, ok := soundM.(soundMsg); ok {
		if string(sm) != "splash" {
			t.Errorf("expected sound 'splash', got '%s'", string(sm))
		}
	} else {
		t.Errorf("expected command to return soundMsg, got %T", soundM)
	}

	// Verify that soundMsg updates the model
	m2, cmd2 := res.Update(soundM)
	res2 := m2.(model)
	if res2.sound != "splash" {
		t.Errorf("expected model.sound to be 'splash', got '%s'", res2.sound)
	}
	if cmd2 != nil {
		t.Error("expected nil command after soundMsg")
	}

	// Verify that another key press clears the sound
	m3, _ := res2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	res3 := m3.(model)
	if res3.sound != "" {
		t.Error("expected model.sound to be cleared after move")
	}
}

func Test_Door_Sound(t *testing.T) {
	grid2 := [][]Tile{{empty}}
	grid1 := [][]Tile{
		{empty, doorTile{baseTile: baseTile{kind: doorKind, sound: "creak"}, targetGrid: grid2, targetX: 0, targetY: 0}},
	}

	m := model{
		x:    0,
		y:    0,
		grid: grid1,
	}

	// Move right onto door
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	newModel, cmd := m.Update(msg)
	res := newModel.(model)

	if res.x != 0 || res.y != 0 {
		t.Errorf("expected teleport to (0,0), got (%d,%d)", res.x, res.y)
	}
	if len(res.grid) != len(grid2) || len(res.grid[0]) != len(grid2[0]) {
		t.Error("expected grid to be grid2")
	}

	if cmd == nil {
		t.Fatal("expected a command for door sound")
	}

	soundM := cmd()
	if sm, ok := soundM.(soundMsg); ok {
		if string(sm) != "creak" {
			t.Errorf("expected sound 'creak', got '%s'", string(sm))
		}
	} else {
		t.Errorf("expected command to return soundMsg, got %T", soundM)
	}
}


func Test_Push_Box(t *testing.T) {
	tests := []struct {
		name      string
		grid      [][]Tile
		startX    int
		startY    int
		key       string
		wantX     int
		wantY     int
		boxStartX int
		boxStartY int
		boxEndX   int
		boxEndY   int
		canMove   bool
	}{
		{
			name: "push box right into empty space",
			grid: [][]Tile{
				{empty, box, empty, empty},
			},
			startX:    0,
			startY:    0,
			key:       "l",
			wantX:     1,
			wantY:     0,
			boxStartX: 1,
			boxStartY: 0,
			boxEndX:   2,
			boxEndY:   0,
			canMove:   true,
		},
		{
			name: "push box right into wall",
			grid: [][]Tile{
				{empty, box, wall},
			},
			startX:    0,
			startY:    0,
			key:       "l",
			wantX:     0,
			wantY:     0,
			boxStartX: 1,
			boxStartY: 0,
			boxEndX:   1, // should not move
			boxEndY:   0,
			canMove:   false,
		},
		{
			name: "push box right into another box",
			grid: [][]Tile{
				{empty, box, box, empty},
			},
			startX:    0,
			startY:    0,
			key:       "l",
			wantX:     0,
			wantY:     0,
			boxStartX: 1,
			boxStartY: 0,
			boxEndX:   1, // should not move
			boxEndY:   0,
			canMove:   false,
		},
		{
			name: "push box right into boundary",
			grid: [][]Tile{
				{empty, box},
			},
			startX:    0,
			startY:    0,
			key:       "l",
			wantX:     0,
			wantY:     0,
			boxStartX: 1,
			boxStartY: 0,
			boxEndX:   1, // should not move
			boxEndY:   0,
			canMove:   false,
		},
		{
			name: "push box down into empty space",
			grid: [][]Tile{
				{empty},
				{box},
				{empty},
			},
			startX:    0,
			startY:    0,
			key:       "j",
			wantX:     0,
			wantY:     1,
			boxStartX: 0,
			boxStartY: 1,
			boxEndX:   0,
			boxEndY:   2,
			canMove:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Deep copy grid to avoid side effects between tests
			gridCopy := make([][]Tile, len(tt.grid))
			for i := range tt.grid {
				gridCopy[i] = make([]Tile, len(tt.grid[i]))
				copy(gridCopy[i], tt.grid[i])
			}

			m := model{x: tt.startX, y: tt.startY, grid: gridCopy}
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			newModel, _ := m.Update(msg)
			res := newModel.(model)

			if res.x != tt.wantX || res.y != tt.wantY {
				t.Errorf("%s: player moved to (%d,%d), want (%d,%d)", tt.name, res.x, res.y, tt.wantX, tt.wantY)
			}

			// Check box final position
			if tt.canMove {
				if res.grid[tt.boxStartY][tt.boxStartX].Kind() == boxKind {
					t.Errorf("%s: box still at original position (%d,%d)", tt.name, tt.boxStartX, tt.boxStartY)
				}
				if res.grid[tt.boxEndY][tt.boxEndX].Kind() != boxKind {
					t.Errorf("%s: box NOT at target position (%d,%d)", tt.name, tt.boxEndX, tt.boxEndY)
				}
			} else {
				if res.grid[tt.boxStartY][tt.boxStartX].Kind() != boxKind {
					t.Errorf("%s: box moved from (%d,%d) but it shouldn't have", tt.name, tt.boxStartX, tt.boxStartY)
				}
			}
		})
	}
}
