<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let members = [];
  let picking = '';
  let action = 'warn';
  let reason = '';
  let duration = 60;
  let submitting = false;
  let result = '';

  async function load() {
    if (!$currentGuildId) return;
    members = await guildApi($currentGuildId).get('/members');
  }
  async function submit() {
    if (!picking) return;
    submitting = true;
    result = '';
    try {
      const m = members.find(x => x.id === picking);
      await guildApi($currentGuildId).post('/moderation/actions', {
        target_id: m.id, target_name: m.username,
        action, reason, duration_sec: (action === 'mute' ? duration * 60 : 0),
      });
      result = `Applied ${action} to ${m.username}.`;
      reason = ''; picking = '';
    } catch (e) {
      result = 'Error: ' + e.message;
    } finally { submitting = false; }
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

<Panel title="Manual moderation" subtitle="Apply a warn, mute, kick, or ban directly from the dashboard.">
  <div class="form">
    <label>Member
      <select bind:value={picking}>
        <option value="">Select…</option>
        {#each members as m}<option value={m.id}>{m.username}</option>{/each}
      </select>
    </label>
    <label>Action
      <select bind:value={action}>
        <option value="warn">Warn</option>
        <option value="mute">Mute</option>
        <option value="kick">Kick</option>
        <option value="ban">Ban</option>
      </select>
    </label>
    {#if action === 'mute'}
      <label>Duration (minutes)<input type="number" min="1" max="10080" bind:value={duration} /></label>
    {/if}
    <label class="full">Reason<input bind:value={reason} placeholder="Optional, shown in log" /></label>
  </div>
  <button class="btn primary" on:click={submit} disabled={!picking || submitting}>
    {submitting ? 'Applying…' : 'Apply'}
  </button>
  {#if result}<p class="result">{result}</p>{/if}
</Panel>

<style>
  .form { display:grid; grid-template-columns:repeat(3, 1fr); gap:14px; margin-bottom:14px; }
  .form .full { grid-column: 1 / -1; }
  .form label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .form input, .form select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .btn { padding:10px 20px; background:#5865f2; color:#fff; border:none; border-radius:8px; font-size:14px; font-weight:600; cursor:pointer; }
  .btn:disabled { opacity:0.5; cursor:not-allowed; }
  .btn:hover:not(:disabled) { filter:brightness(0.95); }
  .result { margin-top:14px; color:#059669; }
</style>
