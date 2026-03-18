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

	go func() {
		err := game.Run(tea.WithInput(inR), tea.WithOutput(outW))
		if err != nil {
			js.Global().Get("console").Call("error", "Game crashed: ", err.Error())
		}
		close(c)
	}()

	// Wait indefinitely for the core game loop to complete or user to exit
	<-c
}

