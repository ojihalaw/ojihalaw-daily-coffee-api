package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
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

func NewJWTMaker(viper *viper.Viper) *JWTMaker {
	return &JWTMaker{
		AccessSecret:  []byte(viper.GetString("CMS_JWT_ACCESS_SECRET")),
		RefreshSecret: []byte(viper.GetString("CMS_JWT_REFRESH_SECRET")),
		Issuer:        viper.GetString("CMS_JWT_ISSUER"),
		AccessTTL:     viper.GetDuration("CMS_JWT_ACCESS_TTL"),
		RefreshTTL:    viper.GetDuration("CMS_JWT_REFRESH_TTL"),
	}
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

func (m *JWTMaker) VerifyAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func (m *JWTMaker) VerifyRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return m.RefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}
