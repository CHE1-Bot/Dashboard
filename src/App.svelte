<script>
  import { onMount } from 'svelte';
  import Home from './pages/Home.svelte';
  import About from './pages/About.svelte';
  import Premium from './pages/Premium.svelte';
  import Documentation from './pages/Documentation.svelte';
  import Contact from './pages/Contact.svelte';
  import TOS from './pages/TOS.svelte';
  import PrivacyPolicy from './pages/PrivacyPolicy.svelte';
  import Dashboard from './pages/Dashboard.svelte';
  import Login from './pages/Login.svelte';
  import GuildSelect from './pages/GuildSelect.svelte';

  let route = parse(window.location.hash);

  function parse(hash) {
    return (hash || '').replace(/^#\/?/, '');
  }

  function update() { route = parse(window.location.hash); }

  onMount(() => {
    window.addEventListener('hashchange', update);
    return () => window.removeEventListener('hashchange', update);
  });

  $: first = route.split('/')[0];
</script>

{#if first === 'login'}
  <Login />
{:else if first === 'servers' || route === 'dashboard/servers'}
  <GuildSelect />
{:else if first === 'dashboard'}
  <Dashboard />
{:else if first === 'about'}
  <About />
{:else if first === 'premium'}
  <Premium />
{:else if first === 'documentation'}
  <Documentation />
{:else if first === 'contact'}
  <Contact />
{:else if first === 'tos'}
  <TOS />
{:else if first === 'privacy'}
  <PrivacyPolicy />
{:else}
  <Home />
{/if}
