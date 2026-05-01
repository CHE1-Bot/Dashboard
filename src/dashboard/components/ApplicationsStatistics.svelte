<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let data = null;

  async function load() {
    if (!$currentGuildId) return;
    data = await guildApi($currentGuildId).get('/applications/stats');
  }
  onMount(load);
  $: if ($currentGuildId) load();

  function fmtHours(h) {
    if (!h) return '—';
    if (h < 1) return Math.round(h * 60) + ' min';
    if (h < 48) return h.toFixed(1) + ' h';
    return (h / 24).toFixed(1) + ' d';
  }
</script>

{#if !data}<p>Loading…</p>{:else}
  <div class="cards">
    <StatCard label="Pending" value={String(data.totals.pending)} icon="fa-hourglass-half" />
    <StatCard label="Accepted" value={String(data.totals.accepted)} icon="fa-check" />
    <StatCard label="Rejected" value={String(data.totals.rejected)} icon="fa-xmark" />
    <StatCard label="Avg review time" value={fmtHours(data.totals.avg_review_hours)} icon="fa-stopwatch" />
  </div>

  <Panel title="Submissions (last 7 days)">
    <div class="bars">
      {#each data.submissions_per_day as p}
        <div class="bar"><div class="fill" style={`height:${Math.max(4, p.value * 14)}px`}></div><span>{p.date.slice(5)}</span></div>
      {/each}
    </div>
  </Panel>

  <Panel title="By form">
    {#if !data.by_form.length}<p class="empty">No applications yet.</p>{:else}
      <table>
        <thead><tr><th>Form</th><th>Pending</th><th>Accepted</th><th>Rejected</th><th>Total</th><th>Accept rate</th></tr></thead>
        <tbody>
          {#each data.by_form as f}
            {@const total = f.pending + f.accepted + f.rejected}
            {@const reviewed = f.accepted + f.rejected}
            <tr>
              <td><strong>{f.name || '—'}</strong></td>
              <td>{f.pending}</td>
              <td>{f.accepted}</td>
              <td>{f.rejected}</td>
              <td><strong>{total}</strong></td>
              <td>{reviewed > 0 ? Math.round(f.accepted / reviewed * 100) + '%' : '—'}</td>
            </tr>
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
