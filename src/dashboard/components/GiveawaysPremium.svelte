<script>
  import { onMount } from 'svelte';
  import { currentGuildId } from '../../lib/stores.js';
  import { guildApi } from '../../lib/api.js';
  import Panel from '../../lib/Panel.svelte';

  let data = null;

  async function load() {
    if (!$currentGuildId) return;
    data = await guildApi($currentGuildId).get('/giveaways/premium');
  }
  onMount(load);
  $: if ($currentGuildId) load();
</script>

{#if !data}<p>Loading…</p>{:else}
  <Panel title="Giveaway Premium"
         subtitle={data.premium ? 'Premium is active for this server.' : 'Upgrade to unlock weekly & monthly giveaways.'}>
    <div slot="actions">
      <a class="btn primary" href={data.upgrade_url}>
        {data.premium ? 'Manage subscription' : 'Upgrade to Premium'}
      </a>
    </div>

    <div class="plans">
      <div class="plan">
        <div class="plan-head">
          <div class="plan-title">Free</div>
          <div class="plan-price">$0<span>/server</span></div>
        </div>
        <ul>
          {#each data.free_includes as b}
            <li><i class="fa-solid fa-check"></i> {b}</li>
          {/each}
        </ul>
      </div>
      <div class="plan premium {data.premium ? 'active' : ''}">
        <div class="plan-head">
          <div class="plan-title"><i class="fa-solid fa-star"></i> Premium</div>
          <div class="plan-price">$5<span>/server/mo</span></div>
        </div>
        <ul>
          {#each data.benefits as b}
            <li><i class="fa-solid fa-check"></i> {b}</li>
          {/each}
        </ul>
        {#if data.premium}<div class="active-pill">Currently active</div>{/if}
      </div>
    </div>
  </Panel>

  <Panel title="Frequencies" subtitle="Pick a cadence when you create a giveaway.">
    <table class="tiers">
      <thead><tr><th>Frequency</th><th>Window</th><th>Plan</th><th></th></tr></thead>
      <tbody>
        {#each data.tiers as t}
          <tr class:locked={!t.available}>
            <td>
              <strong>{t.label}</strong>
              <div class="desc">{t.description}</div>
            </td>
            <td>{t.hours}h</td>
            <td>
              {#if t.premium_only}
                <span class="tag premium"><i class="fa-solid fa-star"></i> Premium</span>
              {:else}
                <span class="tag free">Free</span>
              {/if}
            </td>
            <td>
              {#if t.available}
                <span class="ok"><i class="fa-solid fa-circle-check"></i> Available</span>
              {:else}
                <span class="locked-tag"><i class="fa-solid fa-lock"></i> Upgrade required</span>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </Panel>
{/if}

<style>
  .plans { display:grid; grid-template-columns:repeat(2, 1fr); gap:16px; }
  .plan { border:1px solid #e5e7eb; border-radius:12px; padding:18px 20px; }
  .plan.premium { border-color:#5b21b6; background:linear-gradient(180deg, #faf5ff 0%, #fff 60%); position:relative; }
  .plan.premium.active { border-color:#10b981; }
  .plan-head { display:flex; justify-content:space-between; align-items:baseline; margin-bottom:14px; }
  .plan-title { font-size:18px; font-weight:700; color:#0f172a; }
  .plan.premium .plan-title { color:#5b21b6; }
  .plan-price { font-size:20px; font-weight:700; color:#0f172a; }
  .plan-price span { font-size:12px; color:#94a3b8; font-weight:500; }
  .plan ul { list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:8px; }
  .plan li { display:flex; align-items:center; gap:8px; color:#1f2937; font-size:14px; }
  .plan li i { color:#10b981; }
  .plan.premium li i { color:#5b21b6; }
  .active-pill { position:absolute; top:14px; right:14px; padding:4px 10px; background:#10b981; color:#fff; border-radius:999px; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px; }

  .tiers { width:100%; border-collapse:collapse; font-size:14px; margin-top:8px; }
  .tiers th, .tiers td { text-align:left; padding:12px; border-bottom:1px solid #f1f5f9; }
  .tiers th { color:#6b7280; font-weight:600; font-size:12px; text-transform:uppercase; }
  .tiers .desc { color:#64748b; font-size:12px; margin-top:2px; }
  .tiers tr.locked td { opacity:0.6; }
  .tag { padding:3px 10px; border-radius:999px; font-size:11px; font-weight:700; text-transform:uppercase; letter-spacing:0.5px; display:inline-flex; align-items:center; gap:4px; }
  .tag.free { background:#dcfce7; color:#166534; }
  .tag.premium { background:#ede9fe; color:#5b21b6; }
  .ok { color:#10b981; font-weight:600; font-size:13px; }
  .locked-tag { color:#9a3412; font-weight:600; font-size:13px; }

  .btn { display:inline-block; padding:10px 20px; background:#5865f2; color:#fff; border-radius:8px; font-weight:600; text-decoration:none; font-size:14px; }
  .btn:hover { filter:brightness(0.95); }
</style>
