package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Task kinds accepted by the Bot's actions.Handlers (internal/actions).
// The Worker receives these via POST /api/v1/tasks and broadcasts task.created
// events on the WS hub; the Bot executes the Discord side-effect.
const (
	KindSendMessage          = "send_message"
	KindSendApplicationPanel = "send_application_panel"
	KindSendTicketPanel      = "send_ticket_panel"
	KindSendGiveawayPanel    = "send_giveaway_panel"

	// Soft kinds used by the dashboard to record an intent even if there is no
	// Worker wired up. The bot ignores unknown kinds.
	KindGiveawayEnd    = "giveaway.end"
	KindGiveawayReroll = "giveaway.reroll"
	KindModKick        = "moderation.kick"
	KindModBan         = "moderation.ban"
	KindModUnban       = "moderation.unban"
	KindModMute        = "moderation.mute"
	KindModUnmute      = "moderation.unmute"
	KindModWarn        = "moderation.warn"
)

var workerClient = &http.Client{Timeout: 10 * time.Second}

// forwardWorker posts an action to the Worker if WORKER_URL is set. Non-blocking:
// the network call runs in a goroutine so handlers return promptly. When a
// worker is not configured we only record the intent in dash_history.
//
// Payload wire shape matches the Worker's REST contract: the Worker accepts
// either {kind,input,created_by} or {event,guild_id,payload} — we send the
// dashboard shape and the Worker normalizes it to {kind, input{guild_id,payload}}.
func forwardWorker(kind, gid string, payload any) {
	// Always record the intent locally so the UI shows the action even without
	// a Worker running.
	recordHistoryAsync(gid, "worker."+kind, summarizePayload(payload))

	if workerURL == "" {
		return
	}
	body, err := json.Marshal(map[string]any{
		"event":    kind,
		"guild_id": gid,
		"payload":  payload,
	})
	if err != nil {
		return
	}

	go func(b []byte) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		url := strings.TrimRight(workerURL, "/") + "/api/v1/tasks"
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if workerKey != "" {
			req.Header.Set("Authorization", "Bearer "+workerKey)
		}
		resp, err := workerClient.Do(req)
		if err != nil {
			logger.Warn("worker forward failed", "kind", kind, "guild", gid, "err", err.Error())
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			logger.Warn("worker forward rejected", "kind", kind, "guild", gid, "status", resp.StatusCode, "body", string(snippet))
		}
	}(body)
}

// summarizePayload extracts a readable string for the history log.
func summarizePayload(p any) string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	if len(b) > 200 {
		return string(b[:200]) + "…"
	}
	return string(b)
}

// recordHistoryAsync enqueues a dash_history row without holding any lock on
// the in-memory Store. Safe to call from request handlers.
func recordHistoryAsync(gid, event, detail string) {
	if db != nil && db.Pool != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_, _ = db.Pool.Exec(ctx, `
				INSERT INTO dash_history (guild_id, actor, event, detail)
				VALUES ($1, 'dashboard', $2, $3)`, gid, event, detail)
		}()
		return
	}
	store.mu.Lock()
	store.history[gid] = append([]HistoryEvent{{
		ID: store.nextID(), GuildID: gid, Actor: "dashboard",
		Event: event, Detail: detail,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}}, store.history[gid]...)
	store.mu.Unlock()
}

// taskPayload is a helper for consumers that want to build input in a typed
// way. Most callers just pass map[string]any.
func taskPayload(fields ...any) map[string]any {
	m := map[string]any{}
	for i := 0; i+1 < len(fields); i += 2 {
		key, _ := fields[i].(string)
		if key != "" {
			m[key] = fields[i+1]
		}
	}
	return m
}

// asJSONError writes an error log line when the worker rejects a task.
func asJSONError(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf(`{"error":%q}`, err.Error())
}
