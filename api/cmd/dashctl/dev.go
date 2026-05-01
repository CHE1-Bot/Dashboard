package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// runDev starts Vite (5173) and the Go API (8080) in parallel.
// Quitting one tears down the other. Mirrors scripts/dev.mjs exactly.
func runDev(_ []string) error {
	fmt.Println("CHE1 dashboard — development mode")
	fmt.Println("  Web: http://localhost:5173")
	fmt.Println("  API: http://localhost:8080")
	fmt.Println("  Mode: APP_ENV=development DEV_MODE=true")
	fmt.Println()

	// Preflight: if the API port is already taken, report who's holding it
	// before spawning anything. Saves the cryptic
	//   "listen tcp :8080: bind: Only one usage…"
	// + the dance where Vite spins up just to be killed.
	if portInUse("127.0.0.1", 8080) {
		hint := holderHint(8080)
		fmt.Fprintln(os.Stderr, "dashctl: port 8080 is already in use.")
		if hint != "" {
			fmt.Fprintln(os.Stderr, "         ", hint)
		} else {
			fmt.Fprintln(os.Stderr, "         Run `npm start` from another terminal? Stop it first.")
		}
		return fmt.Errorf("api port unavailable")
	}
	if portInUse("127.0.0.1", 5173) {
		fmt.Fprintln(os.Stderr, "dashctl: port 5173 is already in use — Vite will pick a different port.")
	}

	type proc struct {
		name string
		cmd  *exec.Cmd
	}

	apiCmd := exec.Command("go", "run", ".")
	apiCmd.Dir = "api"
	apiCmd.Stdout = prefixed("api", os.Stdout)
	apiCmd.Stderr = prefixed("api", os.Stderr)
	apiCmd.Env = append(os.Environ(), "APP_ENV=development")
	platformSetpgid(apiCmd)

	webCmd := npmCommand("run", "dev:web")
	webCmd.Stdout = prefixed("web", os.Stdout)
	webCmd.Stderr = prefixed("web", os.Stderr)
	webCmd.Env = os.Environ()
	platformSetpgid(webCmd)

	procs := []proc{{"api", apiCmd}, {"web", webCmd}}
	for _, p := range procs {
		if err := p.cmd.Start(); err != nil {
			return fmt.Errorf("start %s: %w", p.name, err)
		}
	}

	// Track the first child that exits so we can tear down the rest.
	exits := make(chan proc, len(procs))
	var wg sync.WaitGroup
	wg.Add(len(procs))
	for _, p := range procs {
		p := p
		go func() {
			defer wg.Done()
			_ = p.cmd.Wait()
			exits <- p
		}()
	}

	// Forward Ctrl+C / SIGTERM to children, then wait briefly.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	shutdown := func(reason string) {
		fmt.Fprintf(os.Stderr, "\ndashctl: %s — stopping siblings\n", reason)
		for _, p := range procs {
			if p.cmd.Process != nil {
				_ = killProcessGroup(p.cmd)
			}
		}
	}

	select {
	case s := <-sig:
		shutdown("received " + s.String())
	case p := <-exits:
		shutdown(p.name + " exited")
	}

	// Give children up to 3 seconds to exit cleanly.
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}
	return nil
}

// npmCommand resolves npm correctly on Windows where it ships as npm.cmd.
func npmCommand(args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		full := append([]string{"/c", "npm"}, args...)
		return exec.Command("cmd", full...)
	}
	return exec.Command("npm", args...)
}
