<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let settings = null;
  let channels = [];
  let roles = [];
  let tickets = [];
  let categories = [];
  let saving = false;
  let toast = '';
  let filter = 'all';
  let categoryFilter = 'all';
  let claimFilter = 'all';
  let search = '';

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [settings, channels, roles, tickets, categories] = await Promise.all([
      api.get('/tickets/settings'),
      api.get('/channels'),
      api.get('/roles'),
      api.get('/tickets'),
      api.get('/tickets/categories'),
    ]);
    // Defensive defaults so older settings rows still bind
    settings.support_role_ids = settings.support_role_ids || [];
    settings.ping_role_ids = settings.ping_role_ids || [];
  }

  async function save() {
    saving = true;
    try {
      settings = await guildApi($currentGuildId).patch('/tickets/settings', settings);
      toast = 'Saved'; setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }

  function toggle(field, id) {
    const s = new Set(settings[field] || []);
    if (s.has(id)) s.delete(id); else s.add(id);
    settings[field] = [...s];
  }

  async function setStatus(t, status) {
    await guildApi($currentGuildId).patch(`/tickets/${t.id}`, { status });
    await load();
  }
  async function claim(t) {
    try { await guildApi($currentGuildId).post(`/tickets/${t.id}/claim`, {}); }
    catch (e) { alert(e.message); return; }
    await load();
  }
  async function unclaim(t) {
    await guildApi($currentGuildId).post(`/tickets/${t.id}/unclaim`, {});
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();

  $: shown = tickets.filter(t =>
    (filter === 'all' || t.status === filter) &&
    (categoryFilter === 'all' || String(t.category_id || '') === categoryFilter) &&
    (claimFilter === 'all'
      || (claimFilter === 'claimed' && t.claimed_by)
      || (claimFilter === 'unclaimed' && !t.claimed_by)) &&
    (!search || [t.subject, t.username].join(' ').toLowerCase().includes(search.toLowerCase()))
  );

  $: openCount = tickets.filter(t => t.status === 'open').length;
  $: closedCount = tickets.filter(t => t.status === 'closed').length;
  $: claimedCount = tickets.filter(t => t.claimed_by).length;
</script>

{#if !settings}<p>Loading…</p>{:else}
  <Panel title="Ticket settings" subtitle="Behavior shared by all panels and categories.">
    <div slot="actions">
      {#if toast}<span class="toast">{toast}</span>{/if}
      <button class="btn primary" on:click={save} disabled={saving}>Save</button>
    </div>

    <div class="grid">
      <label>Default category channel
        <select bind:value={settings.category_id}>
          <option value="">— auto —</option>
          {#each channels.filter(c => c.type === 'category') as c}<option value={c.id}>{c.name}</option>{/each}
        </select>
      </label>
      <label>Transcript channel
        <select bind:value={settings.transcript_channel_id}>
          <option value="">— do not post —</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
      <label>Naming pattern
        <input bind:value={settings.naming_pattern} placeholder="ticket-{'{user}'}" />
        <small>Variables: <code>{'{user}'}</code> <code>{'{id}'}</code> <code>{'{category}'}</code></small>
      </label>
      <label>Max open per user
        <input type="number" min="1" max="25" bind:value={settings.max_open_per_user} />
      </label>
      <label>Auto-close after (hours, 0 = off)
        <input type="number" min="0" max="720" bind:value={settings.autoclose_hours} />
      </label>
      <label>Auto-close warning at (hours)
        <input type="number" min="0" max="720" bind:value={settings.autoclose_warning_hours} />
      </label>
    </div>

    <div class="toggles">
      <label class="row"><input type="checkbox" bind:checked={settings.transcripts_on} /> Save transcripts on close</label>
      <label class="row"><input type="checkbox" bind:checked={settings.close_confirm} /> Require confirmation before closing</label>
      <label class="row"><input type="checkbox" bind:checked={settings.claim_required} /> Require staff to claim before responding</label>
      <label class="row"><input type="checkbox" bind:checked={settings.user_can_close} /> Let the ticket opener close their own ticket</label>
      <label class="row"><input type="checkbox" bind:checked={settings.use_threads} /> Use threads instead of channels</label>
    </div>

    <div class="role-block">
      <div class="label">Default support roles</div>
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={settings.support_role_ids?.includes(r.id)} on:change={() => toggle('support_role_ids', r.id)} />{r.name}</label>
      {/each}
    </div>
    <div class="role-block">
      <div class="label">Ping on open</div>
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={settings.ping_role_ids?.includes(r.id)} on:change={() => toggle('ping_role_ids', r.id)} />{r.name}</label>
      {/each}
    </div>
  </Panel>

  <Panel title="Tickets" subtitle="{openCount} open · {claimedCount} claimed · {closedCount} closed">
    <div slot="actions" class="filters">
      <input placeholder="Search subject or user" bind:value={search} />
      <select bind:value={filter}>
        <option value="all">All statuses</option>
        <option value="open">Open</option>
        <option value="closed">Closed</option>
      </select>
      <select bind:value={claimFilter}>
        <option value="all">All</option>
        <option value="claimed">Claimed</option>
        <option value="unclaimed">Unclaimed</option>
      </select>
      <select bind:value={categoryFilter}>
        <option value="all">All categories</option>
        {#each categories as c}<option value={String(c.id)}>{c.name}</option>{/each}
      </select>
    </div>
    <table>
      <thead>
        <tr>
          <th>#</th><th>User</th><th>Subject</th><th>Category</th>
          <th>Claimed</th><th>Last message</th><th>Status</th><th></th>
        </tr>
      </thead>
      <tbody>
        {#each shown as t}
          <tr>
            <td>#{t.id}</td>
            <td>{t.username}</td>
            <td>{t.subject}</td>
            <td>{t.category_name || '—'}</td>
            <td>
              {#if t.claimed_by}
                <span class="claim">{t.claimed_by_name}</span>
              {:else}
                <span class="muted">—</span>
              {/if}
            </td>
            <td><span class="muted">{t.last_message_at ? relativeTime(t.last_message_at) : '—'}</span></td>
            <td><span class={'status ' + t.status}>{t.status}</span></td>
            <td class="actions">
              {#if t.status === 'open'}
                {#if t.claimed_by}
                  <button class="btn" on:click={() => unclaim(t)}>Unclaim</button>
                {:else}
                  <button class="btn" on:click={() => claim(t)}>Claim</button>
                {/if}
                <button class="btn" on:click={() => setStatus(t, 'closed')}>Close</button>
              {:else}
                <button class="btn" on:click={() => setStatus(t, 'open')}>Reopen</button>
              {/if}
              <a class="btn link" href={`/api/guilds/${$currentGuildId}/tickets/${t.id}/transcript`} target="_blank" rel="noopener">Transcript</a>
            </td>
          </tr>
        {:else}
          <tr><td colspan="8" class="empty">No tickets match.</td></tr>
        {/each}
      </tbody>
    </table>
  </Panel>
{/if}

<style>
  .grid { display:grid; grid-template-columns:repeat(3, 1fr); gap:14px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid input, .grid select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .grid small { color:#94a3b8; font-size:11px; }
  .toggles { margin-top:18px; display:grid; grid-template-columns:repeat(2, 1fr); gap:10px; }
  .toggles .row { display:flex; align-items:center; gap:8px; font-size:13px; color:#334155; padding:8px 12px; background:#f8fafc; border-radius:8px; }
  .role-block { margin-top:18px; }
  .role-block .label { font-size:13px; color:#475569; margin-bottom:8px; font-weight:600; }
  .chip { display:inline-flex; align-items:center; gap:6px; padding:6px 12px; border:1px solid #e5e7eb; border-radius:999px; font-size:13px; margin:0 6px 6px 0; }
  .filters { display:flex; gap:8px; align-items:center; }
  .filters input, .filters select { padding:6px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .status { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; }
  .status.open { background:#dbeafe; color:#1e40af; }
  .status.closed { background:#e5e7eb; color:#475569; }
  .claim { display:inline-flex; align-items:center; gap:6px; padding:2px 8px; border-radius:999px; background:#dcfce7; color:#166534; font-size:11px; font-weight:600; }
  .muted { color:#94a3b8; }
  .empty { text-align:center; color:#94a3b8; padding:20px; }
  .actions { display:flex; gap:6px; flex-wrap:wrap; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:12px; font-weight:600; cursor:pointer; text-decoration:none; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.link { background:transparent; color:#475569; padding:6px 8px; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
</style>
