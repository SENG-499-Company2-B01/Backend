package helper

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT_INFO struct {
    Email     	string 		`json:"email"`
    IsAdmin  	bool    	`json:"isAdmin"`
    ExpiredAt 	time.Time 	`json:"expiry"`
}

// VerifyJWT - decrypts JWT to make sure its valid and also checks whether it is past its expiry or not
func VerifyJWT(tokenString string) (bool, JWT_INFO, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
        _, ok := token.Method.(*jwt.SigningMethodHMAC)
        if !ok {
            return nil, fmt.Errorf("Invalid Token")
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    }

	// parse claims
    jwtToken, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return false, JWT_INFO{}, err
	}

	// verify claims
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		email, present := claims["email"]
		if !present {
			return false, JWT_INFO{}, fmt.Errorf("Claims not found") 
		}

		isAdmin, present := claims["isAdmin"]
		if !present {
			return false, JWT_INFO{}, fmt.Errorf("Claims not found") 
		}

		e, present := claims["expiry"]
		if !present {
			return false, JWT_INFO{}, fmt.Errorf("Claims not found") 
		}

		expiry, err := time.Parse(time.RFC3339, e.(string))
		if err != nil {
			fmt.Println(err)
		}

		if time.Now().Before(expiry) {
			return true, JWT_INFO{
				IsAdmin: isAdmin.(bool),
				Email: email.(string),
				ExpiredAt: expiry,
			}, nil
		}
	}

	// jwtToken.
	return false, JWT_INFO{}, nil
}

func CleanJWT(tokenString string) (string, error) {
	const bearerPrefix = "Bearer "

	if strings.HasPrefix(tokenString, bearerPrefix) {
		return strings.TrimPrefix(tokenString, bearerPrefix), nil
	}

	return "", fmt.Errorf("invalid token format: doesn't start with 'Bearer '")
}