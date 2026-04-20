<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let stats = null;
  let history = [];
  let alerts = [];

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [stats, history, alerts] = await Promise.all([
      api.get('/overview'), api.get('/history'), api.get('/alerts'),
    ]);
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !stats}<p>Loading…</p>{:else}
  <div class="cards">
    <StatCard label="Members" value={stats.members.toLocaleString()} icon="fa-users" />
    <StatCard label="Online" value={stats.online_members.toLocaleString()} icon="fa-signal" />
    <StatCard label="Messages today" value={stats.messages_today.toLocaleString()} icon="fa-comments" />
    <StatCard label="Commands today" value={stats.commands_today.toLocaleString()} icon="fa-terminal" />
    <StatCard label="Open tickets" value={String(stats.open_tickets)} icon="fa-ticket" />
    <StatCard label="Active giveaways" value={String(stats.active_giveaways)} icon="fa-gift" />
    <StatCard label="Mod actions (7d)" value={String(stats.mod_actions_week)} icon="fa-shield" />
  </div>

  <div class="two">
    <Panel title="Recent activity">
      {#if history.length === 0}<p class="empty">No activity yet.</p>{:else}
        <ul class="feed">
          {#each history.slice(0, 8) as h}
            <li><strong>{h.actor}</strong> · {h.event} · <em>{h.detail}</em><span class="time">{relativeTime(h.created_at)}</span></li>
          {/each}
        </ul>
      {/if}
    </Panel>
    <Panel title="Active alerts">
      {#if alerts.length === 0}<p class="empty">No alerts.</p>{:else}
        <ul class="alerts">
          {#each alerts.slice(0, 5) as a}
            <li class={'sev-' + a.severity}>
              <strong>{a.title}</strong>
              <span>{a.detail}</span>
              <span class="time">{relativeTime(a.created_at)}</span>
            </li>
          {/each}
        </ul>
      {/if}
    </Panel>
  </div>
{/if}

<style>
  .cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(200px, 1fr)); gap:14px; margin-bottom:20px; }
  .two { display:grid; grid-template-columns:repeat(2, 1fr); gap:20px; }
  @media (max-width: 860px) { .two { grid-template-columns:1fr; } }
  .feed { list-style:none; margin:0; padding:0; }
  .feed li { padding:10px 0; border-bottom:1px solid #f1f5f9; font-size:13px; color:#475569; display:flex; flex-wrap:wrap; gap:6px; align-items:baseline; }
  .feed li:last-child { border-bottom:0; }
  .feed em { color:#111827; font-style:normal; }
  .time { margin-left:auto; color:#94a3b8; font-size:12px; }
  .alerts { list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:10px; }
  .alerts li { padding:10px 12px; border-radius:8px; display:grid; gap:4px; background:#f8fafc; }
  .alerts li.sev-warn { background:#fef3c7; }
  .alerts li.sev-error { background:#fee2e2; }
  .empty { color:#94a3b8; margin:0; }
</style>
