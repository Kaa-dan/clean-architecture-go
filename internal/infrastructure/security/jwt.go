package security

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
	"github.com/kaa-dan/clean-architecture-go/pkg/errors"
)

type JWTManager struct {
	secretKey   string
	expiryHours int
}

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, expiryHours int) *JWTManager {
	return &JWTManager{
		secretKey:   secretKey,
		expiryHours: expiryHours,
	}
}

func (j *JWTManager) GenerateToken(user *entities.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID.Hex(),
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "user-api",
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
