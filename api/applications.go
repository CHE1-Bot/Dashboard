package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// handleApplications dispatches all routes under
//   /api/guilds/{gid}/applications/...
//
// Routes:
//   GET    /                         — list submissions
//   GET    /forms                    — list forms
//   POST   /forms                    — create form
//   GET    /forms/{id}               — fetch one form
//   PATCH  /forms/{id}               — update form
//   DELETE /forms/{id}               — delete form
//   GET    /stats                    — review queue + outcome aggregates
//   PATCH  /{id}                     — update submission (legacy: status only)
//   POST   /{id}/accept              — accept (DM template + role grant)
//   POST   /{id}/reject              — reject (DM template + reason)
func handleApplications(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		if r.Method == "GET" {
			store.mu.RLock()
			out := append([]Application(nil), store.apps[gid]...)
			store.mu.RUnlock()
			// Decorate with form name if missing
			if len(out) > 0 {
				store.mu.RLock()
				formByID := map[int64]string{}
				for _, f := range store.appForms[gid] {
					formByID[f.ID] = f.Name
				}
				store.mu.RUnlock()
				for i := range out {
					if out[i].FormName == "" {
						out[i].FormName = formByID[out[i].FormID]
					}
				}
			}
			sort.SliceStable(out, func(i, j int) bool { return out[i].CreatedAt > out[j].CreatedAt })
			writeJSON(w, http.StatusOK, out)
			return
		}
		if r.Method == "POST" {
			// Accept a brand-new submission (e.g. from the Bot when a user
			// finishes a modal flow). The dashboard normally just reviews.
			var a Application
			if err := readJSON(r, &a); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			a.ID = store.nextID()
			a.GuildID = gid
			if a.Status == "" {
				a.Status = "pending"
			}
			a.CreatedAt = time.Now().UTC().Format(time.RFC3339)
			store.apps[gid] = append([]Application{a}, store.apps[gid]...)
			recordHistoryAsync(gid, "applications.submit", a.Username+" → "+a.FormName)
			writeJSON(w, http.StatusOK, a)
			return
		}
	}
	switch rem[0] {
	case "forms":
		handleApplicationForms(w, r, gid, rem[1:])
		return
	case "stats":
		handleApplicationStats(w, r, gid)
		return
	default:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		if len(rem) == 2 && rem[1] == "accept" && r.Method == "POST" {
			handleApplicationDecision(w, r, gid, id, "accepted")
			return
		}
		if len(rem) == 2 && rem[1] == "reject" && r.Method == "POST" {
			handleApplicationDecision(w, r, gid, id, "rejected")
			return
		}
		if r.Method == "PATCH" {
			var patch Application
			if err := readJSON(r, &patch); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			for i, a := range store.apps[gid] {
				if a.ID == id {
					if patch.Status != "" {
						store.apps[gid][i].Status = patch.Status
					}
					if patch.ReviewNote != "" {
						store.apps[gid][i].ReviewNote = patch.ReviewNote
					}
					writeJSON(w, http.StatusOK, store.apps[gid][i])
					return
				}
			}
			writeErr(w, http.StatusNotFound, "application not found")
			return
		}
		if r.Method == "DELETE" {
			store.mu.Lock()
			defer store.mu.Unlock()
			out := store.apps[gid][:0]
			for _, a := range store.apps[gid] {
				if a.ID != id {
					out = append(out, a)
				}
			}
			store.apps[gid] = out
			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
			return
		}
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleApplicationForms(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	switch {
	case r.Method == "GET" && len(rem) == 0:
		store.mu.RLock()
		defer store.mu.RUnlock()
		out := append([]ApplicationForm(nil), store.appForms[gid]...)
		sort.SliceStable(out, func(i, j int) bool { return out[i].Name < out[j].Name })
		writeJSON(w, http.StatusOK, out)
	case r.Method == "POST" && len(rem) == 0:
		var f ApplicationForm
		if err := readJSON(r, &f); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.TrimSpace(f.Name) == "" {
			writeErr(w, http.StatusBadRequest, "name is required")
			return
		}
		f.ID = 0
		f.GuildID = gid
		applyFormDefaults(&f)
		store.mu.Lock()
		defer store.mu.Unlock()
		f.ID = store.nextID()
		f.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		store.appForms[gid] = append(store.appForms[gid], f)
		recordHistoryAsync(gid, "applications.form.create", f.Name)
		writeJSON(w, http.StatusOK, f)
	case r.Method == "GET" && len(rem) == 1:
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.RLock()
		defer store.mu.RUnlock()
		for _, f := range store.appForms[gid] {
			if f.ID == id {
				writeJSON(w, http.StatusOK, f)
				return
			}
		}
		writeErr(w, http.StatusNotFound, "form not found")
	case r.Method == "PATCH" && len(rem) == 1:
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		var patch ApplicationForm
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		applyFormDefaults(&patch)
		store.mu.Lock()
		defer store.mu.Unlock()
		for i, f := range store.appForms[gid] {
			if f.ID == id {
				patch.ID = f.ID
				patch.GuildID = gid
				patch.CreatedAt = f.CreatedAt
				store.appForms[gid][i] = patch
				recordHistoryAsync(gid, "applications.form.update", patch.Name)
				writeJSON(w, http.StatusOK, store.appForms[gid][i])
				return
			}
		}
		writeErr(w, http.StatusNotFound, "form not found")
	case r.Method == "DELETE" && len(rem) == 1:
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		removed := ""
		out := store.appForms[gid][:0]
		for _, f := range store.appForms[gid] {
			if f.ID == id {
				removed = f.Name
				continue
			}
			out = append(out, f)
		}
		store.appForms[gid] = out
		if removed != "" {
			recordHistoryAsync(gid, "applications.form.delete", removed)
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

func applyFormDefaults(f *ApplicationForm) {
	if f.Color == "" {
		f.Color = "#5865f2"
	}
	if f.AcceptDMTemplate == "" {
		f.AcceptDMTemplate = "Congrats {user}! Your application for **{form}** in **{guild}** has been accepted."
	}
	if f.RejectDMTemplate == "" {
		f.RejectDMTemplate = "Thanks for applying to **{form}** in **{guild}**, {user}. Unfortunately your application has been declined.\n\nReason: {reason}"
	}
	for i := range f.Questions {
		if f.Questions[i].ID == "" {
			f.Questions[i].ID = fmt.Sprintf("q%d", i+1)
		}
		if f.Questions[i].Type == "" {
			f.Questions[i].Type = "short"
		}
	}
}

// handleApplicationDecision applies an accept/reject. Records the reviewer,
// sets ReviewedAt, copies an optional reason from the request body, and
// forwards to the Worker so the Bot can DM the applicant + grant the role.
func handleApplicationDecision(w http.ResponseWriter, r *http.Request, gid string, id int64, decision string) {
	u := requireUser(w, r)
	if u == nil {
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	_ = readJSON(r, &body)

	store.mu.Lock()
	defer store.mu.Unlock()
	for i, a := range store.apps[gid] {
		if a.ID != id {
			continue
		}
		if a.Status != "pending" {
			writeErr(w, http.StatusConflict, "application already "+a.Status)
			return
		}
		// Find the form so we can attach DM template + accepted role.
		var form *ApplicationForm
		for j := range store.appForms[gid] {
			if store.appForms[gid][j].ID == a.FormID {
				form = &store.appForms[gid][j]
				break
			}
		}
		store.apps[gid][i].Status = decision
		store.apps[gid][i].ReviewedBy = u.ID
		store.apps[gid][i].ReviewedByName = u.Username
		store.apps[gid][i].ReviewedAt = time.Now().UTC().Format(time.RFC3339)
		if body.Reason != "" {
			store.apps[gid][i].ReviewNote = body.Reason
		}

		updated := store.apps[gid][i]
		event := "applications.accept"
		if decision == "rejected" {
			event = "applications.reject"
		}
		recordHistoryAsync(gid, event,
			updated.Username+" — "+updated.FormName+
				" by "+u.Username+
				condStr(body.Reason != "", " ("+body.Reason+")", ""))

		// Build the worker payload — includes the rendered DM template so
		// the Bot just sends what we hand it.
		dm := ""
		grantRole := ""
		if form != nil {
			tpl := form.AcceptDMTemplate
			if decision == "rejected" {
				tpl = form.RejectDMTemplate
			}
			dm = renderTemplate(tpl, map[string]string{
				"user":   updated.Username,
				"form":   updated.FormName,
				"guild":  guildName(gid),
				"reason": body.Reason,
			})
			if decision == "accepted" {
				grantRole = form.AcceptedRoleID
			}
		}
		forwardWorker("applications."+decision, gid, map[string]any{
			"application_id": id,
			"user_id":        updated.UserID,
			"form_id":        updated.FormID,
			"reviewer":       u.Username,
			"reason":         body.Reason,
			"dm":             dm,
			"grant_role_id":  grantRole,
		})

		writeJSON(w, http.StatusOK, updated)
		return
	}
	writeErr(w, http.StatusNotFound, "application not found")
}

// handleApplicationStats — counts and per-form aggregates for the stats page.
func handleApplicationStats(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	apps := append([]Application(nil), store.apps[gid]...)
	forms := append([]ApplicationForm(nil), store.appForms[gid]...)
	store.mu.RUnlock()

	pending, accepted, rejected := 0, 0, 0
	var sumReviewSec int64
	var reviewedN int
	for _, a := range apps {
		switch a.Status {
		case "pending":
			pending++
		case "accepted":
			accepted++
		case "rejected":
			rejected++
		}
		if a.ReviewedAt != "" && a.CreatedAt != "" {
			c, e1 := time.Parse(time.RFC3339, a.CreatedAt)
			r, e2 := time.Parse(time.RFC3339, a.ReviewedAt)
			if e1 == nil && e2 == nil && r.After(c) {
				sumReviewSec += int64(r.Sub(c).Seconds())
				reviewedN++
			}
		}
	}
	avgHours := 0.0
	if reviewedN > 0 {
		avgHours = (float64(sumReviewSec) / float64(reviewedN)) / 3600
	}

	type bucket struct {
		FormID   int64  `json:"form_id"`
		Name     string `json:"name"`
		Pending  int    `json:"pending"`
		Accepted int    `json:"accepted"`
		Rejected int    `json:"rejected"`
	}
	buckets := map[int64]*bucket{}
	for _, f := range forms {
		buckets[f.ID] = &bucket{FormID: f.ID, Name: f.Name}
	}
	for _, a := range apps {
		b, ok := buckets[a.FormID]
		if !ok {
			b = &bucket{FormID: a.FormID, Name: a.FormName}
			buckets[a.FormID] = b
		}
		switch a.Status {
		case "pending":
			b.Pending++
		case "accepted":
			b.Accepted++
		case "rejected":
			b.Rejected++
		}
	}
	byForm := make([]bucket, 0, len(buckets))
	for _, b := range buckets {
		byForm = append(byForm, *b)
	}
	sort.Slice(byForm, func(i, j int) bool { return byForm[i].Name < byForm[j].Name })

	type pt struct {
		Date  string `json:"date"`
		Value int    `json:"value"`
	}
	now := time.Now().UTC()
	subs := make([]pt, 0, 7)
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		c := 0
		for _, a := range apps {
			if strings.HasPrefix(a.CreatedAt, day) {
				c++
			}
		}
		subs = append(subs, pt{Date: day, Value: c})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"totals": map[string]any{
			"pending":          pending,
			"accepted":         accepted,
			"rejected":         rejected,
			"avg_review_hours": avgHours,
		},
		"by_form":            byForm,
		"submissions_per_day": subs,
	})
}

func renderTemplate(tpl string, vars map[string]string) string {
	out := tpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	return out
}

func guildName(gid string) string {
	store.mu.RLock()
	defer store.mu.RUnlock()
	if g := store.guilds[gid]; g != nil {
		return g.Name
	}
	return "the server"
}

func condStr(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
