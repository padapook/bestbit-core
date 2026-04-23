package auth

import (
	"os"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/padapook/bestbit-core/internal/account/model"

	"log"
)

func getJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatalln("jw secret key is empty")
	}
	return []byte(secret)
}

func getShareTokenSecret() []byte {
	secret := os.Getenv("SHARE_TOKEN_SECRET_KEY")
	if secret == "" {
		log.Fatalln("share secret key is empty")
	}
	return []byte(secret)
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	AccountID string `json:"account_id"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateTokens(user *model.User) (*TokenDetails, error) {
	accessExpirationTime := time.Now().Add(1 * time.Hour)
	accessClaims := &Claims{
		AccountID: user.AccountId,
		Username:  user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bestbit-core",
			Subject:   user.AccountId,
		},
	}
	log.Println("'jwtSecret")
	log.Println("'shareTokenSecret")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(getJwtSecret())
	if err != nil {
		return nil, err
	}

	// ให้ exp refresh tk 1 วัน
	refreshExpirationTime := time.Now().Add(1 * 24 * time.Hour)
	refreshClaims := &Claims{
		AccountID: user.AccountId,
		Username:  user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bestbit-core",
			Subject:   user.AccountId,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(getJwtSecret())
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GenerateShareToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		AccountID: user.AccountId,
		Username:  user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bestbit-core",
			Subject:   user.AccountId,
			Audience:  jwt.ClaimStrings{"share-session"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getShareTokenSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateShareToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getShareTokenSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		isShareToken := false
		for _, aud := range claims.Audience {
			if aud == "share-session" {
				isShareToken = true
				break
			}
		}

		if !isShareToken {
			return nil, errors.New("invalid token type, expected share token")
		}

		return claims, nil
	}

	return nil, errors.New("invalid share token")
}
