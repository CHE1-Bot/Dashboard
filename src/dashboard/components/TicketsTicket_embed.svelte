<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let embed = null;
  let saving = false;
  let toast = '';

  async function load() {
    if (!$currentGuildId) return;
    embed = await guildApi($currentGuildId).get('/tickets/embed');
  }
  async function save() {
    saving = true;
    try {
      embed = await guildApi($currentGuildId).patch('/tickets/embed', embed);
      toast = 'Saved'; setTimeout(() => toast = '', 2000);
    } finally { saving = false; }
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !embed}<p>Loading…</p>{:else}
  <div class="split">
    <Panel title="Ticket embed">
      <div slot="actions">
        {#if toast}<span class="toast">{toast}</span>{/if}
        <button class="btn primary" on:click={save} disabled={saving}>Save</button>
      </div>
      <div class="grid">
        <label>Title<input bind:value={embed.title} /></label>
        <label>Color<input type="color" bind:value={embed.color} /></label>
        <label class="full">Description<textarea rows="5" bind:value={embed.description}></textarea></label>
        <label>Footer<input bind:value={embed.footer} /></label>
        <label>Thumbnail URL<input bind:value={embed.thumbnail} /></label>
      </div>
    </Panel>

    <Panel title="Preview">
      <div class="preview" style="--c:{embed.color}">
        {#if embed.thumbnail}<img class="thumb" src={embed.thumbnail} alt="" />{/if}
        <h3>{embed.title || '(no title)'}</h3>
        <p>{embed.description || '(no description)'}</p>
        {#if embed.footer}<small>{embed.footer}</small>{/if}
      </div>
    </Panel>
  </div>
{/if}

<style>
  .split { display:grid; grid-template-columns:1fr 1fr; gap:20px; }
  @media (max-width: 860px) { .split { grid-template-columns:1fr; } }
  .grid { display:grid; grid-template-columns:1fr 120px; gap:14px; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid .full { grid-column: 1 / -1; }
  .grid input, .grid textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  .grid input[type="color"] { padding:0; height:38px; }
  .preview { background:#2b2d31; color:#fff; border-left:4px solid var(--c); padding:14px 16px; border-radius:6px; }
  .preview h3 { margin:0 0 6px; font-size:16px; }
  .preview p { margin:0 0 8px; color:#dbdee1; white-space:pre-wrap; }
  .preview small { color:#a3a6aa; font-size:12px; }
  .thumb { float:right; max-width:80px; max-height:80px; border-radius:4px; margin-left:12px; }
  .btn { padding:8px 16px; border:none; border-radius:8px; cursor:pointer; font-weight:600; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .toast { color:#059669; font-size:13px; margin-right:8px; }
</style>
