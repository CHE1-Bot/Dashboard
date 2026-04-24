package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Discord OAuth2 permission bit: MANAGE_GUILD (0x20) is the minimum perm to
// appear on the dashboard's "you can manage this server" list.
const discordPermManageGuild int64 = 0x20

// partialGuild is the shape returned by GET /users/@me/guilds.
type partialGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions string `json:"permissions"`
	Approximate int    `json:"approximate_member_count,omitempty"`
}

// discordClient batches OAuth/Bot API calls with a timeout.
var discordClient = &http.Client{Timeout: 10 * time.Second}

// botGuildCache caches the set of guild IDs the bot is in so we don't hammer
// the Discord Bot API on every login. The set refreshes every 60 seconds.
type botCache struct {
	mu      sync.RWMutex
	guilds  map[string]struct{}
	fetched time.Time
}

var botPresence = &botCache{guilds: map[string]struct{}{}}

// botInGuild returns whether the CHE1 bot is a member of the given guild. It
// uses DISCORD_BOT_TOKEN to enumerate guilds. Without a bot token it always
// returns false.
func botInGuild(ctx context.Context, gid string) bool {
	if botToken == "" {
		return false
	}
	botPresence.mu.RLock()
	fresh := time.Since(botPresence.fetched) < 60*time.Second
	_, in := botPresence.guilds[gid]
	botPresence.mu.RUnlock()
	if fresh {
		return in
	}
	if err := refreshBotGuilds(ctx); err != nil {
		return in
	}
	botPresence.mu.RLock()
	defer botPresence.mu.RUnlock()
	_, in = botPresence.guilds[gid]
	return in
}

func refreshBotGuilds(ctx context.Context) error {
	// Enumerate up to 200 guilds; CHE1 deployments are smaller than that in
	// practice. Paginate via ?after= if we ever outgrow this.
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/v10/users/@me/guilds?limit=200", nil)
	req.Header.Set("Authorization", "Bot "+botToken)
	resp, err := discordClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord bot %d: %s", resp.StatusCode, string(b))
	}
	var gs []partialGuild
	if err := json.NewDecoder(resp.Body).Decode(&gs); err != nil {
		return err
	}
	set := make(map[string]struct{}, len(gs))
	for _, g := range gs {
		set[g.ID] = struct{}{}
	}
	botPresence.mu.Lock()
	botPresence.guilds = set
	botPresence.fetched = time.Now()
	botPresence.mu.Unlock()
	return nil
}

// exchangeCode swaps an OAuth code for an access token.
func exchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSec)
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://discord.com/api/oauth2/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := discordClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token exchange %d: %s", resp.StatusCode, string(b))
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", err
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("empty access_token")
	}
	return tok.AccessToken, nil
}

// fetchMe queries /users/@me with a user access token.
func fetchMe(ctx context.Context, token string) (*User, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/v10/users/@me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := discordClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("users/@me %d: %s", resp.StatusCode, string(b))
	}
	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

// fetchMyGuilds queries /users/@me/guilds and filters to guilds the user can
// actually administer (owner OR MANAGE_GUILD bit set).
func fetchMyGuilds(ctx context.Context, token string) ([]Guild, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://discord.com/api/v10/users/@me/guilds?with_counts=true", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := discordClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("users/@me/guilds %d: %s", resp.StatusCode, string(b))
	}
	var raw []partialGuild
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	out := make([]Guild, 0, len(raw))
	for _, g := range raw {
		perms, _ := strconv.ParseInt(g.Permissions, 10, 64)
		if !g.Owner && perms&discordPermManageGuild == 0 {
			continue
		}
		out = append(out, Guild{
			ID:          g.ID,
			Name:        g.Name,
			Icon:        g.Icon,
			Owner:       g.Owner,
			Permissions: perms,
			BotPresent:  botInGuild(ctx, g.ID),
			MemberCount: g.Approximate,
		})
	}
	return out, nil
}
