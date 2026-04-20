<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import { formatDate, humanBytes } from '../../lib/ui.js';
  import Panel from '../../lib/Panel.svelte';

  let backups = [];
  let loading = true;
  let creating = false;
  let newLabel = '';

  async function load() {
    if (!$currentGuildId) return;
    loading = true;
    backups = await guildApi($currentGuildId).get('/backup');
    loading = false;
  }
  onMount(load);
  $: if ($currentGuildId) load();

  async function create() {
    creating = true;
    try {
      await guildApi($currentGuildId).post('/backup', { label: newLabel || 'Manual backup' });
      newLabel = '';
      await load();
    } finally { creating = false; }
  }
  async function restore(b) {
    if (!confirm(`Restore backup "${b.label}"? This overwrites current settings.`)) return;
    await guildApi($currentGuildId).post(`/backup/${b.id}/restore`, {});
    alert('Backup restored.');
  }
  async function remove(b) {
    if (!confirm(`Delete backup "${b.label}"?`)) return;
    await guildApi($currentGuildId).del('/backup/' + b.id);
    await load();
  }
</script>

<Panel title="Backups" subtitle="Snapshots of your server's bot settings.">
  <div slot="actions">
    <input placeholder="Label (optional)" bind:value={newLabel} />
    <button class="btn primary" on:click={create} disabled={creating}>Create backup</button>
  </div>
  {#if loading}<p>Loading…</p>{:else}
    <table>
      <thead><tr><th>Label</th><th>Size</th><th>Created</th><th>By</th><th></th></tr></thead>
      <tbody>
        {#each backups as b}
          <tr>
            <td>{b.label}</td>
            <td>{humanBytes(b.size_bytes)}</td>
            <td>{formatDate(b.created_at)}</td>
            <td>{b.created_by}</td>
            <td class="actions">
              <button class="btn" on:click={() => restore(b)}>Restore</button>
              <button class="btn danger" on:click={() => remove(b)}>Delete</button>
            </td>
          </tr>
        {:else}
          <tr><td colspan="5" class="empty">No backups yet.</td></tr>
        {/each}
      </tbody>
    </table>
  {/if}
</Panel>

<style>
  table { width:100%; border-collapse:collapse; font-size:14px; }
  th, td { text-align:left; padding:10px 12px; border-bottom:1px solid #f1f5f9; }
  th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .empty { color:#94a3b8; text-align:center; padding:20px; }
  .actions { display:flex; gap:6px; }
  input { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:13px; }
  .btn { padding:6px 12px; border:none; border-radius:8px; font-size:13px; font-weight:600; cursor:pointer; background:#e2e8f0; color:#0f172a; }
  .btn.primary { background:#5865f2; color:#fff; }
  .btn.danger { background:#fee2e2; color:#b91c1c; }
  .btn:hover { filter:brightness(0.95); }
</style>
