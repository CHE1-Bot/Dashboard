<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let forms = [];
  let loading = true;
  let draft = { name: '', fields: [{ label: 'Subject', type: 'short', required: true, placeholder: '' }] };

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    forms = await guildApi($currentGuildId).get('/tickets/forms');
    loading = false;
  }
  async function create() {
    if (!draft.name) return;
    await guildApi($currentGuildId).post('/tickets/forms', draft);
    draft = { name: '', fields: [{ label: 'Subject', type: 'short', required: true, placeholder: '' }] };
    await load();
  }
  async function remove(f) {
    if (!confirm(`Delete form "${f.name}"?`)) return;
    await guildApi($currentGuildId).del('/tickets/forms/' + f.id);
    await load();
  }
  function addField() { draft.fields = [...draft.fields, { label: '', type: 'short', required: false, placeholder: '' }]; }
  function removeField(i) { draft.fields = draft.fields.filter((_, idx) => idx !== i); }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="New form" subtitle="Modal questions shown when a user opens a ticket.">
  <div class="form">
    <label>Name<input bind:value={draft.name} placeholder="e.g. Bug report" /></label>
    <div class="label">Fields <button class="btn" on:click={addField}>+ Add field</button></div>
    {#each draft.fields as f, i}
      <div class="field-row">
        <input bind:value={f.label} placeholder="Question" />
        <select bind:value={f.type}>
          <option value="short">Short</option>
          <option value="paragraph">Paragraph</option>
        </select>
        <input bind:value={f.placeholder} placeholder="Placeholder" />
        <label class="req"><input type="checkbox" bind:checked={f.required} />Required</label>
        <button class="btn danger" on:click={() => removeField(i)}>×</button>
      </div>
    {/each}
    <button class="btn primary" on:click={create}>Create form</button>
  </div>
</Panel>

<Panel title="Existing forms">
  {#if loading}<p>Loading…</p>{:else if forms.length === 0}<p class="empty">No forms.</p>{:else}
    <table>
      <thead><tr><th>Name</th><th>Fields</th><th></th></tr></thead>
      <tbody>
        {#each forms as f}
          <tr>
            <td>{f.name}</td>
            <td>{f.fields.length}</td>
            <td><button class="btn danger" on:click={() => remove(f)}>Delete</button></td>
          </tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  .form { display:flex; flex-direction:column; gap:10px; }
  .form > label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; max-width:400px; }
  .form input, .form select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .label { font-size:13px; color:#475569; display:flex; align-items:center; justify-content:space-between; }
  .field-row { display:grid; grid-template-columns:2fr 140px 1fr auto auto; gap:8px; align-items:center; }
  .field-row .req { font-size:12px; color:#475569; display:flex; align-items:center; gap:4px; }
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .btn { padding:6px 12px; background:#e2e8f0; color:#0f172a; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.primary { background:#5865f2; color:#fff; align-self:flex-start; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
  .empty { color:#94a3b8; }
</style>
