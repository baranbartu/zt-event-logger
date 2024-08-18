package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"zt-event-logger/pkg/db"
	"zt-event-logger/pkg/events"

	"github.com/gin-gonic/gin"
)

// EventRequest represents the payload for a ZeroTier Central event hook.
// It includes details about the network, device, and user associated with the event.
// This struct is used to log the event information in the system.
type EventRequest struct {
	Network string `json:"network" binding:"required"` // Network is the identifier of the ZeroTier network.
	Device  string `json:"device" binding:"required"`  // Device is the name or identifier of the device involved in the event.
	UserID  int    `json:"userID" binding:"required"`  // UserID is the unique identifier of the user associated with the event.
}

type config struct {
	dbFileLocation string
	preSharedKey   string
}

func generateConfig() (*config, error) {
	c := &config{}

	dbFileLocation := os.Getenv("DB_FILE_LOCATION")
	if dbFileLocation == "" {
		return nil, errors.New("env var DB_FILE_LOCATION is mandatory that should be provided")
	}
	c.dbFileLocation = dbFileLocation

	preSharedKey := os.Getenv("PRE_SHARED_KEY")
	c.preSharedKey = preSharedKey

	return c, nil
}

func main() {
	config, err := generateConfig()
	if err != nil {
		panic(err)
	}

	dbClient, err := db.NewSQLiteClient(config.dbFileLocation)
	if err != nil {
		panic(err)
	}

	eventProcessor, err := events.NewProcessor(dbClient)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// POST endpoint to receive ZeroTier event hooks
	router.POST("/events/receive", func(c *gin.Context) {
		rawPayload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// fetches the signature from the header
		signature := c.GetHeader("X-ZTC-Signature")
		psk := config.preSharedKey

		hookBase, err := eventProcessor.Process(rawPayload, events.WithSignatureInfo(signature, psk))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with received data
		c.JSON(http.StatusOK, gin.H{
			"message":   fmt.Sprintf("Event received and logged successfully"),
			"hook_id":   hookBase.HookID,
			"org_id":    hookBase.OrgID,
			"hook_type": hookBase.HookType,
		})
	})

	// Start the server on port PORT - default 8080
	router.Run(":8080")
}
