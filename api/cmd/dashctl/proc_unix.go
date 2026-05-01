//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// platformSetpgid puts the child in its own process group on Unix so
// killProcessGroup can take down its descendants too.
func platformSetpgid(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup sends SIGTERM to the whole group; SIGKILL after 1s
// is handled by the caller's overall timeout.
func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return cmd.Process.Signal(syscall.SIGTERM)
	}
	return syscall.Kill(-pgid, syscall.SIGTERM)
}
