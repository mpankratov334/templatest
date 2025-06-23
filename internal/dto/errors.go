package dto

import "github.com/pkg/errors"

var (
	ErrNotFound  = errors.New("object not found")
	ErrInvalidID = errors.New("invalid id")
)
