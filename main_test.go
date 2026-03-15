package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_Character_Movement(t *testing.T) {
	grid := [][]tile{
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
	grid := [][]tile{
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
	m := model{x: 0, y: 0, grid: [][]tile{{empty}}}
	
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