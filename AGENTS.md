# AGENTS.md

Guidelines for AI agents working on the stayinalive codebase.

## Project Overview

stayinalive is a macOS terminal TUI that runs Conway's Game of Life with disco-themed visuals and prevents the Mac from sleeping. It uses Bubble Tea V2 for the TUI framework and Lip Gloss V2 for cell styling.

## Language and Dependencies

- **Language**: Go (1.25.5+)
- **Module name**: `stayinalive`
- **Key dependencies**:
  - `charm.land/bubbletea/v2` (v2.0.1): TUI framework (Model-View-Update pattern)
  - `charm.land/lipgloss/v2` (v2.0.0): Terminal styling and colors

## File Layout

All source files are in the repository root under package `main`.

| File             | Responsibility                                     |
| ---------------- | -------------------------------------------------- |
| `main.go`        | Entry point, CLI flag parsing, caffeinate lifecycle |
| `game.go`        | Conway's Game of Life engine (pure logic, no UI)   |
| `ui.go`          | Bubble Tea model, Init, Update, View               |
| `disco.go`       | Disco color palettes, cell rendering               |
| `caffeinate.go`  | macOS caffeinate subprocess start/stop              |

## Architecture

The app follows Bubble Tea's Model-View-Update pattern:

1. **Model** (`ui.go`): holds grid state, terminal dimensions, generation counter, pause state, BPM, background brightness flag, and caffeinate process handle.
2. **Update** (`ui.go`): handles window resize, key presses, tick events, and background color detection messages.
3. **View** (`ui.go`): renders the grid cell-by-cell using `disco.RenderCell` as a full-screen display.
4. **Game engine** (`game.go`): pure logic. `Grid.Tick()` returns a new grid (double-buffered). Toroidal wrapping on all edges.
5. **Rendering** (`disco.go`): color cycling via `(generation + x + y) % len(palette)`. Dead-cell color adapts to terminal background brightness.

## Conventions

### Build and Verification

```bash
go build ./...    # Must succeed with exit code 0
go vet ./...      # Must pass with no warnings
```

Always run both commands after making changes. There are no unit tests currently.

### Code Style

- Standard Go formatting (`gofmt`).
- All exported functions and types have doc comments.
- No third-party dependencies beyond Bubble Tea and Lip Gloss.

### Bubble Tea V2 Patterns

- `tea.Cmd` is `func() tea.Msg`. Functions like `tea.RequestBackgroundColor` are already of this type. Pass them by name (no parentheses) to `tea.Batch`.
- `Init()` uses `tea.Batch` to combine multiple startup commands.
- `View()` returns `tea.View` (not a plain string). Set `AltScreen = true` and `WindowTitle = "stayinalive"`.
- Use `tea.Tick` for recurring timer events. Return a new tick command from the tick message handler.

### Grid and Rendering

- Grid dimensions are derived from terminal size: `width / 2` columns (cells are ~2 chars wide), `height` rows (the grid fills the entire terminal).
- Alive cells: full-block character `U+2588`, styled with disco palette colors and bold.
- Dead cells: middle dot `U+00B7`, styled with adaptive foreground color (`#333333` on dark, `#CCCCCC` on light backgrounds).
- `RenderCell` and `DeadCellStyle` accept a `darkBg bool` parameter for background-adaptive rendering.

### Terminal Background Detection

- `Init()` sends `tea.RequestBackgroundColor` at startup.
- `Update()` handles `tea.BackgroundColorMsg` and stores `msg.IsDark()` in the model's `darkBg` field.
- Default is `darkBg = false` (light-background styling) for terminals that do not respond to the query.

### macOS Caffeinate

- `caffeinate -di` is spawned at startup to prevent display and idle sleep.
- The process is killed explicitly on quit and via `defer` in `main()`.
- If caffeinate fails to start, the app continues with a warning (non-fatal).

## Key Constraints

- **macOS only**: relies on `/usr/bin/caffeinate`.
- **No new files without reason**: the five-file structure is intentional. Do not add files unless a feature genuinely requires it.
- **No new dependencies**: avoid adding modules beyond the existing two.
- **Cosmetic defaults**: when the terminal does not support background color queries, dead cells default to light-background styling. This is acceptable.
