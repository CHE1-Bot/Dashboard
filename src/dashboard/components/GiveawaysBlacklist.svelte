<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let entries = [];
  let members = [];
  let loading = true;
  let draft = { user_id: '', reason: '' };

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [entries, members] = await Promise.all([api.get('/giveaways/blacklist'), api.get('/members')]);
    loading = false;
  }
  async function add() {
    if (!draft.user_id) return;
    const m = members.find(x => x.id === draft.user_id);
    await guildApi($currentGuildId).post('/giveaways/blacklist', {
      user_id: draft.user_id, username: m?.username || '?', reason: draft.reason,
    });
    draft = { user_id: '', reason: '' };
    await load();
  }
  async function remove(e) {
    if (!confirm(`Remove ${e.username} from blacklist?`)) return;
    await guildApi($currentGuildId).del('/giveaways/blacklist/' + e.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Giveaway blacklist" subtitle="These users cannot win giveaways.">
  <div slot="actions">
    <select bind:value={draft.user_id}>
      <option value="">Pick user…</option>
      {#each members as m}<option value={m.id}>{m.username}</option>{/each}
    </select>
    <input placeholder="Reason" bind:value={draft.reason} />
    <button class="btn primary" on:click={add} disabled={!draft.user_id}>Add</button>
  </div>
  {#if loading}<p>Loading…</p>{:else if entries.length === 0}<p class="empty">Nobody is blacklisted.</p>{:else}
    <table>
      <thead><tr><th>User</th><th>Reason</th><th>Added</th><th></th></tr></thead>
      <tbody>
        {#each entries as e}
          <tr>
            <td>{e.username}</td>
            <td>{e.reason || '—'}</td>
            <td>{formatDate(e.created_at)}</td>
            <td><button class="btn danger" on:click={() => remove(e)}>Remove</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  select, input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
