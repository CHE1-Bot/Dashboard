<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let staff = [];
  let members = [];
  let roles = [];
  let picking = '';
  let loading = true;

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [staff, members, roles] = await Promise.all([api.get('/tickets/staff'), api.get('/members'), api.get('/roles')]);
    loading = false;
  }
  async function add() {
    if (!picking) return;
    const m = members.find(x => x.id === picking);
    await guildApi($currentGuildId).post('/tickets/staff', {
      user_id: m.id, username: m.username, roles: m.roles,
    });
    picking = '';
    await load();
  }
  async function remove(s) {
    if (!confirm(`Remove ${s.username} from ticket staff?`)) return;
    await guildApi($currentGuildId).del('/tickets/staff/' + s.user_id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: candidates = members.filter(m => !staff.some(s => s.user_id === m.id));
</script>

<Panel title="Ticket staff" subtitle="Who can claim, reply to, and close tickets.">
  <div slot="actions">
    <select bind:value={picking}>
      <option value="">Add member…</option>
      {#each candidates as m}<option value={m.id}>{m.username}</option>{/each}
    </select>
    <button class="btn primary" on:click={add} disabled={!picking}>Add</button>
  </div>
  {#if loading}<p>Loading…</p>{:else}
    <table>
      <thead><tr><th>User</th><th>Roles</th><th></th></tr></thead>
      <tbody>
        {#each staff as s}
          <tr>
            <td>{s.username}</td>
            <td>{s.roles.map(rid => roles.find(r => r.id === rid)?.name || rid).join(', ')}</td>
            <td><button class="btn danger" on:click={() => remove(s)}>Remove</button></td>
          </tr>
        {:else}<tr><td colspan="3" class="empty">No staff assigned.</td></tr>{/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .empty { text-align:center; color:#94a3b8; padding:20px; }
  select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; min-width:160px; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
</style>
