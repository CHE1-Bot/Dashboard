<script>
  import { onMount } from 'svelte';
  import { guilds, currentGuildId, loadGuilds } from '../lib/stores.js';
  import { api } from '../lib/api.js';
  import { guildIconUrl } from '../lib/ui.js';

  let loading = true;
  let error = '';

  onMount(async () => {
    try { await loadGuilds(); } catch (e) { error = e.message; }
    loading = false;
  });

  function pick(g) {
    currentGuildId.set(g.id);
    location.hash = '#/dashboard/overview';
  }

  async function invite(g) {
    const { invite_url } = await api.get(`/guilds/${g.id}/enable`);
    window.open(invite_url, '_blank', 'width=520,height=720');
  }
</script>

<div class="wrap">
  <div class="inner">
    <h1>Select a server</h1>
    <p class="sub">Choose the Discord server you want to manage.</p>

    {#if loading}
      <div class="state">Loading your servers…</div>
    {:else if error}
      <div class="state err">{error}</div>
    {:else if $guilds.length === 0}
      <div class="state">You're not in any servers we can see.</div>
    {:else}
      <div class="grid">
        {#each $guilds as g (g.id)}
          <div class="card" class:disabled={!g.bot_present}>
            <div class="icon">
              {#if guildIconUrl(g)}
                <img src={guildIconUrl(g)} alt="" />
              {:else}
                {g.name.slice(0, 2).toUpperCase()}
              {/if}
            </div>
            <div class="name">{g.name}</div>
            <div class="meta">{g.member_count} members</div>
            {#if g.bot_present}
              <button class="btn primary" on:click={() => pick(g)}>Manage</button>
            {:else}
              <button class="btn" on:click={() => invite(g)}>Invite Bot</button>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .wrap { min-height:100vh; background:#0f172a; padding:60px 20px; }
  .inner { max-width:960px; margin:0 auto; color:#e2e8f0; }
  h1 { margin:0 0 6px; font-size:26px; color:#fff; }
  .sub { margin:0 0 28px; color:#94a3b8; }
  .grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(220px, 1fr)); gap:16px; }
  .card { background:#fff; color:#111827; border-radius:12px; padding:20px; text-align:center; display:flex; flex-direction:column; align-items:center; gap:8px; }
  .card.disabled { opacity:0.85; }
  .icon { width:64px; height:64px; border-radius:50%; background:#e5e7eb; display:grid; place-items:center; font-weight:700; color:#475569; overflow:hidden; }
  .icon img { width:100%; height:100%; object-fit:cover; }
  .name { font-weight:600; }
  .meta { color:#6b7280; font-size:12px; }
  .btn { margin-top:6px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; padding:8px 16px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn:hover { filter:brightness(0.95); }
  .state { padding:40px; text-align:center; background:rgba(255,255,255,0.04); border-radius:12px; color:#cbd5e1; }
  .state.err { color:#fca5a5; }
</style>
