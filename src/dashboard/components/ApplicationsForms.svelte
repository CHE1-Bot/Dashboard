<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let forms = [];
  let channels = [];
  let roles = [];
  let loading = true;
  let editingId = null;
  let draft = emptyDraft();

  function emptyDraft() {
    return {
      name: '', description: '', emoji: '📝', color: '#5865f2',
      questions: [
        { id: 'q1', label: 'Why do you want to apply?', type: 'paragraph', required: true, placeholder: '' },
      ],
      submission_channel_id: '', accepted_role_id: '',
      required_role_id: '', blocked_role_ids: [],
      cooldown_hours: 24, account_age_days: 7,
      accept_dm_template:
        'Congrats {user}! Your application for **{form}** in **{guild}** has been accepted.',
      reject_dm_template:
        'Thanks for applying to **{form}** in **{guild}**, {user}. Unfortunately your application has been declined.\n\nReason: {reason}',
      enabled: true,
    };
  }

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [forms, channels, roles] = await Promise.all([
      api.get('/applications/forms'),
      api.get('/channels'),
      api.get('/roles'),
    ]);
    loading = false;
  }

  function startEdit(f) {
    editingId = f.id;
    draft = { ...emptyDraft(), ...f, blocked_role_ids: f.blocked_role_ids || [] };
  }
  function cancel() { editingId = null; draft = emptyDraft(); }

  async function save() {
    if (!draft.name.trim()) { alert('Name is required'); return; }
    if (!draft.questions.length) { alert('Add at least one question'); return; }
    const api = guildApi($currentGuildId);
    if (editingId) await api.patch('/applications/forms/' + editingId, draft);
    else await api.post('/applications/forms', draft);
    cancel();
    await load();
  }

  async function remove(f) {
    if (!confirm(`Delete form "${f.name}"?`)) return;
    await guildApi($currentGuildId).del('/applications/forms/' + f.id);
    await load();
  }

  function addQuestion() {
    draft.questions = [...draft.questions, {
      id: 'q' + (draft.questions.length + 1),
      label: 'New question', type: 'short', required: false, placeholder: '',
    }];
  }
  function removeQuestion(i) {
    draft.questions = draft.questions.filter((_, idx) => idx !== i);
  }
  function moveQuestion(i, dir) {
    const j = i + dir;
    if (j < 0 || j >= draft.questions.length) return;
    const arr = [...draft.questions];
    [arr[i], arr[j]] = [arr[j], arr[i]];
    draft.questions = arr;
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title={editingId ? 'Edit form' : 'New application form'}
       subtitle="Like Appy: pick questions, where submissions go, and what happens on accept/reject.">
  <div class="grid">
    <label>Name<input bind:value={draft.name} placeholder="Moderator Application" /></label>
    <label>Emoji<input bind:value={draft.emoji} maxlength="4" /></label>
    <label>Color<input type="color" bind:value={draft.color} /></label>
    <label class="full">Description<textarea rows="2" bind:value={draft.description}
      placeholder="Shown to applicants before they start the form."></textarea></label>

    <label>Submission channel
      <select bind:value={draft.submission_channel_id}>
        <option value="">— pick a channel —</option>
        {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
      </select>
    </label>
    <label>Role granted on accept
      <select bind:value={draft.accepted_role_id}>
        <option value="">— none —</option>
        {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
      </select>
    </label>
    <label>Required role to apply
      <select bind:value={draft.required_role_id}>
        <option value="">— anyone —</option>
        {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
      </select>
    </label>
    <label>Cooldown between applications (hours)
      <input type="number" min="0" max="2160" bind:value={draft.cooldown_hours} />
    </label>
    <label>Minimum account age (days)
      <input type="number" min="0" max="3650" bind:value={draft.account_age_days} />
    </label>
    <label class="row"><input type="checkbox" bind:checked={draft.enabled} /> Form is enabled</label>
  </div>

  <div class="questions">
    <div class="qhead">
      <strong>Questions</strong>
      <button class="btn" on:click={addQuestion}><i class="fa-solid fa-plus"></i> Add question</button>
    </div>
    {#each draft.questions as q, i (i)}
      <div class="qcard">
        <div class="qrow">
          <input class="qid" bind:value={q.id} placeholder="id" />
          <input class="qlabel" bind:value={q.label} placeholder="Question text" />
          <select bind:value={q.type}>
            <option value="short">Short answer</option>
            <option value="paragraph">Paragraph</option>
            <option value="choice">Choice</option>
            <option value="scale">Scale</option>
          </select>
          <label class="req"><input type="checkbox" bind:checked={q.required} /> Required</label>
          <div class="qmove">
            <button class="btn ghost" aria-label="Move up" title="Move up" on:click={() => moveQuestion(i, -1)} disabled={i === 0}><i class="fa-solid fa-arrow-up"></i></button>
            <button class="btn ghost" aria-label="Move down" title="Move down" on:click={() => moveQuestion(i, 1)} disabled={i === draft.questions.length - 1}><i class="fa-solid fa-arrow-down"></i></button>
            <button class="btn danger" aria-label="Remove question" title="Remove question" on:click={() => removeQuestion(i)} disabled={draft.questions.length === 1}><i class="fa-solid fa-xmark"></i></button>
          </div>
        </div>
        {#if q.type === 'short' || q.type === 'paragraph'}
          <input bind:value={q.placeholder} placeholder="Placeholder (optional)" />
        {:else if q.type === 'choice'}
          <input value={(q.choices || []).join(', ')}
                 on:input={(e) => q.choices = e.currentTarget.value.split(',').map(s => s.trim()).filter(Boolean)}
                 placeholder="Options, comma-separated (e.g. Yes, No, Maybe)" />
        {:else if q.type === 'scale'}
          <div class="scale">
            <label>Min<input type="number" bind:value={q.min} min="0" max="10" /></label>
            <label>Max<input type="number" bind:value={q.max} min="1" max="10" /></label>
          </div>
        {/if}
      </div>
    {/each}
  </div>

  <div class="dmblock">
    <div class="dmtitle"><strong>DM templates</strong> <small>Variables: <code>{'{user}'}</code> <code>{'{form}'}</code> <code>{'{guild}'}</code> <code>{'{reason}'}</code></small></div>
    <label>On accept<textarea rows="2" bind:value={draft.accept_dm_template}></textarea></label>
    <label>On reject<textarea rows="3" bind:value={draft.reject_dm_template}></textarea></label>
  </div>

  <div class="row-actions">
    <button class="btn primary" on:click={save}>{editingId ? 'Save changes' : 'Create form'}</button>
    {#if editingId}<button class="btn" on:click={cancel}>Cancel</button>{/if}
  </div>
</Panel>

<Panel title="Forms">
  {#if loading}<p>Loading…</p>{:else if forms.length === 0}<p class="empty">No forms yet.</p>{:else}
    <table>
      <thead><tr><th></th><th>Name</th><th>Questions</th><th>Channel</th><th>Status</th><th></th></tr></thead>
      <tbody>
        {#each forms as f}
          <tr class:disabled={!f.enabled}>
            <td><span class="emoji" style={`background:${f.color}`}>{f.emoji || '📝'}</span></td>
            <td>
              <strong>{f.name}</strong>
              <div class="muted">{f.description || ''}</div>
            </td>
            <td>{f.questions?.length || 0}</td>
            <td>#{channels.find(c => c.id === f.submission_channel_id)?.name || '—'}</td>
            <td>
              {#if f.enabled}<span class="tag enabled">enabled</span>
              {:else}<span class="tag disabled-tag">off</span>{/if}
            </td>
            <td class="actions">
              <button class="btn" on:click={() => startEdit(f)}>Edit</button>
              <button class="btn danger" on:click={() => remove(f)}>Delete</button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .grid { display:grid; grid-template-columns:repeat(3, 1fr); gap:12px; }
  .grid .full { grid-column: 1 / -1; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid label.row { flex-direction:row; align-items:center; padding:8px 12px; background:#f8fafc; border-radius:8px; }
  .grid input, .grid select, .grid textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; font-family:inherit; }
  .grid input[type="color"] { padding:0; height:38px; }

  .questions { margin-top:18px; }
  .qhead { display:flex; align-items:center; justify-content:space-between; margin-bottom:10px; }
  .qcard { padding:12px 14px; border:1px solid #e5e7eb; border-radius:10px; margin-bottom:10px; display:flex; flex-direction:column; gap:8px; background:#fafbfc; }
  .qrow { display:grid; grid-template-columns:120px 1fr 140px auto auto; gap:8px; align-items:center; }
  .qid { font-family:Menlo, monospace; font-size:12px; }
  .qlabel { font-weight:600; }
  .qrow input, .qrow select { padding:6px 8px; border:1px solid #e5e7eb; border-radius:6px; font-size:13px; }
  .req { font-size:12px; color:#475569; display:flex; align-items:center; gap:4px; }
  .qmove { display:flex; gap:4px; }
  .scale { display:flex; gap:10px; }
  .scale label { display:flex; flex-direction:column; font-size:12px; gap:4px; }

  .dmblock { margin-top:18px; padding-top:18px; border-top:1px solid #f1f5f9; }
  .dmtitle { display:flex; align-items:center; justify-content:space-between; margin-bottom:10px; color:#475569; font-size:13px; }
  .dmblock label { display:flex; flex-direction:column; gap:6px; font-size:13px; color:#475569; margin-bottom:10px; }
  .dmblock textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-family:inherit; font-size:13px; }
  small code { background:#f1f5f9; padding:1px 5px; border-radius:4px; font-size:11px; }

  .row-actions { display:flex; gap:8px; margin-top:14px; }
  .btn { padding:8px 14px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; display:inline-flex; align-items:center; gap:6px; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn.ghost { background:transparent; color:#64748b; padding:6px 8px; }
  .btn:hover:not(:disabled) { filter:brightness(0.95); }
  .btn:disabled { opacity:0.4; cursor:not-allowed; }

  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; vertical-align:middle; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  tr.disabled td { opacity:0.55; }
  .emoji { display:inline-grid; place-items:center; width:34px; height:34px; border-radius:8px; color:#fff; font-size:16px; }
  .muted { color:#94a3b8; font-size:11px; margin-top:2px; }
  .tag { padding:2px 8px; border-radius:6px; font-size:11px; font-weight:600; text-transform:uppercase; letter-spacing:0.5px; }
  .tag.enabled { background:#dcfce7; color:#166534; }
  .tag.disabled-tag { background:#fee2e2; color:#b91c1c; }
  .actions { display:flex; gap:6px; }
  .empty { color:#94a3b8; }
</style>
