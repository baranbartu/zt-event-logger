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
