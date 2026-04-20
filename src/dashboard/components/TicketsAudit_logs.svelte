<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let events = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    events = await guildApi($currentGuildId).get('/tickets/audit');
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Ticket audit log" subtitle="Every create, update, close and delete performed on tickets.">
  {#if loading}<p>Loading…</p>{:else if events.length === 0}<p class="empty">No ticket events.</p>{:else}
    <table>
      <thead><tr><th>When</th><th>Actor</th><th>Event</th><th>Detail</th></tr></thead>
      <tbody>
        {#each events as e}
          <tr><td>{formatDate(e.created_at)}</td><td>{e.actor}</td><td><code>{e.event}</code></td><td>{e.detail}</td></tr>
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
  .empty { color:#94a3b8; }
</style>
