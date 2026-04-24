# CHE1 Dashboard — Production Deployment

This repo ships a single-binary Svelte + Go dashboard that is designed to run
alongside [`CHE1-Bot/Bot`](https://github.com/CHE1-Bot/Bot) and
[`CHE1-Bot/Worker`](https://github.com/CHE1-Bot/Worker).

## Architecture at a glance

```
 Browser  ──HTTPS──►  Dashboard BFF  ──┬─ Postgres (shared schema) ─┬──►  Bot
                       (this repo)     │                            │
                                       └────► Worker ──WS── Bot ◄───┘
```

- **Postgres** is shared by all three services. The dashboard migrates its
  own extension tables (`dash_*`) as well as the shared schema (`tickets`,
  `user_levels`, `mod_logs`, `applications`, `application_forms`, `giveaways`)
  on boot.
- **Worker** receives actions from the dashboard over REST
  (`POST /api/v1/tasks` with `{event, guild_id, payload}`) and broadcasts
  `task.created` on its WebSocket hub.
- **Bot** subscribes to the Worker WebSocket and executes the Discord side
  effects (send message, deploy panels, etc.). It also calls the dashboard
  back at `GET /api/bot/guilds/:gid/applications/forms`.

## Environment

| Variable                 | Required? | Purpose                                                          |
|--------------------------|-----------|------------------------------------------------------------------|
| `DATABASE_URL`           | yes       | Shared Postgres. Schema migrates on boot.                        |
| `DISCORD_CLIENT_ID`      | yes       | OAuth2 app client id.                                            |
| `DISCORD_CLIENT_SECRET`  | yes       | OAuth2 app client secret.                                        |
| `DISCORD_REDIRECT_URI`   | yes       | Must match the callback configured in Discord Dev Portal.        |
| `DISCORD_BOT_TOKEN`      | yes       | Used to check which guilds the bot is in (Manage vs Invite).     |
| `WORKER_URL`             | no        | Worker REST base, e.g. `http://worker:8081`.                     |
| `WORKER_API_KEY`         | if WORKER | Bearer token, must equal Worker `INBOUND_API_KEY`.               |
| `DASHBOARD_API_KEY`      | yes       | Bearer token the Bot uses to call `/api/bot/*`.                  |
| `FRONTEND_URL`           | yes       | Where the browser loads the SPA; used for CORS + redirects.      |
| `ALLOWED_ORIGINS`        | no        | Comma-separated CORS allowlist (defaults to `FRONTEND_URL`).     |
| `SESSION_SECURE_COOKIES` | prod      | Set to `true` behind HTTPS so cookies are flagged `Secure`.      |
| `DEV_MODE`               | no        | Never enable in production. Bypasses Discord OAuth.              |

Copy `.env.example` to `.env` and fill in real values. `docker-compose.yml`
reads from this file.

## Deploy with Docker

```bash
cp .env.example .env
# edit .env
docker compose up -d --build        # dashboard + postgres
# or the full stack (also launches prebuilt worker + bot images):
docker compose --profile full up -d
```

The dashboard publishes on `http://localhost:8080`. Put a TLS-terminating
reverse proxy (Caddy, Nginx, Cloudflare) in front of it for production, and
set `SESSION_SECURE_COOKIES=true`.

## Deploy without Docker

```bash
# 1. Build the SPA
npm ci && npm run build

# 2. Copy the SPA next to the Go binary so it gets embedded
rm -rf api/dist && cp -R dist api/dist

# 3. Build the single binary
cd api && go build -trimpath -ldflags "-s -w" -o che1-dashboard ./...

# 4. Run it
./che1-dashboard
```

## Health checks

- `GET /healthz` — always `ok`, for liveness.
- `GET /readyz`  — pings Postgres; `503` while DB is unavailable.

## Discord OAuth setup

1. Create an app at <https://discord.com/developers/applications>.
2. Under **OAuth2 → Redirects**, add `https://your-host/api/auth/callback`.
3. Copy the Client ID / Secret into `DISCORD_CLIENT_ID` / `DISCORD_CLIENT_SECRET`.
4. Under **Bot**, generate a token and set `DISCORD_BOT_TOKEN`.
5. Grant your bot `applications.commands`, `bot`, `Manage Guild`, and any
   module-specific perms (Kick, Ban, Manage Messages, etc.).

## Shared secrets

`WORKER_API_KEY` and `DASHBOARD_API_KEY` must match on **all three services**:

| Secret              | Dashboard env          | Worker env         | Bot env            |
|---------------------|------------------------|--------------------|--------------------|
| `WORKER_API_KEY`    | `WORKER_API_KEY`       | `INBOUND_API_KEY`  | `WORKER_API_KEY`   |
| `DASHBOARD_API_KEY` | `DASHBOARD_API_KEY`    | `DASHBOARD_API_KEY`| `DASHBOARD_API_KEY`|

Generate with `openssl rand -hex 32`.
