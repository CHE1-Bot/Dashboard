<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let data = null;

  async function load() {
    if (!$currentGuildId) return;
    data = await guildApi($currentGuildId).get('/tickets/stats');
  }
  onMount(load);
  $: if ($currentGuildId) load();

  function fmtHours(h) {
    if (!h) return '—';
    if (h < 1) return Math.round(h * 60) + ' min';
    return h.toFixed(1) + ' h';
  }
  function fmtMin(m) {
    if (!m) return '—';
    if (m < 1) return '<1 min';
    if (m < 60) return Math.round(m) + ' min';
    return (m / 60).toFixed(1) + ' h';
  }
</script>

{#if !data}<p>Loading…</p>{:else}
  <div class="cards">
    <StatCard label="Open tickets" value={String(data.totals.open)} icon="fa-ticket" />
    <StatCard label="Closed" value={String(data.totals.closed)} icon="fa-check" />
    <StatCard label="Currently claimed" value={String(data.totals.claimed)} icon="fa-hand" />
    <StatCard label="Avg handle time" value={fmtHours(data.totals.avg_handle_hours)} icon="fa-clock" />
  </div>

  <Panel title="Opens (last 7 days)">
    <div class="bars">
      {#each data.opens_per_day as p}
        <div class="bar"><div class="fill" style={`height:${Math.max(4, p.value * 12)}px`}></div><span>{p.date.slice(5)}</span></div>
      {/each}
    </div>
  </Panel>

  <Panel title="By category">
    {#if !data.by_category.length}<p class="empty">No tickets yet.</p>{:else}
      <table>
        <thead><tr><th>Category</th><th>Open</th><th>Closed</th><th>Total</th></tr></thead>
        <tbody>
          {#each data.by_category as c}
            <tr><td>{c.name}</td><td>{c.open}</td><td>{c.closed}</td><td><strong>{c.open + c.closed}</strong></td></tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </Panel>

  <Panel title="By staff" subtitle="Sorted by claims taken.">
    {#if !data.by_staff.length}<p class="empty">No claims yet.</p>{:else}
      <table>
        <thead><tr><th>Staff</th><th>Claimed</th><th>Closed</th><th>Avg response</th></tr></thead>
        <tbody>
          {#each data.by_staff as s}
            <tr><td>{s.username}</td><td>{s.claimed}</td><td>{s.closed}</td><td>{fmtMin(s.avg_response_min)}</td></tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </Panel>
{/if}

<style>
  .cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(200px, 1fr)); gap:14px; margin-bottom:20px; }
  .bars { display:flex; gap:12px; align-items:flex-end; height:160px; }
  .bar { flex:1; text-align:center; display:flex; flex-direction:column; justify-content:flex-end; gap:6px; font-size:11px; color:#94a3b8; }
  .fill { background:linear-gradient(180deg, #5865f2, #10b981); border-radius:6px 6px 0 0; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .empty { color:#94a3b8; }
</style>
