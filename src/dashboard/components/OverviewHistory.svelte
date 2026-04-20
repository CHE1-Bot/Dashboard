<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let history = [];
  let filter = '';
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    history = await guildApi($currentGuildId).get('/history');
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  $: filtered = filter
    ? history.filter(h => [h.actor, h.event, h.detail].join(' ').toLowerCase().includes(filter.toLowerCase()))
    : history;
</script>

<Panel title="Activity history" subtitle="Every dashboard-originated change, oldest to newest.">
  <div slot="actions"><input placeholder="Filter…" bind:value={filter} /></div>
  {#if loading}<p>Loading…</p>{:else if filtered.length === 0}<p class="empty">No events.</p>{:else}
    <table>
      <thead><tr><th>When</th><th>Actor</th><th>Event</th><th>Detail</th></tr></thead>
      <tbody>
        {#each filtered as h}
          <tr><td>{formatDate(h.created_at)}</td><td>{h.actor}</td><td><code>{h.event}</code></td><td>{h.detail}</td></tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:13px; }
  th, td { text-align:left; padding:8px 10px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:11px; text-transform:uppercase; }
  code { font-size:12px; color:#4f46e5; background:#eef2ff; padding:2px 6px; border-radius:4px; }
  input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  .empty { color:#94a3b8; }
</style>
