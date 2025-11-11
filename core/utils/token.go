package utils

import (
	"backend/web/model"
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TokenUtil struct {
	SecretKey string
}

func NewTokenUtil(secretKey string) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
	}
}

func (t TokenUtil) CreateAccessToken(ctx context.Context, auth *model.Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":    auth.UID,
		"expire": time.Now().Add(time.Minute * 1).UnixMilli(),
	})

	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (t TokenUtil) CreateRefreshToken(ctx context.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (t TokenUtil) ParseToken(ctx context.Context, jwtToken string) (*model.Auth, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	expire := claims["expire"].(float64)
	if int64(expire) < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}

	uid := claims["uid"].(string)
	auth := &model.Auth{
		UID: uid,
	}
	return auth, nil
}
