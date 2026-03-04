# stayinalive

A terminal program that runs Conway's Game of Life with disco visuals while preventing your Mac from going to sleep.

## Project Structure

```
stayinalive/
├── main.go          # Entry point, program setup, CLI flags
├── game.go          # Conway's Game of Life engine (pure logic)
├── ui.go            # Bubble Tea model, Update, View
├── disco.go         # Disco color palettes and cycling effects
├── caffeinate.go    # macOS caffeinate subprocess lifecycle
├── go.mod
└── go.sum
```

## Dependencies

| Package        | Import Path              | Purpose                          |
| -------------- | ------------------------ | -------------------------------- |
| Bubble Tea V2  | `charm.land/bubbletea/v2` | TUI framework (Model-View-Update) |
| Lip Gloss V2   | `charm.land/lipgloss/v2`  | Cell styling, colors, layout     |

That's it. Two deps plus `os/exec` from stdlib for `caffeinate`.

## Component Plan

### 1. `game.go` — Game of Life Engine

Pure logic, no UI concerns.

- **`Grid` type**: 2D `[][]bool` with width/height
- **`NewGrid(w, h int) *Grid`**: allocate grid
- **`Randomize(density float64)`**: seed cells randomly (good default: ~0.3)
- **`Tick() *Grid`**: compute next generation using standard Conway rules (birth on 3 neighbors, survival on 2-3, death otherwise). Returns a new grid (double-buffer to avoid mutation during iteration)
- **`Set(x, y int, alive bool)`** / **`Get(x, y int) bool`**: cell access
- **`CountAlive() int`**: for stats display
- **Wrapping edges**: toroidal topology (cells wrap around) keeps the simulation interesting indefinitely

### 2. `caffeinate.go` — Sleep Prevention

- **`StartCaffeinate() (*exec.Cmd, error)`**: spawn `/usr/bin/caffeinate -di` as a child process. `-d` prevents display sleep, `-i` prevents idle sleep.
- **`StopCaffeinate(cmd *exec.Cmd)`**: send SIGTERM, then Wait(). Called on program exit.
- The caffeinate process inherits the parent's lifecycle — if stayinalive is killed, caffeinate dies too (child process). But explicit cleanup is still good practice.

### 3. `disco.go` — Disco Visual Theme

This is where the personality lives.

- **Color palette**: Define 6-8 disco colors (hot pink `#FF69B4`, gold `#FFD700`, electric blue `#00BFFF`, lime `#39FF14`, violet `#8B00FF`, orange `#FF6600`). Cycle through them per-generation.
- **`CellStyle(generation int, x, y int) lipgloss.Style`**: returns a Lip Gloss style for a live cell. Color shifts based on generation count, creating a rainbow wave across the grid.
- **Dead cell rendering**: dim dot or space. Subtle so alive cells pop.
- **Alive cell rendering**: full block (`█`) or circle (`●`) styled with the current disco color.
- **Color cycling strategy**: `offset = (generation + x + y) % len(palette)` — creates a diagonal wave of color across the grid that shifts each tick.
- **Optional**: use `lipgloss.Lighten`/`lipgloss.Darken` for a pulsing brightness effect synced to an internal beat counter.

### 4. `ui.go` — Bubble Tea Application

The core TUI wiring.

**Model struct fields:**

- `grid` — current game state (`*Grid`)
- `width` — terminal columns
- `height` — terminal rows
- `generation` — tick counter
- `paused` — pause toggle
- `bpm` — current tick speed (default 104)
- `cafCmd` — caffeinate process handle (`*exec.Cmd`)

**Init():**

- Return `tea.Tick` cmd at BPM interval (~577ms for 104 BPM) to drive the simulation
- Return `tea.RequestBackgroundColor` for adaptive styling

**Update() message handling:**

| Message                          | Action                                                      |
| -------------------------------- | ----------------------------------------------------------- |
| `tea.WindowSizeMsg`              | Resize grid to fit terminal (w/2 for cell width, h-3 for status bar) |
| `tea.KeyPressMsg "q"` / `ctrl+c` | Kill caffeinate, return `tea.Quit`                          |
| `tea.KeyPressMsg "space"`        | Toggle pause                                                |
| `tea.KeyPressMsg "r"`           | Re-randomize the grid                                       |
| `tea.KeyPressMsg "+"` / `"-"`   | Speed up / slow down BPM                                    |
| `tickMsg`                        | If not paused: advance generation, return next tick cmd     |
| `tea.BackgroundColorMsg`         | Store dark/light for adaptive colors                        |

**View():**

- Render grid: iterate cells, apply `disco.CellStyle()` per live cell
- Build string row-by-row using styled `Render()` calls
- Status bar at bottom: generation count, alive cell count, BPM, pause state, and a rotating disco tagline (e.g. `"Ah, ha, ha, ha, stayin' alive"`)
- Return `tea.View` with `AltScreen = true` and `WindowTitle = "stayinalive"`

### 5. `main.go` — Entry Point

- Parse optional CLI flags: `--bpm` (default 104), `--density` (default 0.3)
- Call `StartCaffeinate()`
- Create Bubble Tea program: `tea.NewProgram(initialModel)`
- `defer StopCaffeinate(cmd)` for cleanup
- Run and handle exit error

## Key Design Decisions

- **BPM-driven tick rate**: Stayin' Alive is 104 BPM (~577ms/beat). Each beat advances one generation. User can adjust with `+`/`-`.
- **Toroidal grid**: wrapping edges prevent the simulation from dying out at boundaries.
- **Grid auto-sizes to terminal**: uses `tea.WindowSizeMsg` and re-creates the grid on resize.
- **Caffeinate as subprocess**: simplest zero-dep approach. No cgo, no IOKit bindings needed.
- **Double-buffered grid**: `Tick()` returns a new grid rather than mutating in place, avoiding artifacts.

## Controls

| Key          | Action                    |
| ------------ | ------------------------- |
| `q` / `ctrl+c` | Quit (stops caffeinate) |
| `space`      | Pause / resume            |
| `r`          | Randomize grid            |
| `+` / `-`   | Increase / decrease BPM   |
