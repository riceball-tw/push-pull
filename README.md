# Push Pull Puzzle Game

> Push & Pull box to unlock doors! Puzzle game in terminal or webpage through WASM

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). This version includes a WebAssembly (WASM) port that allows the game to run directly in any modern web browser.

## How to Run

### Native Terminal
To run the game natively in your terminal:
```bash
go run ./cmd/push-pull
```

### Web Browser (WASM)
To build and run the browser version:
1. **Build the WASM binary:**
   ```bash
   ./build-wasm.sh
   ```
2. **Start the local server:**
   ```bash
   go run ./wasm/serve.go
   ```
3. **Open in browser:**
   Navigate to [http://localhost:8080](http://localhost:8080).

## How the WASM Port Works

Running a TUI framework like Bubble Tea in a browser requires bridging the gap between OS-level terminal I/O and the web environment.

### 1. Terminal Emulation with [wterm](https://wterm.dev)
Since browsers don't have a native TUI, we use [wterm](https://wterm.dev) — a terminal emulator with a Zig + WASM core (~12 KB). It renders to the DOM (native text selection, copy/paste, screen reader support) and handles ANSI escape codes and keyboard input. The frontend theme uses Catppuccin via CSS custom properties.

### 2. Standard I/O Redirection
In `cmd/wasm/main.go`, we intercept Bubble Tea's input and output:
- **Output:** We wrap `term.write()` in a Go `io.Writer` and pass it to Bubble Tea using `tea.WithOutput()`.
- **Input:** We listen to `onData` events, pipe them into an `io.Pipe`, and pass the reader to Bubble Tea using `tea.WithInput()`.

The Go code (`cmd/wasm/main.go`) calls into a JS wrapper (`wasm/terminal.js`) that provides an xterm-compatible API, so the Go bridge is decoupled from the specific terminal library.

### 3. Build & Patching
Two build steps happen in `./build-wasm.sh`:

**Go WASM build** — Bubble Tea and its dependencies (like `muesli/termenv` or `containerd/console`) rely on Unix/Windows specific syscalls for terminal state (window resizing, raw mode). The script:
- Runs `go mod vendor` to create a temporary copy of dependencies
- Injects a `tea_js.go` stub into the vendored `bubbletea` source providing no-op implementations of OS-specific functions (`listenForResize`, `initInput`, etc.)
- Compiles with `GOOS=js GOARCH=wasm` using the patched vendor files
- Deletes the `vendor` folder afterwards

**JS bundle build** — The wterm frontend is bundled via esbuild:
- `wasm/terminal.js` imports `@wterm/dom`, creates a `WTerm` instance, and wraps it with an xterm-compatible API for the Go bridge
- esbuild outputs `wasm/wterm-bundle.js` (57 KB) with the core JS and CSS
- `wterm.wasm` (13 KB) is copied from `node_modules/@wterm/core` for separate serving and caching

## Project Structure
- `cmd/push-pull`: Native entrypoint.
- `cmd/wasm`: WebAssembly entrypoint (`main.go` bridges Bubble Tea I/O to JS).
- `internal/game`: Core game logic and Bubble Tea model.
- `wasm/`:
  - `index.html` — Page layout with Catppuccin-themed wterm terminal
  - `terminal.js` — wterm setup + xterm-compat wrapper for Go bridge
  - `wterm-bundle.js` — esbuild bundle (includes `@wterm/dom` and `@wterm/core`)
  - `wterm.wasm` — Zig VT parser binary (~12 KB, served separately for caching)
  - `main.wasm` — Compiled Go WASM binary (Bubble Tea game)
  - `wasm_exec.js` — Go WASM runtime
  - `serve.go` — Static file server for development
  - `package.json` — npm dependencies (`@wterm/dom`, `esbuild`)
- `build-wasm.sh`: Automation script for building both Go WASM and JS bundle.
