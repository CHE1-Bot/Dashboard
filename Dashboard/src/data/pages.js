function checkboxItems(items, startIndex = 1) {
  return items
    .map(
      (text, index) => `
        <div class="ticket-overview" id="check-button${startIndex + index}">
          <h3 class="accordion1-texts">${text}</h3>
          <input type="checkbox" id="check${startIndex + index}">
          <label for="check${startIndex + index}" class="button${startIndex + index}"></label>
        </div>`
    )
    .join('');
}

function accordionSection(number, title, bodyHtml) {
  return `
    <div class="accordion${number}">
      <div class="accordion-item${number}">
        <div class="accordion-header${number}">
          <h3>${title}</h3>
        </div>
        <div class="accordion-content${number}" id="content${number}">
          ${bodyHtml}
        </div>
      </div>
    </div>
  `;
}

function settingsPage(title, sections) {
  return sections.map((section, index) => accordionSection(index + 1, section.header, section.body)).join('');
}

const sidebarLinks = [
  { path: 'overview', icon: 'fa-chart-line', label: 'Overview', id: 'Overview' },
  { path: 'servers/servers', icon: 'fa-server', label: 'Servers', id: 'Servers' },
  { path: 'tickets/tickets', icon: 'fa-ticket', label: 'Tickets', id: 'Tickets' },
  { path: 'giveaways/giveaway', icon: 'fa-gift', label: 'Giveaways', id: 'Giveaways' },
  { path: 'moderation/moderation', icon: 'fa-shield', label: 'Moderation', id: 'Moderation' },
  { path: 'dashboard', icon: 'fa-gear', label: 'Settings', id: 'Dashboard' },
  { path: 'leveling/leveling', icon: 'fa-ranking-star', label: 'Leveling', id: 'Leveling' },
];

function rightNav(items) {
  return items.map(item => ({ path: item.path, icon: item.icon, label: item.label, id: item.id }));
}

function topicsPage(header, texts, startIndex) {
  return accordionSection(1, header, checkboxItems(texts, startIndex));
}

export function getRoute(hash) {
  const value = hash.startsWith('#/') ? hash.slice(2) : hash.startsWith('#') ? hash.slice(1) : hash;
  return value || 'dashboard';
}

export const routes = {
  dashboard: {
    pageTitle: 'Dashboard',
    sidebarActive: 'Dashboard',
    contentHtml: settingsPage('Dashboard', [
      {
        header: 'Quick Actions',
        body: checkboxItems(['Enable All Features', 'Auto Setup', 'Backup Data', 'Reset Settings'], 1),
      },
      {
        header: 'System Status',
        body: `
          <div class="ticket-overview" id="check-button5">
            <h3 class="accordion1-texts">Bot Online</h3>
            <input type="checkbox" id="check5" checked disabled>
            <label for="check5" class="button5"></label>
          </div>
          <div class="ticket-overview" id="check-button6">
            <h3 class="accordion1-texts">Database Connected</h3>
            <input type="checkbox" id="check6" checked disabled>
            <label for="check6" class="button6"></label>
          </div>
          <div class="ticket-overview" id="check-button7">
            <h3 class="accordion1-texts">API Active</h3>
            <input type="checkbox" id="check7" checked disabled>
            <label for="check7" class="button7"></label>
          </div>
          <div class="ticket-overview" id="check-button8">
            <h3 class="accordion1-texts">All Systems Go</h3>
            <input type="checkbox" id="check8" checked disabled>
            <label for="check8" class="button8"></label>
          </div>
        `,
      },
    ]),
    rightNav: [],
  },

  overview: {
    pageTitle: 'Overview',
    sidebarActive: 'Overview',
    contentHtml: settingsPage('Overview', [
      {
        header: 'Statistics',
        body: checkboxItems(['Show Server Stats', 'Display User Count', 'Enable Activity Logs', 'Show Bot Uptime'], 1),
      },
      {
        header: 'Dashboard Settings',
        body: checkboxItems(['Auto Refresh', 'Dark Mode', 'Notifications', 'Compact View'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
      { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
      { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
    ]),
  },

  'overview/alerts': {
    pageTitle: 'Alerts',
    sidebarActive: 'Overview',
    contentHtml: settingsPage('Alerts', [
      {
        header: 'Notification Settings',
        body: checkboxItems(['Email Alerts', 'Discord Notifications', 'Push Notifications', 'SMS Alerts'], 1),
      },
      {
        header: 'Alert Types',
        body: checkboxItems(['Server Issues', 'Security Alerts', 'User Reports', 'System Updates'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
      { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
      { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
    ]),
  },

  'overview/history': {
    pageTitle: 'History',
    sidebarActive: 'Overview',
    contentHtml: settingsPage('History', [
      {
        header: 'Recent Activities',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Activity Log</h3>
            <p>View recent server activities and changes.</p>
            <button class="btn">View Log</button>
          </div>
        `,
      },
      {
        header: 'Event History',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Past Events</h3>
            <p>Review historical events and their outcomes.</p>
            <button class="btn">View Events</button>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
      { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
      { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
    ]),
  },

  'overview/settings': {
    pageTitle: 'Settings',
    sidebarActive: 'Overview',
    contentHtml: settingsPage('Settings', [
      {
        header: 'General Settings',
        body: `
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Enable Dashboard</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Auto Refresh</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
        `,
      },
      {
        header: 'Privacy Settings',
        body: `
          <div class="ticket-overview" id="check-button3">
            <h3 class="accordion1-texts">Share Data</h3>
            <input type="checkbox" id="check3">
            <label for="check3" class="button3"></label>
          </div>
          <div class="ticket-overview" id="check-button4">
            <h3 class="accordion1-texts">Analytics</h3>
            <input type="checkbox" id="check4">
            <label for="check4" class="button4"></label>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
      { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
      { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
    ]),
  },

  'overview/statistics': {
    pageTitle: 'Statistics',
    sidebarActive: 'Overview',
    contentHtml: settingsPage('Statistics', [
      {
        header: 'Server Stats',
        body: checkboxItems(['Member Count', 'Channel Count', 'Role Count', 'Bot Uptime'], 1),
      },
      {
        header: 'Activity Stats',
        body: checkboxItems(['Messages Today', 'Active Users', 'Tickets Created', 'Giveaways Won'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'overview/statistics', icon: 'fa-chart-bar', label: 'Statistics', id: 'Statistics' },
      { path: 'overview/settings', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'overview/alerts', icon: 'fa-bell', label: 'Alerts', id: 'Alerts' },
      { path: 'overview/history', icon: 'fa-history', label: 'History', id: 'History' },
    ]),
  },

  'servers/servers': {
    pageTitle: 'Servers',
    sidebarActive: 'Servers',
    contentHtml: settingsPage('Servers', [
      {
        header: 'Server Management',
        body: checkboxItems(['Auto Join', 'Welcome Messages', 'Server Logs', 'Backup Settings'], 1),
      },
      {
        header: 'Permissions',
        body: checkboxItems(['Role Management', 'Channel Permissions', 'Bot Commands', 'Admin Only'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
      { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
      { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
      { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
    ]),
  },

  'servers/add-server': {
    pageTitle: 'Add Server',
    sidebarActive: 'Servers',
    contentHtml: settingsPage('Add Server', [
      {
        header: 'Server Details',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Server Name</h3>
            <input type="text" placeholder="Enter server name">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Server ID</h3>
            <input type="text" placeholder="Enter server ID">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Invite Link</h3>
            <input type="url" placeholder="Enter invite link">
          </div>
        `,
      },
      {
        header: 'Configuration',
        body: `
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Enable Bot</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Auto Setup</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
          <button class="btn">Add Server</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
      { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
      { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
      { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
    ]),
  },

  'servers/backup': {
    pageTitle: 'Backup',
    sidebarActive: 'Servers',
    contentHtml: settingsPage('Backup', [
      {
        header: 'Create Backup',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Backup Name</h3>
            <input type="text" placeholder="Enter backup name">
          </div>
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Include Roles</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Include Channels</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
          <button class="btn">Create Backup</button>
        `,
      },
      {
        header: 'Restore Backup',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Select Backup</h3>
            <select>
              <option>Backup 1 - 2023-10-01</option>
              <option>Backup 2 - 2023-09-15</option>
            </select>
          </div>
          <button class="btn">Restore Backup</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
      { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
      { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
      { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
    ]),
  },

  'servers/manage-roles': {
    pageTitle: 'Manage Roles',
    sidebarActive: 'Servers',
    contentHtml: settingsPage('Manage Roles', [
      {
        header: 'Create Role',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Role Name</h3>
            <input type="text" placeholder="Enter role name">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Color</h3>
            <input type="color">
          </div>
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Mentionable</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <button class="btn">Create Role</button>
        `,
      },
      {
        header: 'Edit Roles',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Select Role</h3>
            <select>
              <option>Admin</option>
              <option>Moderator</option>
              <option>Member</option>
            </select>
          </div>
          <button class="btn">Edit Role</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
      { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
      { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
      { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
    ]),
  },

  'servers/permissions': {
    pageTitle: 'Permissions',
    sidebarActive: 'Servers',
    contentHtml: settingsPage('Permissions', [
      {
        header: 'Role Permissions',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Select Role</h3>
            <select>
              <option>Admin</option>
              <option>Moderator</option>
              <option>Member</option>
            </select>
          </div>
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Manage Server</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Kick Members</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
          <div class="ticket-overview" id="check-button3">
            <h3 class="accordion1-texts">Ban Members</h3>
            <input type="checkbox" id="check3">
            <label for="check3" class="button3"></label>
          </div>
        `,
      },
      {
        header: 'Channel Permissions',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Select Channel</h3>
            <select>
              <option>#general</option>
              <option>#announcements</option>
              <option>#moderation</option>
            </select>
          </div>
          <div class="ticket-overview" id="check-button4">
            <h3 class="accordion1-texts">Read Messages</h3>
            <input type="checkbox" id="check4">
            <label for="check4" class="button4"></label>
          </div>
          <div class="ticket-overview" id="check-button5">
            <h3 class="accordion1-texts">Send Messages</h3>
            <input type="checkbox" id="check5">
            <label for="check5" class="button5"></label>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'servers/add-server', icon: 'fa-plus', label: 'Add Server', id: 'AddServer' },
      { path: 'servers/manage-roles', icon: 'fa-users-cog', label: 'Manage Roles', id: 'ManageRoles' },
      { path: 'servers/permissions', icon: 'fa-key', label: 'Permissions', id: 'Permissions' },
      { path: 'servers/backup', icon: 'fa-save', label: 'Backup', id: 'Backup' },
    ]),
  },

  'tickets/tickets': {
    pageTitle: 'Tickets',
    sidebarActive: 'Tickets',
    contentHtml: `
      ${accordionSection(1, 'General', checkboxItems(['Enable /open', 'Allow user to close tickets', 'Enable Rating', 'Confirm Close'], 1))}
      ${accordionSection(2, 'Close', `
        <div class="ticket-overview" id="check-button5">
          <h3>Allow user to close tickets</h3>
          <input type="checkbox" id="check5">
          <label for="check5" class="button5"></label>
        </div>
        <div class="ticket-overview" id="check-button6">
          <input type="checkbox" id="check6">
          <label for="check6" class="button6"></label>
        </div>
        <div class="ticket-overview" id="check-button7">
          <input type="checkbox" id="check7">
          <label for="check7" class="button7"></label>
        </div>
        <div class="ticket-overview" id="check-button8">
          <input type="checkbox" id="check8">
          <label for="check8" class="button8"></label>
        </div>
      `)}
      ${accordionSection(3, 'Ticket Auto Close', `
        <div class="ticket-overview" id="check-button9">
          <input type="checkbox" id="check9">
          <label for="check9" class="button9"></label>
        </div>
        <div class="ticket-overview" id="check-button10">
          <input type="checkbox" id="check10">
          <label for="check10" class="button10"></label>
        </div>
        <div class="ticket-overview" id="check-button11">
          <input type="checkbox" id="check11">
          <label for="check11" class="button11"></label>
        </div>
        <div class="ticket-overview" id="check-button12">
          <input type="checkbox" id="check12">
          <label for="check12" class="button12"></label>
        </div>
      `)}
      ${accordionSection(4, 'Transcript', `
        <div class="dropdown">
          <button class="dropdown-button">
            Transcript Channel
            <span class="arrow">▼</span>
          </button>
          <div class="dropdown-content">
            <a href="#">1</a>
            <a href="#">2</a>
            <a href="#">3</a>
          </div>
        </div>
        <div class="ticket-overview" id="check-button13">
          <input type="checkbox" id="check13">
          <label for="check13" class="button13"></label>
        </div>
        <div class="ticket-overview" id="check-button14">
          <input type="checkbox" id="check14">
          <label for="check14" class="button14"></label>
        </div>
        <div class="ticket-overview" id="check-button15">
          <input type="checkbox" id="check15">
          <label for="check15" class="button15"></label>
        </div>
        <div class="ticket-overview" id="check-button16">
          <input type="checkbox" id="check16">
          <label for="check16" class="button16"></label>
        </div>
      `)}`,
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/audit-logs': {
    pageTitle: 'Audit Logs',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Audit Logs', [
      {
        header: 'Log Settings',
        body: checkboxItems(['Enable Audit Logs', 'Log User Actions', 'Log Moderation', 'Log Server Changes'], 1),
      },
      {
        header: 'Log Channels',
        body: checkboxItems(['Separate Channels', 'Embed Format', 'Timestamp Logs', 'Archive Old Logs'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/forms': {
    pageTitle: 'Forms',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Forms', [
      {
        header: 'Form Settings',
        body: checkboxItems(['Enable Forms', 'Require Email', 'Custom Fields', 'Form Validation'], 1),
      },
      {
        header: 'Submission',
        body: checkboxItems(['Auto Respond', 'Store Responses', 'Notify Staff', 'Export Data'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/panels': {
    pageTitle: 'Panels',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Panels', [
      {
        header: 'Panel Settings',
        body: checkboxItems(['Enable Panels', 'Custom Colors', 'Panel Layout', 'Interactive Elements'], 1),
      },
      {
        header: 'Panel Management',
        body: checkboxItems(['Create Panels', 'Edit Panels', 'Delete Panels', 'Panel Permissions'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/tags': {
    pageTitle: 'Tags',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Tags', [
      {
        header: 'Tag Settings',
        body: checkboxItems(['Enable Tags', 'Custom Tags', 'Tag Colors', 'Tag Permissions'], 1),
      },
      {
        header: 'Tag Management',
        body: checkboxItems(['Create Tags', 'Edit Tags', 'Delete Tags', 'Tag Categories'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/ticket-embed': {
    pageTitle: 'Ticket Embed',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Ticket Embed', [
      {
        header: 'Embed Settings',
        body: checkboxItems(['Enable Embed', 'Custom Title', 'Embed Color', 'Embed Footer'], 1),
      },
      {
        header: 'Embed Content',
        body: checkboxItems(['Description', 'Fields', 'Thumbnail', 'Author'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/ticket-staff': {
    pageTitle: 'Ticket Staff',
    sidebarActive: 'Tickets',
    contentHtml: settingsPage('Ticket Staff', [
      {
        header: 'Staff Roles',
        body: checkboxItems(['Admin Role', 'Moderator Role', 'Support Role', 'Helper Role'], 1),
      },
      {
        header: 'Staff Permissions',
        body: checkboxItems(['Close Tickets', 'Transfer Tickets', 'Add Notes', 'View Transcripts'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'tickets/transcripts': {
    pageTitle: 'Transcripts',
    sidebarActive: 'Tickets',
    contentHtml: `
      <div class="transcripts">
        <div class="search-section">
          <div class="dropdown">
            <button class="dropdown-button">
              Sorting
              <span class="arrow">▼</span>
            </button>
            <div class="dropdown-content">
              <a href="#" data-sort="newest">Newest First</a>
              <a href="#" data-sort="oldest">Oldest First</a>
              <a href="#" data-sort="id">By ID</a>
            </div>
          </div>
          <div class="search-input">
            <input type="text" id="ticket-id" class="ticket-id" placeholder="Enter ticket ID">
            <button id="search-btn" class="search-btn">Search</button>
          </div>
        </div>
        <div class="transcripts-list" id="transcripts-list">
          <p>Use the search bar to find transcripts.</p>
        </div>
      </div>
    `,
    rightNav: rightNav([
      { path: 'tickets/tickets', icon: 'fa-gears', label: 'Settings', id: 'Settings' },
      { path: 'tickets/ticket-embed', icon: 'fa-regular fa-message', label: 'Embed', id: 'TicketEmbed' },
      { path: 'tickets/ticket-staff', icon: 'fa-people-group', label: 'Ticket Staff', id: 'TicketStaff' },
      { path: 'tickets/transcripts', icon: 'fa-file', label: 'Transcripts', id: 'Transcripts' },
      { path: 'tickets/panels', icon: 'fa-arrow-pointer', label: 'Panels', id: 'Panels' },
      { path: 'tickets/forms', icon: 'fa-file-lines', label: 'Forms', id: 'Forms' },
      { path: 'tickets/tags', icon: 'fa-tag', label: 'Tags', id: 'Tags' },
      { path: 'tickets/audit-logs', icon: 'fa-file-medical', label: 'Audit logs', id: 'AuditLogs' },
    ]),
  },

  'giveaways/giveaway': {
    pageTitle: 'Giveaways',
    sidebarActive: 'Giveaways',
    contentHtml: settingsPage('Giveaways', [
      {
        header: 'General Settings',
        body: checkboxItems(['Enable Giveaways', 'Require Role', 'Auto Announce', 'DM Winners'], 1),
      },
      {
        header: 'Advanced Options',
        body: checkboxItems(['Multiple Winners', 'Blacklist Roles', 'Sponsor Mode', 'Premium Features'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
      { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
      { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
      { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
    ]),
  },

  'giveaways/active-giveaways': {
    pageTitle: 'Active Giveaways',
    sidebarActive: 'Giveaways',
    contentHtml: settingsPage('Active Giveaways', [
      {
        header: 'Current Giveaways',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Giveaway #1</h3>
            <p>Prize: Discord Nitro</p>
            <p>Ends: 2023-10-05</p>
            <p>Participants: 150</p>
            <button class="btn">End Early</button>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Giveaway #2</h3>
            <p>Prize: $50 Gift Card</p>
            <p>Ends: 2023-10-10</p>
            <p>Participants: 75</p>
            <button class="btn">End Early</button>
          </div>
        `,
      },
      {
        header: 'Manage Giveaways',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Reroll Winners</h3>
            <select>
              <option>Giveaway #1</option>
              <option>Giveaway #2</option>
            </select>
            <button class="btn">Reroll</button>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Delete Giveaway</h3>
            <select>
              <option>Giveaway #1</option>
              <option>Giveaway #2</option>
            </select>
            <button class="btn">Delete</button>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
      { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
      { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
      { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
    ]),
  },

  'giveaways/blacklist': {
    pageTitle: 'Blacklist',
    sidebarActive: 'Giveaways',
    contentHtml: settingsPage('Blacklist', [
      {
        header: 'Add to Blacklist',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">User ID</h3>
            <input type="text" placeholder="Enter user ID">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Reason</h3>
            <textarea placeholder="Reason for blacklisting"></textarea>
          </div>
          <button class="btn">Add to Blacklist</button>
        `,
      },
      {
        header: 'Current Blacklist',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Blacklisted Users</h3>
            <ul>
              <li>User123 - Spam</li>
              <li>User456 - Cheating</li>
              <li>User789 - Harassment</li>
            </ul>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Remove User</h3>
            <input type="text" placeholder="User ID to remove">
            <button class="btn">Remove</button>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
      { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
      { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
      { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
    ]),
  },

  'giveaways/create-giveaway': {
    pageTitle: 'Create Giveaway',
    sidebarActive: 'Giveaways',
    contentHtml: settingsPage('Create Giveaway', [
      {
        header: 'Giveaway Details',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Prize</h3>
            <input type="text" placeholder="Enter prize description">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Duration</h3>
            <input type="number" placeholder="Hours" min="1">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Winners</h3>
            <input type="number" value="1" min="1">
          </div>
        `,
      },
      {
        header: 'Requirements',
        body: `
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Require Role</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Required Role</h3>
            <input type="text" placeholder="Role name">
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Blacklist Previous Winners</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
          <button class="btn">Create Giveaway</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
      { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
      { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
      { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
    ]),
  },

  'giveaways/premium': {
    pageTitle: 'Premium',
    sidebarActive: 'Giveaways',
    contentHtml: settingsPage('Premium', [
      {
        header: 'Premium Features',
        body: `
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Unlimited Giveaways</h3>
            <input type="checkbox" id="check1" checked disabled>
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Custom Emojis</h3>
            <input type="checkbox" id="check2" checked disabled>
            <label for="check2" class="button2"></label>
          </div>
          <div class="ticket-overview" id="check-button3">
            <h3 class="accordion1-texts">Priority Support</h3>
            <input type="checkbox" id="check3" checked disabled>
            <label for="check3" class="button3"></label>
          </div>
        `,
      },
      {
        header: 'Subscription',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Current Plan</h3>
            <p>Premium Monthly</p>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Expires</h3>
            <p>2023-11-01</p>
          </div>
          <button class="btn">Renew Subscription</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'giveaways/create-giveaway', icon: 'fa-plus-circle', label: 'Create Giveaway', id: 'CreateGiveaway' },
      { path: 'giveaways/active-giveaways', icon: 'fa-list', label: 'Active Giveaways', id: 'ActiveGiveaways' },
      { path: 'giveaways/premium', icon: 'fa-star', label: 'Premium', id: 'Premium' },
      { path: 'giveaways/blacklist', icon: 'fa-ban', label: 'Blacklist', id: 'Blacklist' },
    ]),
  },

  'moderation/moderation': {
    pageTitle: 'Moderation',
    sidebarActive: 'Moderation',
    contentHtml: settingsPage('Moderation', [
      {
        header: 'Auto Moderation',
        body: checkboxItems(['Spam Filter', 'Link Blocker', 'Caps Lock Detection', 'Word Filter'], 1),
      },
      {
        header: 'Manual Moderation',
        body: checkboxItems(['Kick Users', 'Ban Users', 'Mute Users', 'Warn System'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
      { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
      { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
      { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
    ]),
  },

  'moderation/auto-mod': {
    pageTitle: 'Auto Mod',
    sidebarActive: 'Moderation',
    contentHtml: settingsPage('Auto Mod', [
      {
        header: 'Spam Detection',
        body: `
          <div class="ticket-overview" id="check-button1">
            <h3 class="accordion1-texts">Enable Spam Filter</h3>
            <input type="checkbox" id="check1">
            <label for="check1" class="button1"></label>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Spam Threshold</h3>
            <input type="number" min="1" max="10" value="5">
          </div>
        `,
      },
      {
        header: 'Word Filters',
        body: `
          <div class="ticket-overview" id="check-button2">
            <h3 class="accordion1-texts">Enable Word Filter</h3>
            <input type="checkbox" id="check2">
            <label for="check2" class="button2"></label>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Banned Words</h3>
            <textarea placeholder="Enter banned words, one per line"></textarea>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
      { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
      { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
      { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
    ]),
  },

  'moderation/manual-actions': {
    pageTitle: 'Manual Actions',
    sidebarActive: 'Moderation',
    contentHtml: settingsPage('Manual Actions', [
      {
        header: 'Ban User',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">User ID</h3>
            <input type="text" placeholder="Enter user ID">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Reason</h3>
            <textarea placeholder="Reason for ban"></textarea>
          </div>
          <button class="btn">Ban User</button>
        `,
      },
      {
        header: 'Kick User',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">User ID</h3>
            <input type="text" placeholder="Enter user ID">
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Reason</h3>
            <textarea placeholder="Reason for kick"></textarea>
          </div>
          <button class="btn">Kick User</button>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
      { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
      { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
      { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
    ]),
  },

  'moderation/reports': {
    pageTitle: 'Reports',
    sidebarActive: 'Moderation',
    contentHtml: settingsPage('Reports', [
      {
        header: 'Open Reports',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Report #1</h3>
            <p>User reported for spam.</p>
            <button class="btn">Review</button>
          </div>
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Report #2</h3>
            <p>User reported for harassment.</p>
            <button class="btn">Review</button>
          </div>
        `,
      },
      {
        header: 'Closed Reports',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Report #3</h3>
            <p>Resolved: User warned.</p>
            <button class="btn">View Details</button>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
      { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
      { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
      { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
    ]),
  },

  'moderation/logs': {
    pageTitle: 'Logs',
    sidebarActive: 'Moderation',
    contentHtml: settingsPage('Logs', [
      {
        header: 'Moderation Logs',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Recent Actions</h3>
            <ul>
              <li>User banned: Spam violation</li>
              <li>User kicked: Harassment</li>
              <li>Message deleted: Inappropriate content</li>
            </ul>
          </div>
        `,
      },
      {
        header: 'System Logs',
        body: `
          <div class="ticket-overview">
            <h3 class="accordion1-texts">Bot Activity</h3>
            <ul>
              <li>Bot started: 2023-10-01</li>
              <li>Auto mod triggered: 5 times</li>
              <li>Server joined: New server</li>
            </ul>
          </div>
        `,
      },
    ]),
    rightNav: rightNav([
      { path: 'moderation/auto-mod', icon: 'fa-robot', label: 'Auto Mod', id: 'AutoMod' },
      { path: 'moderation/manual-actions', icon: 'fa-hand-paper', label: 'Manual Actions', id: 'ManualActions' },
      { path: 'moderation/reports', icon: 'fa-flag', label: 'Reports', id: 'Reports' },
      { path: 'moderation/logs', icon: 'fa-file-alt', label: 'Logs', id: 'Logs' },
    ]),
  },

  'leveling/leveling': {
    pageTitle: 'Leveling',
    sidebarActive: 'Leveling',
    contentHtml: settingsPage('Leveling', [
      {
        header: 'Leveling System',
        body: checkboxItems(['Enable Leveling', 'XP per Message', 'Level Up Messages', 'Role Rewards'], 1),
      },
      {
        header: 'Advanced Settings',
        body: checkboxItems(['Multiplier Channels', 'Ignore Bots', 'Leaderboard', 'Reset Levels'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'leveling/leveling', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'leveling/leaderboard', icon: 'fa-trophy', label: 'Leaderboard', id: 'Leaderboard' },
      { path: 'leveling/rewards', icon: 'fa-gift', label: 'Rewards', id: 'Rewards' },
      { path: 'leveling/stats', icon: 'fa-chart-line', label: 'Stats', id: 'Stats' },
    ]),
  },

  'leveling/leaderboard': {
    pageTitle: 'Leaderboard',
    sidebarActive: 'Leveling',
    contentHtml: settingsPage('Leaderboard', [
      {
        header: 'Top Users',
        body: checkboxItems(['Display Top 10', 'Show Levels', 'Include XP', 'Update Real-time'], 1),
      },
      {
        header: 'Filters',
        body: checkboxItems(['Global Leaderboard', 'Server Specific', 'Monthly Reset', 'Weekly Reset'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'leveling/leveling', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'leveling/leaderboard', icon: 'fa-trophy', label: 'Leaderboard', id: 'Leaderboard' },
      { path: 'leveling/rewards', icon: 'fa-gift', label: 'Rewards', id: 'Rewards' },
      { path: 'leveling/stats', icon: 'fa-chart-line', label: 'Stats', id: 'Stats' },
    ]),
  },

  'leveling/rewards': {
    pageTitle: 'Rewards',
    sidebarActive: 'Leveling',
    contentHtml: settingsPage('Rewards', [
      {
        header: 'Role Rewards',
        body: checkboxItems(['Enable Role Rewards', 'Auto Assign', 'Stack Roles', 'Remove Old Roles'], 1),
      },
      {
        header: 'Custom Rewards',
        body: checkboxItems(['XP Multipliers', 'Special Titles', 'Emote Access', 'Channel Access'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'leveling/leveling', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'leveling/leaderboard', icon: 'fa-trophy', label: 'Leaderboard', id: 'Leaderboard' },
      { path: 'leveling/rewards', icon: 'fa-gift', label: 'Rewards', id: 'Rewards' },
      { path: 'leveling/stats', icon: 'fa-chart-line', label: 'Stats', id: 'Stats' },
    ]),
  },

  'leveling/stats': {
    pageTitle: 'Stats',
    sidebarActive: 'Leveling',
    contentHtml: settingsPage('Leveling Stats', [
      {
        header: 'Leveling System',
        body: checkboxItems(['Enable Leveling', 'XP per Message', 'Level Up Messages', 'Role Rewards'], 1),
      },
      {
        header: 'Advanced Settings',
        body: checkboxItems(['Multiplier Channels', 'Ignore Bots', 'Leaderboard', 'Reset Levels'], 5),
      },
    ]),
    rightNav: rightNav([
      { path: 'leveling/leveling', icon: 'fa-cog', label: 'Settings', id: 'Settings' },
      { path: 'leveling/leaderboard', icon: 'fa-trophy', label: 'Leaderboard', id: 'Leaderboard' },
      { path: 'leveling/rewards', icon: 'fa-gift', label: 'Rewards', id: 'Rewards' },
      { path: 'leveling/stats', icon: 'fa-chart-line', label: 'Stats', id: 'Stats' },
    ]),
  },
};

export { sidebarLinks };
