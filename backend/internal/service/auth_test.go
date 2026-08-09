package service

import (
	"testing"
	"time"

	"github.com/entaku0818/aso-compass/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

// signWith は指定した鍵で有効なクレームのトークンを組み立てる。
func signWith(t *testing.T, secret string) string {
	t.Helper()

	claims := &Claims{
		UserID:  "00000000-0000-0000-0000-000000000000",
		Email:   "attacker@example.invalid",
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

// JWT_SECRET が設定されていれば、公開リポジトリに含まれる既定値で署名された
// トークンは受け付けてはならない。
func TestValidateTokenRejectsPublicDefaultSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-real-secret-set-in-production")

	svc := NewAuthService(nil)
	forged := signWith(t, "default-secret-change-in-production")

	if _, err := svc.ValidateToken(forged); err != model.ErrUnauthorized {
		t.Fatalf("既定シークレットで署名されたトークンが受理された (err=%v)", err)
	}
}

func TestValidateTokenAcceptsOwnToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "a-real-secret-set-in-production")

	svc := NewAuthService(nil)
	user := &model.User{ID: "user-1", Email: "user@example.com", IsAdmin: false}

	token, err := svc.GenerateTokenForUser(user)
	if err != nil {
		t.Fatalf("GenerateTokenForUser: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("自前で発行したトークンが拒否された: %v", err)
	}
	if claims.UserID != user.ID || claims.Email != user.Email {
		t.Fatalf("クレームが一致しない: %+v", claims)
	}
}

// 鍵を回すと既存トークンは失効する (macOSアプリは自動再アクティベートで復帰する)。
func TestValidateTokenRejectsTokenFromRotatedSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "old-secret")
	oldToken, err := NewAuthService(nil).GenerateTokenForUser(
		&model.User{ID: "user-1", Email: "user@example.com"})
	if err != nil {
		t.Fatalf("GenerateTokenForUser: %v", err)
	}

	t.Setenv("JWT_SECRET", "new-secret")
	if _, err := NewAuthService(nil).ValidateToken(oldToken); err != model.ErrUnauthorized {
		t.Fatalf("回転前のトークンが受理された (err=%v)", err)
	}
}
