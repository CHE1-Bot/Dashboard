# CHE1 Dashboard — Development & Production

## Switching modes

The whole stack runs in one of two modes, picked by `APP_ENV`:

| Mode          | When to use                                  | Loads             |
|---------------|----------------------------------------------|-------------------|
| `development` | Local work. Default for `npm run dev`.       | `api/.env.development` then `api/.env` |
| `production`  | Real deploy. Default for `npm start`.        | `api/.env.production` then `api/.env`  |

In `development`, `DEV_MODE=true` is the default and — when `DISCORD_CLIENT_ID`
is unset — `/api/auth/login` bypasses Discord and signs you in as a seeded demo
user. The dashboard shows a yellow **DEV MODE** badge in the top bar.

In `production`, `DEV_MODE=false`, real OAuth is required, and the badge turns
green (or hides if you don't want it).

### npm scripts

```bash
npm run dev        # Vite (5173) + Go API (8080), APP_ENV=development
npm run dev:web    # Vite only (use if you run `go run ./api` yourself)
npm run build      # vite build → dist/
npm start          # build SPA, embed it, build Go binary, run with APP_ENV=production
npm run gen-pages  # regenerate src/dashboard/data/pages.js from Go route defs
```

These npm scripts are thin shims around `dashctl`, a Go CLI that lives in
[api/cmd/dashctl/](api/cmd/dashctl/). It replaces what used to be Node-authored
launcher scripts. You can invoke it directly too:

```bash
cd api && go run ./cmd/dashctl dev
cd api && go run ./cmd/dashctl start
cd api && go run ./cmd/dashctl gen-pages
```

`dashctl` walks up from the cwd to find `package.json`, then runs all commands
relative to that project root.

### Per-mode env files

- [api/.env.development](api/.env.development) — committed, safe defaults.
  Edit when you want to point local dev at a real DB / Worker / OAuth app.
- `api/.env.production` — gitignored. Copy
  [api/.env.production.example](api/.env.production.example) to
  `api/.env.production` and fill in real secrets.
- `api/.env` is still loaded last as a generic override.

The frontend learns its mode from `GET /api/meta`
(`{ app_env, dev_mode, oauth_enabled, ... }`).

---

## Production deployment

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
