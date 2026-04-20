
<script>
import { onMount } from 'svelte';

onMount(() => {
  document.querySelectorAll('.accordion-header1, .accordion-header2, .accordion-header3, .accordion-header4, .accordion-header5')
    .forEach(header => {
      header.addEventListener('click', () => {
        header.parentElement.classList.toggle('active');
      });
    });

  document.querySelectorAll('.dropdown-button').forEach(button => {
    button.addEventListener('click', () => {
      button.parentElement.classList.toggle('active');
    });
  });
});
</script>

<div class="accordion1">
    <div class="accordion-item1">
        <div class="accordion-header1">
                        <h3>Leveling System</h3>
                    </div>
        <div class="accordion-content1" id="content1">
                            <div class="ticket-overview" id="check-button1">
                                <h3 class="accordion1-texts">Enable Leveling</h3>
                                <input type="checkbox" id="check1">
                                <label for="check1" class="button1"></label>
                            </div>
</div>
</div>
</div>
<div class="accordion2">
    <div class="accordion-item2">
        <div class="accordion-header2">
                        <h3>Advanced Settings</h3>
                    </div>
        <div class="accordion-content2" id="content2">
                            <div class="ticket-overview" id="check-button5">
                                <h3 class="accordion1-texts">Multiplier Channels</h3>
                                <input type="checkbox" id="check5">
                                <label for="check5" class="button5"></label>
                            </div>
</div>
</div>
</div>
