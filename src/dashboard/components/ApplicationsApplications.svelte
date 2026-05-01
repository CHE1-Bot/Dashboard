<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate, relativeTime } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let apps = [];
  let forms = [];
  let loading = true;
  let statusFilter = 'pending';
  let formFilter = 'all';
  let search = '';
  let active = null; // currently expanded application

  let actioning = false;
  let reasonDraft = '';

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    const api = guildApi($currentGuildId);
    [apps, forms] = await Promise.all([
      api.get('/applications'),
      api.get('/applications/forms'),
    ]);
    loading = false;
    if (active) {
      active = apps.find(a => a.id === active.id) || null;
    }
  }

  async function decide(a, decision) {
    if (actioning) return;
    actioning = true;
    try {
      await guildApi($currentGuildId).post(`/applications/${a.id}/${decision}`, {
        reason: reasonDraft.trim(),
      });
      reasonDraft = '';
      await load();
    } catch (e) {
      alert(e.message || 'Failed');
    } finally { actioning = false; }
  }

  async function remove(a) {
    if (!confirm(`Delete this application from ${a.username}?`)) return;
    await guildApi($currentGuildId).del('/applications/' + a.id);
    if (active?.id === a.id) active = null;
    await load();
  }

  function findForm(formId) { return forms.find(f => f.id === formId); }
  function questionLabel(formId, qid) {
    return findForm(formId)?.questions?.find(q => q.id === qid)?.label || qid;
  }

  onMount(load);
  $: if ($currentGuildId) load();

  $: shown = apps.filter(a =>
    (statusFilter === 'all' || a.status === statusFilter) &&
    (formFilter === 'all' || String(a.form_id) === formFilter) &&
    (!search || (a.username + ' ' + (a.form_name || '')).toLowerCase().includes(search.toLowerCase()))
  );

  $: pending  = apps.filter(a => a.status === 'pending').length;
  $: accepted = apps.filter(a => a.status === 'accepted').length;
  $: rejected = apps.filter(a => a.status === 'rejected').length;
</script>

<Panel title="Review queue"
       subtitle="{pending} pending · {accepted} accepted · {rejected} rejected">
  <div slot="actions" class="filters">
    <input placeholder="Search applicant" bind:value={search} />
    <select bind:value={statusFilter}>
      <option value="pending">Pending</option>
      <option value="accepted">Accepted</option>
      <option value="rejected">Rejected</option>
      <option value="all">All</option>
    </select>
    <select bind:value={formFilter}>
      <option value="all">All forms</option>
      {#each forms as f}<option value={String(f.id)}>{f.emoji ? f.emoji + ' ' : ''}{f.name}</option>{/each}
    </select>
  </div>

  {#if loading}<p>Loading…</p>
  {:else if shown.length === 0}
    <div class="empty">
      <i class="fa-solid fa-inbox"></i>
      <p>Nothing in the queue right now.</p>
    </div>
  {:else}
    <div class="cards">
      {#each shown as a}
        <div class="card" class:open={active?.id === a.id}
             on:click={() => active = active?.id === a.id ? null : a}
             on:keydown={(e) => e.key === 'Enter' && (active = active?.id === a.id ? null : a)}
             role="button" tabindex="0">
          <div class="row top">
            <div class="who">
              <div class="avatar">{(a.username || '?').slice(0, 2).toUpperCase()}</div>
              <div>
                <div class="name">{a.username}</div>
                <div class="form-tag">
                  {#if findForm(a.form_id)?.emoji}{findForm(a.form_id).emoji}{/if}
                  {a.form_name || 'Application'}
                </div>
              </div>
            </div>
            <span class={'status ' + a.status}>{a.status}</span>
          </div>

          <div class="meta">
            <span><i class="fa-regular fa-clock"></i> {relativeTime(a.created_at)}</span>
            {#if a.reviewed_by_name}
              <span><i class="fa-solid fa-gavel"></i> by {a.reviewed_by_name}</span>
            {/if}
          </div>

          {#if active?.id === a.id}
            <div class="answers" on:click|stopPropagation>
              {#each Object.entries(a.answers || {}) as [k, v]}
                <div class="answer">
                  <div class="q">{questionLabel(a.form_id, k)}</div>
                  <div class="ans">{v}</div>
                </div>
              {/each}

              {#if a.review_note}
                <div class="answer note">
                  <div class="q">Reviewer note</div>
                  <div class="ans">{a.review_note}</div>
                </div>
              {/if}

              {#if a.status === 'pending'}
                <textarea
                  rows="2"
                  bind:value={reasonDraft}
                  placeholder="Optional message to send with the decision (DM template uses {'{reason}'})"
                ></textarea>
                <div class="actions">
                  <button class="btn accept" disabled={actioning} on:click={() => decide(a, 'accept')}>
                    <i class="fa-solid fa-check"></i> Accept
                  </button>
                  <button class="btn reject" disabled={actioning} on:click={() => decide(a, 'reject')}>
                    <i class="fa-solid fa-xmark"></i> Reject
                  </button>
                  <button class="btn ghost" on:click={() => remove(a)}>
                    <i class="fa-regular fa-trash-can"></i>
                  </button>
                </div>
              {:else}
                <div class="reviewed">
                  <i class="fa-regular fa-clock"></i> reviewed {formatDate(a.reviewed_at)}
                  <button class="btn ghost small" on:click|stopPropagation={() => remove(a)}>
                    <i class="fa-regular fa-trash-can"></i> Delete
                  </button>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</Panel>

<style>
  .filters { display:flex; gap:8px; align-items:center; }
  .filters input, .filters select { padding:6px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  .empty { text-align:center; padding:40px 20px; color:#94a3b8; }
  .empty i { font-size:32px; margin-bottom:10px; display:block; }
  .cards { display:flex; flex-direction:column; gap:12px; }
  .card { border:1px solid #e5e7eb; border-radius:12px; padding:16px 18px; cursor:pointer; transition:border-color 0.12s, box-shadow 0.12s; }
  .card:hover { border-color:#5865f2; }
  .card.open { border-color:#5865f2; box-shadow:0 0 0 3px rgba(88,101,242,0.12); cursor:default; }
  .row.top { display:flex; justify-content:space-between; align-items:center; gap:12px; }
  .who { display:flex; align-items:center; gap:12px; }
  .avatar { width:42px; height:42px; border-radius:50%; background:#5865f2; color:#fff; display:grid; place-items:center; font-weight:700; font-size:14px; }
  .name { font-weight:700; color:#0f172a; font-size:15px; }
  .form-tag { font-size:12px; color:#64748b; margin-top:2px; }
  .meta { font-size:11px; color:#94a3b8; display:flex; gap:14px; margin-top:8px; }
  .status { padding:3px 10px; border-radius:999px; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px; }
  .status.pending { background:#fef3c7; color:#92400e; }
  .status.accepted { background:#dcfce7; color:#166534; }
  .status.rejected { background:#fee2e2; color:#b91c1c; }
  .answers { margin-top:14px; padding-top:14px; border-top:1px solid #f1f5f9; display:flex; flex-direction:column; gap:12px; }
  .answer .q { font-size:11px; text-transform:uppercase; letter-spacing:0.5px; color:#94a3b8; font-weight:600; }
  .answer .ans { color:#1f2937; font-size:14px; white-space:pre-wrap; margin-top:3px; }
  .answer.note .q { color:#92400e; }
  .answer.note .ans { background:#fef3c7; padding:8px 12px; border-radius:8px; }
  textarea { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; font-family:inherit; resize:vertical; }
  .actions { display:flex; gap:8px; }
  .btn { display:inline-flex; align-items:center; gap:6px; padding:9px 16px; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; }
  .btn.accept { background:#10b981; color:#fff; }
  .btn.reject { background:#ef4444; color:#fff; }
  .btn.ghost { background:transparent; color:#64748b; padding:9px 12px; }
  .btn.ghost.small { padding:4px 10px; font-size:12px; }
  .btn:hover:not(:disabled) { filter:brightness(0.95); }
  .btn:disabled { opacity:0.5; cursor:not-allowed; }
  .reviewed { display:flex; gap:10px; align-items:center; color:#64748b; font-size:12px; }
</style>
