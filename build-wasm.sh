#!/bin/bash
set -e

echo "Creating temporary vendor folder to apply WASM patch..."
go mod vendor

PATCH_FILE="vendor/github.com/charmbracelet/bubbletea/tea_js.go"
echo "Applying WASM build patch to vendored bubbletea source..."
cat << 'EOF' > "$PATCH_FILE"
//go:build js || wasm

package tea

import "os"

const suspendSupported = false
func suspendProcess() {}
func openInputTTY() (*os.File, error) { return nil, nil }
func (p *Program) initInput() (err error) { return nil }
func (p *Program) listenForResize(done chan struct{}) {}
EOF

echo "Copying wasm_exec.js..."
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ./wasm/ 2>/dev/null || cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./wasm/

echo "Compiling via GOOS=js GOARCH=wasm..."
GOOS=js GOARCH=wasm go build -mod=vendor -o ./wasm/main.wasm ./cmd/wasm

echo "Cleaning up temporary vendor folder..."
rm -rf vendor

echo "Done! Run 'go run ./wasm/serve.go' to start the local server."