<script>
  export let label = '';
  export let value = '';
  export let icon = 'fa-chart-line';
  export let trend = null;
</script>
<div class="stat">
  <div class="ico"><i class={'fa-solid ' + icon}></i></div>
  <div>
    <div class="label">{label}</div>
    <div class="value">{value}</div>
    {#if trend !== null}<div class="trend" class:up={trend >= 0}>{trend >= 0 ? '▲' : '▼'} {Math.abs(trend)}%</div>{/if}
  </div>
</div>
<style>
  .stat { display:flex; gap:14px; align-items:center; background:#fff; border:1px solid #e5e7eb; border-radius:12px; padding:16px; }
  .ico { width:44px; height:44px; border-radius:10px; background:#eef2ff; color:#4f46e5; display:grid; place-items:center; font-size:18px; }
  .label { font-size:12px; color:#6b7280; text-transform:uppercase; letter-spacing:0.05em; }
  .value { font-size:24px; font-weight:700; color:#111827; }
  .trend { font-size:12px; color:#ef4444; }
  .trend.up { color:#10b981; }
</style>
