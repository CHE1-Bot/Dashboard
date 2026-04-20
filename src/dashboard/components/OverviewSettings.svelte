<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let settings = null;
  let channels = [];
  let saving = false;
  let toast = '';

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [settings, channels] = await Promise.all([api.get('/settings'), api.get('/channels')]);
  }
  async function save() {
    saving = true;
    try {
      settings = await guildApi($currentGuildId).patch('/settings', settings);
      toast = 'Saved'; setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !settings}<p>Loading…</p>{:else}
  <Panel title="Notifications & logging" subtitle="How and where Che1 reports what's happening.">
    <div slot="actions">
      {#if toast}<span class="toast">{toast}</span>{/if}
      <button class="btn primary" on:click={save} disabled={saving}>Save</button>
    </div>
    <div class="grid">
      <label>Log channel
        <select bind:value={settings.log_channel_id}>
          <option value="">— none —</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
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
    </div>
  </Panel>

  <Panel title="Premium" subtitle="Premium unlocks unlimited giveaways, extra automod rules and more.">
    <div class="premium">
      <div>
        Status: <strong class={settings.premium ? 'ok' : 'no'}>{settings.premium ? 'Active' : 'Inactive'}</strong>
      </div>
      <a class="btn primary" href="#/premium">{settings.premium ? 'Manage' : 'Upgrade'}</a>
    </div>
  </Panel>
{/if}

<style>
  .grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:16px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid input, .grid select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .premium { display:flex; justify-content:space-between; align-items:center; }
  .ok { color:#059669; } .no { color:#b45309; }
  .btn { padding:8px 16px; border:none; border-radius:8px; cursor:pointer; font-weight:600; background:#e2e8f0; color:#0f172a; text-decoration:none; display:inline-block; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
</style>
