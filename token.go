package main

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "time"
    _ "github.com/lib/pq"
)

//Creates the JWT token for the user's session
func CreateJWTToken(username string) (string, error) {
  userToken := jwt.New(jwt.SigningMethodHS256)
  userClaims := userToken.Claims.(jwt.MapClaims)

  userClaims["name"] = fmt.Sprintf("%s", username)
  userClaims["exp"] = time.Now().Add(time.Minute * 10).Unix()

  tokenString, err2 := userToken.SignedString(SecretKey)

  if err2 != nil {
    return "", err2
  }

  return tokenString, nil
}
