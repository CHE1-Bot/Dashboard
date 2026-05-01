<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let channels = [];
  let roles = [];
  let premium = null;
  let creating = false;
  let created = null;
  let error = '';

  let draft = {
    prize: '',
    channel_id: '',
    winner_count: 1,
    frequency: 'daily',
    recurring: false,
    required_role_id: '',
  };

  async function load() {
    if (!$currentGuildId) return;
    const api = guildApi($currentGuildId);
    [channels, roles, premium] = await Promise.all([
      api.get('/channels'),
      api.get('/roles'),
      api.get('/giveaways/premium'),
    ]);
  }

  function tier(key) { return premium?.tiers?.find(t => t.key === key); }
  function locked(key) { return tier(key)?.premium_only && !premium?.premium; }

  async function submit() {
    error = '';
    if (!draft.prize || !draft.channel_id) { error = 'Prize and channel are required'; return; }
    if (locked(draft.frequency)) { error = `${tier(draft.frequency)?.label} giveaways require Premium`; return; }
    creating = true;
    try {
      created = await guildApi($currentGuildId).post('/giveaways', {
        prize: draft.prize,
        channel_id: draft.channel_id,
        winner_count: Number(draft.winner_count) || 1,
        frequency: draft.frequency,
        recurring: !!draft.recurring,
        required_role_id: draft.required_role_id,
      });
      draft = { ...draft, prize: '', required_role_id: '' };
    } catch (e) {
      error = e.message || 'Failed to create giveaway';
    } finally { creating = false; }
  }

  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !premium}<p>Loading…</p>{:else}
  <Panel title="Create giveaway"
         subtitle={premium.premium ? 'Premium server — all frequencies unlocked.' : 'Free plan — daily giveaways are included.'}>
    <div class="freq">
      {#each premium.tiers as t}
        <button
          class="tier"
          class:active={draft.frequency === t.key}
          class:locked={t.premium_only && !premium.premium}
          on:click={() => { if (!(t.premium_only && !premium.premium)) draft.frequency = t.key; }}
          type="button"
        >
          <div class="tier-head">
            <span class="tier-label">{t.label}</span>
            {#if t.premium_only}
              <span class="badge {premium.premium ? 'unlocked' : ''}">
                <i class="fa-solid {premium.premium ? 'fa-star' : 'fa-lock'}"></i>
                {premium.premium ? 'Premium' : 'Locked'}
              </span>
            {:else}
              <span class="badge free">Free</span>
            {/if}
          </div>
          <div class="tier-hours">{t.hours}h window</div>
          <div class="tier-desc">{t.description}</div>
        </button>
      {/each}
    </div>

    {#if locked(draft.frequency)}
      <div class="upgrade">
        <i class="fa-solid fa-star"></i>
        {tier(draft.frequency)?.label} giveaways are part of Premium.
        <a href="#/dashboard/giveaways/premium">Upgrade →</a>
      </div>
    {/if}

    <div class="grid">
      <label class="full">Prize<input bind:value={draft.prize} placeholder="e.g. Nitro Classic" /></label>
      <label>Channel
        <select bind:value={draft.channel_id}>
          <option value="">—</option>
          {#each channels.filter(c => c.type === 'text') as c}<option value={c.id}>#{c.name}</option>{/each}
        </select>
      </label>
      <label>Winners<input type="number" min="1" max="20" bind:value={draft.winner_count} /></label>
      <label>Required role
        <select bind:value={draft.required_role_id}>
          <option value="">— none —</option>
          {#each roles.filter(r => r.name !== '@everyone') as r}<option value={r.id}>{r.name}</option>{/each}
        </select>
      </label>
      <label class="full row">
        <input type="checkbox" bind:checked={draft.recurring} disabled={locked(draft.frequency)} />
        Recurring — automatically start a new {tier(draft.frequency)?.label.toLowerCase()} giveaway when this one ends
      </label>
    </div>

    {#if error}<div class="error">{error}</div>{/if}

    <button class="btn primary" on:click={submit}
            disabled={creating || !draft.prize || !draft.channel_id || locked(draft.frequency)}>
      {creating ? 'Creating…' : `Create ${tier(draft.frequency)?.label.toLowerCase()} giveaway`}
    </button>

    {#if created}
      <div class="success">
        Created <strong>{created.prize}</strong> ·
        {tier(created.frequency)?.label || 'Daily'} ·
        ends {new Date(created.ends_at).toLocaleString()}
      </div>
    {/if}
  </Panel>
{/if}

<style>
  .freq { display:grid; grid-template-columns:repeat(3, 1fr); gap:12px; margin-bottom:16px; }
  .tier { text-align:left; padding:14px 16px; border:2px solid #e5e7eb; border-radius:12px; background:#fff; cursor:pointer; transition:all 0.15s; }
  .tier:hover:not(.locked) { border-color:#5865f2; }
  .tier.active { border-color:#5865f2; box-shadow:0 0 0 3px rgba(88,101,242,0.15); }
  .tier.locked { opacity:0.6; cursor:not-allowed; background:#f8fafc; }
  .tier-head { display:flex; justify-content:space-between; align-items:center; }
  .tier-label { font-weight:700; font-size:15px; color:#0f172a; }
  .tier-hours { font-size:12px; color:#64748b; margin:6px 0 8px; }
  .tier-desc { font-size:12px; color:#475569; }
  .badge { font-size:10px; font-weight:700; padding:3px 8px; border-radius:999px; text-transform:uppercase; letter-spacing:0.5px; display:inline-flex; align-items:center; gap:4px; }
  .badge.free { background:#dcfce7; color:#166534; }
  .badge:not(.free):not(.unlocked) { background:#fef3c7; color:#92400e; }
  .badge.unlocked { background:#ede9fe; color:#5b21b6; }
  .upgrade { padding:10px 14px; background:#fff7ed; color:#9a3412; border:1px solid #fed7aa; border-radius:8px; margin-bottom:14px; font-size:13px; display:flex; align-items:center; gap:10px; }
  .upgrade a { color:#9a3412; font-weight:700; margin-left:auto; }
  .grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:14px; margin-bottom:16px; }
  .grid .full { grid-column: 1 / -1; }
  .grid label { display:flex; flex-direction:column; font-size:13px; color:#475569; gap:6px; }
  .grid label.row { flex-direction:row; align-items:center; gap:8px; padding:8px 12px; background:#f8fafc; border-radius:8px; }
  .grid input, .grid select { padding:8px 10px; border:1px solid #e5e7eb; border-radius:8px; font-size:14px; }
  .btn { padding:10px 20px; background:#5865f2; color:#fff; border:none; border-radius:8px; font-size:14px; font-weight:600; cursor:pointer; }
  .btn:disabled { opacity:0.5; cursor:not-allowed; }
  .btn:hover:not(:disabled) { filter:brightness(0.95); }
  .success { margin-top:14px; padding:10px 14px; background:#d1fae5; color:#065f46; border-radius:8px; font-size:14px; }
  .error { margin-bottom:14px; padding:10px 14px; background:#fee2e2; color:#b91c1c; border-radius:8px; font-size:14px; }
</style>
