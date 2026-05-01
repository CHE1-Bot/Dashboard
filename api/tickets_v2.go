package main

import (
	"fmt"
	"html"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// handleTicketCategories — CRUD for TicketsBot v2-style categories.
//   GET    /tickets/categories         — list, ordered by position
//   POST   /tickets/categories         — create
//   PATCH  /tickets/categories/{id}    — update
//   DELETE /tickets/categories/{id}    — delete
func handleTicketCategories(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	switch {
	case r.Method == "GET" && len(rem) == 0:
		store.mu.RLock()
		defer store.mu.RUnlock()
		out := append([]TicketCategory(nil), store.categories[gid]...)
		sort.SliceStable(out, func(i, j int) bool { return out[i].Position < out[j].Position })
		writeJSON(w, http.StatusOK, out)
	case r.Method == "POST" && len(rem) == 0:
		var c TicketCategory
		if err := readJSON(r, &c); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.TrimSpace(c.Name) == "" {
			writeErr(w, http.StatusBadRequest, "name is required")
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		c.ID = store.nextID()
		c.GuildID = gid
		if c.NamingPattern == "" {
			c.NamingPattern = "ticket-{user}"
		}
		if c.MaxPerUser <= 0 {
			c.MaxPerUser = 1
		}
		if c.Color == "" {
			c.Color = "#5865f2"
		}
		c.Position = len(store.categories[gid])
		store.categories[gid] = append(store.categories[gid], c)
		recordHistory(gid, "tickets.category.create", c.Name)
		writeJSON(w, http.StatusOK, c)
	case r.Method == "PATCH" && len(rem) == 1:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		var patch TicketCategory
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		for i, c := range store.categories[gid] {
			if c.ID == id {
				patch.ID = c.ID
				patch.GuildID = gid
				store.categories[gid][i] = patch
				recordHistory(gid, "tickets.category.update", patch.Name)
				writeJSON(w, http.StatusOK, store.categories[gid][i])
				return
			}
		}
		writeErr(w, http.StatusNotFound, "category not found")
	case r.Method == "DELETE" && len(rem) == 1:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.categories[gid][:0]
		removed := ""
		for _, c := range store.categories[gid] {
			if c.ID == id {
				removed = c.Name
				continue
			}
			out = append(out, c)
		}
		store.categories[gid] = out
		if removed != "" {
			recordHistory(gid, "tickets.category.delete", removed)
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

// handleTicketSnippets — quick-reply CRUD.
func handleTicketSnippets(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	switch {
	case r.Method == "GET" && len(rem) == 0:
		store.mu.RLock()
		defer store.mu.RUnlock()
		out := append([]TicketSnippet(nil), store.snippets[gid]...)
		sort.SliceStable(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
		writeJSON(w, http.StatusOK, out)
	case r.Method == "POST" && len(rem) == 0:
		var s TicketSnippet
		if err := readJSON(r, &s); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		s.Name = strings.TrimSpace(s.Name)
		if s.Name == "" || s.Content == "" {
			writeErr(w, http.StatusBadRequest, "name and content required")
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		// Enforce uniqueness per guild
		for _, existing := range store.snippets[gid] {
			if strings.EqualFold(existing.Name, s.Name) {
				writeErr(w, http.StatusConflict, "snippet name already exists")
				return
			}
		}
		s.ID = store.nextID()
		s.GuildID = gid
		s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		store.snippets[gid] = append(store.snippets[gid], s)
		recordHistory(gid, "tickets.snippet.create", s.Name)
		writeJSON(w, http.StatusOK, s)
	case r.Method == "PATCH" && len(rem) == 1:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		var patch TicketSnippet
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		for i, s := range store.snippets[gid] {
			if s.ID == id {
				if patch.Name != "" {
					store.snippets[gid][i].Name = strings.TrimSpace(patch.Name)
				}
				if patch.Content != "" {
					store.snippets[gid][i].Content = patch.Content
				}
				store.snippets[gid][i].UpdatedAt = time.Now().UTC().Format(time.RFC3339)
				writeJSON(w, http.StatusOK, store.snippets[gid][i])
				return
			}
		}
		writeErr(w, http.StatusNotFound, "snippet not found")
	case r.Method == "DELETE" && len(rem) == 1:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.snippets[gid][:0]
		for _, s := range store.snippets[gid] {
			if s.ID != id {
				out = append(out, s)
			}
		}
		store.snippets[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

// handleTicketClaim assigns the current user to a ticket. TicketsBot v2 lets
// only one staffer claim a ticket at a time; subsequent claims fail unless the
// existing claimer first unclaims.
func handleTicketClaim(w http.ResponseWriter, r *http.Request, gid string, id int64) {
	u := requireUser(w, r)
	if u == nil {
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	for i, t := range store.tickets[gid] {
		if t.ID != id {
			continue
		}
		if t.Status != "open" {
			writeErr(w, http.StatusConflict, "ticket is not open")
			return
		}
		if t.ClaimedBy != "" && t.ClaimedBy != u.ID {
			writeErr(w, http.StatusConflict, "already claimed by "+t.ClaimedByName)
			return
		}
		store.tickets[gid][i].ClaimedBy = u.ID
		store.tickets[gid][i].ClaimedByName = u.Username
		store.tickets[gid][i].ClaimedAt = time.Now().UTC().Format(time.RFC3339)
		recordHistory(gid, "tickets.claim", fmt.Sprintf("Ticket #%d → %s", id, u.Username))
		forwardWorker("tickets.claim", gid, store.tickets[gid][i])
		writeJSON(w, http.StatusOK, store.tickets[gid][i])
		return
	}
	writeErr(w, http.StatusNotFound, "ticket not found")
}

func handleTicketUnclaim(w http.ResponseWriter, r *http.Request, gid string, id int64) {
	if requireUser(w, r) == nil {
		return
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	for i, t := range store.tickets[gid] {
		if t.ID != id {
			continue
		}
		store.tickets[gid][i].ClaimedBy = ""
		store.tickets[gid][i].ClaimedByName = ""
		store.tickets[gid][i].ClaimedAt = ""
		recordHistory(gid, "tickets.unclaim", fmt.Sprintf("Ticket #%d", id))
		forwardWorker("tickets.unclaim", gid, store.tickets[gid][i])
		writeJSON(w, http.StatusOK, store.tickets[gid][i])
		return
	}
	writeErr(w, http.StatusNotFound, "ticket not found")
}

// handleTicketTranscript returns an HTML transcript for a ticket. In a real
// deployment the Bot writes this on close; we return whatever's in
// dash_ticket_transcripts (Postgres) or generate a synthetic one in memory.
func handleTicketTranscript(w http.ResponseWriter, r *http.Request, gid string, id int64) {
	store.mu.RLock()
	var t *Ticket
	for i := range store.tickets[gid] {
		if store.tickets[gid][i].ID == id {
			t = &store.tickets[gid][i]
			break
		}
	}
	store.mu.RUnlock()
	if t == nil {
		writeErr(w, http.StatusNotFound, "ticket not found")
		return
	}
	htmlOut := renderTranscriptHTML(t)
	writeJSON(w, http.StatusOK, map[string]any{
		"id":      id,
		"ticket":  t,
		"html":    htmlOut,
		"content": stripHTML(htmlOut),
	})
}

func renderTranscriptHTML(t *Ticket) string {
	var sb strings.Builder
	sb.WriteString(`<!doctype html><html><head><meta charset="utf-8"><title>Ticket #`)
	sb.WriteString(strconv.FormatInt(t.ID, 10))
	sb.WriteString(`</title>`)
	sb.WriteString(`<style>
		body{font:14px/1.5 -apple-system,Segoe UI,Arial,sans-serif;background:#36393f;color:#dcddde;margin:0;padding:24px;}
		.header{background:#2f3136;border-radius:8px;padding:16px 20px;margin-bottom:18px;}
		.header h1{margin:0 0 6px 0;font-size:18px;color:#fff;}
		.meta{font-size:12px;color:#b9bbbe;}
		.meta span{margin-right:12px;}
		.message{display:flex;gap:14px;padding:10px 6px;border-bottom:1px solid #40444b;}
		.avatar{width:40px;height:40px;border-radius:50%;background:#5865f2;display:grid;place-items:center;color:#fff;font-weight:700;}
		.body .author{color:#fff;font-weight:600;margin-right:8px;}
		.body .ts{color:#72767d;font-size:11px;}
		.body .text{margin-top:4px;color:#dcddde;}
	</style></head><body>`)
	sb.WriteString(`<div class="header"><h1>` + html.EscapeString(t.Subject) + `</h1>`)
	sb.WriteString(`<div class="meta">`)
	sb.WriteString(`<span>Ticket #` + strconv.FormatInt(t.ID, 10) + `</span>`)
	sb.WriteString(`<span>Opened by ` + html.EscapeString(t.Username) + `</span>`)
	if t.CategoryName != "" {
		sb.WriteString(`<span>Category: ` + html.EscapeString(t.CategoryName) + `</span>`)
	}
	if t.ClaimedByName != "" {
		sb.WriteString(`<span>Claimed by ` + html.EscapeString(t.ClaimedByName) + `</span>`)
	}
	sb.WriteString(`<span>Status: ` + html.EscapeString(t.Status) + `</span>`)
	sb.WriteString(`</div></div>`)

	// Synthetic message preview — real bot writes the full transcript on close.
	sb.WriteString(`<div class="message"><div class="avatar">` + initials(t.Username) + `</div>`)
	sb.WriteString(`<div class="body"><span class="author">` + html.EscapeString(t.Username) + `</span>`)
	sb.WriteString(`<span class="ts">` + html.EscapeString(t.CreatedAt) + `</span>`)
	sb.WriteString(`<div class="text">` + html.EscapeString(t.Subject) + `</div></div></div>`)
	if t.ClaimedByName != "" {
		sb.WriteString(`<div class="message"><div class="avatar" style="background:#10b981">` + initials(t.ClaimedByName) + `</div>`)
		sb.WriteString(`<div class="body"><span class="author">` + html.EscapeString(t.ClaimedByName) + `</span>`)
		sb.WriteString(`<span class="ts">` + html.EscapeString(t.ClaimedAt) + `</span>`)
		sb.WriteString(`<div class="text">Hi! I've claimed this ticket and will help you out.</div></div></div>`)
	}
	sb.WriteString(`</body></html>`)
	return sb.String()
}

func initials(name string) string {
	if name == "" {
		return "?"
	}
	if len(name) >= 2 {
		return strings.ToUpper(name[:2])
	}
	return strings.ToUpper(name[:1])
}

func stripHTML(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		switch r {
		case '<':
			in = true
		case '>':
			in = false
		default:
			if !in {
				b.WriteRune(r)
			}
		}
	}
	return strings.TrimSpace(b.String())
}

// handleTicketStats returns aggregate metrics that the Statistics page binds to.
//
// Shape:
//   { totals: {open, closed, claimed, avg_handle_hours},
//     by_category: [{name, open, closed}],
//     by_staff: [{user_id, username, claimed, closed, avg_response_min}],
//     opens_per_day: [{date, value}] }
func handleTicketStats(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	tickets := append([]Ticket(nil), store.tickets[gid]...)
	categories := append([]TicketCategory(nil), store.categories[gid]...)
	store.mu.RUnlock()

	totals := map[string]any{}
	open, closed, claimed := 0, 0, 0
	var totalHandleSec int64
	var handleN int
	for _, t := range tickets {
		switch t.Status {
		case "open":
			open++
		case "closed":
			closed++
		}
		if t.ClaimedBy != "" {
			claimed++
		}
		if t.ClosedAt != "" && t.CreatedAt != "" {
			c, e1 := time.Parse(time.RFC3339, t.CreatedAt)
			x, e2 := time.Parse(time.RFC3339, t.ClosedAt)
			if e1 == nil && e2 == nil && x.After(c) {
				totalHandleSec += int64(x.Sub(c).Seconds())
				handleN++
			}
		}
	}
	avgHandleHours := 0.0
	if handleN > 0 {
		avgHandleHours = (float64(totalHandleSec) / float64(handleN)) / 3600.0
	}
	totals["open"] = open
	totals["closed"] = closed
	totals["claimed"] = claimed
	totals["avg_handle_hours"] = avgHandleHours

	// By category
	type catBucket struct {
		Name   string `json:"name"`
		Open   int    `json:"open"`
		Closed int    `json:"closed"`
	}
	catMap := map[int64]*catBucket{}
	for _, c := range categories {
		catMap[c.ID] = &catBucket{Name: c.Name}
	}
	uncat := &catBucket{Name: "Uncategorized"}
	for _, t := range tickets {
		bucket := uncat
		if b, ok := catMap[t.CategoryID]; ok {
			bucket = b
		}
		if t.Status == "open" {
			bucket.Open++
		} else if t.Status == "closed" {
			bucket.Closed++
		}
	}
	byCategory := []catBucket{}
	for _, c := range categories {
		if b := catMap[c.ID]; b != nil {
			byCategory = append(byCategory, *b)
		}
	}
	if uncat.Open > 0 || uncat.Closed > 0 {
		byCategory = append(byCategory, *uncat)
	}

	// By staff
	type staffBucket struct {
		UserID         string  `json:"user_id"`
		Username       string  `json:"username"`
		Claimed        int     `json:"claimed"`
		Closed         int     `json:"closed"`
		AvgResponseMin float64 `json:"avg_response_min"`
	}
	staffMap := map[string]*staffBucket{}
	for _, t := range tickets {
		if t.ClaimedBy == "" {
			continue
		}
		b := staffMap[t.ClaimedBy]
		if b == nil {
			b = &staffBucket{UserID: t.ClaimedBy, Username: t.ClaimedByName}
			staffMap[t.ClaimedBy] = b
		}
		b.Claimed++
		if t.Status == "closed" {
			b.Closed++
		}
		if t.ClaimedAt != "" && t.CreatedAt != "" {
			c, e1 := time.Parse(time.RFC3339, t.CreatedAt)
			x, e2 := time.Parse(time.RFC3339, t.ClaimedAt)
			if e1 == nil && e2 == nil && x.After(c) {
				diff := x.Sub(c).Minutes()
				if b.AvgResponseMin == 0 {
					b.AvgResponseMin = diff
				} else {
					b.AvgResponseMin = (b.AvgResponseMin + diff) / 2
				}
			}
		}
	}
	byStaff := make([]staffBucket, 0, len(staffMap))
	for _, b := range staffMap {
		byStaff = append(byStaff, *b)
	}
	sort.Slice(byStaff, func(i, j int) bool { return byStaff[i].Claimed > byStaff[j].Claimed })

	// Opens per day (last 7 days)
	type pt struct {
		Date  string `json:"date"`
		Value int    `json:"value"`
	}
	opens := make([]pt, 0, 7)
	now := time.Now().UTC()
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		count := 0
		for _, t := range tickets {
			if strings.HasPrefix(t.CreatedAt, day) {
				count++
			}
		}
		opens = append(opens, pt{Date: day, Value: count})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"totals":        totals,
		"by_category":   byCategory,
		"by_staff":      byStaff,
		"opens_per_day": opens,
	})
}
