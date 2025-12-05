package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func generateToken(userID string) (string, error) {
	// Define a sample secret key
	var jwtKey = []byte("my_secret_key")

	// Create a map to hold the claims
	claims := &jwt.StandardClaims{

		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "kolendar_app",
		Subject:   userID,
	}

	// Create a new token object, specifying the algorithm and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkToken(tokenString string) (*jwt.StandardClaims, error) {
	// Define a sample secret key
	var jwtKey = []byte("my_secret_key")

	claims := &jwt.StandardClaims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
