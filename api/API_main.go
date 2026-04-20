package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ---------------- Models (mirror Bot schema) ----------------

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email,omitempty"`
}

type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int64  `json:"permissions"`
	BotPresent  bool   `json:"bot_present"`
	MemberCount int    `json:"member_count"`
}

type Role struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    int    `json:"color"`
	Position int    `json:"position"`
	Managed  bool   `json:"managed"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // text, voice, category, announcement
}

type Member struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname,omitempty"`
	Avatar   string   `json:"avatar,omitempty"`
	Roles    []string `json:"roles"`
	JoinedAt string   `json:"joined_at"`
}

type GuildSettings struct {
	Prefix            string            `json:"prefix"`
	Language          string            `json:"language"`
	Timezone          string            `json:"timezone"`
	WelcomeChannelID  string            `json:"welcome_channel_id"`
	WelcomeMessage    string            `json:"welcome_message"`
	AutoRoleIDs       []string          `json:"auto_role_ids"`
	LogChannelID      string            `json:"log_channel_id"`
	ModuleEnabled     map[string]bool   `json:"module_enabled"`
	FeatureFlags      map[string]string `json:"feature_flags"`
	Premium           bool              `json:"premium"`
}

type Ticket struct {
	ID            int64  `json:"id"`
	GuildID       string `json:"guild_id"`
	ChannelID     string `json:"channel_id"`
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	Subject       string `json:"subject"`
	Status        string `json:"status"` // open | closed
	TranscriptURL string `json:"transcript_url,omitempty"`
	AssignedTo    string `json:"assigned_to,omitempty"`
	Tags          []string `json:"tags"`
	CreatedAt     string `json:"created_at"`
	ClosedAt      string `json:"closed_at,omitempty"`
}

type TicketPanel struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Color     string `json:"color"`
	Buttons   []TicketPanelButton `json:"buttons"`
}

type TicketPanelButton struct {
	Label    string `json:"label"`
	Style    string `json:"style"`
	FormID   int64  `json:"form_id,omitempty"`
	Category string `json:"category"`
}

type TicketForm struct {
	ID      int64            `json:"id"`
	GuildID string           `json:"guild_id"`
	Name    string           `json:"name"`
	Fields  []TicketFormField `json:"fields"`
}

type TicketFormField struct {
	Label       string `json:"label"`
	Type        string `json:"type"` // short | paragraph
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
}

type TicketTag struct {
	ID      int64  `json:"id"`
	GuildID string `json:"guild_id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

type TicketEmbed struct {
	GuildID     string `json:"guild_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Footer      string `json:"footer"`
	Thumbnail   string `json:"thumbnail"`
}

type TicketStaff struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type TicketSettings struct {
	CategoryID      string   `json:"category_id"`
	SupportRoleIDs  []string `json:"support_role_ids"`
	TranscriptsOn   bool     `json:"transcripts_on"`
	CloseConfirm    bool     `json:"close_confirm"`
	MaxOpenPerUser  int      `json:"max_open_per_user"`
	NamingPattern   string   `json:"naming_pattern"`
}

type UserLevel struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	Level     int    `json:"level"`
	XP        int    `json:"xp"`
	TotalXP   int    `json:"total_xp"`
	UpdatedAt string `json:"updated_at"`
}

type LevelingSettings struct {
	Enabled         bool     `json:"enabled"`
	XPPerMessage    int      `json:"xp_per_message"`
	CooldownSec     int      `json:"cooldown_sec"`
	AnnounceChannel string   `json:"announce_channel"`
	AnnounceMessage string   `json:"announce_message"`
	NoXPChannels    []string `json:"no_xp_channels"`
	NoXPRoles       []string `json:"no_xp_roles"`
	Multiplier      float64  `json:"multiplier"`
}

type LevelReward struct {
	ID      int64  `json:"id"`
	GuildID string `json:"guild_id"`
	Level   int    `json:"level"`
	RoleID  string `json:"role_id"`
	RoleName string `json:"role_name"`
	Stackable bool `json:"stackable"`
}

type ModLog struct {
	ID          int64  `json:"id"`
	GuildID     string `json:"guild_id"`
	ModeratorID string `json:"moderator_id"`
	ModName     string `json:"moderator_name"`
	TargetID    string `json:"target_id"`
	TargetName  string `json:"target_name"`
	Action      string `json:"action"` // kick | ban | mute | warn | unban | unmute
	Reason      string `json:"reason"`
	DurationSec int    `json:"duration_sec"`
	CreatedAt   string `json:"created_at"`
}

type AutoModRule struct {
	ID       int64  `json:"id"`
	GuildID  string `json:"guild_id"`
	Name     string `json:"name"`
	Trigger  string `json:"trigger"` // spam | caps | links | invites | words | mention
	Action   string `json:"action"`  // delete | warn | mute | kick | ban
	Enabled  bool   `json:"enabled"`
	Config   map[string]any `json:"config"`
}

type Report struct {
	ID          int64  `json:"id"`
	GuildID     string `json:"guild_id"`
	ReporterID  string `json:"reporter_id"`
	ReporterName string `json:"reporter_name"`
	TargetID    string `json:"target_id"`
	TargetName  string `json:"target_name"`
	Reason      string `json:"reason"`
	Status      string `json:"status"` // open | resolved | dismissed
	CreatedAt   string `json:"created_at"`
}

type Application struct {
	ID        int64           `json:"id"`
	GuildID   string          `json:"guild_id"`
	UserID    string          `json:"user_id"`
	Username  string          `json:"username"`
	RoleID    string          `json:"role_id"`
	RoleName  string          `json:"role_name"`
	Answers   map[string]any  `json:"answers"`
	Status    string          `json:"status"` // pending | accepted | rejected
	CreatedAt string          `json:"created_at"`
}

type ApplicationForm struct {
	GuildID  string `json:"guild_id"`
	RoleID   string `json:"role_id"`
	URL      string `json:"url"`
}

type Giveaway struct {
	ID        int64    `json:"id"`
	GuildID   string   `json:"guild_id"`
	ChannelID string   `json:"channel_id"`
	MessageID string   `json:"message_id"`
	Prize     string   `json:"prize"`
	Winners   []string `json:"winners"`
	WinnerCount int    `json:"winner_count"`
	EndsAt    string   `json:"ends_at"`
	Status    string   `json:"status"` // running | ended
	Entrants  int      `json:"entrants"`
	HostedBy  string   `json:"hosted_by"`
	CreatedAt string   `json:"created_at"`
}

type BlacklistEntry struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

type Alert struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	Severity  string `json:"severity"` // info | warn | error
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
	Dismissed bool   `json:"dismissed"`
}

type HistoryEvent struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	Actor     string `json:"actor"`
	Event     string `json:"event"`
	Detail    string `json:"detail"`
	CreatedAt string `json:"created_at"`
}

type Backup struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	Label     string `json:"label"`
	SizeBytes int    `json:"size_bytes"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type Permission struct {
	ID        int64    `json:"id"`
	GuildID   string   `json:"guild_id"`
	RoleID    string   `json:"role_id"`
	RoleName  string   `json:"role_name"`
	Modules   []string `json:"modules"` // which dashboard modules this role can access
}

type Session struct {
	ID          string
	UserID      string
	AccessToken string
	Expires     time.Time
}

type Stats struct {
	Members          int            `json:"members"`
	OnlineMembers    int            `json:"online_members"`
	MessagesToday    int            `json:"messages_today"`
	CommandsToday    int            `json:"commands_today"`
	OpenTickets      int            `json:"open_tickets"`
	ActiveGiveaways  int            `json:"active_giveaways"`
	ModActionsWeek   int            `json:"mod_actions_week"`
	MessagesPerDay   []TimeSeriesPt `json:"messages_per_day"`
	JoinsPerDay      []TimeSeriesPt `json:"joins_per_day"`
}

type TimeSeriesPt struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

// ---------------- Store ----------------

type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	users    map[string]*User

	guilds     map[string]*Guild              // guild id -> guild
	userGuilds map[string][]string            // user id -> []guild id
	botGuilds  map[string]bool                // guild id -> bot present
	roles      map[string][]Role
	channels   map[string][]Channel
	members    map[string][]Member

	settings    map[string]*GuildSettings
	tickets     map[string][]Ticket
	ticketIDSeq int64
	panels      map[string][]TicketPanel
	forms       map[string][]TicketForm
	tags        map[string][]TicketTag
	embeds      map[string]*TicketEmbed
	staff       map[string][]TicketStaff
	ticketSet   map[string]*TicketSettings

	userLevels  map[string][]UserLevel
	levelingSet map[string]*LevelingSettings
	rewards     map[string][]LevelReward

	modLogs    map[string][]ModLog
	modLogSeq  int64
	autoMod    map[string][]AutoModRule
	reports    map[string][]Report

	apps      map[string][]Application
	appForms  map[string][]ApplicationForm

	giveaways    map[string][]Giveaway
	giveawaySeq  int64
	blacklist    map[string][]BlacklistEntry

	alerts    map[string][]Alert
	history   map[string][]HistoryEvent
	backups   map[string][]Backup
	perms     map[string][]Permission
	seq       int64
}

func newStore() *Store {
	return &Store{
		sessions: map[string]*Session{},
		users:    map[string]*User{},

		guilds:     map[string]*Guild{},
		userGuilds: map[string][]string{},
		botGuilds:  map[string]bool{},
		roles:      map[string][]Role{},
		channels:   map[string][]Channel{},
		members:    map[string][]Member{},

		settings:  map[string]*GuildSettings{},
		tickets:   map[string][]Ticket{},
		panels:    map[string][]TicketPanel{},
		forms:     map[string][]TicketForm{},
		tags:      map[string][]TicketTag{},
		embeds:    map[string]*TicketEmbed{},
		staff:     map[string][]TicketStaff{},
		ticketSet: map[string]*TicketSettings{},

		userLevels:  map[string][]UserLevel{},
		levelingSet: map[string]*LevelingSettings{},
		rewards:     map[string][]LevelReward{},

		modLogs: map[string][]ModLog{},
		autoMod: map[string][]AutoModRule{},
		reports: map[string][]Report{},

		apps:     map[string][]Application{},
		appForms: map[string][]ApplicationForm{},

		giveaways: map[string][]Giveaway{},
		blacklist: map[string][]BlacklistEntry{},

		alerts:  map[string][]Alert{},
		history: map[string][]HistoryEvent{},
		backups: map[string][]Backup{},
		perms:   map[string][]Permission{},
	}
}

func (s *Store) nextID() int64 {
	s.seq++
	return s.seq
}

// seed creates demo data for dev mode so the UI isn't empty
func (s *Store) seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	devUser := &User{ID: "100000000000000001", Username: "DevUser", Discriminator: "0001", Avatar: ""}
	s.users[devUser.ID] = devUser

	gA := &Guild{ID: "200000000000000001", Name: "Che1 Test Server", Icon: "", Owner: true, Permissions: 8, BotPresent: true, MemberCount: 248}
	gB := &Guild{ID: "200000000000000002", Name: "Demo Community", Icon: "", Owner: false, Permissions: 8, BotPresent: true, MemberCount: 1820}
	gC := &Guild{ID: "200000000000000003", Name: "Friends & Co.", Icon: "", Owner: false, Permissions: 8, BotPresent: false, MemberCount: 42}
	s.guilds[gA.ID] = gA
	s.guilds[gB.ID] = gB
	s.guilds[gC.ID] = gC
	s.userGuilds[devUser.ID] = []string{gA.ID, gB.ID, gC.ID}
	s.botGuilds[gA.ID] = true
	s.botGuilds[gB.ID] = true

	for _, gid := range []string{gA.ID, gB.ID} {
		s.seedGuild(gid)
	}
}

func (s *Store) seedGuild(gid string) {
	s.roles[gid] = []Role{
		{ID: gid + "-r1", Name: "@everyone", Color: 0, Position: 0},
		{ID: gid + "-r2", Name: "Member", Color: 0x3498db, Position: 1},
		{ID: gid + "-r3", Name: "Moderator", Color: 0xe67e22, Position: 2},
		{ID: gid + "-r4", Name: "Admin", Color: 0xe74c3c, Position: 3},
		{ID: gid + "-r5", Name: "Level 10", Color: 0x2ecc71, Position: 4},
	}
	s.channels[gid] = []Channel{
		{ID: gid + "-c1", Name: "general", Type: "text"},
		{ID: gid + "-c2", Name: "announcements", Type: "announcement"},
		{ID: gid + "-c3", Name: "mod-logs", Type: "text"},
		{ID: gid + "-c4", Name: "tickets", Type: "category"},
		{ID: gid + "-c5", Name: "giveaways", Type: "text"},
		{ID: gid + "-c6", Name: "Lounge", Type: "voice"},
	}
	s.members[gid] = []Member{
		{ID: "u1", Username: "alice", Nickname: "Alice", Roles: []string{gid + "-r2", gid + "-r3"}, JoinedAt: isoDaysAgo(120)},
		{ID: "u2", Username: "bob", Roles: []string{gid + "-r2"}, JoinedAt: isoDaysAgo(90)},
		{ID: "u3", Username: "carol", Roles: []string{gid + "-r2", gid + "-r4"}, JoinedAt: isoDaysAgo(300)},
		{ID: "u4", Username: "dave", Roles: []string{gid + "-r2"}, JoinedAt: isoDaysAgo(40)},
	}
	s.settings[gid] = &GuildSettings{
		Prefix: "!", Language: "en", Timezone: "UTC",
		WelcomeChannelID: gid + "-c1",
		WelcomeMessage:   "Welcome {user} to {server}!",
		AutoRoleIDs:      []string{gid + "-r2"},
		LogChannelID:     gid + "-c3",
		ModuleEnabled: map[string]bool{
			"tickets": true, "moderation": true, "giveaways": true, "leveling": true, "applications": true,
		},
		FeatureFlags: map[string]string{},
		Premium:      false,
	}
	s.ticketSet[gid] = &TicketSettings{
		CategoryID: gid + "-c4", SupportRoleIDs: []string{gid + "-r3"},
		TranscriptsOn: true, CloseConfirm: true, MaxOpenPerUser: 1, NamingPattern: "ticket-{user}",
	}
	s.embeds[gid] = &TicketEmbed{
		GuildID: gid, Title: "Support", Description: "Click a button below to open a ticket.",
		Color: "#3498db", Footer: "Che1 Support", Thumbnail: "",
	}
	s.levelingSet[gid] = &LevelingSettings{
		Enabled: true, XPPerMessage: 15, CooldownSec: 60,
		AnnounceChannel: gid + "-c1", AnnounceMessage: "GG {user}, you leveled up to **{level}**!",
		Multiplier: 1.0,
	}
	s.staff[gid] = []TicketStaff{
		{UserID: "u1", Username: "alice", Roles: []string{gid + "-r3"}},
	}

	// tickets
	for i := 1; i <= 4; i++ {
		id := s.nextID()
		status := "open"
		closed := ""
		if i%2 == 0 {
			status = "closed"
			closed = isoDaysAgo(i)
		}
		s.tickets[gid] = append(s.tickets[gid], Ticket{
			ID: id, GuildID: gid, ChannelID: gid + "-c4-" + strconv.Itoa(i),
			UserID: "u" + strconv.Itoa(i), Username: []string{"alice", "bob", "carol", "dave"}[i-1],
			Subject: []string{"Account issue", "Billing question", "Bug report", "Feature request"}[i-1],
			Status:  status, Tags: []string{},
			CreatedAt: isoDaysAgo(i + 1), ClosedAt: closed,
		})
	}
	s.panels[gid] = []TicketPanel{{
		ID: s.nextID(), GuildID: gid, ChannelID: gid + "-c1",
		Title: "Support", Message: "Need help? Open a ticket below.", Color: "#3498db",
		Buttons: []TicketPanelButton{{Label: "Open Ticket", Style: "primary", Category: "general"}},
	}}
	s.forms[gid] = []TicketForm{{
		ID: s.nextID(), GuildID: gid, Name: "General",
		Fields: []TicketFormField{
			{Label: "Subject", Type: "short", Required: true},
			{Label: "Describe your issue", Type: "paragraph", Required: true},
		},
	}}
	s.tags[gid] = []TicketTag{
		{ID: s.nextID(), GuildID: gid, Name: "bug", Color: "#e74c3c"},
		{ID: s.nextID(), GuildID: gid, Name: "billing", Color: "#f1c40f"},
	}

	// leveling
	names := []string{"alice", "bob", "carol", "dave", "eve", "frank", "grace", "heidi"}
	for i, n := range names {
		s.userLevels[gid] = append(s.userLevels[gid], UserLevel{
			UserID: "u" + strconv.Itoa(i+1), Username: n,
			Level: 20 - i*2, XP: 3000 - i*200, TotalXP: 50000 - i*5000,
			UpdatedAt: isoDaysAgo(i),
		})
	}
	s.rewards[gid] = []LevelReward{
		{ID: s.nextID(), GuildID: gid, Level: 5, RoleID: gid + "-r2", RoleName: "Member"},
		{ID: s.nextID(), GuildID: gid, Level: 10, RoleID: gid + "-r5", RoleName: "Level 10"},
	}

	// mod logs
	acts := []string{"warn", "mute", "kick", "ban"}
	for i := 0; i < 6; i++ {
		s.modLogs[gid] = append(s.modLogs[gid], ModLog{
			ID: s.nextID(), GuildID: gid, ModeratorID: "u1", ModName: "alice",
			TargetID: "u" + strconv.Itoa((i%4)+1), TargetName: names[i%len(names)],
			Action: acts[i%len(acts)], Reason: "demo action " + strconv.Itoa(i),
			CreatedAt: isoDaysAgo(i),
		})
	}
	s.autoMod[gid] = []AutoModRule{
		{ID: s.nextID(), GuildID: gid, Name: "No invites", Trigger: "invites", Action: "delete", Enabled: true, Config: map[string]any{}},
		{ID: s.nextID(), GuildID: gid, Name: "Spam", Trigger: "spam", Action: "mute", Enabled: true, Config: map[string]any{"threshold": 5, "window_sec": 10}},
		{ID: s.nextID(), GuildID: gid, Name: "Caps lock", Trigger: "caps", Action: "warn", Enabled: false, Config: map[string]any{"percent": 70}},
	}
	s.reports[gid] = []Report{
		{ID: s.nextID(), GuildID: gid, ReporterID: "u2", ReporterName: "bob", TargetID: "u3", TargetName: "carol", Reason: "spam", Status: "open", CreatedAt: isoDaysAgo(1)},
	}

	// applications
	s.appForms[gid] = []ApplicationForm{
		{GuildID: gid, RoleID: gid + "-r3", URL: "https://forms.example/mod"},
	}
	s.apps[gid] = []Application{
		{ID: s.nextID(), GuildID: gid, UserID: "u4", Username: "dave", RoleID: gid + "-r3", RoleName: "Moderator",
			Answers: map[string]any{"why": "I love helping"}, Status: "pending", CreatedAt: isoDaysAgo(2)},
	}

	// giveaways
	s.giveaways[gid] = []Giveaway{
		{ID: s.nextID(), GuildID: gid, ChannelID: gid + "-c5", MessageID: "msg1",
			Prize: "Nitro Classic (1 month)", WinnerCount: 1, EndsAt: isoDaysFromNow(2), Status: "running",
			Entrants: 37, HostedBy: "alice", CreatedAt: isoDaysAgo(1)},
		{ID: s.nextID(), GuildID: gid, ChannelID: gid + "-c5", MessageID: "msg2",
			Prize: "$10 Steam Gift Card", WinnerCount: 1, EndsAt: isoDaysAgo(2), Status: "ended",
			Winners: []string{"bob"}, Entrants: 112, HostedBy: "carol", CreatedAt: isoDaysAgo(5)},
	}
	s.blacklist[gid] = []BlacklistEntry{}

	// alerts + history
	s.alerts[gid] = []Alert{
		{ID: s.nextID(), GuildID: gid, Severity: "warn", Title: "High spam rate",
			Detail: "Anti-spam triggered 12 times in the last hour.", CreatedAt: isoDaysAgo(0)},
		{ID: s.nextID(), GuildID: gid, Severity: "info", Title: "Backup completed",
			Detail: "Nightly backup finished successfully.", CreatedAt: isoDaysAgo(1)},
	}
	s.history[gid] = []HistoryEvent{
		{ID: s.nextID(), GuildID: gid, Actor: "alice", Event: "settings.update", Detail: "Changed welcome message", CreatedAt: isoDaysAgo(0)},
		{ID: s.nextID(), GuildID: gid, Actor: "carol", Event: "moderation.ban", Detail: "Banned spammer", CreatedAt: isoDaysAgo(1)},
	}
	s.backups[gid] = []Backup{
		{ID: s.nextID(), GuildID: gid, Label: "Nightly", SizeBytes: 184320, CreatedAt: isoDaysAgo(0), CreatedBy: "system"},
		{ID: s.nextID(), GuildID: gid, Label: "Before rename", SizeBytes: 172144, CreatedAt: isoDaysAgo(7), CreatedBy: "alice"},
	}
	s.perms[gid] = []Permission{
		{ID: s.nextID(), GuildID: gid, RoleID: gid + "-r4", RoleName: "Admin", Modules: []string{"tickets", "moderation", "giveaways", "leveling", "settings"}},
		{ID: s.nextID(), GuildID: gid, RoleID: gid + "-r3", RoleName: "Moderator", Modules: []string{"tickets", "moderation"}},
	}
}

// ---------------- Globals ----------------

var (
	store       = newStore()
	devMode     = envBool("DEV_MODE", true)
	clientID    = os.Getenv("DISCORD_CLIENT_ID")
	clientSec   = os.Getenv("DISCORD_CLIENT_SECRET")
	redirectURI = getenv("DISCORD_REDIRECT_URI", "http://localhost:8080/api/auth/callback")
	workerURL   = os.Getenv("WORKER_URL")
	workerKey   = os.Getenv("WORKER_API_KEY")
	frontendURL = getenv("FRONTEND_URL", "http://localhost:5173")
)

func main() {
	if devMode {
		store.seed()
		log.Println("DEV_MODE=true: seeded demo data (DevUser + 3 guilds)")
	}

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/api/auth/login", handleLogin)
	mux.HandleFunc("/api/auth/callback", handleCallback)
	mux.HandleFunc("/api/auth/logout", handleLogout)
	mux.HandleFunc("/api/auth/me", handleMe)

	// Guilds
	mux.HandleFunc("/api/guilds", handleGuilds)
	mux.HandleFunc("/api/guilds/", handleGuildScoped)

	// Static front-end (optional): serve built Vite output when present
	mux.Handle("/", http.FileServer(http.Dir("../dist/")))

	addr := getenv("ADDR", ":8080")
	log.Printf("CHE1 dashboard API listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

// ---------------- HTTP helpers ----------------

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = frontendURL
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func envBool(k string, d bool) bool {
	v := strings.ToLower(os.Getenv(k))
	if v == "" {
		return d
	}
	return v == "1" || v == "true" || v == "yes"
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func isoDaysAgo(n int) string {
	return time.Now().Add(-time.Duration(n) * 24 * time.Hour).UTC().Format(time.RFC3339)
}

func isoDaysFromNow(n int) string {
	return time.Now().Add(time.Duration(n) * 24 * time.Hour).UTC().Format(time.RFC3339)
}

// ---------------- Auth ----------------

const sidCookie = "che1_sid"

func setSession(w http.ResponseWriter, sid string) {
	http.SetCookie(w, &http.Cookie{
		Name: sidCookie, Value: sid, Path: "/", HttpOnly: true,
		SameSite: http.SameSiteLaxMode, MaxAge: 60 * 60 * 24 * 7,
	})
}

func clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: sidCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
}

func currentUser(r *http.Request) *User {
	c, err := r.Cookie(sidCookie)
	if err != nil {
		return nil
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	s := store.sessions[c.Value]
	if s == nil || time.Now().After(s.Expires) {
		return nil
	}
	return store.users[s.UserID]
}

func requireUser(w http.ResponseWriter, r *http.Request) *User {
	u := currentUser(r)
	if u == nil {
		writeErr(w, http.StatusUnauthorized, "not authenticated")
		return nil
	}
	return u
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if devMode {
		// Auto-login as DevUser
		store.mu.Lock()
		sid := randHex(24)
		store.sessions[sid] = &Session{ID: sid, UserID: "100000000000000001", Expires: time.Now().Add(7 * 24 * time.Hour)}
		store.mu.Unlock()
		setSession(w, sid)
		http.Redirect(w, r, frontendURL+"/#/dashboard/servers/servers", http.StatusFound)
		return
	}
	if clientID == "" {
		writeErr(w, http.StatusInternalServerError, "DISCORD_CLIENT_ID not configured")
		return
	}
	state := randHex(16)
	http.SetCookie(w, &http.Cookie{Name: "che1_state", Value: state, Path: "/", HttpOnly: true, MaxAge: 600})
	u := url.Values{}
	u.Set("client_id", clientID)
	u.Set("redirect_uri", redirectURI)
	u.Set("response_type", "code")
	u.Set("scope", "identify guilds")
	u.Set("state", state)
	http.Redirect(w, r, "https://discord.com/oauth2/authorize?"+u.Encode(), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if devMode {
		handleLogin(w, r)
		return
	}
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	sc, err := r.Cookie("che1_state")
	if err != nil || sc.Value != state {
		writeErr(w, http.StatusBadRequest, "invalid state")
		return
	}
	// Exchange code
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSec)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)
	req, _ := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	defer resp.Body.Close()
	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	// Fetch user
	user, err := discordGet[User](tok.AccessToken, "/users/@me")
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	guilds, err := discordGet[[]Guild](tok.AccessToken, "/users/@me/guilds")
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	// Store
	store.mu.Lock()
	store.users[user.ID] = &user
	ids := make([]string, 0, len(guilds))
	for _, g := range guilds {
		copy := g
		store.guilds[g.ID] = &copy
		ids = append(ids, g.ID)
	}
	store.userGuilds[user.ID] = ids
	sid := randHex(24)
	store.sessions[sid] = &Session{ID: sid, UserID: user.ID, AccessToken: tok.AccessToken, Expires: time.Now().Add(7 * 24 * time.Hour)}
	store.mu.Unlock()
	setSession(w, sid)
	http.Redirect(w, r, frontendURL+"/#/dashboard/servers/servers", http.StatusFound)
}

func discordGet[T any](token, path string) (T, error) {
	var out T
	req, _ := http.NewRequest("GET", "https://discord.com/api/v10"+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return out, fmt.Errorf("discord %d: %s", resp.StatusCode, string(b))
	}
	return out, json.NewDecoder(resp.Body).Decode(&out)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(sidCookie)
	if err == nil {
		store.mu.Lock()
		delete(store.sessions, c.Value)
		store.mu.Unlock()
	}
	clearSession(w)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if u == nil {
		writeErr(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

// ---------------- Guilds list ----------------

func handleGuilds(w http.ResponseWriter, r *http.Request) {
	u := requireUser(w, r)
	if u == nil {
		return
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	ids := store.userGuilds[u.ID]
	out := make([]Guild, 0, len(ids))
	for _, id := range ids {
		if g, ok := store.guilds[id]; ok {
			g.BotPresent = store.botGuilds[id]
			out = append(out, *g)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

// ---------------- Guild-scoped router ----------------
// URL shape: /api/guilds/{guildID}/{resource...}

func handleGuildScoped(w http.ResponseWriter, r *http.Request) {
	u := requireUser(w, r)
	if u == nil {
		return
	}
	rest := strings.TrimPrefix(r.URL.Path, "/api/guilds/")
	parts := strings.Split(rest, "/")
	if len(parts) < 1 || parts[0] == "" {
		writeErr(w, http.StatusBadRequest, "missing guild id")
		return
	}
	gid := parts[0]
	// authorize
	store.mu.RLock()
	ok := false
	for _, id := range store.userGuilds[u.ID] {
		if id == gid {
			ok = true
			break
		}
	}
	store.mu.RUnlock()
	if !ok {
		writeErr(w, http.StatusForbidden, "no access to this guild")
		return
	}
	if len(parts) == 1 {
		// /api/guilds/{id}
		store.mu.RLock()
		g := store.guilds[gid]
		store.mu.RUnlock()
		if g == nil {
			writeErr(w, http.StatusNotFound, "guild not found")
			return
		}
		writeJSON(w, http.StatusOK, g)
		return
	}
	res := parts[1]
	rem := parts[2:]
	switch res {
	case "enable":
		handleEnableBot(w, r, gid)
	case "channels":
		handleChannels(w, r, gid)
	case "roles":
		handleRoles(w, r, gid)
	case "members":
		handleMembers(w, r, gid)
	case "settings":
		handleSettings(w, r, gid)
	case "overview":
		handleOverview(w, r, gid)
	case "alerts":
		handleAlerts(w, r, gid, rem)
	case "history":
		handleHistory(w, r, gid)
	case "backup":
		handleBackup(w, r, gid, rem)
	case "permissions":
		handlePermissions(w, r, gid, rem)
	case "tickets":
		handleTickets(w, r, gid, rem)
	case "moderation":
		handleModeration(w, r, gid, rem)
	case "giveaways":
		handleGiveaways(w, r, gid, rem)
	case "leveling":
		handleLeveling(w, r, gid, rem)
	case "applications":
		handleApplications(w, r, gid, rem)
	default:
		writeErr(w, http.StatusNotFound, "unknown resource: "+res)
	}
}

// ---------------- Per-guild meta ----------------

func handleEnableBot(w http.ResponseWriter, r *http.Request, gid string) {
	switch r.Method {
	case "POST":
		store.mu.Lock()
		store.botGuilds[gid] = true
		if store.guilds[gid] != nil {
			store.guilds[gid].BotPresent = true
		}
		if store.settings[gid] == nil {
			store.mu.Unlock()
			store.mu.Lock()
		}
		store.mu.Unlock()
		store.mu.Lock()
		if _, ok := store.settings[gid]; !ok {
			store.seedGuild(gid)
		}
		store.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "invite_url": botInvite(gid)})
	case "DELETE":
		store.mu.Lock()
		delete(store.botGuilds, gid)
		if store.guilds[gid] != nil {
			store.guilds[gid].BotPresent = false
		}
		store.mu.Unlock()
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	case "GET":
		writeJSON(w, http.StatusOK, map[string]any{"invite_url": botInvite(gid)})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

func botInvite(gid string) string {
	if clientID == "" {
		return "https://discord.com/oauth2/authorize?client_id=CLIENT_ID&scope=bot&permissions=8&guild_id=" + gid
	}
	return "https://discord.com/oauth2/authorize?client_id=" + clientID + "&scope=bot&permissions=8&guild_id=" + gid
}

func handleChannels(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.channels[gid])
}

func handleRoles(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.roles[gid])
}

func handleMembers(w http.ResponseWriter, r *http.Request, gid string) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	store.mu.RLock()
	defer store.mu.RUnlock()
	all := store.members[gid]
	if q == "" {
		writeJSON(w, http.StatusOK, all)
		return
	}
	out := make([]Member, 0)
	for _, m := range all {
		if strings.Contains(strings.ToLower(m.Username), q) || strings.Contains(strings.ToLower(m.Nickname), q) {
			out = append(out, m)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func handleSettings(w http.ResponseWriter, r *http.Request, gid string) {
	switch r.Method {
	case "GET":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.settings[gid])
	case "PATCH":
		var patch GuildSettings
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		store.settings[gid] = &patch
		recordHistory(gid, "settings.update", "Settings updated")
		store.mu.Unlock()
		writeJSON(w, http.StatusOK, patch)
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

func recordHistory(gid, event, detail string) {
	store.history[gid] = append([]HistoryEvent{{
		ID: store.nextID(), GuildID: gid, Actor: "dashboard", Event: event, Detail: detail, CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}}, store.history[gid]...)
}

// ---------------- Overview ----------------

func handleOverview(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	var s Stats
	s.Members = 0
	if g := store.guilds[gid]; g != nil {
		s.Members = g.MemberCount
	}
	s.OnlineMembers = s.Members * 40 / 100
	s.MessagesToday = 1200 + int(time.Now().Unix()%500)
	s.CommandsToday = 80 + int(time.Now().Unix()%40)
	for _, t := range store.tickets[gid] {
		if t.Status == "open" {
			s.OpenTickets++
		}
	}
	for _, g := range store.giveaways[gid] {
		if g.Status == "running" {
			s.ActiveGiveaways++
		}
	}
	s.ModActionsWeek = len(store.modLogs[gid])
	for i := 6; i >= 0; i-- {
		s.MessagesPerDay = append(s.MessagesPerDay, TimeSeriesPt{
			Date: time.Now().AddDate(0, 0, -i).Format("2006-01-02"), Value: 800 + (i*120)%500,
		})
		s.JoinsPerDay = append(s.JoinsPerDay, TimeSeriesPt{
			Date: time.Now().AddDate(0, 0, -i).Format("2006-01-02"), Value: 5 + (i*3)%10,
		})
	}
	writeJSON(w, http.StatusOK, s)
}

func handleAlerts(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	switch {
	case r.Method == "GET":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.alerts[gid])
	case r.Method == "DELETE" && len(rem) == 1:
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.alerts[gid][:0]
		for _, a := range store.alerts[gid] {
			if a.ID != id {
				out = append(out, a)
			}
		}
		store.alerts[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

func handleHistory(w http.ResponseWriter, r *http.Request, gid string) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	writeJSON(w, http.StatusOK, store.history[gid])
}

func handleBackup(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.backups[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var body struct{ Label string }
		_ = readJSON(r, &body)
		if body.Label == "" {
			body.Label = "Manual backup"
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		b := Backup{ID: store.nextID(), GuildID: gid, Label: body.Label, SizeBytes: 140000 + int(time.Now().Unix()%50000),
			CreatedAt: time.Now().UTC().Format(time.RFC3339), CreatedBy: "dashboard"}
		store.backups[gid] = append([]Backup{b}, store.backups[gid]...)
		recordHistory(gid, "backup.create", body.Label)
		writeJSON(w, http.StatusOK, b)
		return
	}
	if r.Method == "POST" && len(rem) == 2 && rem[1] == "restore" {
		recordHistory(gid, "backup.restore", "Backup "+rem[0])
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.backups[gid][:0]
		for _, b := range store.backups[gid] {
			if b.ID != id {
				out = append(out, b)
			}
		}
		store.backups[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handlePermissions(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	switch {
	case r.Method == "GET":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.perms[gid])
	case r.Method == "POST":
		var p Permission
		if err := readJSON(r, &p); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		p.ID = store.nextID()
		p.GuildID = gid
		store.perms[gid] = append(store.perms[gid], p)
		recordHistory(gid, "permissions.grant", p.RoleName)
		writeJSON(w, http.StatusOK, p)
	case r.Method == "DELETE" && len(rem) == 1:
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.perms[gid][:0]
		for _, p := range store.perms[gid] {
			if p.ID != id {
				out = append(out, p)
			}
		}
		store.perms[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		writeErr(w, http.StatusMethodNotAllowed, "")
	}
}

// ---------------- Tickets ----------------

func handleTickets(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.tickets[gid])
			return
		}
		if r.Method == "POST" {
			var t Ticket
			if err := readJSON(r, &t); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			t.ID = store.nextID()
			t.GuildID = gid
			t.Status = "open"
			t.CreatedAt = time.Now().UTC().Format(time.RFC3339)
			store.tickets[gid] = append([]Ticket{t}, store.tickets[gid]...)
			recordHistory(gid, "tickets.create", t.Subject)
			forwardWorker("tickets.create", gid, t)
			writeJSON(w, http.StatusOK, t)
			return
		}
	}
	switch rem[0] {
	case "settings":
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.ticketSet[gid])
			return
		}
		if r.Method == "PATCH" {
			var s TicketSettings
			if err := readJSON(r, &s); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			store.ticketSet[gid] = &s
			writeJSON(w, http.StatusOK, s)
			return
		}
	case "panels":
		handleTicketPanels(w, r, gid, rem[1:])
		return
	case "forms":
		handleTicketForms(w, r, gid, rem[1:])
		return
	case "tags":
		handleTicketTags(w, r, gid, rem[1:])
		return
	case "embed":
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.embeds[gid])
			return
		}
		if r.Method == "PATCH" {
			var e TicketEmbed
			if err := readJSON(r, &e); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			e.GuildID = gid
			store.embeds[gid] = &e
			writeJSON(w, http.StatusOK, e)
			return
		}
	case "staff":
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.staff[gid])
			return
		}
		if r.Method == "POST" {
			var s TicketStaff
			if err := readJSON(r, &s); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			store.staff[gid] = append(store.staff[gid], s)
			writeJSON(w, http.StatusOK, s)
			return
		}
		if r.Method == "DELETE" && len(rem) == 2 {
			store.mu.Lock()
			defer store.mu.Unlock()
			out := store.staff[gid][:0]
			for _, s := range store.staff[gid] {
				if s.UserID != rem[1] {
					out = append(out, s)
				}
			}
			store.staff[gid] = out
			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
			return
		}
	case "audit":
		store.mu.RLock()
		defer store.mu.RUnlock()
		out := make([]HistoryEvent, 0)
		for _, h := range store.history[gid] {
			if strings.HasPrefix(h.Event, "tickets.") {
				out = append(out, h)
			}
		}
		writeJSON(w, http.StatusOK, out)
		return
	default:
		// ticket by id
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		if len(rem) == 2 && rem[1] == "transcript" && r.Method == "GET" {
			writeJSON(w, http.StatusOK, map[string]any{
				"id": id, "content": "Transcript for ticket " + rem[0] + "\n\n(server-generated)",
			})
			return
		}
		if r.Method == "PATCH" {
			var patch Ticket
			if err := readJSON(r, &patch); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			for i, t := range store.tickets[gid] {
				if t.ID == id {
					if patch.Status != "" {
						store.tickets[gid][i].Status = patch.Status
						if patch.Status == "closed" {
							store.tickets[gid][i].ClosedAt = time.Now().UTC().Format(time.RFC3339)
						}
					}
					if patch.AssignedTo != "" {
						store.tickets[gid][i].AssignedTo = patch.AssignedTo
					}
					if patch.Tags != nil {
						store.tickets[gid][i].Tags = patch.Tags
					}
					recordHistory(gid, "tickets.update", "Ticket #"+rem[0])
					forwardWorker("tickets.update", gid, store.tickets[gid][i])
					writeJSON(w, http.StatusOK, store.tickets[gid][i])
					return
				}
			}
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			for _, t := range store.tickets[gid] {
				if t.ID == id {
					writeJSON(w, http.StatusOK, t)
					return
				}
			}
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleTicketPanels(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.panels[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var p TicketPanel
		if err := readJSON(r, &p); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		p.ID = store.nextID()
		p.GuildID = gid
		store.panels[gid] = append(store.panels[gid], p)
		writeJSON(w, http.StatusOK, p)
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.panels[gid][:0]
		for _, p := range store.panels[gid] {
			if p.ID != id {
				out = append(out, p)
			}
		}
		store.panels[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleTicketForms(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.forms[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var f TicketForm
		if err := readJSON(r, &f); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		f.ID = store.nextID()
		f.GuildID = gid
		store.forms[gid] = append(store.forms[gid], f)
		writeJSON(w, http.StatusOK, f)
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.forms[gid][:0]
		for _, f := range store.forms[gid] {
			if f.ID != id {
				out = append(out, f)
			}
		}
		store.forms[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleTicketTags(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.tags[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var t TicketTag
		if err := readJSON(r, &t); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		t.ID = store.nextID()
		t.GuildID = gid
		store.tags[gid] = append(store.tags[gid], t)
		writeJSON(w, http.StatusOK, t)
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.tags[gid][:0]
		for _, t := range store.tags[gid] {
			if t.ID != id {
				out = append(out, t)
			}
		}
		store.tags[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

// ---------------- Moderation ----------------

func handleModeration(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		writeErr(w, http.StatusNotFound, "")
		return
	}
	switch rem[0] {
	case "logs":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.modLogs[gid])
	case "automod":
		handleAutoMod(w, r, gid, rem[1:])
	case "reports":
		handleReports(w, r, gid, rem[1:])
	case "actions":
		if r.Method != "POST" {
			writeErr(w, http.StatusMethodNotAllowed, "")
			return
		}
		var body struct {
			TargetID    string `json:"target_id"`
			TargetName  string `json:"target_name"`
			Action      string `json:"action"`
			Reason      string `json:"reason"`
			DurationSec int    `json:"duration_sec"`
		}
		if err := readJSON(r, &body); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		if body.Action == "" || body.TargetID == "" {
			writeErr(w, http.StatusBadRequest, "action and target_id required")
			return
		}
		u := currentUser(r)
		store.mu.Lock()
		defer store.mu.Unlock()
		ml := ModLog{
			ID: store.nextID(), GuildID: gid,
			ModeratorID: u.ID, ModName: u.Username,
			TargetID: body.TargetID, TargetName: body.TargetName,
			Action: body.Action, Reason: body.Reason, DurationSec: body.DurationSec,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}
		store.modLogs[gid] = append([]ModLog{ml}, store.modLogs[gid]...)
		recordHistory(gid, "moderation."+body.Action, body.TargetName)
		forwardWorker("moderation.action", gid, ml)
		writeJSON(w, http.StatusOK, ml)
	default:
		writeErr(w, http.StatusNotFound, "")
	}
}

func handleAutoMod(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.autoMod[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var rule AutoModRule
		if err := readJSON(r, &rule); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		rule.ID = store.nextID()
		rule.GuildID = gid
		store.autoMod[gid] = append(store.autoMod[gid], rule)
		writeJSON(w, http.StatusOK, rule)
		return
	}
	if r.Method == "PATCH" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		var patch AutoModRule
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		for i, rule := range store.autoMod[gid] {
			if rule.ID == id {
				patch.ID = id
				patch.GuildID = gid
				store.autoMod[gid][i] = patch
				writeJSON(w, http.StatusOK, patch)
				return
			}
		}
		writeErr(w, http.StatusNotFound, "")
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.autoMod[gid][:0]
		for _, rule := range store.autoMod[gid] {
			if rule.ID != id {
				out = append(out, rule)
			}
		}
		store.autoMod[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleReports(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.reports[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var rep Report
		if err := readJSON(r, &rep); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		rep.ID = store.nextID()
		rep.GuildID = gid
		rep.Status = "open"
		rep.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		store.reports[gid] = append([]Report{rep}, store.reports[gid]...)
		writeJSON(w, http.StatusOK, rep)
		return
	}
	if r.Method == "PATCH" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		var patch Report
		if err := readJSON(r, &patch); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		for i, rep := range store.reports[gid] {
			if rep.ID == id {
				if patch.Status != "" {
					store.reports[gid][i].Status = patch.Status
				}
				writeJSON(w, http.StatusOK, store.reports[gid][i])
				return
			}
		}
		writeErr(w, http.StatusNotFound, "")
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

// ---------------- Giveaways ----------------

func handleGiveaways(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.giveaways[gid])
			return
		}
		if r.Method == "POST" {
			var g Giveaway
			if err := readJSON(r, &g); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			g.ID = store.nextID()
			g.GuildID = gid
			g.Status = "running"
			g.CreatedAt = time.Now().UTC().Format(time.RFC3339)
			if g.WinnerCount <= 0 {
				g.WinnerCount = 1
			}
			store.giveaways[gid] = append([]Giveaway{g}, store.giveaways[gid]...)
			recordHistory(gid, "giveaways.create", g.Prize)
			forwardWorker("giveaways.create", gid, g)
			writeJSON(w, http.StatusOK, g)
			return
		}
	}
	switch rem[0] {
	case "active":
		store.mu.RLock()
		defer store.mu.RUnlock()
		out := make([]Giveaway, 0)
		for _, g := range store.giveaways[gid] {
			if g.Status == "running" {
				out = append(out, g)
			}
		}
		writeJSON(w, http.StatusOK, out)
		return
	case "blacklist":
		handleBlacklist(w, r, gid, rem[1:])
		return
	case "premium":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, map[string]any{
			"premium":      store.settings[gid] != nil && store.settings[gid].Premium,
			"benefits":     []string{"Unlimited active giveaways", "Bonus-entry roles", "Custom embeds", "Role-weighted winners"},
			"upgrade_url":  "/#/premium",
		})
		return
	default:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
			return
		}
		if r.Method == "POST" && len(rem) == 2 && rem[1] == "end" {
			store.mu.Lock()
			defer store.mu.Unlock()
			for i, g := range store.giveaways[gid] {
				if g.ID == id {
					store.giveaways[gid][i].Status = "ended"
					if len(store.giveaways[gid][i].Winners) == 0 {
						store.giveaways[gid][i].Winners = []string{"bob"}
					}
					forwardWorker("giveaways.end", gid, store.giveaways[gid][i])
					writeJSON(w, http.StatusOK, store.giveaways[gid][i])
					return
				}
			}
			writeErr(w, http.StatusNotFound, "")
			return
		}
		if r.Method == "POST" && len(rem) == 2 && rem[1] == "reroll" {
			store.mu.Lock()
			defer store.mu.Unlock()
			for i, g := range store.giveaways[gid] {
				if g.ID == id {
					store.giveaways[gid][i].Winners = []string{"alice"}
					forwardWorker("giveaways.reroll", gid, store.giveaways[gid][i])
					writeJSON(w, http.StatusOK, store.giveaways[gid][i])
					return
				}
			}
			writeErr(w, http.StatusNotFound, "")
			return
		}
		if r.Method == "DELETE" {
			store.mu.Lock()
			defer store.mu.Unlock()
			out := store.giveaways[gid][:0]
			for _, g := range store.giveaways[gid] {
				if g.ID != id {
					out = append(out, g)
				}
			}
			store.giveaways[gid] = out
			writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
			return
		}
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleBlacklist(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.blacklist[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var e BlacklistEntry
		if err := readJSON(r, &e); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		e.ID = store.nextID()
		e.GuildID = gid
		e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		store.blacklist[gid] = append(store.blacklist[gid], e)
		writeJSON(w, http.StatusOK, e)
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.blacklist[gid][:0]
		for _, e := range store.blacklist[gid] {
			if e.ID != id {
				out = append(out, e)
			}
		}
		store.blacklist[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

// ---------------- Leveling ----------------

func handleLeveling(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		writeErr(w, http.StatusNotFound, "")
		return
	}
	switch rem[0] {
	case "settings":
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.levelingSet[gid])
			return
		}
		if r.Method == "PATCH" {
			var s LevelingSettings
			if err := readJSON(r, &s); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			store.levelingSet[gid] = &s
			writeJSON(w, http.StatusOK, s)
			return
		}
	case "leaderboard":
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.userLevels[gid])
		return
	case "rewards":
		handleRewards(w, r, gid, rem[1:])
		return
	case "stats":
		store.mu.RLock()
		defer store.mu.RUnlock()
		lvls := store.userLevels[gid]
		total, avg, top := 0, 0, 0
		for _, u := range lvls {
			total += u.TotalXP
			if u.Level > top {
				top = u.Level
			}
		}
		if len(lvls) > 0 {
			avg = total / len(lvls)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"tracked_users":  len(lvls),
			"total_xp":       total,
			"average_xp":     avg,
			"top_level":      top,
			"top_users":      lvls,
		})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

func handleRewards(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if r.Method == "GET" && len(rem) == 0 {
		store.mu.RLock()
		defer store.mu.RUnlock()
		writeJSON(w, http.StatusOK, store.rewards[gid])
		return
	}
	if r.Method == "POST" && len(rem) == 0 {
		var rw LevelReward
		if err := readJSON(r, &rw); err != nil {
			writeErr(w, http.StatusBadRequest, err.Error())
			return
		}
		store.mu.Lock()
		defer store.mu.Unlock()
		rw.ID = store.nextID()
		rw.GuildID = gid
		store.rewards[gid] = append(store.rewards[gid], rw)
		writeJSON(w, http.StatusOK, rw)
		return
	}
	if r.Method == "DELETE" && len(rem) == 1 {
		id, _ := strconv.ParseInt(rem[0], 10, 64)
		store.mu.Lock()
		defer store.mu.Unlock()
		out := store.rewards[gid][:0]
		for _, rw := range store.rewards[gid] {
			if rw.ID != id {
				out = append(out, rw)
			}
		}
		store.rewards[gid] = out
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

// ---------------- Applications ----------------

func handleApplications(w http.ResponseWriter, r *http.Request, gid string, rem []string) {
	if len(rem) == 0 {
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.apps[gid])
			return
		}
	}
	switch rem[0] {
	case "forms":
		if r.Method == "GET" {
			store.mu.RLock()
			defer store.mu.RUnlock()
			writeJSON(w, http.StatusOK, store.appForms[gid])
			return
		}
		if r.Method == "POST" {
			var f ApplicationForm
			if err := readJSON(r, &f); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			f.GuildID = gid
			store.mu.Lock()
			defer store.mu.Unlock()
			store.appForms[gid] = append(store.appForms[gid], f)
			writeJSON(w, http.StatusOK, f)
			return
		}
	default:
		id, err := strconv.ParseInt(rem[0], 10, 64)
		if err != nil {
			writeErr(w, http.StatusBadRequest, "bad id")
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
					writeJSON(w, http.StatusOK, store.apps[gid][i])
					return
				}
			}
			writeErr(w, http.StatusNotFound, "")
			return
		}
	}
	writeErr(w, http.StatusMethodNotAllowed, "")
}

// ---------------- Worker forward ----------------

// forwardWorker sends an action event to the Worker service if configured.
// In dev/without WORKER_URL it's a no-op.
func forwardWorker(event, gid string, payload any) {
	if workerURL == "" {
		return
	}
	body, _ := json.Marshal(map[string]any{"event": event, "guild_id": gid, "payload": payload})
	req, _ := http.NewRequest("POST", strings.TrimRight(workerURL, "/")+"/api/v1/tasks", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	if workerKey != "" {
		req.Header.Set("Authorization", "Bearer "+workerKey)
	}
	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("worker forward error: %v", err)
			return
		}
		resp.Body.Close()
	}()
}
