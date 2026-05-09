#!/usr/bin/env bash
set -e
GOOS=js GOARCH=wasm go build -o docs/sudoku.wasm .
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" docs/wasm_exec.js
