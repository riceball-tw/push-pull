# Push Pull - Bubble Tea TUI with WASM Port

Push Pull is a terminal-based puzzle game built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). This version includes a WebAssembly (WASM) port that allows the game to run directly in any modern web browser.

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

### 1. Terminal Emulation with xterm.js
Since browsers don't have a native TUI, we use [xterm.js](https://xtermjs.org/) to provide a terminal interface in the DOM. It handles rendering ANSI escape codes and capturing keyboard input.

### 2. Standard I/O Redirection
In `cmd/wasm/main.go`, we intercept Bubble Tea's input and output:
- **Output:** We wrap `xterm.write()` in a Go `io.Writer` and pass it to Bubble Tea using `tea.WithOutput()`.
- **Input:** We listen to `onData` events from xterm.js, pipe them into an `io.Pipe`, and pass the reader to Bubble Tea using `tea.WithInput()`.

### 3. Build & Patching
Bubble Tea and its dependencies (like `muesli/termenv` or `containerd/console`) often rely on Unix/Windows specific syscalls for terminal state (like window resizing or raw mode). 

The `./build-wasm.sh` script automates the build process without cluttering your project:
- **Temporary Vendoring**: It runs `go mod vendor` to create a temporary copy of dependencies.
- **The WASM Stub**: It injects a `tea_js.go` stub into the vendored `bubbletea` source. This provides no-op implementations of OS-specific functions (`listenForResize`, `initInput`, etc.) that would otherwise prevent WASM compilation.
- **WASM Compilation**: It runs the build using the patched vendor files.
- **Cleanup**: It deletes the `vendor` folder immediately after the build is finished.

This approach ensures a "zero-maintenance" workflow where no patched code needs to be stored in your repository.

## Project Structure
- `cmd/push-pull`: Native entrypoint.
- `cmd/wasm`: WebAssembly entrypoint.
- `internal/game`: Core game logic and Bubble Tea model.
- `wasm/`: Web assets (HTML, JS, CSS) and the compiled `.wasm` binary.
- `build-wasm.sh`: Automation script for building the web version.
