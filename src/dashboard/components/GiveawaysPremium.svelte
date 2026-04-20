<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let data = null;

  async function load() {
    if (!$currentGuildId) return;
    data = await guildApi($currentGuildId).get('/giveaways/premium');
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !data}<p>Loading…</p>{:else}
  <Panel title="Giveaway Premium" subtitle={data.premium ? 'Active for this server' : 'Upgrade to unlock premium features'}>
    <ul class="benefits">
      {#each data.benefits as b}<li><i class="fa-solid fa-check"></i> {b}</li>{/each}
    </ul>
    <a class="btn primary" href={data.upgrade_url}>{data.premium ? 'Manage subscription' : 'Upgrade to Premium'}</a>
  </Panel>
{/if}

<style>
  .benefits { list-style:none; margin:0 0 20px; padding:0; display:flex; flex-direction:column; gap:10px; }
  .benefits li { display:flex; align-items:center; gap:10px; color:#111827; }
  .benefits i { color:#10b981; }
  .btn { display:inline-block; padding:10px 20px; background:#5865f2; color:#fff; border-radius:8px; font-weight:600; text-decoration:none; }
  .btn:hover { filter:brightness(0.95); }
</style>
