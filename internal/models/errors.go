package models

import (
	"errors"
)

var (
	// error for when record we are looking for doesn't exist in db
	ErrNoRecord = errors.New("models: no matching record found")

	// error for when a user tries to login with incorrect email or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// error for when a user tries to signup with an email address that already exists
	ErrDuplicateEmail = errors.New("models: duplicate email")
)