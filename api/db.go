package main

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgx connection pool. When Pool is nil the API falls back to the
// in-memory Store so the dashboard still boots for pure-dev work.
type DB struct {
	Pool *pgxpool.Pool
}

// openDB opens a pgx pool and runs migrations. Returns nil (not error) when
// DATABASE_URL is empty so the caller can decide to continue in memory mode.
func openDB(ctx context.Context, url string) (*DB, error) {
	if url == "" {
		return nil, nil
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	d := &DB{Pool: pool}
	if err := d.Migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return d, nil
}

// Migrate runs idempotent DDL. It covers both the shared schema (tickets,
// user_levels, mod_logs, applications, application_forms, giveaways) from
// the Bot repo and dashboard-owned extension tables (settings, panels, etc.)
// that don't exist in Bot/Worker.
func (d *DB) Migrate(ctx context.Context) error {
	stmts := []string{
		// ---- Shared schema (matches CHE1-Bot/Bot schema.sql exactly) ----
		`CREATE TABLE IF NOT EXISTS tickets (
			id             BIGSERIAL PRIMARY KEY,
			guild_id       TEXT NOT NULL,
			channel_id     TEXT NOT NULL,
			user_id        TEXT NOT NULL,
			subject        TEXT NOT NULL,
			status         TEXT NOT NULL,
			transcript_url TEXT,
			opened_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			closed_at      TIMESTAMPTZ
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tickets_guild_status ON tickets(guild_id, status)`,

		`CREATE TABLE IF NOT EXISTS user_levels (
			guild_id   TEXT NOT NULL,
			user_id    TEXT NOT NULL,
			xp         BIGINT NOT NULL DEFAULT 0,
			level      BIGINT NOT NULL DEFAULT 0,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			PRIMARY KEY (guild_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS mod_logs (
			id           BIGSERIAL PRIMARY KEY,
			guild_id     TEXT NOT NULL,
			moderator_id TEXT NOT NULL,
			target_id    TEXT NOT NULL,
			action       TEXT NOT NULL,
			reason       TEXT,
			created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_mod_logs_guild_created ON mod_logs(guild_id, created_at DESC)`,

		`CREATE TABLE IF NOT EXISTS applications (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			user_id    TEXT NOT NULL,
			role       TEXT NOT NULL,
			answers    JSONB NOT NULL,
			status     TEXT NOT NULL DEFAULT 'pending',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS application_forms (
			guild_id TEXT NOT NULL,
			role     TEXT NOT NULL,
			url      TEXT NOT NULL,
			PRIMARY KEY (guild_id, role)
		)`,

		`CREATE TABLE IF NOT EXISTS giveaways (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			message_id TEXT NOT NULL,
			prize      TEXT NOT NULL,
			winners    TEXT[],
			ends_at    TIMESTAMPTZ NOT NULL,
			status     TEXT NOT NULL
		)`,

		// ---- Dashboard-owned tables (prefixed dash_) ----
		`CREATE TABLE IF NOT EXISTS dash_sessions (
			sid          TEXT PRIMARY KEY,
			user_id      TEXT NOT NULL,
			access_token TEXT,
			expires_at   TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_sessions_user ON dash_sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_sessions_expires ON dash_sessions(expires_at)`,

		`CREATE TABLE IF NOT EXISTS dash_users (
			id             TEXT PRIMARY KEY,
			username       TEXT NOT NULL,
			discriminator  TEXT,
			avatar         TEXT,
			email          TEXT,
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS dash_guild_settings (
			guild_id          TEXT PRIMARY KEY,
			prefix            TEXT NOT NULL DEFAULT '!',
			language          TEXT NOT NULL DEFAULT 'en',
			timezone          TEXT NOT NULL DEFAULT 'UTC',
			welcome_channel   TEXT,
			welcome_message   TEXT,
			auto_role_ids     JSONB NOT NULL DEFAULT '[]'::JSONB,
			log_channel       TEXT,
			modules           JSONB NOT NULL DEFAULT '{}'::JSONB,
			feature_flags     JSONB NOT NULL DEFAULT '{}'::JSONB,
			premium           BOOLEAN NOT NULL DEFAULT FALSE,
			updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS dash_leveling_settings (
			guild_id         TEXT PRIMARY KEY,
			enabled          BOOLEAN NOT NULL DEFAULT TRUE,
			xp_per_message   INTEGER NOT NULL DEFAULT 15,
			cooldown_sec     INTEGER NOT NULL DEFAULT 60,
			announce_channel TEXT,
			announce_message TEXT NOT NULL DEFAULT '',
			no_xp_channels   JSONB NOT NULL DEFAULT '[]'::JSONB,
			no_xp_roles      JSONB NOT NULL DEFAULT '[]'::JSONB,
			multiplier       REAL NOT NULL DEFAULT 1.0
		)`,

		`CREATE TABLE IF NOT EXISTS dash_level_rewards (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			level      INTEGER NOT NULL,
			role_id    TEXT NOT NULL,
			role_name  TEXT NOT NULL,
			stackable  BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_level_rewards_guild ON dash_level_rewards(guild_id)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_settings (
			guild_id           TEXT PRIMARY KEY,
			category_id        TEXT,
			support_role_ids   JSONB NOT NULL DEFAULT '[]'::JSONB,
			transcripts_on     BOOLEAN NOT NULL DEFAULT TRUE,
			close_confirm      BOOLEAN NOT NULL DEFAULT TRUE,
			max_open_per_user  INTEGER NOT NULL DEFAULT 1,
			naming_pattern     TEXT NOT NULL DEFAULT 'ticket-{user}'
		)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_panels (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			title      TEXT NOT NULL,
			message    TEXT NOT NULL,
			color      TEXT NOT NULL DEFAULT '#3498db',
			buttons    JSONB NOT NULL DEFAULT '[]'::JSONB
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_ticket_panels_guild ON dash_ticket_panels(guild_id)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_forms (
			id        BIGSERIAL PRIMARY KEY,
			guild_id  TEXT NOT NULL,
			name      TEXT NOT NULL,
			fields    JSONB NOT NULL DEFAULT '[]'::JSONB
		)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_tags (
			id        BIGSERIAL PRIMARY KEY,
			guild_id  TEXT NOT NULL,
			name      TEXT NOT NULL,
			color     TEXT NOT NULL DEFAULT '#64748b'
		)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_embed (
			guild_id     TEXT PRIMARY KEY,
			title        TEXT NOT NULL DEFAULT 'Support',
			description  TEXT NOT NULL DEFAULT '',
			color        TEXT NOT NULL DEFAULT '#3498db',
			footer       TEXT NOT NULL DEFAULT '',
			thumbnail    TEXT NOT NULL DEFAULT ''
		)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_staff (
			guild_id  TEXT NOT NULL,
			user_id   TEXT NOT NULL,
			username  TEXT NOT NULL,
			roles     JSONB NOT NULL DEFAULT '[]'::JSONB,
			PRIMARY KEY (guild_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS dash_automod_rules (
			id        BIGSERIAL PRIMARY KEY,
			guild_id  TEXT NOT NULL,
			name      TEXT NOT NULL,
			trigger   TEXT NOT NULL,
			action    TEXT NOT NULL,
			enabled   BOOLEAN NOT NULL DEFAULT TRUE,
			config    JSONB NOT NULL DEFAULT '{}'::JSONB
		)`,

		`CREATE TABLE IF NOT EXISTS dash_reports (
			id             BIGSERIAL PRIMARY KEY,
			guild_id       TEXT NOT NULL,
			reporter_id    TEXT NOT NULL,
			reporter_name  TEXT NOT NULL,
			target_id      TEXT NOT NULL,
			target_name    TEXT NOT NULL,
			reason         TEXT NOT NULL,
			status         TEXT NOT NULL DEFAULT 'open',
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS dash_giveaway_meta (
			giveaway_id   BIGINT PRIMARY KEY,
			winner_count  INTEGER NOT NULL DEFAULT 1,
			entrants      INTEGER NOT NULL DEFAULT 0,
			hosted_by     TEXT,
			required_role TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS dash_giveaway_blacklist (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			user_id    TEXT NOT NULL,
			username   TEXT NOT NULL,
			reason     TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS dash_alerts (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			severity   TEXT NOT NULL DEFAULT 'info',
			title      TEXT NOT NULL,
			detail     TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			dismissed  BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_alerts_guild ON dash_alerts(guild_id, dismissed)`,

		`CREATE TABLE IF NOT EXISTS dash_history (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			actor      TEXT NOT NULL,
			event      TEXT NOT NULL,
			detail     TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_dash_history_guild_created ON dash_history(guild_id, created_at DESC)`,

		`CREATE TABLE IF NOT EXISTS dash_backups (
			id          BIGSERIAL PRIMARY KEY,
			guild_id    TEXT NOT NULL,
			label       TEXT NOT NULL,
			size_bytes  BIGINT NOT NULL DEFAULT 0,
			payload     JSONB,
			created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_by  TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS dash_permissions (
			id         BIGSERIAL PRIMARY KEY,
			guild_id   TEXT NOT NULL,
			role_id    TEXT NOT NULL,
			role_name  TEXT NOT NULL,
			modules    JSONB NOT NULL DEFAULT '[]'::JSONB
		)`,

		`CREATE TABLE IF NOT EXISTS dash_mod_log_meta (
			mod_log_id       BIGINT PRIMARY KEY,
			moderator_name   TEXT,
			target_name      TEXT,
			duration_sec     INTEGER NOT NULL DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS dash_application_meta (
			application_id BIGINT PRIMARY KEY,
			username       TEXT,
			role_id        TEXT,
			role_name      TEXT
		)`,

		`CREATE TABLE IF NOT EXISTS dash_ticket_meta (
			ticket_id   BIGINT PRIMARY KEY,
			username    TEXT,
			assigned_to TEXT,
			tags        JSONB NOT NULL DEFAULT '[]'::JSONB
		)`,
	}

	conn, err := d.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, stmt := range stmts {
		if _, err := conn.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

// ErrNotFound is returned by DB helpers when a row does not exist.
var ErrNotFound = errors.New("not found")

// QueryRow is a thin helper so callers can branch on pgx.ErrNoRows.
func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.Pool.QueryRow(ctx, sql, args...)
}
