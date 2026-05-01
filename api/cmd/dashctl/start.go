package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

// runStart performs a full production build + run, mirroring scripts/start.mjs:
//  1. vite build           → ./dist
//  2. copy dist            → api/dist (so go:embed picks it up)
//  3. go build -trimpath   → api/che1-dashboard(.exe)
//  4. exec the binary with APP_ENV=production
func runStart(_ []string) error {
	fmt.Println("CHE1 dashboard — production build")

	binName := "che1-dashboard"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// 1. Vite build
	if err := runStep("vite build", npmCommand("run", "build")); err != nil {
		return err
	}

	// 2. Copy dist → api/dist (replace, so old hashed assets don't linger)
	apiDist := filepath.Join("api", "dist")
	if err := os.RemoveAll(apiDist); err != nil {
		return fmt.Errorf("clean api/dist: %w", err)
	}
	if err := copyTree("dist", apiDist); err != nil {
		return fmt.Errorf("copy dist: %w", err)
	}
	fmt.Println("✓ copied dist → api/dist")

	// 3. Go build — restrict to the server package so we don't accidentally
	// build dashctl itself.
	build := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", binName, ".")
	build.Dir = "api"
	if err := runStep("go build", build); err != nil {
		return err
	}
	fmt.Printf("✓ built api/%s\n", binName)

	// 4. Run with APP_ENV=production. Forward signals so Ctrl+C reaches the
	// child cleanly, matching what scripts/start.mjs did.
	binary := filepath.Join("api", binName)
	fmt.Printf("\n→ launching %s (APP_ENV=production)\n\n", binary)
	srv := exec.Command(binary)
	srv.Stdout = os.Stdout
	srv.Stderr = os.Stderr
	srv.Stdin = os.Stdin
	srv.Env = append(os.Environ(), "APP_ENV=production")
	platformSetpgid(srv)
	if err := srv.Start(); err != nil {
		return fmt.Errorf("launch %s: %w", binary, err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	done := make(chan error, 1)
	go func() { done <- srv.Wait() }()

	select {
	case s := <-sig:
		fmt.Fprintf(os.Stderr, "\ndashctl: received %s — shutting down\n", s)
		_ = killProcessGroup(srv)
		<-done
	case err := <-done:
		if err != nil {
			return err
		}
	}
	return nil
}

// runStep prints what it's doing, runs the command synchronously with stdio
// passthrough, and surfaces non-zero exits as errors.
func runStep(name string, c *exec.Cmd) error {
	fmt.Printf("\n→ %s: %s\n", name, c.String())
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", name, err)
	}
	return nil
}

// copyTree recursively copies src into dst.
func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		out := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(out, info.Mode())
		}
		// Regular file
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		o, err := os.Create(out)
		if err != nil {
			return err
		}
		defer o.Close()
		if _, err := io.Copy(o, in); err != nil {
			return err
		}
		return os.Chmod(out, info.Mode())
	})
}
