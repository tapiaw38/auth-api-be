package domain

import "time"

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

func (r RefreshToken) Usable(now time.Time) bool {
	return r.RevokedAt == nil && now.Before(r.ExpiresAt)
}
