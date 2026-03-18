package main

import (
	"io"
	"os"
	"syscall/js"

	"push-pull/internal/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type xtermWriter struct {
	term js.Value
}

func (w *xtermWriter) Write(p []byte) (n int, err error) {
	w.term.Call("write", string(p))
	return len(p), nil
}

func main() {
	os.Setenv("TERM", "xterm-256color")

	// Force lipgloss to use full 24-bit true color.
	// In WASM, termenv cannot auto-detect the terminal's color profile,
	// so without this it falls back to a limited palette (or no color),
	// causing similar hex values to map to the same color.
	lipgloss.SetColorProfile(termenv.TrueColor)

	c := make(chan struct{})

	// Ensure term is globally available from index.html
	term := js.Global().Get("term")
	if term.IsUndefined() || term.IsNull() {
		js.Global().Get("console").Call("error", "xterm.js is not initialized globally as 'term'")
		return
	}

	inR, inW := io.Pipe()

	// Register onData callback in xterm.js to capture keystrokes
	onDataCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		inW.Write([]byte(args[0].String()))
		return nil
	})
	defer onDataCallback.Release()
	term.Call("onData", onDataCallback)

	outW := &xtermWriter{term: term}

	p, err := game.NewProgram(tea.WithInput(inR), tea.WithOutput(outW))
	if err != nil {
		js.Global().Get("console").Call("error", "Failed to init program: ", err.Error())
		return
	}

	sendWindowSize := func(cols, rows int) {
		if cols <= 0 || rows <= 0 {
			return
		}
		go p.Send(tea.WindowSizeMsg{Width: cols, Height: rows})
	}

	// Seed initial size so Bubble Tea knows how wide/tall the xterm surface is.
	sendWindowSize(term.Get("cols").Int(), term.Get("rows").Int())

	// Keep Bubble Tea in sync with xterm.js resizing (fit addon updates cols/rows).
	onResize := js.FuncOf(func(this js.Value, args []js.Value) any {
		evt := term
		if len(args) > 0 {
			evt = args[0]
		}
		sendWindowSize(evt.Get("cols").Int(), evt.Get("rows").Int())
		return nil
	})
	defer onResize.Release()
	term.Call("onResize", onResize)

	go func() {
		if _, err := p.Run(); err != nil {
			js.Global().Get("console").Call("error", "Game crashed: ", err.Error())
		}
		close(c)
	}()

	// Wait indefinitely for the core game loop to complete or user to exit
	<-c
}
