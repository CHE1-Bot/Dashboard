<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let tags = [];
  let loading = true;
  let draft = { name: '', color: '#3498db' };

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    tags = await guildApi($currentGuildId).get('/tickets/tags');
    loading = false;
  }
  async function create() {
    if (!draft.name) return;
    await guildApi($currentGuildId).post('/tickets/tags', draft);
    draft = { name: '', color: '#3498db' };
    await load();
  }
  async function remove(t) {
    if (!confirm(`Delete tag "${t.name}"?`)) return;
    await guildApi($currentGuildId).del('/tickets/tags/' + t.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Ticket tags" subtitle="Labels staff can apply to tickets for filtering.">
  <div slot="actions">
    <input placeholder="Tag name" bind:value={draft.name} />
    <input type="color" bind:value={draft.color} />
    <button class="btn primary" on:click={create}>Add</button>
  </div>
  {#if loading}<p>Loading…</p>{:else if tags.length === 0}<p class="empty">No tags yet.</p>{:else}
    <div class="tags">
      {#each tags as t}
        <span class="tag" style="background:{t.color}">
          {t.name}
          <button on:click={() => remove(t)}>×</button>
        </span>
      {/each}
    </div>
  {/if}
</Panel>

<style>
  input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  input[type="color"] { padding:0; width:40px; height:36px; }
  .tags { display:flex; flex-wrap:wrap; gap:8px; }
  .tag { display:inline-flex; align-items:center; gap:6px; padding:4px 10px; border-radius:999px; color:#fff; font-size:13px; font-weight:600; }
  .tag button { background:rgba(0,0,0,0.2); color:#fff; border:none; width:18px; height:18px; border-radius:50%; cursor:pointer; font-size:11px; line-height:1; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
