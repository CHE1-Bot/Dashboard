<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let stats = null;

  async function load() {
    if (!$currentGuildId) return;
    stats = await guildApi($currentGuildId).get('/overview');
  }
  onMount(load);
  $: if ($currentGuildId) load();

  function barHeight(arr, v) {
    const max = Math.max(...arr.map(p => p.value), 1);
    return Math.round((v / max) * 100);
  }
</script>

{#if !stats}<p>Loading…</p>{:else}
  <Panel title="Messages (7 days)">
    <div class="chart">
      {#each stats.messages_per_day as p}
        <div class="bar-col">
          <div class="bar" style="height:{barHeight(stats.messages_per_day, p.value)}%" title={`${p.date}: ${p.value}`}></div>
          <small>{p.date.slice(5)}</small>
        </div>
      {/each}
    </div>
  </Panel>

  <Panel title="Joins (7 days)">
    <div class="chart">
      {#each stats.joins_per_day as p}
        <div class="bar-col">
          <div class="bar green" style="height:{barHeight(stats.joins_per_day, p.value)}%" title={`${p.date}: ${p.value}`}></div>
          <small>{p.date.slice(5)}</small>
        </div>
      {/each}
    </div>
  </Panel>

  <Panel title="Snapshot">
    <dl>
      <div><dt>Members</dt><dd>{stats.members.toLocaleString()}</dd></div>
      <div><dt>Online now</dt><dd>{stats.online_members.toLocaleString()}</dd></div>
      <div><dt>Messages today</dt><dd>{stats.messages_today.toLocaleString()}</dd></div>
      <div><dt>Commands today</dt><dd>{stats.commands_today.toLocaleString()}</dd></div>
      <div><dt>Open tickets</dt><dd>{stats.open_tickets}</dd></div>
      <div><dt>Active giveaways</dt><dd>{stats.active_giveaways}</dd></div>
      <div><dt>Mod actions (7d)</dt><dd>{stats.mod_actions_week}</dd></div>
    </dl>
  </Panel>
{/if}

<style>
  .chart { display:flex; align-items:flex-end; gap:12px; height:200px; }
  .bar-col { flex:1; display:flex; flex-direction:column; align-items:center; gap:6px; }
  .bar { width:100%; background:#5865f2; border-radius:6px 6px 0 0; min-height:4px; }
  .bar.green { background:#10b981; }
  small { color:#94a3b8; font-size:11px; }
  dl { display:grid; grid-template-columns:repeat(auto-fill, minmax(180px, 1fr)); gap:12px; margin:0; }
  dl > div { background:#f8fafc; border-radius:8px; padding:12px; }
  dt { font-size:11px; color:#6b7280; text-transform:uppercase; letter-spacing:0.05em; }
  dd { margin:4px 0 0; font-size:20px; font-weight:700; color:#111827; }
</style>
