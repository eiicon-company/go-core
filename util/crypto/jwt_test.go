package crypto

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJwtParseWithExpiry(t *testing.T) {
	SetSecret("0123456789abcdef")
	t.Cleanup(func() { SetSecret(jwtSecretDefault) })

	token, err := JwtToken(map[string]any{"pk": 1}, JwtTokenOption{Expires: 15 * time.Minute})
	if err != nil {
		t.Fatalf("JwtToken: %v", err)
	}

	data, expiresAt, err := JwtParseWithExpiry(token)
	if err != nil {
		t.Fatalf("JwtParseWithExpiry: %v", err)
	}
	if got := data.(map[string]any)["pk"]; got != float64(1) {
		t.Errorf("data pk = %v, want 1", got)
	}
	if remaining := time.Until(expiresAt); remaining <= 14*time.Minute || remaining > 15*time.Minute {
		t.Errorf("remaining lifetime = %s, want about 15m", remaining)
	}

	if _, _, err := JwtParseWithExpiry(token + "x"); err == nil {
		t.Error("a token with a broken signature must be rejected")
	}

	noExp, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"data": map[string]any{"pk": 1}}).SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatalf("sign a token without exp: %v", err)
	}
	if _, _, err := JwtParseWithExpiry(noExp); err == nil {
		t.Error("a token without exp must be rejected, as its lifetime is unknown")
	}
	if _, err := JwtParse(noExp); err != nil {
		t.Errorf("JwtParse must keep accepting a token without exp: %v", err)
	}
}

func TestJwtTokenWithExpiresAt(t *testing.T) {
	SetSecret("0123456789abcdef")
	t.Cleanup(func() { SetSecret(jwtSecretDefault) })

	deadline := time.Now().Add(7 * time.Minute).Truncate(time.Second)
	token, err := JwtToken(map[string]any{"pk": 1}, JwtTokenOption{Expires: time.Hour, ExpiresAt: deadline})
	if err != nil {
		t.Fatalf("JwtToken: %v", err)
	}

	_, expiresAt, err := JwtParseWithExpiry(token)
	if err != nil {
		t.Fatalf("JwtParseWithExpiry: %v", err)
	}
	if !expiresAt.Equal(deadline) {
		t.Errorf("exp = %s, want the given deadline %s", expiresAt, deadline)
	}
}

func TestJwtParseWithTimes(t *testing.T) {
	SetSecret("0123456789abcdef")
	t.Cleanup(func() { SetSecret(jwtSecretDefault) })

	issuedAt := time.Now().Add(-time.Hour).Truncate(time.Second)
	token, err := JwtToken(map[string]any{"pk": 1}, JwtTokenOption{Expires: 15 * time.Minute, IssuedAt: issuedAt})
	if err != nil {
		t.Fatalf("JwtToken: %v", err)
	}

	data, gotIssuedAt, expiresAt, err := JwtParseWithTimes(token)
	if err != nil {
		t.Fatalf("JwtParseWithTimes: %v", err)
	}
	if got := data.(map[string]any)["pk"]; got != float64(1) {
		t.Errorf("data pk = %v, want 1", got)
	}
	if !gotIssuedAt.Equal(issuedAt) {
		t.Errorf("iat = %s, want the given IssuedAt %s", gotIssuedAt, issuedAt)
	}
	if remaining := time.Until(expiresAt); remaining <= 14*time.Minute || remaining > 15*time.Minute {
		t.Errorf("remaining lifetime = %s, want about 15m counted from now", remaining)
	}

	reissued, err := JwtToken(data, JwtTokenOption{IssuedAt: gotIssuedAt, ExpiresAt: expiresAt})
	if err != nil {
		t.Fatalf("reissue: %v", err)
	}
	_, reissuedIat, reissuedExp, err := JwtParseWithTimes(reissued)
	if err != nil {
		t.Fatalf("parse the reissued token: %v", err)
	}
	if !reissuedIat.Equal(gotIssuedAt) || !reissuedExp.Equal(expiresAt) {
		t.Errorf("reissued iat, exp = %s, %s, want the original %s, %s", reissuedIat, reissuedExp, gotIssuedAt, expiresAt)
	}
}

func TestJwtParseWithTimesRefuses(t *testing.T) {
	SetSecret("0123456789abcdef")
	t.Cleanup(func() { SetSecret(jwtSecretDefault) })

	// sign signs a token; a nil iat or exp leaves that claim out.
	sign := func(method jwt.SigningMethod, iat, exp any) string {
		t.Helper()
		claims := jwt.MapClaims{jwtDataClaim: 1}
		if iat != nil {
			claims["iat"] = iat
		}
		if exp != nil {
			claims["exp"] = exp
		}
		signed, err := jwt.NewWithClaims(method, claims).SignedString([]byte(jwtSecret))
		if err != nil {
			t.Fatalf("sign: %v", err)
		}
		return signed
	}
	now := time.Now()
	exp := now.Add(time.Hour).Unix()
	refused := map[string]string{
		"no iat":            sign(jwt.SigningMethodHS256, nil, exp),
		"iat not a number":  sign(jwt.SigningMethodHS256, "yesterday", exp),
		"iat in the future": sign(jwt.SigningMethodHS256, now.Add(10*time.Minute).Unix(), exp),
		"no exp":            sign(jwt.SigningMethodHS256, now.Unix(), nil),
		"HS512":             sign(jwt.SigningMethodHS512, now.Unix(), exp),
	}
	for name, signed := range refused {
		if _, _, _, err := JwtParseWithTimes(signed); err == nil {
			t.Errorf("%s: must be rejected", name)
		}
	}
}
