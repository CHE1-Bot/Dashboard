<script>
  import Navbar from '../components/Navbar.svelte';
  import Footer from '../components/Footer.svelte';
</script>

<Navbar />

<div class="tos-container">
  <div class="text">
    <div class="TOS">
      <h1 id="TOS">Terms Of Service</h1>
    </div>
    <h3>CHE1 Terms Of Service</h3>
    <p>
      Before using any application of CHE1 bot read the Terms Of Service
      carefully. By using this service you agree to our Terms Of Service and
      guidelines. If you do not agree you must not use this service.
    </p>

    <h3>Privacy</h3>
    <p>
      We process and use your data under
      <a href="https://www.finlex.fi/fi/lainsaadanto/2025/1148#chp_1__sec_1">
        the Finnish law of data processing and sharing</a
      >. For any further question please refer to our
      <a href="#/privacy">Privacy Policy</a>.
    </p>

    <h3>Guidelines</h3>
    <p>1. Do not use this application for anything that can cause harm to this application, you or anyone.</p>
    <p>2. Do not use this application for anything that violates any law applicable in your country or Finland.</p>
    <p>3. Do not try to access data or areas that you are not allowed to view.</p>
    <p>4. Do not try to bypass any level of security that this application has.</p>
    <p>5. Do not use this service to process any sensitive data.</p>
    <p>6. Do not use this service or store any adult media.</p>
    <p>
      7. Do not use this service for anything that violates
      <a href="https://discord.com/terms">Discord's Terms Of Service</a>
      or
      <a href="https://discord.com/guidelines">Discord's Community Guidelines</a>. This includes:
    </p>
    <p>Exploiting or disclosing bugs for unintended advantages.</p>
    <p>
      Violating applicable laws or promoting illegal activities, such as
      hacking, cracking, or distributing counterfeit software.
    </p>
    <p>
      Abusing Discord products, such as engaging in fraudulent Nitro incentives
      or Boosting activities.
    </p>

    <h3>Purchases</h3>
    <p>1. We have the right to cancel your purchase or refund you at any time for any reason.</p>

    <h3>Termination</h3>
    <p>
      You can ask your data to be deleted in our
      <a href="https://discord.gg/C3PBvDA2G">Support Server</a>
      or by contacting us via our gmail.
    </p>
  </div>
</div>

<Footer />

<style>
  .tos-container {
    max-width: 900px;
    margin: 40px auto;
    padding: 40px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 12px;
    color: #111;
    line-height: 1.6;
  }

  #TOS {
    color: #001aff;
  }

  a {
    color: #001aff;
  }
</style>
