<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let rules = [];
  let loading = true;
  let draft = emptyDraft();

  function emptyDraft() { return { name: '', trigger: 'spam', action: 'delete', enabled: true, config: {} }; }

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    rules = await guildApi($currentGuildId).get('/moderation/automod');
    loading = false;
  }
  async function create() {
    if (!draft.name) return;
    await guildApi($currentGuildId).post('/moderation/automod', draft);
    draft = emptyDraft();
    await load();
  }
  async function toggle(r) {
    await guildApi($currentGuildId).patch('/moderation/automod/' + r.id, { ...r, enabled: !r.enabled });
    await load();
  }
  async function remove(r) {
    if (!confirm(`Delete rule "${r.name}"?`)) return;
    await guildApi($currentGuildId).del('/moderation/automod/' + r.id);
    await load();
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="New rule">
  <div class="form">
    <input bind:value={draft.name} placeholder="Rule name" />
    <select bind:value={draft.trigger}>
      <option value="spam">Spam</option>
      <option value="caps">Caps lock</option>
      <option value="links">Links</option>
      <option value="invites">Discord invites</option>
      <option value="words">Banned words</option>
      <option value="mention">Mass mentions</option>
    </select>
    <select bind:value={draft.action}>
      <option value="delete">Delete message</option>
      <option value="warn">Warn</option>
      <option value="mute">Mute</option>
      <option value="kick">Kick</option>
      <option value="ban">Ban</option>
    </select>
    <button class="btn primary" on:click={create}>Add rule</button>
  </div>
</Panel>

<Panel title="Rules">
  {#if loading}<p>Loading…</p>{:else if rules.length === 0}<p class="empty">No rules configured.</p>{:else}
    <table>
      <thead><tr><th>Name</th><th>Trigger</th><th>Action</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {#each rules as r}
          <tr class:off={!r.enabled}>
            <td>{r.name}</td>
            <td>{r.trigger}</td>
            <td>{r.action}</td>
            <td>
              <label class="switch">
                <input type="checkbox" checked={r.enabled} on:change={() => toggle(r)} />
                <span>{r.enabled ? 'On' : 'Off'}</span>
              </label>
            </td>
            <td><button class="btn danger" on:click={() => remove(r)}>Delete</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .form { display:flex; gap:10px; flex-wrap:wrap; }
  input, select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  input { flex:1; min-width:160px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  tr.off td { color:#94a3b8; }
  .switch { display:flex; align-items:center; gap:6px; font-size:13px; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
