package main

import (
	"io"
	"os"
	"syscall/js"

	"push-pull/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

type xtermWriter struct {
	term js.Value
}

func (w *xtermWriter) Write(p []byte) (n int, err error) {
	w.term.Call("write", string(p))
	return len(p), nil
}

func main() {
	// Force full true color support in WASM environment.
	// Without these, lipgloss/termenv cannot detect the terminal profile
	// and may fall back to a limited palette (or no color at all),
	// causing similar hex colors to collapse into the same value.
	os.Setenv("TERM", "xterm-256color")
	os.Setenv("COLORTERM", "truecolor")

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
