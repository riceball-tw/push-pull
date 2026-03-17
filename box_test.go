package main

import (
	"strings"
	"testing"
)

func TestDisplayChar(t *testing.T) {
	tests := []struct {
		name  string
		count int
		want  string
	}{
		{"Initial box state", 0, "箱"},
		{"Count 1", 1, "１"},
		{"Count 9", 9, "９"},
		{"Count 10", 10, "10"},
		{"Count 99", 99, "99"},
		{"Count 100", 100, "＋"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box := b()
			box.count = tt.count
			got := box.DisplayChar()
			if got != tt.want {
				t.Errorf("DisplayChar() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestBoxCountDisplay(t *testing.T) {
	// Setup a grid with a box
	box := b()
	grid := [][]Tile{
		{e, box, e},
	}
	m := model{
		x:    0,
		y:    0,
		grid: grid,
	}

	// Helper to get the rendered view
	getView := func() string {
		return m.View()
	}

	// 1. Initial box state (count 0)
	// Should display "箱" (standard box character)
	view := getView()
	if !strings.Contains(view, "箱") {
		t.Errorf("Initial view should contain '箱', got: %s", view)
	}

	// 2. Push box once (count 1)
	box.count = 1
	view = getView()
	if !strings.Contains(view, "１") {
		t.Errorf("View with count 1 should contain '１', got: %s", view)
	}

	// 3. Push box 9 times (count 9)
	box.count = 9
	view = getView()
	if !strings.Contains(view, "９") {
		t.Errorf("View with count 9 should contain '９', got: %s", view)
	}

	// 4. Push box 10 times (count 10)
	box.count = 10
	view = getView()
	// Now we use regular "10" for consistent 2-column width
	if !strings.Contains(view, "10") {
		t.Errorf("View with count 10 should contain '10', got: %s", view)
	}

	// 5. Push box 100 times
	box.count = 100
	view = getView()
	// Should show "＋" to maintain 2-column width
	if !strings.Contains(view, "＋") {
		t.Errorf("View with count 100 should contain '＋', got: %s", view)
	}
}
