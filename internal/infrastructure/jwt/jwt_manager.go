package jwt

import (
	"fmt"
	"jwt-auth/internal/domain/services"
)

type JWTManagerImpl struct {
	tokens map[string]map[string]interface{}
}

func NewJWTManager() services.JWTManager {
	return &JWTManagerImpl{
		tokens: make(map[string]map[string]interface{}),
	}
}

func (j *JWTManagerImpl) GenerateToken(userID string, claims map[string]interface{}) (string, error) {
	email, _ := claims["email"].(string)
	username, _ := claims["username"].(string)
	if email == "" {
		email = "unknown@example.com"
	}
	token := "access_" + userID + "_" + email
	j.tokens[token] = map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"email":    email,
	}
	return token, nil
}

func (j *JWTManagerImpl) GenerateRefreshToken(userID string) (string, error) {
	token := "refresh_" + userID
	// Store minimal claims for refresh token
	j.tokens[token] = map[string]interface{}{
		"user_id":  userID,
		"username": "",
		"email":    "",
	}
	return token, nil
}

func (j *JWTManagerImpl) ValidateToken(token string) (map[string]interface{}, error) {
	if claims, ok := j.tokens[token]; ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
