<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let settings = null;
  let channels = [];
  let roles = [];
  let tickets = [];
  let saving = false;
  let toast = '';
  let filter = 'all';

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [settings, channels, roles, tickets] = await Promise.all([
      api.get('/tickets/settings'), api.get('/channels'), api.get('/roles'), api.get('/tickets'),
    ]);
  }
  async function save() {
    saving = true;
    try {
      settings = await guildApi($currentGuildId).patch('/tickets/settings', settings);
      toast = 'Saved'; setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }
  function toggleRole(id) {
    const s = new Set(settings.support_role_ids || []);
    if (s.has(id)) s.delete(id); else s.add(id);
    settings.support_role_ids = [...s];
  }
  async function setStatus(t, status) {
    await guildApi($currentGuildId).patch(`/tickets/${t.id}`, { status });
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: shown = filter === 'all' ? tickets : tickets.filter(t => t.status === filter);
</script>

{#if !settings}<p>Loading…</p>{:else}
  <Panel title="Ticket settings">
    <div slot="actions">
      {#if toast}<span class="toast">{toast}</span>{/if}
      <button class="btn primary" on:click={save} disabled={saving}>Save</button>
    </div>
    <div class="grid">
      <label>Category channel
        <select bind:value={settings.category_id}>
          <option value="">— auto —</option>
          {#each channels.filter(c => c.type === 'category') as c}<option value={c.id}>{c.name}</option>{/each}
        </select>
      </label>
      <label>Max open per user
        <input type="number" min="1" max="10" bind:value={settings.max_open_per_user} />
      </label>
      <label>Naming pattern
        <input bind:value={settings.naming_pattern} placeholder="ticket-{'{user}'}" />
      </label>
      <label class="row"><input type="checkbox" bind:checked={settings.transcripts_on} /> Save transcripts when closing</label>
      <label class="row"><input type="checkbox" bind:checked={settings.close_confirm} /> Require close confirmation</label>
    </div>
    <div class="roles">
      <div class="label">Support roles</div>
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={settings.support_role_ids?.includes(r.id)} on:change={() => toggleRole(r.id)} />{r.name}</label>
      {/each}
    </div>
  </Panel>

  <Panel title="Open tickets" subtitle="{tickets.filter(t=>t.status==='open').length} open · {tickets.filter(t=>t.status==='closed').length} closed">
    <div slot="actions">
      <select bind:value={filter}>
        <option value="all">All</option>
        <option value="open">Open</option>
        <option value="closed">Closed</option>
      </select>
    </div>
    <table>
      <thead><tr><th>#</th><th>User</th><th>Subject</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {#each shown as t}
          <tr>
            <td>#{t.id}</td>
            <td>{t.username}</td>
            <td>{t.subject}</td>
            <td><span class={'status ' + t.status}>{t.status}</span></td>
            <td>
              {#if t.status === 'open'}<button class="btn" on:click={() => setStatus(t, 'closed')}>Close</button>
              {:else}<button class="btn" on:click={() => setStatus(t, 'open')}>Reopen</button>{/if}
            </td>
          </tr>
        {:else}
          <tr><td colspan="5" class="empty">No tickets.</td></tr>
        {/each}
      </tbody>
    </table>
  </Panel>
{/if}

<style>
  .grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:14px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid label.row { flex-direction:row; align-items:center; }
  .grid input, .grid select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .roles { margin-top:20px; }
  .roles .label { font-size:13px; color:#475569; margin-bottom:8px; }
  .chip { display:inline-flex; align-items:center; gap:6px; padding:6px 12px; border:1px solid #e5e7eb; border-radius:999px; font-size:13px; margin:0 6px 6px 0; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .status { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; }
  .status.open { background:#dbeafe; color:#1e40af; }
  .status.closed { background:#e5e7eb; color:#475569; }
  .empty { text-align:center; color:#94a3b8; padding:20px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
  select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
</style>
