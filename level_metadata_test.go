package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLevelMetadataInitialization(t *testing.T) {
	l := Level{
		Title:       "Test Level",
		Description: "A test description",
		StartX:      1,
		StartY:      1,
		Grid: [][]Tile{
			{e, e, e},
			{e, e, e},
			{e, e, e},
		},
	}

	m := NewModel(l)

	if m.title != l.Title {
		t.Errorf("expected title to be %q, got %q", l.Title, m.title)
	}
	if m.description != l.Description {
		t.Errorf("expected description to be %q, got %q", l.Description, m.description)
	}
}

func TestLevelMetadataView(t *testing.T) {
	m := model{
		title:       "Gopher World",
		description: "Where gophers play",
		x:           0,
		y:           0,
		grid: [][]Tile{
			{e},
		},
	}

	view := m.View()

	if !strings.Contains(view, m.title) {
		t.Errorf("view does not contain level title %q", m.title)
	}
	if !strings.Contains(view, m.description) {
		t.Errorf("view does not contain level description %q", m.description)
	}
}

func TestLevelMetadataTransition(t *testing.T) {
	level2 := Level{
		Title:       "Level Two",
		Description: "The second level",
		Grid: [][]Tile{
			{e},
		},
		StartX: 0,
		StartY: 0,
	}

	level1 := Level{
		Title:       "Level One",
		Description: "The first level",
		Grid: [][]Tile{
			{e, newDoor(level2)},
		},
		StartX: 0,
		StartY: 0,
	}

	m := NewModel(level1)

	if m.title != "Level One" {
		t.Fatalf("expected initial title 'Level One', got %q", m.title)
	}

	// Move right onto the door
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	newModel, _ := m.Update(msg)
	res := newModel.(model)

	if res.title != "Level Two" {
		t.Errorf("expected title to update to 'Level Two', got %q", res.title)
	}
	if res.description != "The second level" {
		t.Errorf("expected description to update to 'The second level', got %q", res.description)
	}
}

func TestLevelMetadataUndo(t *testing.T) {
	level2 := Level{
		Title:       "Level Two",
		Description: "The second level",
		Grid: [][]Tile{
			{e},
		},
		StartX: 0,
		StartY: 0,
	}

	level1 := Level{
		Title:       "Level One",
		Description: "The first level",
		Grid: [][]Tile{
			{e, newDoor(level2)},
		},
		StartX: 0,
		StartY: 0,
	}

	m := NewModel(level1)

	// Move onto door
	msgL := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	m1, _ := m.Update(msgL)
	res1 := m1.(model)

	if res1.title != "Level Two" {
		t.Fatalf("failed to transition to Level Two")
	}

	// Undo - should no longer go back to Level One
	msgU := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")}
	m2, _ := res1.Update(msgU)
	res2 := m2.(model)

	if res2.title != "Level Two" {
		t.Errorf("expected title to stay 'Level Two' after undo, got %q", res2.title)
	}
	if res2.description != "The second level" {
		t.Errorf("expected description to stay 'The second level' after undo, got %q", res2.description)
	}
}
