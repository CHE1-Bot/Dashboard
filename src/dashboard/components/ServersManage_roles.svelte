<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let roles = [];
  let members = [];
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [roles, members] = await Promise.all([api.get('/roles'), api.get('/members')]);
    loading = false;
  }

  onMount(load);
  $: if ($currentGuildId) load();

  function color(r) { return '#' + r.color.toString(16).padStart(6, '0'); }
  function membersOf(roleId) { return members.filter(m => m.roles.includes(roleId)).length; }
</script>

<Panel title="Roles" subtitle="Roles synced from Discord. Role creation is managed in Discord itself.">
  {#if loading}
    <p>Loading…</p>
  {:else}
    <table>
      <thead><tr><th>Role</th><th>Members</th><th>Position</th><th>Color</th></tr></thead>
      <tbody>
        {#each [...roles].sort((a,b) => b.position - a.position) as r}
          <tr>
            <td><span class="dot" style="background:{color(r)}"></span>{r.name}</td>
            <td>{membersOf(r.id)}</td>
            <td>{r.position}</td>
            <td><code>{color(r)}</code></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .dot { display:inline-block; width:10px; height:10px; border-radius:50%; margin-right:8px; vertical-align:middle; }
  code { font-size:12px; color:#64748b; }
</style>
