<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let rewards = [];
  let roles = [];
  let loading = true;
  let draft = { level: 5, role_id: '', stackable: false };

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [rewards, roles] = await Promise.all([api.get('/leveling/rewards'), api.get('/roles')]);
    loading = false;
  }
  async function create() {
    if (!draft.role_id) return;
    const r = roles.find(x => x.id === draft.role_id);
    await guildApi($currentGuildId).post('/leveling/rewards', {
      level: Number(draft.level) || 1, role_id: r.id, role_name: r.name, stackable: draft.stackable,
    });
    draft = { level: 5, role_id: '', stackable: false };
    await load();
  }
  async function remove(rw) {
    if (!confirm(`Remove reward for level ${rw.level}?`)) return;
    await guildApi($currentGuildId).del('/leveling/rewards/' + rw.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
  $: sorted = [...rewards].sort((a, b) => a.level - b.level);
</script>

<Panel title="Add reward">
  <div class="form">
    <label>Level<input type="number" min="1" max="500" bind:value={draft.level} /></label>
    <label>Role
      <select bind:value={draft.role_id}>
        <option value="">—</option>
        {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
      </select>
    </label>
    <label class="stack"><input type="checkbox" bind:checked={draft.stackable} /> Stackable (keep earlier rewards)</label>
    <button class="btn primary" on:click={create} disabled={!draft.role_id}>Add</button>
  </div>
</Panel>

<Panel title="Rewards">
  {#if loading}<p>Loading…</p>{:else if sorted.length === 0}<p class="empty">No rewards configured.</p>{:else}
    <table>
      <thead><tr><th>Level</th><th>Role</th><th>Stackable</th><th></th></tr></thead>
      <tbody>
        {#each sorted as r}
          <tr>
            <td>{r.level}</td>
            <td>{r.role_name}</td>
            <td>{r.stackable ? 'Yes' : 'No'}</td>
            <td><button class="btn danger" on:click={() => remove(r)}>Remove</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .form { display:grid; grid-template-columns:120px 1fr auto auto; gap:12px; align-items:flex-end; }
  .form label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .form input, .form select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .form .stack { flex-direction:row; align-items:center; gap:6px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
