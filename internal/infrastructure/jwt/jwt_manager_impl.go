package jwt

import (
	"fmt"
	"jwt-auth/internal/domain/services"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtManagerImpl struct {
	secretKey []byte
	accessExpiry time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager() services.JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "feh5tpb9aYtPxbCAxRKHZU967WyH3yjE"
	}
	
	accessExpiry := time.Hour
	if exp := os.Getenv("JWT_ACCESS_EXPIRY"); exp != "" {
		if d, err := time.ParseDuration(exp); err == nil {
			accessExpiry = d
		}
	}
	
	return &jwtManagerImpl{
		secretKey: []byte(secret),
		accessExpiry: accessExpiry,
		refreshExpiry: 24 * time.Hour,
	}
}

func (j *jwtManagerImpl) GenerateToken(userID string, claims map[string]interface{}) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"username": claims["username"],
		"email": claims["email"],
		"iat": now.Unix(),
		"exp": now.Add(j.accessExpiry).Unix(),
		"type": "access",
	})
	return token.SignedString(j.secretKey)
}

func (j *jwtManagerImpl) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iat": now.Unix(),
		"exp": now.Add(j.refreshExpiry).Unix(),
		"type": "refresh",
	})
	return token.SignedString(j.secretKey)
}

func (j *jwtManagerImpl) ValidateToken(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Convert claims to map[string]interface{}
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}