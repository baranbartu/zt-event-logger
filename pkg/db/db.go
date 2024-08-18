package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// DB defines the methods for database operations
type DB interface {
	// Insert inserts an event into the database
	Insert(event *Event) error
}

// Query is used when a granular query structure with network, device (member), or/and user ID are
// needed
type Query struct {
	NetworkID string
	MemberID  string
	UserID    string
}

// QueryOpt is a way to abstract AdditionalQuery construction and func signature for its usages
type QueryOpt func(aq *Query)

// WithNetworkID us used when query should contain NetworkID
func WithNetworkID(networkID string) QueryOpt {
	return func(q *Query) {
		q.NetworkID = networkID
	}
}

// WithMemberID us used when query should contain MemberID
func WithMemberID(memberID string) QueryOpt {
	return func(q *Query) {
		q.MemberID = memberID
	}
}

// WithUserID us used when query should contain UserID
func WithUserID(userID string) QueryOpt {
	return func(q *Query) {
		q.MemberID = userID
	}
}

// Event represents the DB table that stores the webhook events
type Event struct {
	ID            int    `db:"id" json:"id"`
	HookID        string `db:"hook_id" json:"hook_id"`
	OrgID         string `db:"org_id" json:"org_id"`
	HookType      string `db:"hook_type" json:"hook_type"`
	NetworkID     string `db:"network_id" json:"network_id"`
	MemberID      string `db:"member_id" json:"member_id,omitempty"`
	UserID        string `db:"user_id" json:"user_id,omitempty"`
	UserEmail     string `db:"user_email" json:"user_email,omitempty"`
	NetworkConfig JSONB  `json:"network_config,omitempty" gorm:"type:json"`
	OldConfig     JSONB  `json:"old_config,omitempty" gorm:"type:json"`
	NewConfig     JSONB  `json:"new_config,omitempty" gorm:"type:json"`
	Metadata      JSONB  `json:"metadata,omitempty" gorm:"type:json"`
	CreatedAt     string `db:"created_at" json:"created_at"`
}

// JSONB is a custom type to handle JSON data
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}
