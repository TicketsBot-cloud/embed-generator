// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package postgres

import (
	"database/sql"
	"encoding/json"
	"time"
)

type MessageActionSet struct {
	ID        string
	MessageID string
	SetID     string
	Actions   json.RawMessage
}

type SavedMessage struct {
	ID          string
	CreatorID   string
	GuildID     sql.NullString
	UpdatedAt   time.Time
	Name        string
	Description sql.NullString
	Data        json.RawMessage
}

type Session struct {
	TokenHash   string
	UserID      string
	GuildIds    []string
	AccessToken string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type User struct {
	ID               string
	Name             string
	Discriminator    string
	Avatar           sql.NullString
	StripeCustomerID sql.NullString
	StripeEmail      sql.NullString
}

type UserSubscription struct {
	ID       string
	UserID   string
	Status   string
	PriceIds []string
	GuildIds []string
}
