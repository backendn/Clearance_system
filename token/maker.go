package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
    CreateToken(userID int64, role string, duration time.Duration) (string, *Payload, error)
    VerifyToken(token string) (*Payload, error)
}
