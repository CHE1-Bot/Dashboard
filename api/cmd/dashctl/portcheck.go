package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// portInUse reports whether something is already listening on the port.
// host is unused but kept for backwards compatibility.
//
// Strategy:
//  1. Dial 127.0.0.1:port — catches any active listener regardless of which
//     address it bound to (Windows happily lets you bind 127.0.0.1:p when
//     someone else has :p, so binding alone is unreliable).
//  2. If dial fails, try binding 0.0.0.0:p the way the API server will, in
//     case the port is in TIME_WAIT or otherwise can't accept a fresh bind.
func portInUse(host string, port int) bool {
	_ = host
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	if conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond); err == nil {
		_ = conn.Close()
		return true
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	_ = l.Close()
	time.Sleep(50 * time.Millisecond)
	return false
}

// holderHint tries to identify which process is on the port. Best-effort —
// a missing tool returns "" and the caller falls back to a generic message.
func holderHint(port int) string {
	switch runtime.GOOS {
	case "windows":
		return windowsHolder(port)
	default:
		return unixHolder(port)
	}
}

func windowsHolder(port int) string {
	out, err := exec.Command("netstat", "-ano").Output()
	if err != nil {
		return ""
	}
	needle := fmt.Sprintf(":%d", port)
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, needle) || !strings.Contains(line, "LISTENING") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		pid := fields[len(fields)-1]
		name := windowsProcessName(pid)
		if name != "" {
			return fmt.Sprintf("PID %s (%s) — kill with:  taskkill /PID %s /F", pid, name, pid)
		}
		return fmt.Sprintf("PID %s — kill with:  taskkill /PID %s /F", pid, pid)
	}
	return ""
}

func windowsProcessName(pid string) string {
	out, err := exec.Command("tasklist", "/FI", "PID eq "+pid, "/NH", "/FO", "CSV").Output()
	if err != nil {
		return ""
	}
	// tasklist CSV: "image","pid","session","sess#","mem"
	line := strings.TrimSpace(string(out))
	if !strings.HasPrefix(line, `"`) {
		return ""
	}
	end := strings.Index(line[1:], `"`)
	if end <= 0 {
		return ""
	}
	return line[1 : 1+end]
}

func unixHolder(port int) string {
	// Try lsof first; fall back to ss.
	if out, err := exec.Command("lsof", "-iTCP:"+fmt.Sprint(port), "-sTCP:LISTEN", "-Pn", "-Fpc").Output(); err == nil && len(out) > 0 {
		var pid, name string
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if strings.HasPrefix(line, "p") {
				pid = strings.TrimPrefix(line, "p")
			} else if strings.HasPrefix(line, "c") {
				name = strings.TrimPrefix(line, "c")
			}
		}
		if pid != "" {
			return fmt.Sprintf("PID %s (%s) — kill with:  kill -9 %s", pid, name, pid)
		}
	}
	if out, err := exec.Command("ss", "-ltnp", "sport = :"+fmt.Sprint(port)).Output(); err == nil {
		txt := string(out)
		if strings.Contains(txt, "pid=") {
			return strings.TrimSpace(txt)
		}
	}
	return ""
}
