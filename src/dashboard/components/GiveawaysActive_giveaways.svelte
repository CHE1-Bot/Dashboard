<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let list = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    list = await guildApi($currentGuildId).get('/giveaways/active');
    loading = false;
  }
  async function endNow(g) {
    if (!confirm('End this giveaway now?')) return;
    await guildApi($currentGuildId).post(`/giveaways/${g.id}/end`, {});
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Active giveaways" subtitle="{list.length} running">
  {#if loading}<p>Loading…</p>{:else if list.length === 0}<p class="empty">No active giveaways right now.</p>{:else}
    <div class="grid">
      {#each list as g}
        <div class="card">
          <div class="prize">{g.prize}</div>
          <div class="meta"><i class="fa-solid fa-clock"></i> ends {relativeTime(g.ends_at)}</div>
          <div class="meta"><i class="fa-solid fa-user-plus"></i> {g.entrants} entrants</div>
          <div class="meta"><i class="fa-solid fa-trophy"></i> {g.winner_count} winner(s)</div>
          <button class="btn" on:click={() => endNow(g)}>End now</button>
        </div>
      {/each}
    </div>
  {/if}
</Panel>

<style>
  .grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(260px, 1fr)); gap:14px; }
  .card { border:1px solid #e5e7eb; border-radius:12px; padding:16px; display:flex; flex-direction:column; gap:6px; }
  .prize { font-size:16px; font-weight:700; color:#111827; }
  .meta { color:#475569; font-size:13px; }
  .btn { margin-top:10px; padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; align-self:flex-start; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
