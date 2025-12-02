package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Constructor
func NewPayload(userID int64, role string, duration time.Duration) *Payload {
	now := time.Now()
	return &Payload{
		ID:        uuid.New(),
		UserID:    userID,
		Role:      role,
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}
}

//
// Implement jwt.Claims interface
//

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiresAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetIssuer() (string, error) {
	return "", nil // optional
}

func (p *Payload) GetSubject() (string, error) {
	return "", nil // optional
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return []string{}, nil
}

// Your Valid method still works
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
