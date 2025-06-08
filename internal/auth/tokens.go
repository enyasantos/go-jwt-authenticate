package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = jwt.TimeFunc().Add(15 * time.Minute).Unix() // Token valid for 15 minutes

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = jwt.TimeFunc().Add(7 * 24 * time.Hour).Unix() // Token valid for 7 days

	secret := os.Getenv("JWT_SECRET_REFRESH")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString, secretKey string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorExpired)
	}

	clains, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, jwt.NewValidationError("invalid claims", jwt.ValidationErrorClaimsInvalid)
	}

	return token, clains, nil
}
