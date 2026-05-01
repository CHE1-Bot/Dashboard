<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let categories = [];
  let channels = [];
  let roles = [];
  let forms = [];
  let loading = true;
  let editingId = null;
  let draft = emptyDraft();

  function emptyDraft() {
    return {
      name: '', channel_id: '',
      support_role_ids: [], mention_role_ids: [],
      welcome_message: 'Thanks for opening a ticket. Support will be with you shortly.',
      naming_pattern: 'ticket-{user}',
      claim_required: false, max_per_user: 1,
      form_id: 0, color: '#5865f2', emoji: '🎫', disabled: false,
    };
  }

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [categories, channels, roles, forms] = await Promise.all([
      api.get('/tickets/categories'),
      api.get('/channels'),
      api.get('/roles'),
      api.get('/tickets/forms'),
    ]);
    loading = false;
  }

  function startEdit(c) {
    editingId = c.id;
    draft = {
      ...emptyDraft(), ...c,
      support_role_ids: c.support_role_ids || [],
      mention_role_ids: c.mention_role_ids || [],
    };
  }
  function cancel() { editingId = null; draft = emptyDraft(); }

  async function save() {
    if (!draft.name.trim()) { alert('Name is required'); return; }
    const api = guildApi($currentGuildId);
    if (editingId) {
      await api.patch('/tickets/categories/' + editingId, draft);
    } else {
      await api.post('/tickets/categories', draft);
    }
    cancel();
    await load();
  }
  async function remove(c) {
    if (!confirm(`Delete category "${c.name}"? Panels using this category will need to be updated.`)) return;
    await guildApi($currentGuildId).del('/tickets/categories/' + c.id);
    await load();
  }

  function toggle(field, id) {
    const s = new Set(draft[field] || []);
    if (s.has(id)) s.delete(id); else s.add(id);
    draft[field] = [...s];
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title={editingId ? 'Edit category' : 'New category'}
       subtitle="Categories group tickets by purpose. Each panel button opens a ticket in one category.">
  <div class="form">
    <label>Name<input bind:value={draft.name} placeholder="General Support" /></label>
    <label>Emoji<input bind:value={draft.emoji} maxlength="4" /></label>
    <label>Color<input type="color" bind:value={draft.color} /></label>

    <label>Parent category channel
      <select bind:value={draft.channel_id}>
        <option value="">— inherit default —</option>
        {#each channels.filter(c => c.type === 'category') as c}<option value={c.id}>{c.name}</option>{/each}
      </select>
    </label>
    <label>Naming pattern
      <input bind:value={draft.naming_pattern} placeholder="ticket-{'{user}'}" />
    </label>
    <label>Max open per user
      <input type="number" min="1" max="25" bind:value={draft.max_per_user} />
    </label>

    <label class="full">Welcome message
      <textarea rows="2" bind:value={draft.welcome_message}></textarea>
      <small>Variables: <code>{'{user}'}</code> <code>{'{ticket}'}</code></small>
    </label>

    <label>Form (optional)
      <select bind:value={draft.form_id}>
        <option value={0}>None</option>
        {#each forms as f}<option value={f.id}>{f.name}</option>{/each}
      </select>
    </label>
    <label class="row"><input type="checkbox" bind:checked={draft.claim_required} /> Require claim before responding</label>
    <label class="row"><input type="checkbox" bind:checked={draft.disabled} /> Disable this category</label>

    <div class="full">
      <div class="label">Support roles (can see all tickets in this category)</div>
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={draft.support_role_ids.includes(r.id)} on:change={() => toggle('support_role_ids', r.id)} />{r.name}</label>
      {/each}
    </div>
    <div class="full">
      <div class="label">Mention roles when ticket opens</div>
      {#each roles.filter(r => r.name !== '@everyone') as r}
        <label class="chip"><input type="checkbox" checked={draft.mention_role_ids.includes(r.id)} on:change={() => toggle('mention_role_ids', r.id)} />{r.name}</label>
      {/each}
    </div>

    <div class="full row-actions">
      <button class="btn primary" on:click={save}>{editingId ? 'Save changes' : 'Create category'}</button>
      {#if editingId}<button class="btn" on:click={cancel}>Cancel</button>{/if}
    </div>
  </div>
</Panel>

<Panel title="Categories">
  {#if loading}<p>Loading…</p>{:else if categories.length === 0}<p class="empty">No categories yet.</p>{:else}
    <table>
      <thead><tr><th></th><th>Name</th><th>Claim</th><th>Form</th><th>Roles</th><th></th></tr></thead>
      <tbody>
        {#each categories as c}
          <tr class:disabled={c.disabled}>
            <td><span class="emoji" style={`background:${c.color}`}>{c.emoji || '🎫'}</span></td>
            <td>
              <strong>{c.name}</strong>
              {#if c.disabled}<span class="tag">disabled</span>{/if}
              <div class="muted">#{channels.find(ch => ch.id === c.channel_id)?.name || 'inherits'}</div>
            </td>
            <td>{c.claim_required ? 'required' : 'optional'}</td>
            <td>{forms.find(f => f.id === c.form_id)?.name || '—'}</td>
            <td>
              {#each c.support_role_ids || [] as rid}
                <span class="rchip">{roles.find(r => r.id === rid)?.name || rid}</span>
              {/each}
            </td>
            <td class="actions">
              <button class="btn" on:click={() => startEdit(c)}>Edit</button>
              <button class="btn danger" on:click={() => remove(c)}>Delete</button>
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
  .form label.row { flex-direction:row; align-items:center; padding:8px 12px; background:#f8fafc; border-radius:8px; }
  .form input, .form select, .form textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  .form input[type="color"] { padding:0; height:38px; }
  .form small { color:#94a3b8; font-size:11px; }
  .label { font-size:13px; color:#475569; margin-bottom:8px; font-weight:600; }
  .chip { display:inline-flex; align-items:center; gap:6px; padding:6px 12px; border:1px solid #e5e7eb; border-radius:999px; font-size:13px; margin:0 6px 6px 0; }
  .row-actions { display:flex; gap:8px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; vertical-align:middle; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  tr.disabled td { opacity:0.55; }
  .emoji { display:inline-grid; place-items:center; width:34px; height:34px; border-radius:8px; color:#fff; font-size:16px; }
  .muted { color:#94a3b8; font-size:11px; }
  .tag { display:inline-block; margin-left:6px; padding:1px 6px; background:#fee2e2; color:#b91c1c; border-radius:4px; font-size:10px; font-weight:600; }
  .rchip { display:inline-block; padding:2px 8px; background:#f1f5f9; border-radius:6px; font-size:11px; margin:0 4px 4px 0; }
  .actions { display:flex; gap:6px; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
