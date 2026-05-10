# sudoku

Simple sudoku solver and unique-solution sudoku generator using gini sat solver.

wasm backend [https://fbundle.github.io/sudoku](https://fbundle.github.io/sudoku/)

## Features

- Unique-solution board generation via SAT solver
- Implication hint (finds a cell that can be deduced without branching)
- Undo
- Web interface

## Running

**Server mode** (multiplayer, Go backend)
```bash
go run .
# open http://localhost:3000/sudoku/
```

**WASM mode** (single player, runs entirely in the browser)
```bash
./build_wasm.sh
cd docs && python3 -m http.server 3000
# open http://localhost:3000
```
