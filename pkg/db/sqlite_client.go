package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SQLiteClient implements the DB interface for SQLite
type SQLiteClient struct {
	db *gorm.DB
}

// NewSQLiteClient creates a new SQLite client
func NewSQLiteClient(dbFile string) (DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	client := &SQLiteClient{db: db}
	if err := client.initializeSchema(); err != nil {
		return nil, err
	}

	return client, nil
}

// initializeSchema sets up the database schema
func (c *SQLiteClient) initializeSchema() error {
	// Auto-migrate the schema
	if err := c.db.AutoMigrate(&Event{}); err != nil {
		return fmt.Errorf("failed to migrate database schema: %v", err)
	}
	return nil
}

// Insert inserts an event into the database
func (c *SQLiteClient) Insert(event *Event) error {
	if err := c.db.Create(event).Error; err != nil {
		return fmt.Errorf("could not insert event: %v", err)
	}
	return nil
}

// Search makes a `select` query to database and returns all the recorded events based on a
// given query criteria
func (c *SQLiteClient) Search(opts ...QueryOpt) ([]Event, error) {
	// Initialize the base query
	query := c.db.Model(&Event{})

	// Apply the query options
	criteria := &QueryCriteria{}
	for _, opt := range opts {
		opt(criteria)
	}

	// Add conditions to the query based on the criteria
	if criteria.NetworkID != "" {
		query = query.Where("network_id = ?", criteria.NetworkID)
	}
	if criteria.MemberID != "" {
		query = query.Where("member_id = ?", criteria.MemberID)
	}
	if criteria.UserID != "" {
		query = query.Where("user_id = ?", criteria.UserID)
	}

	// Execute the query and retrieve the events
	var events []Event
	if err := query.Find(&events).Error; err != nil {
		return nil, fmt.Errorf("could not query events: %v", err)
	}

	return events, nil
}
