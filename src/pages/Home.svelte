<script>
  import Navbar from '../components/Navbar.svelte';
  import Footer from '../components/Footer.svelte';
</script>

<Navbar />

<div class="page-bg">

<section class="hero">
  <div class="hero-content">
    <h1>Che1 – The next generation of Discord bots</h1>
    <p>
      Che1 is a new and innovative Discord bot that offers moderation, tickets,
      giveaways and powerful utilities to enhance your server.
    </p>

    <div class="hero-buttons">
      <a
        href="https://discord.com/oauth2/authorize?client_id=1479549318825246843"
        target="_blank"
        rel="noopener"
      >
        <button class="btn invite">
          <i class="fa-brands fa-discord"></i>
          Invite Che1
        </button>
      </a>

      <a href="#/dashboard">
        <button class="btn dashboard">
          <i class="fa-solid fa-gear"></i>
          Dashboard
        </button>
      </a>
    </div>
  </div>
</section>

<section class="features">
  <h2 class="features-title">Why Che1?</h2>
  <div class="feature-container">
    <div class="feature">
      <i class="fa-solid fa-users-gear" style="color: rgb(0, 52, 145);"></i>
      <h3>Feature Rich</h3>
      <p>
        Che1 offers a wide range of features to enhance your server, including
        moderation, tickets, giveaways and powerful utilities.
      </p>
    </div>

    <div class="feature">
      <i class="fa-solid fa-arrows-down-to-people" style="color: rgb(0, 52, 145);"></i>
      <h3>Community</h3>
      <p>
        CHE1 is built with the community in mind. We listen to our users and
        implement features that they want to see.
      </p>
    </div>

    <div class="feature">
      <i class="fa-solid fa-money-bills" style="color: rgb(0, 52, 145);"></i>
      <h3>Free to Use</h3>
      <p>
        CHE1 is completely free to use and will always be free to use. We don't
        charge for any of our features, but we do have a premium subscription
        that offers additional features and supports the development of CHE1.
      </p>
    </div>
  </div>
</section>

<section class="features">
  <h2 class="features-title">Features</h2>
  <div class="feature-container">
    <div class="feature">
      <i class="fa-solid fa-shield-alt" style="color: #298344;"></i>
      <h3>Advanced Moderation</h3>
      <p>Che1 provides robust moderation tools to keep your server safe and well-maintained.</p>
      <a href="#/documentation" class="btn" id="Moderation-Button">Learn More</a>
    </div>
    <div class="feature">
      <i class="fa-solid fa-ticket" style="color: #c93e3e;"></i>
      <h3>Ticketing System</h3>
      <p>Effortlessly manage support requests with our intuitive ticketing system.</p>
      <a href="#/documentation" class="btn" id="Ticketing-Button">Learn More</a>
    </div>
    <div class="feature">
      <i class="fa-solid fa-gift" style="color: #001aff;"></i>
      <h3>Giveaways</h3>
      <p>Engage your community with exciting giveaways and contests.</p>
      <a href="#/documentation" class="btn" id="Giveaways-Button">Learn More</a>
    </div>
  </div>
</section>

<section class="get-started-section">
  <div class="get-started" id="Get-Started">
    <h1>Get Started with Che1</h1>
    <p>Setup process is easy and quick. Just follow our simple steps to get started!</p>
    <p>
      For any further questions, feel free to reach out to our
      <a href="https://discord.gg/YAgbpz7QPX">support team</a>.
    </p>
  </div>
  <a href="#/dashboard">
    <button class="btn invite" id="Get-Started-Button">
      Get started with Che1
      <i class="fa-solid fa-arrow-right" style="color: rgb(255, 255, 255);"></i>
    </button>
  </a>
</section>

</div>

<Footer />

<style>
  .page-bg {
    background-image: url('/images/Background2.jpg');
    background-size: cover;
    background-position: center;
    background-attachment: fixed;
    background-repeat: no-repeat;
    color: white;
    min-height: 100vh;
  }

  .hero {
    min-height: auto;
    padding: 60px 20px 40px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
  }

  .hero h1 {
    font-size: 50px;
    margin-bottom: 20px;
  }

  .hero p {
    max-width: 600px;
    margin: auto;
    font-size: 18px;
  }

  .hero-buttons {
    margin-top: 30px;
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 15px;
  }

  .btn {
    padding: 14px 30px;
    border-radius: 30px;
    border: none;
    font-size: 16px;
    display: inline-flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;
  }

  .invite {
    background: white;
    color: black;
  }

  .dashboard {
    background: transparent;
    border: 2px solid white;
    color: white;
  }

  .features {
    padding: 100px 20px;
    text-align: center;
  }

  .features-title {
    font-size: 35px;
  }

  .feature-container {
    display: flex;
    justify-content: center;
    gap: 40px;
    flex-wrap: wrap;
  }

  .feature {
    background: white;
    color: black;
    padding: 30px;
    border-radius: 15px;
    width: 250px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  }

  .feature i {
    font-size: 30px;
    margin-bottom: 15px;
  }

  #Moderation-Button {
    color: #298344;
    cursor: default;
  }
  #Moderation-Button:hover {
    transform: translateY(-5px);
    transition: 0.2s;
  }
  #Ticketing-Button {
    color: #c93e3e;
    cursor: default;
  }
  #Ticketing-Button:hover {
    transform: translateY(-5px);
    transition: 0.2s;
  }
  #Giveaways-Button {
    color: #001aff;
    cursor: default;
  }
  #Giveaways-Button:hover {
    transform: translateY(-5px);
    transition: 0.2s;
  }

  #Get-Started {
    background: white;
    color: black;
    padding: 40px 50px;
    border-radius: 15px;
    width: 100%;
    max-width: 700px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    text-align: center;
  }

  #Get-Started h1 {
    font-size: 22px;
    margin-top: 0;
  }

  #Get-Started a {
    color: #001aff;
  }

  .get-started-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 60px 20px;
  }

  #Get-Started-Button {
    border: #888 solid 2px;
    margin-top: 30px;
    background-color: #00042b;
    color: white;
  }
</style>
