<script>
  import Navbar from '../components/Navbar.svelte';
  import Footer from '../components/Footer.svelte';
</script>

<Navbar />

<section class="page">
  <h1>About Che1</h1>
  <p>
    Che1 is a community-driven Discord bot focused on moderation, tickets,
    giveaways, and leveling. Built with care to be reliable, fast, and free.
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
</style>
