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

type SignUpData struct {
  FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
  Email string `json:"email"`
	Password string `json:"password"`
  PWReentry string `json:"reenter"`
}

//Creates a new row in the users table in the database and creates a JWT token for their first session
func SignUp(w http.ResponseWriter, r *http.Request) {

  jsonInBytes, err := ioutil.ReadAll(r.Body)

  if err != nil {
    log.Print(err)

    return
  }

  var newUser SignUpData

  err = json.Unmarshal(jsonInBytes, &newUser)

  if err != nil {
    log.Print(err)
    http.Error(w, "Invalid Request Body", http.StatusBadRequest)
    return
  }

  errFound, errors := ErrorHandler(newUser, w)

  if errFound {
    marshalledJSON, err1 := json.Marshal(errors)

    if err1 != nil {
      log.Print(err1)
      http.Error(w, "Invalid Request Body", http.StatusBadRequest)
      return
    }

    w.Write(marshalledJSON)
    return
  }

  sqlEmailCheckStmt := `SELECT email FROM users WHERE email=$1`
  var email string

  err2 := db.QueryRow(sqlEmailCheckStmt, newUser.Email).Scan(&email)

  if err2 == nil {
    http.Error(w, "That email is already taken", http.StatusUnauthorized)
    return
  } else if err2 != sql.ErrNoRows {
    log.Println(err2)
    http.Error(w, "Trouble checking through our database", http.StatusInternalServerError)
    return
  }


  sqlStmt := `INSERT INTO users (first_name, last_name, email, encrypted_password) VALUES ($1, $2, $3, $4) RETURNING id`
  id := 0
  hashedPassword, err3 := PasswordHash(newUser.Password)

  if err3 != nil {
    panic(err3)
  }

  err3 = db.QueryRow(sqlStmt, newUser.FirstName, newUser.LastName, newUser.Email, hashedPassword).Scan(&id)

  if err3 != nil {
    panic(err3)
  }

  fmt.Println("New users record id is ", id)

  tokenString, err4 := CreateJWTToken(newUser.Email)

  if err4 != nil {
    log.Println(err4)
    http.Error(w, "Trouble Creating JWT token", http.StatusInternalServerError)
    return
  }

  fmt.Println("New user successfully created!")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(tokenString))
}
