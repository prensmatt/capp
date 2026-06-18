package models

import(
	"errors"
)

var(
	ErrNotFound = errors.New("record not found")
	ErrDuplicateEmail =	errors.New("duplicate email")
	ErrInvalidInput =	errors.New("invalid input")
)