package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// Endpoints under /api/bot/* are consumed by the CHE1 Bot and authenticated
// with Bearer DASHBOARD_API_KEY (matches internal/dashboard.Client in the Bot
// repo). These do NOT use the user session cookie.
func handleBotRoutes(w http.ResponseWriter, r *http.Request) {
	if !botAPIKeyValid(r) {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/bot/")
	parts := strings.Split(rest, "/")
	if len(parts) >= 3 && parts[0] == "guilds" && parts[2] == "applications" && len(parts) >= 4 && parts[3] == "forms" {
		handleBotApplicationForms(w, r, parts[1])
		return
	}
	if len(parts) >= 3 && parts[0] == "guilds" && parts[2] == "settings" {
		handleBotGuildSettings(w, r, parts[1])
		return
	}
	if len(parts) >= 1 && parts[0] == "health" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	writeErr(w, http.StatusNotFound, "unknown bot route")
}

func botAPIKeyValid(r *http.Request) bool {
	if dashboardAPIKey == "" {
		// No key configured: accept any caller. Good for single-host local dev,
		// not safe for prod — flagged in startup logs.
		return true
	}
	got := r.Header.Get("Authorization")
	return got == "Bearer "+dashboardAPIKey
}

// GET /api/bot/guilds/:gid/applications/forms — consumed by Bot's
// internal/dashboard/client.go.
func handleBotApplicationForms(w http.ResponseWriter, r *http.Request, gid string) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "")
		return
	}
	type formOut struct {
		ID    string `json:"id"`
		Role  string `json:"role"`
		URL   string `json:"url"`
		Title string `json:"title"`
	}
	if db != nil && db.Pool != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		rows, err := db.Pool.Query(ctx, `
			SELECT role, url FROM application_forms WHERE guild_id = $1 ORDER BY role`, gid)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()
		out := []formOut{}
		for rows.Next() {
			var f formOut
			if err := rows.Scan(&f.Role, &f.URL); err != nil {
				writeErr(w, http.StatusInternalServerError, err.Error())
				return
			}
			f.ID = gid + ":" + f.Role
			f.Title = f.Role
			out = append(out, f)
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	// Memory fallback
	store.mu.RLock()
	forms := store.appForms[gid]
	store.mu.RUnlock()
	out := make([]formOut, 0, len(forms))
	for _, f := range forms {
		out = append(out, formOut{ID: gid + ":" + f.RoleID, Role: f.RoleID, URL: f.URL, Title: f.RoleID})
	}
	writeJSON(w, http.StatusOK, out)
}

// GET /api/bot/guilds/:gid/settings — a lightweight read-only view the Bot
// can use to respect per-guild prefix, log channel, modules, etc.
func handleBotGuildSettings(w http.ResponseWriter, r *http.Request, gid string) {
	if r.Method != http.MethodGet {
		writeErr(w, http.StatusMethodNotAllowed, "")
		return
	}
	if db != nil && db.Pool != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var prefix, language, timezone, welcomeCh, welcomeMsg, logCh string
		var autoRoles, modules, flags []byte
		var premium bool
		err := db.QueryRow(ctx, `
			SELECT prefix, language, timezone,
			       COALESCE(welcome_channel,''), COALESCE(welcome_message,''),
			       auto_role_ids, COALESCE(log_channel,''), modules, feature_flags, premium
			  FROM dash_guild_settings WHERE guild_id = $1`, gid,
		).Scan(&prefix, &language, &timezone, &welcomeCh, &welcomeMsg, &autoRoles, &logCh, &modules, &flags, &premium)
		if err != nil {
			writeErr(w, http.StatusNotFound, "no settings")
			return
		}
		out := map[string]any{
			"prefix":             prefix,
			"language":           language,
			"timezone":           timezone,
			"welcome_channel_id": welcomeCh,
			"welcome_message":    welcomeMsg,
			"log_channel_id":     logCh,
			"premium":            premium,
		}
		var v any
		if len(autoRoles) > 0 {
			_ = json.Unmarshal(autoRoles, &v)
			out["auto_role_ids"] = v
		}
		if len(modules) > 0 {
			_ = json.Unmarshal(modules, &v)
			out["module_enabled"] = v
		}
		if len(flags) > 0 {
			_ = json.Unmarshal(flags, &v)
			out["feature_flags"] = v
		}
		writeJSON(w, http.StatusOK, out)
		return
	}
	store.mu.RLock()
	s := store.settings[gid]
	store.mu.RUnlock()
	if s == nil {
		writeErr(w, http.StatusNotFound, "no settings")
		return
	}
	writeJSON(w, http.StatusOK, s)
}
