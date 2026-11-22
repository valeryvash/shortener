package config

import "errors"

var (
	ErrEmptyBindAddress = errors.New("empty bind address")
	ErrEmptyRedirectAddress = errors.New("empty redirect address")
)