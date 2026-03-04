# stayinalive

A terminal program that runs Conway's Game of Life with disco visuals while preventing your Mac from going to sleep.

Built with [Bubble Tea V2](https://charm.sh/bubbletea/) and [Lip Gloss V2](https://charm.sh/lipgloss/).

## Features

- **Conway's Game of Life** with toroidal (wrapping) grid that auto-sizes to your terminal
- **Disco color palette**: alive cells cycle through hot pink, gold, electric blue, lime, violet, orange, cyan, and red in a diagonal wave pattern
- **Adaptive background detection**: automatically detects whether your terminal has a dark or light background and adjusts dead-cell styling accordingly
- **macOS sleep prevention**: spawns `caffeinate -di` to prevent display and idle sleep while running
- **BPM-driven simulation**: defaults to 400 BPM for fast, fluid animation; adjustable at runtime with `+`/`-` keys

## Requirements

- macOS (uses `/usr/bin/caffeinate` for sleep prevention)
- Go 1.25.5 or later
- A terminal emulator (works best with one that supports true color and background color queries)

## Installation

```bash
git clone <repo-url> && cd stayinalive
go build -o stayinalive .
```

## Usage

```bash
./stayinalive [flags]
```

### Flags

| Flag        | Default | Description                         |
| ----------- | ------- | ----------------------------------- |
| `--bpm`     | `400`   | Tick speed in beats per minute      |
| `--density` | `0.3`   | Initial cell density (0.0 to 1.0)   |

### Examples

```bash
# Run with defaults (400 BPM, 30% density)
./stayinalive

# Slower simulation with more cells
./stayinalive --bpm 120 --density 0.5

# Slow, sparse simulation
./stayinalive --bpm 30 --density 0.1
```

## Controls

| Key             | Action                      |
| --------------- | --------------------------- |
| `q` / `Ctrl+C`  | Quit (stops caffeinate)    |
| `Space`          | Pause / resume             |
| `r`              | Randomize the grid         |
| `+` / `=`        | Increase BPM (max 600)     |
| `-` / `_`        | Decrease BPM (min 10)      |

## How It Works

The simulation runs in an alternate screen buffer. Each generation advances on a timer driven by the configured BPM. Alive cells are rendered as full-block characters (`U+2588`) styled with rotating disco colors. Dead cells are rendered as middle dots (`U+00B7`) with a color that adapts to your terminal background:

- **Dark terminals**: dead cells render in `#333333` (dim, unobtrusive)
- **Light terminals**: dead cells render in `#CCCCCC` (visible but subtle)

The grid uses toroidal topology: cells wrap around edges, which keeps the simulation active indefinitely.

On startup, the app spawns `caffeinate -di` as a child process. The `-d` flag prevents display sleep and the `-i` flag prevents idle sleep. The process is killed when you quit.

## Terminal Background Detection

At startup, the app sends a background color query to the terminal using Bubble Tea's `RequestBackgroundColor` command. When the terminal responds, the app stores whether the background is dark or light and adjusts dead-cell rendering.

Terminals that do not support background color queries will not respond. In that case, the app defaults to light-background styling (`#CCCCCC` for dead cells). This is cosmetic-only and does not affect functionality.

## Development

```bash
# Build
go build ./...

# Vet
go vet ./...
```

There are no unit tests at this time. Verification is done through build, vet, and manual testing.
