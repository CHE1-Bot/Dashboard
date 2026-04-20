<script>
  import { onMount } from 'svelte';
  import { currentGuild, currentGuildId } from '../../lib/stores.js';
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
      api.get('/settings'), api.get('/channels'), api.get('/roles'),
    ]);
  }

  async function save() {
    saving = true;
    try {
      settings = await guildApi($currentGuildId).patch('/settings', settings);
      toast = 'Saved';
      setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }

  function toggleModule(key) {
    settings.module_enabled = { ...(settings.module_enabled || {}), [key]: !settings.module_enabled?.[key] };
  }

  function toggleAutoRole(id) {
    const s = new Set(settings.auto_role_ids || []);
    if (s.has(id)) s.delete(id); else s.add(id);
    settings.auto_role_ids = [...s];
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !$currentGuild}
  <p>No server selected. <a href="#/servers">Pick a server</a>.</p>
{:else if !settings}
  <p>Loading…</p>
{:else}
  <Panel title={$currentGuild.name} subtitle="General settings">
    <div slot="actions">
      {#if toast}<span class="toast">{toast}</span>{/if}
      <button class="btn primary" on:click={save} disabled={saving}>Save</button>
    </div>

    <div class="grid">
      <label>Prefix<input bind:value={settings.prefix} maxlength="5" /></label>
      <label>Language
        <select bind:value={settings.language}>
          <option value="en">English</option>
          <option value="fi">Suomi</option>
          <option value="de">Deutsch</option>
          <option value="es">Español</option>
        </select>
      </label>
      <label>Timezone<input bind:value={settings.timezone} /></label>
      <label>Log channel
        <select bind:value={settings.log_channel_id}>
          <option value="">— none —</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
      <label>Welcome channel
        <select bind:value={settings.welcome_channel_id}>
          <option value="">— none —</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
      <label class="full">Welcome message
        <textarea rows="3" bind:value={settings.welcome_message}></textarea>
        <small>Placeholders: <code>{'{user}'}</code>, <code>{'{server}'}</code></small>
      </label>
    </div>
  </Panel>

  <Panel title="Modules" subtitle="Turn features on or off for this server">
    <div class="modules">
      {#each ['tickets','moderation','giveaways','leveling','applications'] as key}
        <label class="switch">
          <input type="checkbox" checked={!!settings.module_enabled?.[key]} on:change={() => toggleModule(key)} />
          <span>{key}</span>
        </label>
      {/each}
    </div>
  </Panel>

  <Panel title="Auto-roles" subtitle="Roles assigned to new members on join">
    <div class="roles">
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip">
          <input type="checkbox" checked={settings.auto_role_ids?.includes(r.id)} on:change={() => toggleAutoRole(r.id)} />
          <span style="color:{'#' + r.color.toString(16).padStart(6,'0')}">{r.name}</span>
        </label>
      {/each}
    </div>
  </Panel>
{/if}

<style>
  .grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:16px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid .full { grid-column: 1 / -1; }
  .grid input, .grid select, .grid textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  small { color:#94a3b8; }
  .modules { display:flex; flex-wrap:wrap; gap:10px; }
  .switch { display:flex; align-items:center; gap:8px; padding:8px 14px; border:1px solid #e5e7eb; border-radius:999px; text-transform:capitalize; }
  .roles { display:flex; flex-wrap:wrap; gap:8px; }
  .chip { display:flex; align-items:center; gap:6px; padding:6px 12px; border:1px solid #e5e7eb; border-radius:999px; font-size:13px; }
  .btn { padding:8px 16px; border:none; border-radius:8px; cursor:pointer; font-weight:600; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
</style>
