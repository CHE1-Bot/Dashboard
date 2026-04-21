<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let data = null;

  async function load() {
    if (!$currentGuildId) return;
    data = await guildApi($currentGuildId).get('/leveling/stats');
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !data}<p>Loading…</p>{:else}
  <div class="cards">
    <StatCard label="Tracked users" value={data.tracked_users.toLocaleString()} icon="fa-users" />
    <StatCard label="Top level" value={String(data.top_level)} icon="fa-crown" />
    <StatCard label="Average XP" value={data.average_xp.toLocaleString()} icon="fa-chart-line" />
    <StatCard label="Total XP awarded" value={data.total_xp.toLocaleString()} icon="fa-bolt" />
  </div>

  <Panel title="Level distribution">
    <div class="dist">
      {#each data.top_users.slice(0, 10) as u}
        <div class="row">
          <span class="name">{u.username}</span>
          <div class="bar"><div style="width:{Math.min(100, u.level * 5)}%"></div></div>
          <span class="lvl">Lv {u.level}</span>
        </div>
      {/each}
    </div>
  </Panel>
{/if}

<style>
  .cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(200px, 1fr)); gap:14px; margin-bottom:20px; }
  .dist { display:flex; flex-direction:column; gap:10px; }
  .row { display:grid; grid-template-columns:120px 1fr 60px; gap:12px; align-items:center; font-size:13px; }
  .name { color:#111827; font-weight:500; }
  .bar { height:10px; background:#e5e7eb; border-radius:6px; overflow:hidden; }
  .bar > div { height:100%; background:linear-gradient(90deg, #5865f2, #10b981); border-radius:6px; }
  .lvl { color:#475569; text-align:right; }
</style>
