<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let reports = [];
  let loading = true;
  let filter = 'open';

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    reports = await guildApi($currentGuildId).get('/moderation/reports');
    loading = false;
  }
  async function setStatus(r, status) {
    await guildApi($currentGuildId).patch('/moderation/reports/' + r.id, { status });
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: shown = filter === 'all' ? reports : reports.filter(r => r.status === filter);
</script>

<Panel title="User reports" subtitle="Reports submitted by members through /report or a panel.">
  <div slot="actions">
    <select bind:value={filter}>
      <option value="open">Open</option>
      <option value="resolved">Resolved</option>
      <option value="dismissed">Dismissed</option>
      <option value="all">All</option>
    </select>
  </div>
  {#if loading}<p>Loading…</p>{:else if shown.length === 0}<p class="empty">No reports.</p>{:else}
    <table>
      <thead><tr><th>When</th><th>Reporter</th><th>Target</th><th>Reason</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {#each shown as r}
          <tr>
            <td>{formatDate(r.created_at)}</td>
            <td>{r.reporter_name}</td>
            <td>{r.target_name}</td>
            <td>{r.reason}</td>
            <td><span class={'status ' + r.status}>{r.status}</span></td>
            <td class="actions">
              {#if r.status === 'open'}
                <button class="btn" on:click={() => setStatus(r, 'resolved')}>Resolve</button>
                <button class="btn" on:click={() => setStatus(r, 'dismissed')}>Dismiss</button>
              {:else}
                <button class="btn" on:click={() => setStatus(r, 'open')}>Reopen</button>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .status { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; background:#e5e7eb; color:#475569; }
  .status.open { background:#fee2e2; color:#991b1b; }
  .status.resolved { background:#d1fae5; color:#065f46; }
  .actions { display:flex; gap:6px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:12px; font-weight:600; cursor:pointer; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
