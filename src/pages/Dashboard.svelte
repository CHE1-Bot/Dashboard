<script>
  import { onMount } from 'svelte';
  import { sidebarLinks, routes, getRoute } from '../dashboard/data/pages.js';
  import { user, guilds, currentGuild, currentGuildId, loadMe, loadGuilds, logout } from '../lib/stores.js';
  import { avatarUrl, guildIconUrl } from '../lib/ui.js';

  const componentModules = import.meta.glob('../dashboard/components/*.svelte');

  let hash = window.location.hash || '#/dashboard';
  let currentRoute = getRoute(normalize(hash));
  let route = routes[currentRoute] || routes.dashboard;
  let Component;
  let loading = true;
  let switcherOpen = false;
  let userMenuOpen = false;

  function normalize(h) {
    const raw = h.startsWith('#/') ? h.slice(2) : h.startsWith('#') ? h.slice(1) : h;
    if (raw === 'dashboard' || raw === '') return '#/dashboard';
    if (raw.startsWith('dashboard/')) return '#/' + raw.slice('dashboard/'.length);
    return '#/' + raw;
  }

  function updateRoute() {
    hash = window.location.hash || '#/dashboard';
    currentRoute = getRoute(normalize(hash));
    route = routes[currentRoute] || routes.dashboard;
    loadComponent();
  }

  async function loadComponent() {
    if (!route.componentName) { Component = null; return; }
    const key = `../dashboard/components/${route.componentName}.svelte`;
    const loader = componentModules[key];
    if (!loader) { Component = null; return; }
    try {
      const m = await loader();
      Component = m.default;
    } catch (e) {
      console.error('Component load failed:', e);
      Component = null;
    }
  }

  async function init() {
    try { await loadMe(); } catch (e) {
      loading = false;
      location.hash = '#/login';
      return;
    }
    await loadGuilds();
    if (!$currentGuildId) {
      location.hash = '#/servers';
      loading = false;
      return;
    }
    await loadComponent();
    loading = false;
  }

  onMount(() => {
    window.addEventListener('hashchange', updateRoute);
    init();
    return () => window.removeEventListener('hashchange', updateRoute);
  });

  function pickGuild(id) {
    currentGuildId.set(id);
    switcherOpen = false;
    // reload current component so it refetches for the new guild
    loadComponent();
  }

  async function handleLogout() {
    await logout();
    location.hash = '#/login';
  }

  const dashHref = (path) => '#/dashboard' + (path === 'dashboard' ? '' : '/' + path);
  const isSidebarActive = (id) => route.sidebarActive === id;
  const isRightNavActive = (path) => currentRoute === path;
</script>

{#if loading}
  <div class="loader">Loading…</div>
{:else}
<div class="dashboard">
  <aside class="sidebar">
    <a class="logo" href="#/">
      <img src="/images/Che1logo.png" alt="Che1" />
      <span>Che1</span>
    </a>

    <div class="guild-switcher" class:open={switcherOpen}>
      <button on:click={() => switcherOpen = !switcherOpen}>
        {#if $currentGuild}
          <span class="gicon">
            {#if guildIconUrl($currentGuild)}
              <img src={guildIconUrl($currentGuild)} alt="" />
            {:else}
              {$currentGuild.name.slice(0,2).toUpperCase()}
            {/if}
          </span>
          <span class="gname">{$currentGuild.name}</span>
        {:else}
          <span class="gname">No server selected</span>
        {/if}
        <i class="fa-solid fa-chevron-down"></i>
      </button>
      {#if switcherOpen}
        <div class="menu">
          {#each $guilds.filter(g => g.bot_present) as g}
            <button on:click={() => pickGuild(g.id)} class:active={g.id === $currentGuildId}>
              <span class="gicon">
                {#if guildIconUrl(g)}<img src={guildIconUrl(g)} alt="" />{:else}{g.name.slice(0,2).toUpperCase()}{/if}
              </span>
              <span>{g.name}</span>
            </button>
          {/each}
          <a class="menu-foot" href="#/servers">+ All servers</a>
        </div>
      {/if}
    </div>

    <nav>
      {#each sidebarLinks as link}
        <a href={dashHref(link.path)} class:selected={isSidebarActive(link.id)}>
          <i class={"fa-solid " + link.icon}></i>
          {link.label}
        </a>
      {/each}
    </nav>
  </aside>

  <div class="main">
    <header class="topbar">
      <h2>{route.pageTitle}</h2>
      <div class="user" class:open={userMenuOpen}>
        <button on:click={() => userMenuOpen = !userMenuOpen}>
          <img src={avatarUrl($user)} alt="" />
          <span>{$user?.username || 'User'}</span>
        </button>
        {#if userMenuOpen}
          <div class="user-menu">
            <a href="#/servers">Switch server</a>
            <a href="#/">Home</a>
            <button on:click={handleLogout}>Log out</button>
          </div>
        {/if}
      </div>
    </header>

    <div class="content">
      {#if Component}
        <svelte:component this={Component} />
      {:else if route.contentHtml}
        {@html route.contentHtml}
      {/if}
    </div>
  </div>

  {#if route.rightNav && route.rightNav.length}
    <aside class="right-sidebar">
      <nav>
        {#each route.rightNav as item}
          <a href={dashHref(item.path)} class:selected={isRightNavActive(item.path)}>
            <i class={"fa-solid " + item.icon}></i>
            {item.label}
          </a>
        {/each}
      </nav>
    </aside>
  {/if}
</div>
{/if}

<style>
  :global(body) { margin:0; font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Arial, sans-serif; background:#f4f6fb; }
  .loader { min-height:100vh; display:grid; place-items:center; color:#475569; }

  .dashboard { display:flex; min-height:100vh; }
  .sidebar { width:240px; background:#0f172a; color:#cbd5e1; display:flex; flex-direction:column; padding:20px; gap:18px; }
  .logo { display:flex; align-items:center; gap:10px; text-decoration:none; color:#fff; }
  .logo img { width:40px; height:40px; border-radius:50%; }
  .logo span { font-weight:700; font-size:18px; }

  .guild-switcher { position:relative; }
  .guild-switcher > button { width:100%; display:flex; align-items:center; gap:10px; background:#1e293b; color:#fff; border:none; border-radius:10px; padding:10px 12px; cursor:pointer; text-align:left; }
  .guild-switcher .gicon { width:28px; height:28px; border-radius:8px; background:#334155; display:grid; place-items:center; font-size:11px; font-weight:700; overflow:hidden; flex-shrink:0; }
  .guild-switcher .gicon img { width:100%; height:100%; object-fit:cover; }
  .guild-switcher .gname { flex:1; font-size:13px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
  .guild-switcher .menu { position:absolute; left:0; right:0; top:calc(100% + 6px); background:#1e293b; border-radius:10px; overflow:hidden; z-index:20; box-shadow:0 8px 24px rgba(0,0,0,0.3); }
  .guild-switcher .menu button { width:100%; display:flex; align-items:center; gap:8px; background:transparent; color:#e2e8f0; border:none; padding:8px 12px; cursor:pointer; font-size:13px; text-align:left; }
  .guild-switcher .menu button.active, .guild-switcher .menu button:hover { background:#334155; }
  .guild-switcher .menu .menu-foot { display:block; padding:10px 12px; font-size:12px; color:#94a3b8; border-top:1px solid #334155; text-decoration:none; }

  .sidebar nav { display:flex; flex-direction:column; gap:4px; }
  .sidebar nav a { color:#cbd5e1; text-decoration:none; padding:10px 12px; border-radius:8px; font-size:14px; display:flex; align-items:center; gap:10px; }
  .sidebar nav a:hover, .sidebar nav a.selected { background:#1e293b; color:#fff; }

  .main { flex:1; display:flex; flex-direction:column; min-width:0; }
  .topbar { display:flex; justify-content:space-between; align-items:center; padding:16px 32px; background:#fff; border-bottom:1px solid #e5e7eb; }
  .topbar h2 { margin:0; font-size:18px; color:#111827; }
  .user { position:relative; }
  .user > button { display:flex; align-items:center; gap:10px; background:transparent; border:none; cursor:pointer; color:#111827; font-size:14px; }
  .user img { width:36px; height:36px; border-radius:50%; }
  .user-menu { position:absolute; right:0; top:calc(100% + 6px); background:#fff; border:1px solid #e5e7eb; border-radius:10px; min-width:160px; overflow:hidden; box-shadow:0 8px 24px rgba(0,0,0,0.08); z-index:20; }
  .user-menu a, .user-menu button { display:block; width:100%; padding:10px 14px; text-decoration:none; color:#111827; background:transparent; border:none; text-align:left; cursor:pointer; font-size:14px; }
  .user-menu a:hover, .user-menu button:hover { background:#f1f5f9; }

  .content { padding:28px 32px; }

  .right-sidebar { width:220px; background:#0f172a; color:#cbd5e1; padding:20px; }
  .right-sidebar nav { display:flex; flex-direction:column; gap:4px; }
  .right-sidebar a { color:#cbd5e1; text-decoration:none; padding:10px 12px; border-radius:8px; font-size:14px; display:flex; align-items:center; gap:10px; }
  .right-sidebar a:hover, .right-sidebar a.selected { background:#1e293b; color:#fff; }
</style>
