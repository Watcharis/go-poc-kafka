package util

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var key1 string

func JWTLoadKey() {
	// key1 = viper.GetString("secrets.jwt-key-access")
	key1 = "secretJwtPOCKafka"
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

func VerifyTokenExp(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractTokenExp(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key1), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenExp(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(ctx context.Context, r *http.Request) (*AccessDetails, error) {
	// validate token from request
	token, err := VerifyTokenExp(r)
	if err != nil {
		return nil, err
	}
	// map jwt from token
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		accessUuid, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("user_id not found")
		}

		return &AccessDetails{
			AccessUuid: accessUuid,
		}, nil
	}
	return nil, err
}
