<script>
  import { loginUrl } from '../lib/stores.js';
</script>

<div class="wrap">
  <div class="card">
    <img src="/images/Che1logo.png" alt="Che1" />
    <h1>Sign in to Che1 Dashboard</h1>
    <p>Manage your Discord server's tickets, moderation, giveaways and leveling — all in one place.</p>
    <a class="btn" href={loginUrl()}>
      <i class="fa-brands fa-discord"></i>
      Continue with Discord
    </a>
    <p class="fine">By continuing you agree to the <a href="#/tos">Terms</a> and <a href="#/privacy">Privacy Policy</a>.</p>
  </div>
</div>

<style>
  .wrap { min-height:100vh; display:grid; place-items:center; background:#0f172a; padding:20px; }
  .card { background:#fff; border-radius:16px; padding:40px; max-width:440px; text-align:center; }
  img { width:64px; height:64px; border-radius:50%; margin-bottom:16px; }
  h1 { margin:0 0 8px; font-size:22px; color:#111827; }
  p { color:#475569; font-size:14px; line-height:1.6; }
  .btn { display:inline-flex; align-items:center; gap:10px; background:#5865f2; color:#fff; padding:12px 24px; border-radius:8px; text-decoration:none; font-weight:600; margin:20px 0 12px; }
  .btn:hover { background:#4752c4; }
  .fine { font-size:12px; color:#94a3b8; }
  .fine a { color:#64748b; }
</style>
