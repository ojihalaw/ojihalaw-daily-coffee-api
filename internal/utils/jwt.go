package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	AccessSecret  []byte
	RefreshSecret []byte
	Issuer        string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type AccessClaims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"uid"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func (m *JWTMaker) NewAccessToken(uid, role string, now time.Time) (string, time.Time, error) {
	exp := now.Add(m.AccessTTL)
	claims := &AccessClaims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(m.AccessSecret)
	return s, exp, err
}

func (m *JWTMaker) NewRefreshToken(uid, jti string, now time.Time) (string, time.Time, error) {
	exp := now.Add(m.RefreshTTL)
	claims := &RefreshClaims{
		UserID: uid,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Subject:   uid,
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(m.RefreshSecret)
	return s, exp, err
}
