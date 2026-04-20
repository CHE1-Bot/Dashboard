<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let logs = [];
  let loading = true;
  let filterAction = 'all';
  let filterQuery = '';

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    logs = await guildApi($currentGuildId).get('/moderation/logs');
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  $: shown = logs.filter(l =>
    (filterAction === 'all' || l.action === filterAction) &&
    (!filterQuery || [l.target_name, l.moderator_name, l.reason].join(' ').toLowerCase().includes(filterQuery.toLowerCase()))
  );
</script>

<Panel title="Moderation logs" subtitle="Every warn, mute, kick, ban and revert.">
  <div slot="actions">
    <select bind:value={filterAction}>
      <option value="all">All actions</option>
      <option value="warn">Warn</option>
      <option value="mute">Mute</option>
      <option value="kick">Kick</option>
      <option value="ban">Ban</option>
      <option value="unban">Unban</option>
      <option value="unmute">Unmute</option>
    </select>
    <input placeholder="Filter by user or reason" bind:value={filterQuery} />
  </div>
  {#if loading}<p>Loading…</p>{:else if shown.length === 0}<p class="empty">No log entries.</p>{:else}
    <table>
      <thead><tr><th>When</th><th>Action</th><th>Target</th><th>Moderator</th><th>Reason</th></tr></thead>
      <tbody>
        {#each shown as l}
          <tr>
            <td>{formatDate(l.created_at)}</td>
            <td><span class={'tag ' + l.action}>{l.action}</span></td>
            <td>{l.target_name}</td>
            <td>{l.moderator_name}</td>
            <td>{l.reason || '—'}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  select, input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .tag { padding:2px 8px; border-radius:999px; font-size:11px; font-weight:600; text-transform:uppercase; background:#e5e7eb; color:#475569; }
  .tag.warn { background:#fef3c7; color:#92400e; }
  .tag.mute { background:#e0e7ff; color:#3730a3; }
  .tag.kick { background:#fed7aa; color:#9a3412; }
  .tag.ban { background:#fee2e2; color:#991b1b; }
  .empty { color:#94a3b8; }
</style>
