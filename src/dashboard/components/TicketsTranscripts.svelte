<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let tickets = [];
  let selected = null;
  let transcript = null;
  let query = '';
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    tickets = await guildApi($currentGuildId).get('/tickets');
    loading = false;
  }
  async function open(t) {
    selected = t;
    transcript = null;
    const data = await guildApi($currentGuildId).get(`/tickets/${t.id}/transcript`);
    transcript = data;
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: shown = query
    ? tickets.filter(t => String(t.id) === query || t.username.toLowerCase().includes(query.toLowerCase()))
    : tickets.filter(t => t.status === 'closed');
</script>

<div class="split">
  <Panel title="Closed tickets">
    <div slot="actions"><input placeholder="Search ID or user" bind:value={query} /></div>
    {#if loading}<p>Loading…</p>{:else}
      <ul class="list">
        {#each shown as t}
          <li class:active={selected?.id === t.id}>
            <button on:click={() => open(t)}>
              <strong>#{t.id}</strong> {t.subject}
              <small>{t.username} · {formatDate(t.closed_at || t.created_at)}</small>
            </button>
          </li>
        {:else}<li class="empty">No tickets.</li>{/each}
      </ul>
    {/if}
  </Panel>

  <Panel title={selected ? `Ticket #${selected.id}` : 'Select a ticket'}>
    {#if selected && transcript}
      <pre>{transcript.content}</pre>
    {:else if selected}
      <p>Loading transcript…</p>
    {:else}
      <p class="empty">Pick a ticket from the list to view its transcript.</p>
    {/if}
  </Panel>
</div>

<style>
  .split { display:grid; grid-template-columns:340px 1fr; gap:20px; }
  @media (max-width: 860px) { .split { grid-template-columns:1fr; } }
  input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  .list { list-style:none; margin:0; padding:0; max-height:480px; overflow:auto; }
  .list li { border-bottom:1px solid #f1f5f9; }
  .list button { display:block; width:100%; background:transparent; border:none; text-align:left; padding:10px 12px; font-size:14px; cursor:pointer; color:#111827; }
  .list button:hover, .list li.active button { background:#f1f5f9; }
  .list small { display:block; color:#94a3b8; font-size:12px; margin-top:2px; }
  .list .empty { padding:20px; text-align:center; color:#94a3b8; }
  pre { background:#0f172a; color:#e2e8f0; padding:16px; border-radius:8px; white-space:pre-wrap; word-break:break-word; font-size:13px; line-height:1.6; max-height:600px; overflow:auto; }
  .empty { color:#94a3b8; }
</style>
