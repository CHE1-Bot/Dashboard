package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed all:dist
var embeddedDist embed.FS

// spaHandler serves the Svelte build. It tries the embedded FS first (produced
// by copying the Vite `dist/` tree into api/dist before `go build`), then
// falls back to reading from ../dist on disk for local dev. When the request
// path doesn't match a file, it serves index.html so client-side routing works.
func spaHandler() http.Handler {
	sub, err := fs.Sub(embeddedDist, "dist")
	hasEmbedded := err == nil
	if hasEmbedded {
		// The embed directive always succeeds, but the tree may be empty on
		// first build. Detect that so we can fall through to disk.
		if f, err := sub.Open("index.html"); err != nil {
			hasEmbedded = false
		} else {
			f.Close()
		}
	}

	onDisk := http.Dir("../dist")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		p := strings.TrimPrefix(r.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}

		if hasEmbedded {
			if data, err := openEmbedded(sub, p); err == nil {
				serveBytes(w, r, p, data)
				return
			}
			// SPA fallback
			if data, err := openEmbedded(sub, "index.html"); err == nil {
				serveBytes(w, r, "index.html", data)
				return
			}
		}

		// Disk fallback (local dev without a re-build)
		f, err := onDisk.Open(p)
		if err != nil {
			f, err = onDisk.Open("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			p = "index.html"
		}
		defer f.Close()
		stat, _ := f.Stat()
		if stat != nil && stat.IsDir() {
			f.Close()
			f, err = onDisk.Open("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			p = "index.html"
		}
		data, err := io.ReadAll(f)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		serveBytes(w, r, p, data)
	})
}

func openEmbedded(sub fs.FS, name string) ([]byte, error) {
	f, err := sub.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func serveBytes(w http.ResponseWriter, r *http.Request, name string, data []byte) {
	ext := strings.ToLower(path.Ext(name))
	ct := "application/octet-stream"
	switch ext {
	case ".html":
		ct = "text/html; charset=utf-8"
	case ".js", ".mjs":
		ct = "application/javascript; charset=utf-8"
	case ".css":
		ct = "text/css; charset=utf-8"
	case ".json":
		ct = "application/json; charset=utf-8"
	case ".svg":
		ct = "image/svg+xml"
	case ".png":
		ct = "image/png"
	case ".jpg", ".jpeg":
		ct = "image/jpeg"
	case ".webp":
		ct = "image/webp"
	case ".ico":
		ct = "image/x-icon"
	case ".woff2":
		ct = "font/woff2"
	case ".woff":
		ct = "font/woff"
	case ".map":
		ct = "application/json; charset=utf-8"
	}
	w.Header().Set("Content-Type", ct)
	// Hashed assets under /assets/ can cache aggressively; index.html must not.
	if strings.HasPrefix(r.URL.Path, "/assets/") && name != "index.html" {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	_, _ = w.Write(data)
}

// ensureDistPlaceholder writes a sentinel index.html into api/dist before
// `go build` if the folder is missing so the `embed` directive compiles even
// when the frontend hasn't been built yet. Intended for CI.
func ensureDistPlaceholder() {
	_ = os.MkdirAll("dist", 0o755)
	p := "dist/index.html"
	if _, err := os.Stat(p); err == nil {
		return
	}
	_ = os.WriteFile(p, []byte("<!doctype html><title>CHE1 Dashboard</title><p>Run npm run build && copy dist → api/dist.</p>"), 0o644)
}
