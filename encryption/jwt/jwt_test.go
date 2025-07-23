package jwt

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJw(t *testing.T) {
	type User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		jwt.RegisteredClaims
	}

	secret := "jwt.token"

	expected := User{
		ID:    int(rand.Int31()),
		Email: "jeo@deo.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add((24 * 365) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := JwtEncode(secret, expected)

	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatalf("Token must not be empty")
	}

	actual, err := JwtDecode(secret, token, User{})

	if err != nil {
		t.Fatal(err)
	}

	if actual.ID != expected.ID {
		t.Fatalf("Expected ID to be (%d) but got (%d)", actual.ID, expected.ID)
	}

	if actual.Email != expected.Email {
		t.Fatalf("Expected ID to be (%s) but got (%s)", actual.Email, expected.Email)
	}

	if actual.RegisteredClaims.ExpiresAt.UnixMilli() != expected.RegisteredClaims.ExpiresAt.UnixMilli() {
		t.Fatalf(
			"Expected expires_at to be (%d) but got (%d)",
			actual.RegisteredClaims.ExpiresAt.UnixMilli(),
			expected.RegisteredClaims.ExpiresAt.UnixMilli(),
		)
	}

	if actual.RegisteredClaims.IssuedAt.UnixMilli() != expected.RegisteredClaims.IssuedAt.UnixMilli() {
		t.Fatalf(
			"Expected issued_at to be (%d) but got (%d)",
			actual.RegisteredClaims.IssuedAt.UnixMilli(),
			expected.RegisteredClaims.IssuedAt.UnixMilli(),
		)
	}

	if actual.RegisteredClaims.NotBefore.UnixMilli() != expected.RegisteredClaims.NotBefore.UnixMilli() {
		t.Fatalf(
			"Expected not_before to be (%d) but got (%d)",
			actual.RegisteredClaims.IssuedAt.UnixMilli(),
			expected.RegisteredClaims.IssuedAt.UnixMilli(),
		)
	}
}
