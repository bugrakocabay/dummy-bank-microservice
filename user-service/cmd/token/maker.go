package token

import (
	"time"
)

type Maker interface {
	// CreateToken creates a new token for a specific user id and duration
	CreateToken(userID string, email string, duration time.Duration) (string, error)

	// VerifyToken checks if the given token is valid
	VerifyToken(token string) (*Payload, error)
}
