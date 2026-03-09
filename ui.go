package main

import (
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// tickMsg signals that the simulation should advance one generation.
type tickMsg time.Time

const (
	stagnationWindow    = 20
	stagnationThreshold = 3
)

// model holds the application state for the Bubble Tea program.
type model struct {
	grid       *Grid
	width      int
	height     int
	generation int
	paused     bool
	bpm        int
	density    float64
	darkBg     bool
	cafCmd     *exec.Cmd
	autoReseed bool                  // auto-reseed on stagnation (on by default)
	aliveHist  [stagnationWindow]int // ring buffer of recent alive counts
	histLen    int                   // number of entries filled in aliveHist
	histIdx    int                   // next write position in aliveHist
}

// tickInterval converts a BPM value to a time.Duration per beat.
func tickInterval(bpm int) time.Duration {
	return time.Duration(60000/bpm) * time.Millisecond
}

// isStagnant returns true if the alive count history shows the simulation
// has settled into a static or near-static state.
func (m model) isStagnant() bool {
	if m.histLen < stagnationWindow {
		return false
	}
	minC, maxC := m.aliveHist[0], m.aliveHist[0]
	for i := 1; i < stagnationWindow; i++ {
		if m.aliveHist[i] < minC {
			minC = m.aliveHist[i]
		}
		if m.aliveHist[i] > maxC {
			maxC = m.aliveHist[i]
		}
	}
	return maxC-minC <= stagnationThreshold
}

// newModel creates a model with the given parameters. The grid is nil
// until the first WindowSizeMsg arrives.
func newModel(bpm int, density float64, autoReseed bool, cafCmd *exec.Cmd) model {
	return model{
		bpm:        bpm,
		density:    density,
		autoReseed: autoReseed,
		cafCmd:     cafCmd,
	}
}

// Init starts the tick timer immediately.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(tickInterval(m.bpm), func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		tea.RequestBackgroundColor,
	)
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
		m.height = msg.Height
		if m.height < 1 {
			m.height = 1
		}
		m.grid = NewGrid(m.width, m.height)
		m.grid.Randomize(m.density)
		m.generation = 0
		m.histLen = 0
		m.histIdx = 0
		return m, nil

	case tea.BackgroundColorMsg:
		m.darkBg = msg.IsDark()
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
				m.histLen = 0
				m.histIdx = 0
			}
			return m, nil
		case "+", "=":
			m.bpm += 10
			if m.bpm > 600 {
				m.bpm = 600
			}
			return m, nil
		case "-", "_":
			m.bpm -= 10
			if m.bpm < 10 {
				m.bpm = 10
			}
			return m, nil
		case "a":
			m.autoReseed = !m.autoReseed
			return m, nil
		}

	case tickMsg:
		if !m.paused && m.grid != nil {
			m.grid = m.grid.Tick()
			m.generation++

			// Track alive count for stagnation detection.
			alive := m.grid.CountAlive()
			m.aliveHist[m.histIdx] = alive
			m.histIdx = (m.histIdx + 1) % stagnationWindow
			if m.histLen < stagnationWindow {
				m.histLen++
			}

			// Auto-reseed when stagnant or completely dead.
			if m.autoReseed && (alive == 0 || m.isStagnant()) {
				m.grid.Randomize(m.density)
				m.generation = 0
				m.histLen = 0
				m.histIdx = 0
			}
		}
		return m, tea.Tick(tickInterval(m.bpm), func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// View renders the grid as a full-screen terminal view.
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
		if y > 0 {
			sb.WriteByte('\n')
		}
		for x := 0; x < m.grid.width; x++ {
			sb.WriteString(RenderCell(m.grid.Get(x, y), m.generation, x, y, m.darkBg))
		}
	}

	v := tea.NewView(sb.String())
	v.AltScreen = true
	v.WindowTitle = "stayinalive"
	return v
}
