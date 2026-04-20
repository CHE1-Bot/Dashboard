<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let alerts = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    alerts = await guildApi($currentGuildId).get('/alerts');
    loading = false;
  }
  async function dismiss(a) {
    await guildApi($currentGuildId).del('/alerts/' + a.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Alerts" subtitle="Notifications from the bot's health checks and automated rules.">
  {#if loading}<p>Loading…</p>{:else if alerts.length === 0}<p class="empty">No active alerts. Everything looks healthy.</p>{:else}
    <ul class="list">
      {#each alerts as a}
        <li class={'sev-' + a.severity}>
          <div>
            <strong>{a.title}</strong>
            <p>{a.detail}</p>
            <time>{formatDate(a.created_at)}</time>
          </div>
          <button class="btn" on:click={() => dismiss(a)}>Dismiss</button>
        </li>
      {/each}
    </ul>
  {/if}
</Panel>

<style>
  .list { list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:10px; }
  .list li { display:flex; justify-content:space-between; align-items:flex-start; gap:14px; padding:14px; border-radius:10px; background:#f8fafc; border-left:4px solid #94a3b8; }
  .list .sev-warn { border-color:#f59e0b; background:#fffbeb; }
  .list .sev-error { border-color:#ef4444; background:#fef2f2; }
  p { margin:4px 0; color:#475569; font-size:13px; }
  time { color:#94a3b8; font-size:12px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
