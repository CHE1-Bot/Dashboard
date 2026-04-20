<script>
  export let title = '';
  export let subtitle = '';
</script>

<section class="panel">
  {#if title}
    <header>
      <h2>{title}</h2>
      {#if subtitle}<p class="sub">{subtitle}</p>{/if}
      <div class="actions"><slot name="actions" /></div>
    </header>
  {/if}
  <div class="body"><slot /></div>
</section>

<style>
  .panel { background:#fff; border:1px solid #e5e7eb; border-radius:12px; margin-bottom:20px; box-shadow:0 1px 2px rgba(0,0,0,0.03); }
  header { display:flex; align-items:center; gap:12px; padding:16px 20px; border-bottom:1px solid #f1f5f9; }
  header h2 { margin:0; font-size:16px; color:#111827; flex-shrink:0; }
  .sub { margin:0; color:#6b7280; font-size:13px; }
  .actions { margin-left:auto; display:flex; gap:8px; }
  .body { padding:20px; }
</style>
