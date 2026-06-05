package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {
	
	// CreateHash returns an Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {

	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	
	// Create new claims and set fields
	claims := jwt.RegisteredClaims{
		
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
		}
	
	// Create new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Convert tokenSecret to a byte to be passed to SignedString
	secret := []byte(tokenSecret)

	// Sign token with the secret key
	return token.SignedString(secret)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	// Valide the JWT signature & extract the claims into *jwt.Token struct
	token, err := jwt.ParseWithClaims(tokenString,  &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {

		//return the key so the library can verify the signature
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	} 
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		userID, err := uuid.Parse(claims.Subject)
			return userID, err
	}
	return uuid.UUID{}, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// Get Authorization header
	authHeader := headers.Get("Authorization")

	// Error handling for no headers
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	// Trim the authHeader of bearer prefix
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}