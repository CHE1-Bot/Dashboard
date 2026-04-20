<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let giveaways = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    giveaways = await guildApi($currentGuildId).get('/giveaways');
    loading = false;
  }
  async function endNow(g) {
    if (!confirm('End this giveaway and pick a winner now?')) return;
    await guildApi($currentGuildId).post(`/giveaways/${g.id}/end`, {});
    await load();
  }
  async function reroll(g) {
    await guildApi($currentGuildId).post(`/giveaways/${g.id}/reroll`, {});
    await load();
  }
  async function remove(g) {
    if (!confirm('Delete this giveaway?')) return;
    await guildApi($currentGuildId).del('/giveaways/' + g.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: running = giveaways.filter(g => g.status === 'running');
  $: ended   = giveaways.filter(g => g.status === 'ended');
</script>

<div class="cards">
  <StatCard label="Active" value={String(running.length)} icon="fa-gift" />
  <StatCard label="Ended" value={String(ended.length)} icon="fa-flag-checkered" />
  <StatCard label="Total entrants" value={giveaways.reduce((a,g) => a+g.entrants, 0).toLocaleString()} icon="fa-user-plus" />
</div>

<Panel title="All giveaways">
  <div slot="actions">
    <a class="btn primary" href="#/dashboard/giveaways/create-giveaway">+ New giveaway</a>
  </div>
  {#if loading}<p>Loading…</p>{:else if giveaways.length === 0}<p class="empty">No giveaways yet.</p>{:else}
    <table>
      <thead><tr><th>Prize</th><th>Winners</th><th>Entrants</th><th>Ends</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {#each giveaways as g}
          <tr>
            <td>{g.prize}</td>
            <td>{g.winner_count} {#if g.winners?.length}· {g.winners.join(', ')}{/if}</td>
            <td>{g.entrants}</td>
            <td>{formatDate(g.ends_at)}</td>
            <td><span class={'status ' + g.status}>{g.status}</span></td>
            <td class="actions">
              {#if g.status === 'running'}<button class="btn" on:click={() => endNow(g)}>End</button>{/if}
              {#if g.status === 'ended'}<button class="btn" on:click={() => reroll(g)}>Reroll</button>{/if}
              <button class="btn danger" on:click={() => remove(g)}>Delete</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(200px, 1fr)); gap:14px; margin-bottom:20px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .status { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; }
  .status.running { background:#dbeafe; color:#1e40af; }
  .status.ended { background:#e5e7eb; color:#475569; }
  .actions { display:flex; gap:6px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; text-decoration:none; display:inline-block; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
