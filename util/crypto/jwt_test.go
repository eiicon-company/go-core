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
