<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';
  import StatCard from '../../lib/StatCard.svelte';

  let logs = [];
  let reports = [];
  let rules = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [logs, reports, rules] = await Promise.all([
      api.get('/moderation/logs'), api.get('/moderation/reports'), api.get('/moderation/automod'),
    ]);
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  $: byAction = logs.reduce((a, l) => { a[l.action] = (a[l.action] || 0) + 1; return a; }, {});
</script>

{#if loading}<p>Loading…</p>{:else}
  <div class="cards">
    <StatCard label="Warnings" value={String(byAction.warn || 0)} icon="fa-triangle-exclamation" />
    <StatCard label="Mutes" value={String(byAction.mute || 0)} icon="fa-volume-xmark" />
    <StatCard label="Kicks" value={String(byAction.kick || 0)} icon="fa-user-slash" />
    <StatCard label="Bans" value={String(byAction.ban || 0)} icon="fa-gavel" />
    <StatCard label="Open reports" value={String(reports.filter(r => r.status === 'open').length)} icon="fa-flag" />
    <StatCard label="AutoMod rules" value={String(rules.filter(r => r.enabled).length) + '/' + rules.length} icon="fa-robot" />
  </div>

  <Panel title="Recent moderation activity">
    {#if logs.length === 0}<p class="empty">No recent actions.</p>{:else}
      <ul class="feed">
        {#each logs.slice(0, 10) as l}
          <li>
            <span class={'tag ' + l.action}>{l.action}</span>
            <strong>{l.target_name}</strong> by <strong>{l.moderator_name}</strong>
            {#if l.reason}— <em>{l.reason}</em>{/if}
            <span class="time">{relativeTime(l.created_at)}</span>
          </li>
        {/each}
      </ul>
    {/if}
  </Panel>
{/if}

<style>
  .cards { display:grid; grid-template-columns:repeat(auto-fill, minmax(180px, 1fr)); gap:14px; margin-bottom:20px; }
  .feed { list-style:none; margin:0; padding:0; }
  .feed li { padding:10px 0; border-bottom:1px solid #f1f5f9; font-size:14px; display:flex; align-items:center; gap:8px; flex-wrap:wrap; }
  .feed li:last-child { border-bottom:0; }
  .feed em { color:#64748b; font-style:normal; }
  .time { margin-left:auto; color:#94a3b8; font-size:12px; }
  .tag { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; background:#e5e7eb; color:#475569; }
  .tag.warn { background:#fef3c7; color:#92400e; }
  .tag.mute { background:#e0e7ff; color:#3730a3; }
  .tag.kick { background:#fed7aa; color:#9a3412; }
  .tag.ban { background:#fee2e2; color:#991b1b; }
  .empty { color:#94a3b8; }
</style>
