package config

import "github.com/golang-jwt/jwt/v5"

var SecretKey = []byte("kuncirahasiaJWTKEY12361q2*")
var KunciRahasia = "asdasdasdas*"

type JWTClaim struct {
	UserId int    `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
