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
	// ExpiresAt, when set, is used as exp instead of Expires, so a reissued token keeps an exact deadline
	ExpiresAt time.Time
	// IssuedAt, when set, is used as iat instead of now, so a reissued token keeps its login time
	IssuedAt time.Time
}

// jwtDataClaim is the claim that carries the caller's data
const jwtDataClaim = "data"

// jwtIssuedAtSkew is how far in the future an iat may be, to absorb clock skew between servers
const jwtIssuedAtSkew = time.Minute

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

	now := time.Now()
	issuedAt := now
	if !option.IssuedAt.IsZero() {
		issuedAt = option.IssuedAt
	}
	expiresAt := now.Add(jwtExpires)
	if option.Expires != 0 {
		expiresAt = now.Add(option.Expires)
	}
	if !option.ExpiresAt.IsZero() {
		expiresAt = option.ExpiresAt
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jwtDataClaim: data,
		"exp":        jwt.NewNumericDate(expiresAt),
		"iat":        jwt.NewNumericDate(issuedAt),
	})
	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte(jwtSecret))
}

// JwtParse returns parsed token as map
func JwtParse(encrypted string) (interface{}, error) {
	claims, err := jwtParseClaims(encrypted)
	if err != nil {
		return nil, err
	}

	return claims[jwtDataClaim], nil
}

// JwtParseWithExpiry returns the parsed token data with its exp, so a reissued token can keep the same lifetime
func JwtParseWithExpiry(encrypted string) (interface{}, time.Time, error) {
	claims, err := jwtParseClaims(encrypted)
	if err != nil {
		return nil, time.Time{}, err
	}

	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return nil, time.Time{}, xerrors.New("token has no exp")
	}

	return claims[jwtDataClaim], exp.Time, nil
}

// JwtParseWithTimes returns the parsed HS256 token data with its iat and exp, refusing a token that lacks either.
// Both are whole seconds, as JwtToken signs them.
func JwtParseWithTimes(encrypted string) (interface{}, time.Time, time.Time, error) {
	claims, err := jwtParseClaims(encrypted, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, time.Time{}, time.Time{}, err
	}

	iat, err := claims.GetIssuedAt()
	if err != nil || iat == nil {
		return nil, time.Time{}, time.Time{}, xerrors.New("token has no iat")
	}
	if iat.After(time.Now().Add(jwtIssuedAtSkew)) {
		return nil, time.Time{}, time.Time{}, xerrors.New("token is issued in the future")
	}
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return nil, time.Time{}, time.Time{}, xerrors.New("token has no exp")
	}

	return claims[jwtDataClaim], iat.Time, exp.Time, nil
}

func jwtParseClaims(encrypted string, options ...jwt.ParserOption) (jwt.MapClaims, error) {
	if jwtSecret == jwtSecretDefault {
		return nil, xerrors.New("must be changed default secret")
	}

	token, err := jwt.Parse(encrypted, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	}, options...)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, xerrors.New("invalid token")
	}

	return claims, nil
}
