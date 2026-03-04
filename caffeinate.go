package main

import "os/exec"

// StartCaffeinate starts the macOS caffeinate process to prevent display and idle sleep.
// The -d flag prevents display sleep; the -i flag prevents idle sleep.
func StartCaffeinate() (*exec.Cmd, error) {
	cmd := exec.Command("/usr/bin/caffeinate", "-di")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// StopCaffeinate kills the caffeinate process and reaps the zombie.
func StopCaffeinate(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
}
