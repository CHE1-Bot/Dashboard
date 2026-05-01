// dashctl is the development & build CLI for the CHE1 dashboard. It replaces
// the Node-authored scripts/dev.mjs and scripts/start.mjs.
//
// Subcommands:
//   dashctl dev         - run Vite (5173) + Go API (8080), APP_ENV=development.
//   dashctl start       - production: vite build, embed dist, build Go binary, run.
//   dashctl gen-pages   - regenerate src/dashboard/data/pages.js from the Go
//                         route definitions in this package.
//
// dashctl auto-discovers the project root by walking up from CWD until it
// finds a package.json. All commands run with cwd at that root.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dashctl:", err)
		os.Exit(1)
	}
	if err := os.Chdir(root); err != nil {
		fmt.Fprintln(os.Stderr, "dashctl:", err)
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]
	var runErr error
	switch cmd {
	case "dev":
		runErr = runDev(args)
	case "start":
		runErr = runStart(args)
	case "gen-pages":
		runErr = runGenPages(args)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "dashctl: unknown command %q\n\n", cmd)
		usage()
		os.Exit(2)
	}
	if runErr != nil {
		fmt.Fprintln(os.Stderr, "dashctl:", runErr)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`dashctl — CHE1 dashboard tooling

USAGE
  dashctl <command>

COMMANDS
  dev          Run Vite + Go API in development mode (APP_ENV=development).
  start        Build the SPA, embed it in the Go binary, run with APP_ENV=production.
  gen-pages    Regenerate src/dashboard/data/pages.js from Go route defs.
  help         Show this help.`)
}

// findProjectRoot walks up from CWD until it finds a directory containing
// package.json — that is the dashboard project root. Stops at the filesystem
// root if not found.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("could not locate project root (no package.json found above " + originalDir() + ")")
		}
		dir = parent
	}
}

// originalDir returns the launching cwd for nicer error messages.
func originalDir() string {
	if d, err := os.Getwd(); err == nil {
		return d
	}
	return "."
}
