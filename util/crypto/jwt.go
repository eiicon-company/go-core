package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/xerrors"
)

var (
	jwtSecretDefault = "****************"
	jwtSecret = jwtSecretDefault
	jwtExpires = time.Hour * 24 * 90 // 3 months
	Domain = "localhost"
	SessionKey = "_go_core_key"
	MaxAge = time.Hour * 24 * 75 // 2.5 months
)

type JwtTokenOption struct {
	Expires time.Duration
}

func SetExpires(expires time.Duration) {
	jwtExpires = expires
}

func SetSecret(secret string) {
	jwtSecret = secret
}

func SetDomain(key string) {
	Domain = key
}

func SetSessionKey(key string) {
	SessionKey = key
}

func SetMaxAge(age time.Duration) {
	MaxAge = age
}

func JwtToken(data interface{}, option JwtTokenOption) (string, error) {
	if jwtSecret == jwtSecretDefault {
		return "", xerrors.New("must be changed default secret")
	}

	expires := jwtExpires
	if option.Expires != 0 {
		expires = option.Expires
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(expires).Unix(),
		"iat":  time.Now().Unix(),
	})
	return jwtToken.SignedString([]byte(jwtSecret))
}

func JwtParse(encrypted string) (interface{}, error) {
	if jwtSecret == jwtSecretDefault {
		return nil, xerrors.New("must be changed default secret")
	}

	token, err := jwt.Parse(encrypted, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, xerrors.New("invalid token")
	}

	return claims["data"], nil
}
