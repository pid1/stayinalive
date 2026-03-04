package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	bpm := flag.Int("bpm", 400, "beats per minute (tick speed)")
	density := flag.Float64("density", 0.3, "initial cell density (0.0-1.0)")
	flag.Parse()

	cafCmd, err := StartCaffeinate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not start caffeinate: %v\n", err)
		cafCmd = nil
	}
	defer StopCaffeinate(cafCmd)

	m := newModel(*bpm, *density, cafCmd)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
