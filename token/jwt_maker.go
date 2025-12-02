package token

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("token is invalid")

// JWT Maker
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("secret key too short")
	}
	return &JWTMaker{secretKey}, nil
}

// Create a signed JWT token
func (maker *JWTMaker) CreateToken(
	userID int64,
	role string,
	duration time.Duration,
) (string, *Payload, error) {

	payload := NewPayload(userID, role, duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenString, payload, nil
}

// Verify token and return payload
func (maker *JWTMaker) VerifyToken(tokenString string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		// JWT library auto-calls p.Valid()
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
