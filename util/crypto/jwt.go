package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/xerrors"
)

var (
	// ensure to be changed
	jwtSecretDefault = "****************"
	// jwtSecret must be changed
	jwtSecret = jwtSecretDefault
	// jwtExpires is used as session key
	jwtExpires = time.Hour * 24 * 90 // 3 months
	// Domain is used s cookie domain
	Domain = "localhost"
	// SessionKey is used as session key
	SessionKey = "_go_core_key"
	// MaxAge is session max age
	MaxAge = time.Hour * 24 * 75 // 2.5 months
)

// JwtTokenOption used as a generate token option
type JwtTokenOption struct {
	Expires time.Duration
}

// SetExpires set value
func SetExpires(expires time.Duration) {
	jwtExpires = expires
}

// SetSecret set value
func SetSecret(secret string) {
	jwtSecret = secret
}

// SetDomain set value
func SetDomain(key string) {
	Domain = key
}

// SetSessionKey set value
func SetSessionKey(key string) {
	SessionKey = key
}

// SetMaxAge set value
func SetMaxAge(age time.Duration) {
	MaxAge = age
}

// JwtToken returns jwt token with expires
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
		"exp":  jwt.NewNumericDate(time.Now().Add(expires)),
		"iat":  jwt.NewNumericDate(time.Now()),
	})
	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte(jwtSecret))
}

// JwtParse returns parsed token as map
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
