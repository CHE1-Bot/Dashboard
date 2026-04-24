package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// dbGetSession returns the Session and its User for a given sid, or nil on miss.
// Implemented against Postgres when enabled; otherwise falls through to the
// in-memory map maintained by Store.
func dbGetSession(ctx context.Context, sid string) (*Session, *User) {
	if db == nil || db.Pool == nil {
		return nil, nil
	}
	var s Session
	s.ID = sid
	err := db.QueryRow(ctx, `
		SELECT user_id, COALESCE(access_token, ''), expires_at
		  FROM dash_sessions WHERE sid = $1`, sid,
	).Scan(&s.UserID, &s.AccessToken, &s.Expires)
	if err != nil {
		return nil, nil
	}
	if time.Now().After(s.Expires) {
		_, _ = db.Pool.Exec(ctx, `DELETE FROM dash_sessions WHERE sid = $1`, sid)
		return nil, nil
	}
	u, err := dbGetUser(ctx, s.UserID)
	if err != nil {
		return &s, nil
	}
	return &s, u
}

func dbPutSession(ctx context.Context, s *Session) error {
	if db == nil || db.Pool == nil {
		return nil
	}
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO dash_sessions (sid, user_id, access_token, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (sid) DO UPDATE SET
			access_token = EXCLUDED.access_token,
			expires_at   = EXCLUDED.expires_at`,
		s.ID, s.UserID, s.AccessToken, s.Expires,
	)
	return err
}

func dbDeleteSession(ctx context.Context, sid string) {
	if db == nil || db.Pool == nil {
		return
	}
	_, _ = db.Pool.Exec(ctx, `DELETE FROM dash_sessions WHERE sid = $1`, sid)
}

func dbGetUser(ctx context.Context, id string) (*User, error) {
	var u User
	err := db.QueryRow(ctx, `
		SELECT id, username, COALESCE(discriminator,''), COALESCE(avatar,''), COALESCE(email,'')
		  FROM dash_users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Username, &u.Discriminator, &u.Avatar, &u.Email)
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func dbUpsertUser(ctx context.Context, u *User) error {
	if db == nil || db.Pool == nil {
		return nil
	}
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO dash_users (id, username, discriminator, avatar, email, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (id) DO UPDATE SET
			username      = EXCLUDED.username,
			discriminator = EXCLUDED.discriminator,
			avatar        = EXCLUDED.avatar,
			email         = EXCLUDED.email,
			updated_at    = NOW()`,
		u.ID, u.Username, u.Discriminator, u.Avatar, u.Email,
	)
	return err
}
