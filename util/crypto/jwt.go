package crypto

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret  = "eiiconGRP47;pike"  // #nosec
	jwtExpires = time.Hour * 24 * 90 // 3 months

	// SessionKey is used as session key
	SessionKey = "_go_utils_key"
	// MaxAge is session max age
	MaxAge = time.Hour * 24 * 75 // 2.5 months
)

// JwtToken returns jwt token with expires
func JwtToken(data interface{}) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(jwtExpires).Unix(),
		"iat":  time.Now().Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte(jwtSecret))
}

// JwtParse returns parsed token as map
func JwtParse(encrypted string) (interface{}, error) {
	token, err := jwt.Parse(encrypted, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims["data"], nil
}
