<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let channels = [];
  let roles = [];
  let draft = { prize: '', channel_id: '', winner_count: 1, duration_hours: 24, required_role_id: '' };
  let creating = false;
  let created = null;

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [channels, roles] = await Promise.all([api.get('/channels'), api.get('/roles')]);
  }
  async function submit() {
    if (!draft.prize || !draft.channel_id) return;
    creating = true;
    try {
      const ends = new Date(Date.now() + draft.duration_hours * 3600 * 1000).toISOString();
      created = await guildApi($currentGuildId).post('/giveaways', {
        prize: draft.prize, channel_id: draft.channel_id,
        winner_count: Number(draft.winner_count) || 1,
        ends_at: ends,
      });
      draft = { prize: '', channel_id: '', winner_count: 1, duration_hours: 24, required_role_id: '' };
    } finally { creating = false; }
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Create giveaway" subtitle="The bot will post the embed and announce winners automatically.">
  <div class="grid">
    <label class="full">Prize<input bind:value={draft.prize} placeholder="e.g. Nitro Classic" /></label>
    <label>Channel
      <select bind:value={draft.channel_id}>
        <option value="">—</option>
        {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
      </select>
    </label>
    <label>Winners<input type="number" min="1" max="20" bind:value={draft.winner_count} /></label>
    <label>Duration (hours)<input type="number" min="1" max="8760" bind:value={draft.duration_hours} /></label>
    <label>Required role
      <select bind:value={draft.required_role_id}>
        <option value="">— none —</option>
        {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
      </select>
    </label>
  </div>
  <button class="btn primary" on:click={submit} disabled={creating || !draft.prize || !draft.channel_id}>
    {creating ? 'Creating…' : 'Create giveaway'}
  </button>
  {#if created}
    <div class="success">Created: <strong>{created.prize}</strong> · ends {new Date(created.ends_at).toLocaleString()}</div>
  {/if}
</Panel>

<style>
  .grid { display:grid; grid-template-columns:repeat(3, 1fr); gap:14px; margin-bottom:16px; }
  .grid .full { grid-column: 1 / -1; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid input, .grid select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .btn { padding:10px 20px; background:#5865f2; color:#fff; border:none; border-radius:8px; font-size:14px; font-weight:600; cursor:pointer; }
  .btn:disabled { opacity:0.5; cursor:not-allowed; }
  .btn:hover:not(:disabled) { filter:brightness(0.95); }
  .success { margin-top:14px; padding:10px 14px; background:#d1fae5; color:#065f46; border-radius:8px; font-size:14px; }
</style>
