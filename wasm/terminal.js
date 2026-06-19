import { WTerm } from "@wterm/dom";

/**
 * Initialize wterm with a Catppuccin theme and expose an xterm-compatible
 * API for the Go WASM bridge (cmd/wasm/main.go).
 *
 * The Go side calls:
 *   term.write(str)         — write output
 *   term.onData(callback)   — register input handler (METHOD, not property)
 *   term.onResize(callback) — register resize handler (METHOD, not property)
 *   term.focus()            — focus terminal
 *   term.cols / term.rows   — current grid dimensions
 *
 * wterm uses properties (term.onData = cb) instead of methods.
 * This wrapper bridges the difference so cmd/wasm/main.go needs ZERO changes.
 */
async function initTerminal() {
  const el = document.getElementById("terminal-container");
  if (!el) throw new Error("#terminal-container not found");

  // Create wterm with Catppuccin-compatible settings.
  // The theme is applied via CSS custom properties on the .wterm element
  // (defined in index.html <style>), so we only pass non-styling options.
  const wterm = new WTerm(el, {
    cursorBlink: true,
    autoResize: true,
    wasmUrl: "wterm.wasm",
  });

  // Init loads the Zig WASM binary (~12 KB, embedded as base64 in the bundle).
  // Must complete before writing data or reading cols/rows.
  await wterm.init();

  // ── xterm-compatible wrapper for Go WASM ──────────────────────
  // The Go code (cmd/wasm/main.go) calls JS via syscall/js:
  //   term.Call("onData", cb)     → term.onData(cb)
  //   term.Call("onResize", cb)   → term.onResize(cb)
  //   term.Call("write", str)     → term.write(str)
  //   term.Call("focus")          → term.focus()
  //   term.Get("cols")            → term.cols
  //   term.Get("rows")            → term.rows
  //
  // wterm uses property assignment for onData/onResize, not method calls.
  // This wrapper converts the method-call pattern to property assignment.

  const compat = {
    write: (data) => wterm.write(data),
    focus: () => wterm.focus(),

    onData: function (callback) {
      wterm.onData = callback;
    },

    onResize: function (callback) {
      wterm.onResize = callback;
    },

    get cols() {
      return wterm.cols;
    },
    get rows() {
      return wterm.rows;
    },
  };

  // Expose globally for the Go WASM bridge to find
  window.term = compat;

  // Mobile button callback — also exposed globally (set by Go side)
  window._wasmOnData = null;

  return compat;
}

// Start initialization immediately. Export the promise so index.html
// can await it before loading Go WASM.
const ready = initTerminal();
window._wtermReady = ready;
