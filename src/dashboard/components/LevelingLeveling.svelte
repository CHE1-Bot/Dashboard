<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let settings = null;
  let channels = [];
  let roles = [];
  let saving = false;
  let toast = '';

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [settings, channels, roles] = await Promise.all([
      api.get('/leveling/settings'), api.get('/channels'), api.get('/roles'),
    ]);
  }
  async function save() {
    saving = true;
    try {
      settings = await guildApi($currentGuildId).patch('/leveling/settings', settings);
      toast = 'Saved'; setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }
  function toggleChannel(id) {
    const s = new Set(settings.no_xp_channels || []);
    s.has(id) ? s.delete(id) : s.add(id);
    settings.no_xp_channels = [...s];
  }
  function toggleRole(id) {
    const s = new Set(settings.no_xp_roles || []);
    s.has(id) ? s.delete(id) : s.add(id);
    settings.no_xp_roles = [...s];
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !settings}<p>Loading…</p>{:else}
  <Panel title="Leveling" subtitle="XP and announcement settings.">
    <div slot="actions">
      {#if toast}<span class="toast">{toast}</span>{/if}
      <button class="btn primary" on:click={save} disabled={saving}>Save</button>
    </div>
    <label class="switch">
      <input type="checkbox" bind:checked={settings.enabled} />
      <span>{settings.enabled ? 'Leveling is enabled' : 'Leveling is disabled'}</span>
    </label>

    <div class="grid">
      <label>XP per message<input type="number" min="1" max="1000" bind:value={settings.xp_per_message} /></label>
      <label>Cooldown (seconds)<input type="number" min="0" max="600" bind:value={settings.cooldown_sec} /></label>
      <label>Multiplier<input type="number" step="0.1" min="0.1" max="5" bind:value={settings.multiplier} /></label>
      <label>Announce channel
        <select bind:value={settings.announce_channel}>
          <option value="">— same channel —</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
      <label class="full">Level-up message
        <textarea rows="2" bind:value={settings.announce_message}></textarea>
        <small>Placeholders: <code>{'{user}'}</code>, <code>{'{level}'}</code></small>
      </label>
    </div>
  </Panel>

  <Panel title="Blocked channels" subtitle="No XP is awarded in these channels.">
    <div class="chips">
      {#each channels.filter(c => c.type === 'text') as c}
        <label class="chip"><input type="checkbox" checked={settings.no_xp_channels?.includes(c.id)} on:change={() => toggleChannel(c.id)} />#{c.name}</label>
      {/each}
    </div>
  </Panel>

  <Panel title="Blocked roles" subtitle="Members with these roles don't earn XP.">
    <div class="chips">
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={settings.no_xp_roles?.includes(r.id)} on:change={() => toggleRole(r.id)} />{r.name}</label>
      {/each}
    </div>
  </Panel>
{/if}

<style>
  .switch { display:flex; align-items:center; gap:10px; padding:12px; background:#f8fafc; border-radius:8px; margin-bottom:14px; }
  .grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:14px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid .full { grid-column: 1 / -1; }
  .grid input, .grid select, .grid textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  small { color:#94a3b8; }
  .chips { display:flex; flex-wrap:wrap; gap:8px; }
  .chip { display:flex; align-items:center; gap:6px; padding:6px 12px; border:1px solid #e5e7eb; border-radius:999px; font-size:13px; }
  .btn { padding:8px 16px; border:none; border-radius:8px; cursor:pointer; font-weight:600; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
</style>
