<script>
  import Navbar from '../components/Navbar.svelte';
  import Footer from '../components/Footer.svelte';
</script>

<Navbar />

<section class="page">
  <h1>Contact</h1>
  <p>
    Reach us in our
    <a href="https://discord.gg/YAgbpz7QPX">support Discord</a>
    or email niina.somppi@gmail.com.
  </p>
</section>

<Footer />

<style>
  .page {
    max-width: 900px;
    margin: 40px auto;
    padding: 40px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 12px;
    color: #111;
    line-height: 1.6;
  }

  a {
    color: #001aff;
  }
</style>
