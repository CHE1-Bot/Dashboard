export const sidebarLinks = [
  { path: 'overview', icon: 'fa-chart-line', label: 'Overview', id: 'Overview' },
  { path: 'servers/servers', icon: 'fa-server', label: 'Servers', id: 'Servers' },
  { path: 'tickets/tickets', icon: 'fa-ticket', label: 'Tickets', id: 'Tickets' },
  { path: 'giveaways/giveaway', icon: 'fa-gift', label: 'Giveaways', id: 'Giveaways' },
  { path: 'moderation/moderation', icon: 'fa-shield', label: 'Moderation', id: 'Moderation' },
  { path: 'dashboard', icon: 'fa-gear', label: 'Settings', id: 'Dashboard' },
  { path: 'leveling/leveling', icon: 'fa-ranking-star', label: 'Leveling', id: 'Leveling' },
];

const overviewRightNav = [
  { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
  { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
  { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
  { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
];

const serversRightNav = [
  { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
  { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
  { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
  { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
];

const ticketsRightNav = [
  { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
  { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
  { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
  { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
  { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
  { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
  { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
  { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
];

const giveawaysRightNav = [
  { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
  { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
  { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
  { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
];

const moderationRightNav = [
  { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
  { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
  { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
  { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
];

const levelingRightNav = [
  { path: 'leveling/leveling', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
  { path: 'leveling/leaderboard', icon: 'fa-trophy', label: 'Leaderboard', id: 'Leaderboard' },
  { path: 'leveling/rewards', icon: 'fa-gift', label: 'Rewards', id: 'Rewards' },
  { path: 'leveling/stats', icon: 'fa-chart-line', label: 'Stats', id: 'Stats' },
];

export function getRoute(hash) {
  const value = hash.startsWith('#/') ? hash.slice(2) : hash.startsWith('#') ? hash.slice(1) : hash;
  return value || 'dashboard';
}

export const routes = {
  dashboard: {
    componentName: 'Dashboard',
    pageTitle: 'Dashboard',
    sidebarActive: 'Dashboard',
    contentHtml: '',
    rightNav: [],
  },
  overview: {
    componentName: 'OverviewOverview',
    pageTitle: 'Overview',
    sidebarActive: 'Overview',
    contentHtml: '',
    rightNav: overviewRightNav,
  },
  'overview/statistics': {
    componentName: 'OverviewStatistics',
    pageTitle: 'Statistics',
    sidebarActive: 'Overview',
    contentHtml: '',
    rightNav: overviewRightNav,
  },
  'overview/settings': {
    componentName: 'OverviewSettings',
    pageTitle: 'Settings',
    sidebarActive: 'Overview',
    contentHtml: '',
    rightNav: overviewRightNav,
  },
  'overview/alerts': {
    componentName: 'OverviewAlerts',
    pageTitle: 'Alerts',
    sidebarActive: 'Overview',
    contentHtml: '',
    rightNav: overviewRightNav,
  },
  'overview/history': {
    componentName: 'OverviewHistory',
    pageTitle: 'History',
    sidebarActive: 'Overview',
    contentHtml: '',
    rightNav: overviewRightNav,
  },
  'servers/servers': {
    componentName: 'ServersServers',
    pageTitle: 'Servers',
    sidebarActive: 'Servers',
    contentHtml: '',
    rightNav: serversRightNav,
  },
  'servers/add-server': {
    componentName: 'ServersAdd_server',
    pageTitle: 'Add Server',
    sidebarActive: 'Servers',
    contentHtml: '',
    rightNav: serversRightNav,
  },
  'servers/backup': {
    componentName: 'ServersBackup',
    pageTitle: 'Backup',
    sidebarActive: 'Servers',
    contentHtml: '',
    rightNav: serversRightNav,
  },
  'servers/manage-roles': {
    componentName: 'ServersManage_roles',
    pageTitle: 'Manage Roles',
    sidebarActive: 'Servers',
    contentHtml: '',
    rightNav: serversRightNav,
  },
  'servers/permissions': {
    componentName: 'ServersPermissions',
    pageTitle: 'Permissions',
    sidebarActive: 'Servers',
    contentHtml: '',
    rightNav: serversRightNav,
  },
  'tickets/tickets': {
    componentName: 'TicketsTickets',
    pageTitle: 'Tickets',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/audit-logs': {
    componentName: 'TicketsAudit_logs',
    pageTitle: 'Audit Logs',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/forms': {
    componentName: 'TicketsForms',
    pageTitle: 'Forms',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/panels': {
    componentName: 'TicketsPanels',
    pageTitle: 'Panels',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/tags': {
    componentName: 'TicketsTags',
    pageTitle: 'Tags',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/ticket-embed': {
    componentName: 'TicketsTicket_embed',
    pageTitle: 'Ticket Embed',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/ticket-staff': {
    componentName: 'TicketsTicket_staff',
    pageTitle: 'Ticket Staff',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'tickets/transcripts': {
    componentName: 'TicketsTranscripts',
    pageTitle: 'Transcripts',
    sidebarActive: 'Tickets',
    contentHtml: '',
    rightNav: ticketsRightNav,
  },
  'giveaways/giveaway': {
    componentName: 'GiveawaysGiveaway',
    pageTitle: 'Giveaways',
    sidebarActive: 'Giveaways',
    contentHtml: '',
    rightNav: giveawaysRightNav,
  },
  'giveaways/active-giveaways': {
    componentName: 'GiveawaysActive_giveaways',
    pageTitle: 'Active Giveaways',
    sidebarActive: 'Giveaways',
    contentHtml: '',
    rightNav: giveawaysRightNav,
  },
  'giveaways/blacklist': {
    componentName: 'GiveawaysBlacklist',
    pageTitle: 'Blacklist',
    sidebarActive: 'Giveaways',
    contentHtml: '',
    rightNav: giveawaysRightNav,
  },
  'giveaways/create-giveaway': {
    componentName: 'GiveawaysCreate_giveaway',
    pageTitle: 'Create Giveaway',
    sidebarActive: 'Giveaways',
    contentHtml: '',
    rightNav: giveawaysRightNav,
  },
  'giveaways/premium': {
    componentName: 'GiveawaysPremium',
    pageTitle: 'Premium',
    sidebarActive: 'Giveaways',
    contentHtml: '',
    rightNav: giveawaysRightNav,
  },
  'moderation/moderation': {
    componentName: 'ModerationModeration',
    pageTitle: 'Moderation',
    sidebarActive: 'Moderation',
    contentHtml: '',
    rightNav: moderationRightNav,
  },
  'moderation/auto-mod': {
    componentName: 'ModerationAuto_mod',
    pageTitle: 'Auto Mod',
    sidebarActive: 'Moderation',
    contentHtml: '',
    rightNav: moderationRightNav,
  },
  'moderation/manual-actions': {
    componentName: 'ModerationManual_actions',
    pageTitle: 'Manual Actions',
    sidebarActive: 'Moderation',
    contentHtml: '',
    rightNav: moderationRightNav,
  },
  'moderation/reports': {
    componentName: 'ModerationReports',
    pageTitle: 'Reports',
    sidebarActive: 'Moderation',
    contentHtml: '',
    rightNav: moderationRightNav,
  },
  'moderation/logs': {
    componentName: 'ModerationLogs',
    pageTitle: 'Logs',
    sidebarActive: 'Moderation',
    contentHtml: '',
    rightNav: moderationRightNav,
  },
  'leveling/leveling': {
    componentName: 'LevelingLeveling',
    pageTitle: 'Leveling',
    sidebarActive: 'Leveling',
    contentHtml: '',
    rightNav: levelingRightNav,
  },
  'leveling/leaderboard': {
    componentName: 'LevelingLeaderboard',
    pageTitle: 'Leaderboard',
    sidebarActive: 'Leveling',
    contentHtml: '',
    rightNav: levelingRightNav,
  },
  'leveling/rewards': {
    componentName: 'LevelingRewards',
    pageTitle: 'Rewards',
    sidebarActive: 'Leveling',
    contentHtml: '',
    rightNav: levelingRightNav,
  },
  'leveling/stats': {
    componentName: 'LevelingStats',
    pageTitle: 'Stats',
    sidebarActive: 'Leveling',
    contentHtml: '',
    rightNav: levelingRightNav,
  },
};
