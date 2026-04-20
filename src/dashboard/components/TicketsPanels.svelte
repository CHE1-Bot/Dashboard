<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let panels = [];
  let channels = [];
  let forms = [];
  let loading = true;
  let draft = emptyDraft();

  function emptyDraft() {
    return { channel_id: '', title: '', message: '', color: '#3498db',
      buttons: [{ label: 'Open Ticket', style: 'primary', category: 'general' }] };
  }

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [panels, channels, forms] = await Promise.all([
      api.get('/tickets/panels'), api.get('/channels'), api.get('/tickets/forms'),
    ]);
    loading = false;
  }
  async function create() {
    if (!draft.channel_id || !draft.title) return;
    await guildApi($currentGuildId).post('/tickets/panels', draft);
    draft = emptyDraft();
    await load();
  }
  async function remove(p) {
    if (!confirm('Delete panel?')) return;
    await guildApi($currentGuildId).del('/tickets/panels/' + p.id);
    await load();
  }
  function addButton() {
    draft.buttons = [...draft.buttons, { label: 'New Button', style: 'secondary', category: 'general' }];
  }
  function removeButton(i) { draft.buttons = draft.buttons.filter((_, idx) => idx !== i); }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="New ticket panel" subtitle="Posts a message in a channel with buttons that open tickets.">
  <div class="form">
    <label>Channel
      <select bind:value={draft.channel_id}>
        <option value="">—</option>
        {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
      </select>
    </label>
    <label>Title<input bind:value={draft.title} /></label>
    <label>Color<input type="color" bind:value={draft.color} /></label>
    <label class="full">Message<textarea rows="3" bind:value={draft.message}></textarea></label>

    <div class="full">
      <div class="label">Buttons <button class="btn" on:click={addButton}>+ Add button</button></div>
      {#each draft.buttons as b, i}
        <div class="btn-row">
          <input bind:value={b.label} placeholder="Label" />
          <select bind:value={b.style}>
            <option value="primary">Primary</option>
            <option value="secondary">Secondary</option>
            <option value="success">Success</option>
            <option value="danger">Danger</option>
          </select>
          <input bind:value={b.category} placeholder="Category" />
          <select bind:value={b.form_id}>
            <option value={0}>No form</option>
            {#each forms as f}<option value={f.id}>{f.name}</option>{/each}
          </select>
          <button class="btn danger" on:click={() => removeButton(i)}>×</button>
        </div>
      {/each}
    </div>
    <div class="full">
      <button class="btn primary" on:click={create}>Create panel</button>
    </div>
  </div>
</Panel>

<Panel title="Existing panels">
  {#if loading}<p>Loading…</p>{:else if panels.length === 0}<p class="empty">No panels yet.</p>{:else}
    <table>
      <thead><tr><th>Title</th><th>Channel</th><th>Buttons</th><th></th></tr></thead>
      <tbody>
        {#each panels as p}
          <tr>
            <td>{p.title}</td>
            <td>#{channels.find(c => c.id === p.channel_id)?.name || '?'}</td>
            <td>{p.buttons.length}</td>
            <td><button class="btn danger" on:click={() => remove(p)}>Delete</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .form { display:grid; grid-template-columns:repeat(3, 1fr); gap:12px; }
  .form .full { grid-column: 1 / -1; }
  .form label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .form input, .form select, .form textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  .form input[type="color"] { padding:0; height:38px; }
  .label { font-size:13px; color:#475569; display:flex; align-items:center; justify-content:space-between; margin-bottom:8px; }
  .btn-row { display:grid; grid-template-columns:1fr 140px 140px 140px auto; gap:6px; margin-bottom:8px; }
  .btn-row input, .btn-row select { padding:6px 8px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
