package types

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	FullName string `json:"full_name"`
	Id       int64  `json:"id"`
	jwt.StandardClaims
}
