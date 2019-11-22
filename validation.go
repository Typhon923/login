package main

import (
    "regexp"
    "net/http"
    _ "github.com/lib/pq"
    "golang.org/x/crypto/bcrypt"
)

type Errors struct {
  Errors []ErrorDetails `json:"errors"`
}

type ErrorDetails struct {
  Error string `json:"error"`
  MissingVariable string `json:"variable,omitempty"`
  Message string `json:"message,omitempty"`
}

//Handles all bad request type errors
func ErrorHandler(newUser SignUpData, w http.ResponseWriter) (bool, Errors) {
  var emailRegex string = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])"
  var errorFound bool = false
  var errors Errors
  errDetails := ErrorDetails{Error: "Missing variable", Message: "Missing a variable in request."}

  if newUser.FirstName == "" {
    errDetails.MissingVariable = "firstname"
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }
  if newUser.LastName == "" {
    errDetails.MissingVariable = "lastname"
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }
  if newUser.Email == "" {
    errDetails.MissingVariable = "email"
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }
  if newUser.Password == "" {
    errDetails.MissingVariable = "password"
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }
  if newUser.PWReentry == "" {
    errDetails.MissingVariable = "reenter"
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }
  if newUser.Password != newUser.PWReentry {
    errDetails = ErrorDetails{Error: "Mismatched passwords", Message: "Mismatched password and password reentry."}
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }

  match, _ := regexp.MatchString(emailRegex, newUser.Email)
  if newUser.Email != "" && !match {
    errDetails = ErrorDetails{Error: "Invalid email", Message: "Email address doesn't adhere to standard form."}
    errors.Errors = append(errors.Errors, errDetails)
    errorFound = true
  }

  if errorFound {
    w.WriteHeader(http.StatusBadRequest)
  }

  return errorFound, errors
}

/* Generates a matching bcrypt variation of the user's password.
Returns a string representation of the bcrypt password, and an error object. */
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(bytes), err
}
