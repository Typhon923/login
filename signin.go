package main

import (
    "fmt"
    "log"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"
)

type SignInData struct {
  Email string `json:"email"`
	Password string `json:"password"`
}

//Verifies a username and password from the request is one that exists in the database and creates a JWT token for the session if successful
func SignIn(w http.ResponseWriter, r *http.Request) {

  jsonInBytes, err := ioutil.ReadAll(r.Body)

  if err != nil {
    log.Print(err)
    return
  }

  var signedInUser SignInData

  err = json.Unmarshal(jsonInBytes, &signedInUser)

  if err != nil {
    log.Print(err)
    http.Error(w, "Invalid Request Body", http.StatusBadRequest)
    return
  }

  sqlStmt := `SELECT id, email FROM users WHERE email=$1 AND encrypted_password=$2`

  var (
    email string
    id int
  )

  hashedPassword, hashErr := PasswordHash(signedInUser.Password)

  if hashErr != nil {
    panic(hashErr)
  }

  row := db.QueryRow(sqlStmt, signedInUser.Email, hashedPassword)

  switch err1 := row.Scan(&id, &email); err {
  case sql.ErrNoRows:
    log.Print(err1)
    http.Error(w, "Couldn't find this user in our database", http.StatusUnauthorized)
    return
  }

  tokenString, err2 := CreateJWTToken(signedInUser.Email)

  if err2 != nil {
    log.Println(err2)
    http.Error(w, "Trouble Creating JWT token", http.StatusInternalServerError)
    return
  }

  fmt.Printf("User \"%s\" has logged in", signedInUser.Email)
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(tokenString))
}
