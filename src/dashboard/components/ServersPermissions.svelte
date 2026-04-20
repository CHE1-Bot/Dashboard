<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  const MODULES = ['tickets', 'moderation', 'giveaways', 'leveling', 'settings', 'applications'];

  let perms = [];
  let roles = [];
  let loading = true;
  let pickingRole = '';
  let pickingModules = {};

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [perms, roles] = await Promise.all([api.get('/permissions'), api.get('/roles')]);
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  async function add() {
    if (!pickingRole) return;
    const r = roles.find(x => x.id === pickingRole);
    const chosen = Object.keys(pickingModules).filter(k => pickingModules[k]);
    await guildApi($currentGuildId).post('/permissions', {
      role_id: r.id, role_name: r.name, modules: chosen,
    });
    pickingRole = ''; pickingModules = {};
    await load();
  }
  async function remove(p) {
    if (!confirm(`Revoke ${p.role_name}'s dashboard access?`)) return;
    await guildApi($currentGuildId).del('/permissions/' + p.id);
    await load();
  }
</script>

<Panel title="Dashboard permissions" subtitle="Which roles can access which parts of the dashboard.">
  {#if loading}<p>Loading…</p>{:else}
    <table>
      <thead><tr><th>Role</th><th>Modules</th><th></th></tr></thead>
      <tbody>
        {#each perms as p}
          <tr>
            <td>{p.role_name}</td>
            <td>{p.modules.join(', ') || '—'}</td>
            <td><button class="btn danger" on:click={() => remove(p)}>Remove</button></td>
          </tr>
        {:else}
          <tr><td colspan="3" class="empty">No custom permissions — server admins have full access.</td></tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<Panel title="Grant role access">
  <div class="row">
    <select bind:value={pickingRole}>
      <option value="">Select role…</option>
      {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
    </select>
    <div class="mods">
      {#each MODULES as m}
        <label><input type="checkbox" bind:checked={pickingModules[m]} /> {m}</label>
      {/each}
    </div>
    <button class="btn primary" on:click={add} disabled={!pickingRole}>Grant</button>
  </div>
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .empty { color:#94a3b8; text-align:center; padding:20px; }
  .row { display:flex; flex-wrap:wrap; gap:12px; align-items:flex-start; }
  select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; min-width:180px; }
  .mods { display:flex; flex-wrap:wrap; gap:12px; flex:1; }
  .mods label { display:flex; align-items:center; gap:4px; text-transform:capitalize; font-size:13px; }
  .btn { padding:8px 14px; border:none; border-radius:8px; font-weight:600; cursor:pointer; font-size:13px; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
</style>
