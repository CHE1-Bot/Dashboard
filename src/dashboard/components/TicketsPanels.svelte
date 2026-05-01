<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let panels = [];
  let channels = [];
  let categories = [];
  let loading = true;
  let draft = emptyDraft();
  let editingId = null;

  function emptyDraft() {
    return {
      channel_id: '', title: 'Support', message: 'Click a button below to open a ticket.',
      color: '#5865f2',
      buttons: [{ label: 'Open Ticket', style: 'primary', emoji: '🎫', category_id: 0 }],
    };
  }

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [panels, channels, categories] = await Promise.all([
      api.get('/tickets/panels'),
      api.get('/channels'),
      api.get('/tickets/categories'),
    ]);
    if (categories.length && !draft.buttons[0].category_id) {
      draft.buttons[0].category_id = categories[0].id;
    }
    loading = false;
  }

  async function save() {
    if (!draft.channel_id || !draft.title) return;
    if (!draft.buttons.length) { alert('Add at least one button'); return; }
    if (draft.buttons.some(b => !b.category_id)) { alert('Every button needs a category'); return; }
    const api = guildApi($currentGuildId);
    if (editingId) {
      await api.patch('/tickets/panels/' + editingId, draft);
    } else {
      await api.post('/tickets/panels', draft);
    }
    editingId = null;
    draft = emptyDraft();
    if (categories.length) draft.buttons[0].category_id = categories[0].id;
    await load();
  }

  function edit(p) {
    editingId = p.id;
    draft = JSON.parse(JSON.stringify(p));
    draft.buttons = (draft.buttons || []).map(b => ({
      label: b.label || '', style: b.style || 'primary',
      emoji: b.emoji || '', category_id: b.category_id || (categories[0]?.id ?? 0),
    }));
  }
  function cancel() { editingId = null; draft = emptyDraft(); if (categories.length) draft.buttons[0].category_id = categories[0].id; }

  async function remove(p) {
    if (!confirm(`Delete panel "${p.title}"?`)) return;
    await guildApi($currentGuildId).del('/tickets/panels/' + p.id);
    await load();
  }

  async function deploy(p) {
    if (!confirm(`Deploy panel "${p.title}" to #${channels.find(c => c.id === p.channel_id)?.name}?`)) return;
    await guildApi($currentGuildId).post('/tickets/panels/' + p.id + '/deploy', {});
  }

  function addButton() {
    draft.buttons = [...draft.buttons, {
      label: 'New Button', style: 'secondary', emoji: '',
      category_id: categories[0]?.id ?? 0,
    }];
  }
  function removeButton(i) { draft.buttons = draft.buttons.filter((_, idx) => idx !== i); }

  onMount(load);
  $: if ($currentGuildId) load();
  $: catName = (id) => categories.find(c => c.id === id)?.name || '—';
</script>

<Panel title={editingId ? 'Edit panel' : 'New ticket panel'}
       subtitle="Posts a message in a channel with buttons that open tickets in a category.">
  {#if categories.length === 0}
    <div class="warn">
      You don't have any ticket categories yet — <a href="#/dashboard/tickets/categories">create one first</a>.
    </div>
  {/if}
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
          <input bind:value={b.emoji} placeholder="Emoji" maxlength="4" />
          <select bind:value={b.style}>
            <option value="primary">Primary (blue)</option>
            <option value="secondary">Secondary (grey)</option>
            <option value="success">Success (green)</option>
            <option value="danger">Danger (red)</option>
          </select>
          <select bind:value={b.category_id}>
            <option value={0}>— select category —</option>
            {#each categories as c}<option value={c.id}>{c.name}</option>{/each}
          </select>
          <button class="btn danger" on:click={() => removeButton(i)} disabled={draft.buttons.length === 1}>×</button>
        </div>
      {/each}
    </div>

    <div class="full preview">
      <div class="preview-title">Preview</div>
      <div class="embed" style={`border-left:4px solid ${draft.color}`}>
        <div class="ettl">{draft.title}</div>
        <div class="emsg">{draft.message}</div>
        <div class="ebtns">
          {#each draft.buttons as b}
            <span class={'pbtn ' + (b.style || 'secondary')}>{b.emoji ? b.emoji + ' ' : ''}{b.label || 'Button'}</span>
          {/each}
        </div>
      </div>
    </div>

    <div class="full row-actions">
      <button class="btn primary" on:click={save}>{editingId ? 'Save changes' : 'Create panel'}</button>
      {#if editingId}<button class="btn" on:click={cancel}>Cancel</button>{/if}
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
            <td><span class="dot" style={`background:${p.color || '#5865f2'}`}></span> {p.title}</td>
            <td>#{channels.find(c => c.id === p.channel_id)?.name || '?'}</td>
            <td>
              {#each p.buttons as b}
                <span class="bchip">{b.label} → {catName(b.category_id)}</span>
              {/each}
            </td>
            <td class="actions">
              <button class="btn" on:click={() => deploy(p)}>Deploy</button>
              <button class="btn" on:click={() => edit(p)}>Edit</button>
              <button class="btn danger" on:click={() => remove(p)}>Delete</button>
            </td>
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
  .btn-row { display:grid; grid-template-columns:1fr 80px 160px 1fr auto; gap:6px; margin-bottom:8px; }
  .btn-row input, .btn-row select { padding:6px 8px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  .row-actions { display:flex; gap:8px; }
  .preview { border-top:1px solid #e5e7eb; padding-top:12px; }
  .preview-title { font-size:11px; color:#94a3b8; text-transform:uppercase; margin-bottom:6px; }
  .embed { background:#2f3136; color:#dcddde; padding:14px 16px; border-radius:8px; max-width:520px; }
  .ettl { color:#fff; font-weight:700; margin-bottom:6px; }
  .emsg { white-space:pre-wrap; font-size:14px; color:#dcddde; }
  .ebtns { margin-top:10px; display:flex; gap:6px; flex-wrap:wrap; }
  .pbtn { padding:6px 12px; border-radius:6px; font-size:13px; font-weight:600; }
  .pbtn.primary { background:#5865f2; color:#fff; }
  .pbtn.secondary { background:#4f545c; color:#fff; }
  .pbtn.success { background:#3ba55c; color:#fff; }
  .pbtn.danger { background:#ed4245; color:#fff; }
  .dot { display:inline-block; width:8px; height:8px; border-radius:50%; margin-right:6px; vertical-align:middle; }
  .bchip { display:inline-block; padding:2px 8px; border-radius:6px; background:#f1f5f9; font-size:11px; margin:0 4px 4px 0; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; vertical-align:top; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .actions { display:flex; gap:6px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .btn:disabled { opacity:0.5; cursor:not-allowed; }
  .empty { color:#94a3b8; }
  .warn { background:#fef3c7; color:#92400e; padding:12px 16px; border-radius:8px; margin-bottom:14px; font-size:13px; }
  .warn a { color:#92400e; font-weight:600; }
</style>
