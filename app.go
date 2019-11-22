package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  password = "Wooloo12!@"
  dbname = "snorlax"

  //Need to replace with .env file system in the future
  keyPath = `D:\AtomProjects/auth0_login/controllers/hmac_key`
)

var (
  SecretKey []byte
  db *sql.DB
)

//Sets up the keypath for the JWT tokens. (Need to change for .env file system)
func init() {
  keyData, err := ioutil.ReadFile(keyPath)

  if err != nil {
    log.Fatal(err)
  } else {
    SecretKey = keyData
  }
}

//Connects to both the database and localhost for the testing ground
func main() {
  var err error
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

  log.Print(psqlInfo)

  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  defer db.Close()

  err = db.Ping()

  if err != nil {
    panic(err)
  } else {
    fmt.Println("Successfully connected!")
  }

  http.HandleFunc("/hi", Hi)

  http.HandleFunc("/signup", SignUp)

  http.HandleFunc("/signin", SignIn)

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func Hi(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi")
}
