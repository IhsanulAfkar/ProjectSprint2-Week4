package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type StaffClaims struct {
	jwt.RegisteredClaims
	Id 	string `json:"id"`
	Role    string `json:"role"`
}
type ReturnToken struct {
	Id 	string `json:"id"`
	Role    string `json:"role"`
}
func SignJWT(id string, role string) string {
	// expiredIn := 28800 // 8 hours
	exp := time.Now().Add(time.Hour * 8)
	claims := StaffClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer: "Cat Socials",
		},
		Id: id,
		Role: role,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ""
	}
	return signedToken

}
func ParseToken(jwtToken string) (*ReturnToken, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}
	parsedToken, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return nil, errors.New("unable to parse claims")
	}
	
	// id:=fmt.Sprint(parsedToken)
	return &ReturnToken{
		Id: fmt.Sprint(parsedToken["id"]),
		Role: fmt.Sprint(parsedToken["role"]),
	}, nil
}