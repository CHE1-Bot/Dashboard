<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let snippets = [];
  let loading = true;
  let editingId = null;
  let draft = { name: '', content: '' };

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    snippets = await guildApi($currentGuildId).get('/tickets/snippets');
    loading = false;
  }
  function startEdit(s) {
    editingId = s.id;
    draft = { name: s.name, content: s.content };
  }
  function cancel() { editingId = null; draft = { name: '', content: '' }; }
  async function save() {
    if (!draft.name.trim() || !draft.content.trim()) { alert('Name and content are required'); return; }
    const api = guildApi($currentGuildId);
    if (editingId) await api.patch('/tickets/snippets/' + editingId, draft);
    else await api.post('/tickets/snippets', draft);
    cancel();
    await load();
  }
  async function remove(s) {
    if (!confirm(`Delete snippet /${s.name}?`)) return;
    await guildApi($currentGuildId).del('/tickets/snippets/' + s.id);
    await load();
  }
  function copy(s) { navigator.clipboard?.writeText(s.content); }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title={editingId ? 'Edit snippet' : 'New snippet'}
       subtitle="Pre-written replies staff can fire with /snippet send <name>.">
  <div class="form">
    <label>Name (no spaces)
      <input bind:value={draft.name} placeholder="welcome" pattern="[A-Za-z0-9_-]+" />
    </label>
    <label class="full">Content
      <textarea rows="4" bind:value={draft.content} placeholder="Hi {'{user}'}, how can we help?"></textarea>
      <small>Variables: <code>{'{user}'}</code> <code>{'{ticket}'}</code></small>
    </label>
    <div class="full row-actions">
      <button class="btn primary" on:click={save}>{editingId ? 'Save changes' : 'Create snippet'}</button>
      {#if editingId}<button class="btn" on:click={cancel}>Cancel</button>{/if}
    </div>
  </div>
</Panel>

<Panel title="Snippets">
  {#if loading}<p>Loading…</p>{:else if snippets.length === 0}<p class="empty">No snippets yet.</p>{:else}
    <table>
      <thead><tr><th>Name</th><th>Content</th><th>Updated</th><th></th></tr></thead>
      <tbody>
        {#each snippets as s}
          <tr>
            <td><code class="name">/{s.name}</code></td>
            <td class="content">{s.content}</td>
            <td class="muted">{relativeTime(s.updated_at)}</td>
            <td class="actions">
              <button class="btn" on:click={() => copy(s)}>Copy</button>
              <button class="btn" on:click={() => startEdit(s)}>Edit</button>
              <button class="btn danger" on:click={() => remove(s)}>Delete</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .form { display:grid; grid-template-columns:repeat(2, 1fr); gap:12px; }
  .form .full { grid-column: 1 / -1; }
  .form label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .form input, .form textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  .form small { color:#94a3b8; font-size:11px; }
  .row-actions { display:flex; gap:8px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; vertical-align:top; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .name { background:#f1f5f9; padding:2px 8px; border-radius:6px; font-size:13px; color:#0f172a; }
  .content { color:#334155; max-width:480px; white-space:pre-wrap; }
  .muted { color:#94a3b8; font-size:12px; }
  .actions { display:flex; gap:6px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
