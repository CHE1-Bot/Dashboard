//go:build windows

package main

import (
	"fmt"
	"os/exec"
)

// platformSetpgid is a no-op on Windows — process trees are killed via taskkill /T.
func platformSetpgid(cmd *exec.Cmd) {}

// killProcessGroup uses taskkill /T /F so child processes (e.g. npm spawning
// node spawning vite) all die together.
func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	pid := cmd.Process.Pid
	tk := exec.Command("taskkill", "/PID", fmt.Sprint(pid), "/T", "/F")
	// Discard output — taskkill is noisy and we don't care about its return.
	_ = tk.Run()
	return nil
}
