<script>
  let quickOpen = false;
  let statusOpen = false;
  let enableAll = false;
  let autoSetup = false;
  let backupData = false;
  let resetSettings = false;
  let botOnline = true;
  let databaseConnected = true;
  let apiActive = true;
  let allSystemsGo = true;
</script>

<div class="accordion1">
  <div class="accordion-item1" class:active={quickOpen}>
    <button type="button" class="accordion-header1" on:click={() => quickOpen = !quickOpen} on:keydown={(event) => (event.key === 'Enter' || event.key === ' ') && (quickOpen = !quickOpen)}>
      <h3>Quick Actions</h3>
    </button>
    <div class="accordion-content1" id="content1">
      <div class="ticket-overview" id="check-button1">
        <h3 class="accordion1-texts">Enable All Features</h3>
        <input type="checkbox" id="check1" bind:checked={enableAll}>
        <label for="check1" class="button1"></label>
      </div>
      <div class="ticket-overview" id="check-button2">
        <h3 class="accordion1-texts">Auto Setup</h3>
        <input type="checkbox" id="check2" bind:checked={autoSetup}>
        <label for="check2" class="button2"></label>
      </div>
      <div class="ticket-overview" id="check-button3">
        <h3 class="accordion1-texts">Backup Data</h3>
        <input type="checkbox" id="check3" bind:checked={backupData}>
        <label for="check3" class="button3"></label>
      </div>
      <div class="ticket-overview" id="check-button4">
        <h3 class="accordion1-texts">Reset Settings</h3>
        <input type="checkbox" id="check4" bind:checked={resetSettings}>
        <label for="check4" class="button4"></label>
      </div>
</div>
</div>
</div>
<div class="accordion2">
  <div class="accordion-item2" class:active={statusOpen}>
    <button type="button" class="accordion-header2" on:click={() => statusOpen = !statusOpen} on:keydown={(event) => (event.key === 'Enter' || event.key === ' ') && (statusOpen = !statusOpen)}>
      <h3>System Status</h3>
    </button>
    <div class="accordion-content2" id="content2">
      <div class="ticket-overview" id="check-button5">
        <h3 class="accordion1-texts">Bot Online</h3>
        <input type="checkbox" id="check5" bind:checked={botOnline} disabled>
        <label for="check5" class="button5"></label>
      </div>
      <div class="ticket-overview" id="check-button6">
        <h3 class="accordion1-texts">Database Connected</h3>
        <input type="checkbox" id="check6" bind:checked={databaseConnected} disabled>
        <label for="check6" class="button6"></label>
      </div>
      <div class="ticket-overview" id="check-button7">
        <h3 class="accordion1-texts">API Active</h3>
        <input type="checkbox" id="check7" bind:checked={apiActive} disabled>
        <label for="check7" class="button7"></label>
      </div>
      <div class="ticket-overview" id="check-button8">
        <h3 class="accordion1-texts">All Systems Go</h3>
        <input type="checkbox" id="check8" bind:checked={allSystemsGo} disabled>
        <label for="check8" class="button8"></label>
      </div>
</div>
</div>
</div>
