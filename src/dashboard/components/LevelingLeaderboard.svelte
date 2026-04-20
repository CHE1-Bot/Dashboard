<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let rows = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    rows = await guildApi($currentGuildId).get('/leveling/leaderboard');
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  $: sorted = [...rows].sort((a, b) => b.total_xp - a.total_xp);
</script>

<Panel title="Leaderboard" subtitle="Ranked by total XP.">
  {#if loading}<p>Loading…</p>{:else if sorted.length === 0}<p class="empty">No levels tracked yet.</p>{:else}
    <table>
      <thead><tr><th>#</th><th>User</th><th>Level</th><th>XP</th><th>Total XP</th></tr></thead>
      <tbody>
        {#each sorted as u, i}
          <tr>
            <td>{i + 1}</td>
            <td>{u.username}</td>
            <td><strong>{u.level}</strong></td>
            <td>{u.xp.toLocaleString()}</td>
            <td>{u.total_xp.toLocaleString()}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  tr:nth-child(1) td:first-child { color:#f59e0b; font-weight:700; }
  tr:nth-child(2) td:first-child { color:#94a3b8; font-weight:700; }
  tr:nth-child(3) td:first-child { color:#b45309; font-weight:700; }
  .empty { color:#94a3b8; }
</style>
