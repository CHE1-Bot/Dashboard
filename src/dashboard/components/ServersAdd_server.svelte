<script>
  import { onMount } from 'svelte';
  import { guilds, currentGuildId, loadGuilds } from '../../lib/stores.js';
  import { api } from '../../lib/api.js';
  import { guildIconUrl } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let loading = true;
  onMount(async () => { await loadGuilds(); loading = false; });

  async function invite(g) {
    const { invite_url } = await api.get(`/guilds/${g.id}/enable`);
    window.open(invite_url, '_blank', 'width=520,height=720');
    setTimeout(async () => { await loadGuilds(); }, 3000);
  }

  function manage(g) {
    currentGuildId.set(g.id);
    location.hash = '#/dashboard/overview';
  }
</script>

<Panel title="Your servers" subtitle="Invite Che1 to a server or switch between the ones it already manages.">
  {#if loading}
    <p>Loading…</p>
  {:else}
    <div class="grid">
      {#each $guilds as g (g.id)}
        <div class="card">
          <div class="icon">
            {#if guildIconUrl(g)}<img src={guildIconUrl(g)} alt="" />{:else}{g.name.slice(0,2).toUpperCase()}{/if}
          </div>
          <div class="info">
            <div class="name">{g.name}</div>
            <div class="meta">{g.member_count} members {#if g.owner}· owner{/if}</div>
          </div>
          {#if g.bot_present}
            <button class="btn" on:click={() => manage(g)}>Manage</button>
          {:else}
            <button class="btn primary" on:click={() => invite(g)}>Invite Bot</button>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</Panel>

<style>
  .grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(260px, 1fr)); gap:12px; }
  .card { display:flex; align-items:center; gap:12px; padding:12px; border:1px solid #e5e7eb; border-radius:10px; }
  .icon { width:44px; height:44px; border-radius:10px; background:#e5e7eb; display:grid; place-items:center; font-weight:700; color:#475569; overflow:hidden; flex-shrink:0; font-size:13px; }
  .icon img { width:100%; height:100%; object-fit:cover; }
  .info { flex:1; min-width:0; }
  .name { font-weight:600; font-size:14px; color:#111827; white-space:nowrap; overflow:hidden; text-overflow:ellipsis; }
  .meta { font-size:12px; color:#6b7280; }
  .btn { padding:6px 12px; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
</style>
