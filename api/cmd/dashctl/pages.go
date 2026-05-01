package main

// pages.go is the source of truth for the dashboard's client-side routing
// table. `dashctl gen-pages` emits src/dashboard/data/pages.js from the
// definitions below.
//
// Adding or renaming a page: edit only this file, then run
//   dashctl gen-pages

// SidebarLink is one entry in the left rail.
type SidebarLink struct {
	Path  string // hash path, e.g. "tickets/tickets"
	Icon  string // FontAwesome class, e.g. "fa-ticket"
	Label string
	ID    string // matches Route.SidebarActive
}

// RightNavItem is one entry in a right-rail submenu.
type RightNavItem struct {
	Path  string
	Icon  string
	Label string
	ID    string
}

// Section bundles a sidebar tab with its submenu and routes.
type Section struct {
	JSVar    string         // JS variable holding the right-nav array
	Sidebar  SidebarLink    // entry inserted into sidebarLinks
	RightNav []RightNavItem // entries in the section's right rail
	Routes   []RouteEntry   // pages under this section
}

// RouteEntry maps a hash path to a Svelte component + page title.
type RouteEntry struct {
	Path          string
	Component     string
	PageTitle     string
	SidebarActive string // defaults to Section.Sidebar.ID
}

// rootRoutes lists routes that don't belong to a section's right-rail (e.g.
// the bare /dashboard landing page).
var rootRoutes = []RouteEntry{
	{Path: "dashboard", Component: "Dashboard", PageTitle: "Dashboard", SidebarActive: "Dashboard"},
}

// sections is the spine of the dashboard. Edit here, regen, done.
var sections = []Section{
	{
		JSVar: "overviewRightNav",
		Sidebar: SidebarLink{Path: "overview", Icon: "fa-chart-line", Label: "Overview", ID: "Overview"},
		RightNav: []RightNavItem{
			{Path: "overview/statistics", Icon: "fa-chart-bar", Label: "Statistics", ID: "Statistics"},
			{Path: "overview/settings", Icon: "fa-cog", Label: "Settings", ID: "Settings"},
			{Path: "overview/alerts", Icon: "fa-bell", Label: "Alerts", ID: "Alerts"},
			{Path: "overview/history", Icon: "fa-history", Label: "History", ID: "History"},
		},
		Routes: []RouteEntry{
			{Path: "overview", Component: "OverviewOverview", PageTitle: "Overview"},
			{Path: "overview/statistics", Component: "OverviewStatistics", PageTitle: "Statistics"},
			{Path: "overview/settings", Component: "OverviewSettings", PageTitle: "Settings"},
			{Path: "overview/alerts", Component: "OverviewAlerts", PageTitle: "Alerts"},
			{Path: "overview/history", Component: "OverviewHistory", PageTitle: "History"},
		},
	},
	{
		JSVar: "serversRightNav",
		Sidebar: SidebarLink{Path: "servers/servers", Icon: "fa-server", Label: "Servers", ID: "Servers"},
		RightNav: []RightNavItem{
			{Path: "servers/add-server", Icon: "fa-plus", Label: "Add Server", ID: "AddServer"},
			{Path: "servers/manage-roles", Icon: "fa-users-cog", Label: "Manage Roles", ID: "ManageRoles"},
			{Path: "servers/permissions", Icon: "fa-key", Label: "Permissions", ID: "Permissions"},
			{Path: "servers/backup", Icon: "fa-save", Label: "Backup", ID: "Backup"},
		},
		Routes: []RouteEntry{
			{Path: "servers/servers", Component: "ServersServers", PageTitle: "Servers"},
			{Path: "servers/add-server", Component: "ServersAdd_server", PageTitle: "Add Server"},
			{Path: "servers/manage-roles", Component: "ServersManage_roles", PageTitle: "Manage Roles"},
			{Path: "servers/permissions", Component: "ServersPermissions", PageTitle: "Permissions"},
			{Path: "servers/backup", Component: "ServersBackup", PageTitle: "Backup"},
		},
	},
	{
		JSVar: "ticketsRightNav",
		Sidebar: SidebarLink{Path: "tickets/tickets", Icon: "fa-ticket", Label: "Tickets", ID: "Tickets"},
		RightNav: []RightNavItem{
			{Path: "tickets/tickets", Icon: "fa-gears", Label: "Settings", ID: "Settings"},
			{Path: "tickets/categories", Icon: "fa-layer-group", Label: "Categories", ID: "Categories"},
			{Path: "tickets/panels", Icon: "fa-arrow-pointer", Label: "Panels", ID: "Panels"},
			{Path: "tickets/forms", Icon: "fa-file-lines", Label: "Forms", ID: "Forms"},
			{Path: "tickets/snippets", Icon: "fa-comment-dots", Label: "Snippets", ID: "Snippets"},
			{Path: "tickets/tags", Icon: "fa-tag", Label: "Tags", ID: "Tags"},
			{Path: "tickets/ticket-staff", Icon: "fa-people-group", Label: "Staff", ID: "TicketStaff"},
			{Path: "tickets/transcripts", Icon: "fa-file", Label: "Transcripts", ID: "Transcripts"},
			{Path: "tickets/statistics", Icon: "fa-chart-bar", Label: "Statistics", ID: "Statistics"},
			{Path: "tickets/audit-logs", Icon: "fa-file-medical", Label: "Audit logs", ID: "AuditLogs"},
		},
		Routes: []RouteEntry{
			{Path: "tickets/tickets", Component: "TicketsTickets", PageTitle: "Tickets"},
			{Path: "tickets/categories", Component: "TicketsCategories", PageTitle: "Categories"},
			{Path: "tickets/panels", Component: "TicketsPanels", PageTitle: "Panels"},
			{Path: "tickets/forms", Component: "TicketsForms", PageTitle: "Forms"},
			{Path: "tickets/snippets", Component: "TicketsSnippets", PageTitle: "Snippets"},
			{Path: "tickets/tags", Component: "TicketsTags", PageTitle: "Tags"},
			{Path: "tickets/ticket-embed", Component: "TicketsTicket_embed", PageTitle: "Ticket Embed"},
			{Path: "tickets/ticket-staff", Component: "TicketsTicket_staff", PageTitle: "Ticket Staff"},
			{Path: "tickets/transcripts", Component: "TicketsTranscripts", PageTitle: "Transcripts"},
			{Path: "tickets/statistics", Component: "TicketsStatistics", PageTitle: "Statistics"},
			{Path: "tickets/audit-logs", Component: "TicketsAudit_logs", PageTitle: "Audit Logs"},
		},
	},
	{
		JSVar: "giveawaysRightNav",
		Sidebar: SidebarLink{Path: "giveaways/giveaway", Icon: "fa-gift", Label: "Giveaways", ID: "Giveaways"},
		RightNav: []RightNavItem{
			{Path: "giveaways/create-giveaway", Icon: "fa-plus-circle", Label: "Create Giveaway", ID: "CreateGiveaway"},
			{Path: "giveaways/active-giveaways", Icon: "fa-list", Label: "Active Giveaways", ID: "ActiveGiveaways"},
			{Path: "giveaways/premium", Icon: "fa-star", Label: "Premium", ID: "Premium"},
			{Path: "giveaways/blacklist", Icon: "fa-ban", Label: "Blacklist", ID: "Blacklist"},
		},
		Routes: []RouteEntry{
			{Path: "giveaways/giveaway", Component: "GiveawaysGiveaway", PageTitle: "Giveaways"},
			{Path: "giveaways/active-giveaways", Component: "GiveawaysActive_giveaways", PageTitle: "Active Giveaways"},
			{Path: "giveaways/blacklist", Component: "GiveawaysBlacklist", PageTitle: "Blacklist"},
			{Path: "giveaways/create-giveaway", Component: "GiveawaysCreate_giveaway", PageTitle: "Create Giveaway"},
			{Path: "giveaways/premium", Component: "GiveawaysPremium", PageTitle: "Premium"},
		},
	},
	{
		JSVar: "moderationRightNav",
		Sidebar: SidebarLink{Path: "moderation/moderation", Icon: "fa-shield", Label: "Moderation", ID: "Moderation"},
		RightNav: []RightNavItem{
			{Path: "moderation/auto-mod", Icon: "fa-robot", Label: "Auto Mod", ID: "AutoMod"},
			{Path: "moderation/manual-actions", Icon: "fa-hand-paper", Label: "Manual Actions", ID: "ManualActions"},
			{Path: "moderation/reports", Icon: "fa-flag", Label: "Reports", ID: "Reports"},
			{Path: "moderation/logs", Icon: "fa-file-alt", Label: "Logs", ID: "Logs"},
		},
		Routes: []RouteEntry{
			{Path: "moderation/moderation", Component: "ModerationModeration", PageTitle: "Moderation"},
			{Path: "moderation/auto-mod", Component: "ModerationAuto_mod", PageTitle: "Auto Mod"},
			{Path: "moderation/manual-actions", Component: "ModerationManual_actions", PageTitle: "Manual Actions"},
			{Path: "moderation/reports", Component: "ModerationReports", PageTitle: "Reports"},
			{Path: "moderation/logs", Component: "ModerationLogs", PageTitle: "Logs"},
		},
	},
	{
		JSVar: "applicationsRightNav",
		Sidebar: SidebarLink{Path: "applications/applications", Icon: "fa-clipboard-check", Label: "Applications", ID: "Applications"},
		RightNav: []RightNavItem{
			{Path: "applications/applications", Icon: "fa-inbox", Label: "Review queue", ID: "ReviewQueue"},
			{Path: "applications/forms", Icon: "fa-list-check", Label: "Forms", ID: "Forms"},
			{Path: "applications/statistics", Icon: "fa-chart-bar", Label: "Statistics", ID: "Statistics"},
		},
		Routes: []RouteEntry{
			{Path: "applications/applications", Component: "ApplicationsApplications", PageTitle: "Applications"},
			{Path: "applications/forms", Component: "ApplicationsForms", PageTitle: "Application Forms"},
			{Path: "applications/statistics", Component: "ApplicationsStatistics", PageTitle: "Application Statistics"},
		},
	},
	{
		// Note: SidebarActive is "Dashboard" rather than the section ID below
		// because the legacy layout chose to highlight the gear icon.
		JSVar: "",
		Sidebar: SidebarLink{Path: "dashboard", Icon: "fa-gear", Label: "Settings", ID: "Dashboard"},
	},
	{
		JSVar: "levelingRightNav",
		Sidebar: SidebarLink{Path: "leveling/leveling", Icon: "fa-ranking-star", Label: "Leveling", ID: "Leveling"},
		RightNav: []RightNavItem{
			{Path: "leveling/leveling", Icon: "fa-cog", Label: "Settings", ID: "Settings"},
			{Path: "leveling/leaderboard", Icon: "fa-trophy", Label: "Leaderboard", ID: "Leaderboard"},
			{Path: "leveling/rewards", Icon: "fa-gift", Label: "Rewards", ID: "Rewards"},
			{Path: "leveling/stats", Icon: "fa-chart-line", Label: "Stats", ID: "Stats"},
		},
		Routes: []RouteEntry{
			{Path: "leveling/leveling", Component: "LevelingLeveling", PageTitle: "Leveling"},
			{Path: "leveling/leaderboard", Component: "LevelingLeaderboard", PageTitle: "Leaderboard"},
			{Path: "leveling/rewards", Component: "LevelingRewards", PageTitle: "Rewards"},
			{Path: "leveling/stats", Component: "LevelingStats", PageTitle: "Stats"},
		},
	},
}
