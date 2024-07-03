package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// This method generates a new JWT for the user
func GenerateToken(email string, userId string) (string, error) {

	// Retrieve the secret key from the ENV file
	secretKey := os.Getenv("SECRET_JWT_KEY")

	// First parameter: Signing Method
	// Second parameter: Additional user data for the JWT
	// Third parameter: Expires at (2 hours from NOW)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	// Return the produced token as Signed
	// Sign Key = Secret Key only known to us
	return token.SignedString([]byte(secretKey))
}

// This method verifies a JWT token
func VerifyToken(token string) (string, error) {

	// Retrieve the secret key from the ENV file
	secretKey := os.Getenv("SECRET_JWT_KEY")

	// Parse the received token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// Check the signing method type
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}

		// Return success
		return []byte(secretKey), nil
	})

	// Check if present error
	if err != nil {
		return "", errors.New("could not parse token")
	}

	// Check the token validity
	if !parsedToken.Valid {
		return "", errors.New("not valid token")
	}

	// Check the token data and expires at
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("not valid token")
	}

	// Access the token data
	fmt.Println("USER ID:", claims["userId"])
	userId := claims["userId"].(string)

	// Success and return the userId and nil for error
	return userId, nil
}
