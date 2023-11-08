package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no mathing record found")
var ErrInvalidCredential = errors.New("models: Invalid Credential")
var ErrDuplicateEmail = errors.New("models: Email already Exist")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

//CHANGE STARTS HERE
type Share struct {
	ID             int
	User_id        int
	Shared_user_id int
}

//CHANGE TILL HERE
