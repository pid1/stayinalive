package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// tickMsg signals that the simulation should advance one generation.
type tickMsg time.Time

// model holds the application state for the Bubble Tea program.
type model struct {
	grid       *Grid
	width      int
	height     int
	generation int
	paused     bool
	bpm        int
	density    float64
	cafCmd     *exec.Cmd
}

// tickInterval converts a BPM value to a time.Duration per beat.
func tickInterval(bpm int) time.Duration {
	return time.Duration(60000/bpm) * time.Millisecond
}

// newModel creates a model with the given parameters. The grid is nil
// until the first WindowSizeMsg arrives.
func newModel(bpm int, density float64, cafCmd *exec.Cmd) model {
	return model{
		bpm:     bpm,
		density: density,
		cafCmd:  cafCmd,
	}
}

// Init starts the tick timer immediately.
func (m model) Init() tea.Cmd {
	return tea.Tick(tickInterval(m.bpm), func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages: window resizes, key presses, and tick events.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// Each cell renders as ~2 chars wide due to the full-block character.
		m.width = msg.Width / 2
		if m.width < 1 {
			m.width = 1
		}
		// Reserve 3 lines for the status bar.
		m.height = msg.Height - 3
		if m.height < 1 {
			m.height = 1
		}
		m.grid = NewGrid(m.width, m.height)
		m.grid.Randomize(m.density)
		m.generation = 0
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			StopCaffeinate(m.cafCmd)
			return m, tea.Quit
		case " ":
			m.paused = !m.paused
			return m, nil
		case "r":
			if m.grid != nil {
				m.grid.Randomize(m.density)
				m.generation = 0
			}
			return m, nil
		case "+", "=":
			m.bpm += 10
			if m.bpm > 300 {
				m.bpm = 300
			}
			return m, nil
		case "-", "_":
			m.bpm -= 10
			if m.bpm < 10 {
				m.bpm = 10
			}
			return m, nil
		}

	case tickMsg:
		if !m.paused && m.grid != nil {
			m.grid = m.grid.Tick()
			m.generation++
		}
		return m, tea.Tick(tickInterval(m.bpm), func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// View renders the grid and status bar as a terminal view.
func (m model) View() tea.View {
	if m.grid == nil {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		v.WindowTitle = "stayinalive"
		return v
	}

	var sb strings.Builder

	// Render the grid.
	for y := 0; y < m.grid.height; y++ {
		for x := 0; x < m.grid.width; x++ {
			sb.WriteString(RenderCell(m.grid.Get(x, y), m.generation, x, y))
		}
		sb.WriteByte('\n')
	}

	// Build the status bar.
	pauseLabel := "PLAYING"
	if m.paused {
		pauseLabel = "PAUSED"
	}
	sb.WriteString(fmt.Sprintf("Gen: %d | Alive: %d | BPM: %d | %s\n",
		m.generation, m.grid.CountAlive(), m.bpm, pauseLabel))
	sb.WriteString(DiscoTagline(m.generation))
	sb.WriteByte('\n')
	sb.WriteString("q: quit | space: pause | r: randomize | +/-: speed")

	v := tea.NewView(sb.String())
	v.AltScreen = true
	v.WindowTitle = "stayinalive"
	return v
}
