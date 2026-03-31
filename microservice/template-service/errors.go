package main

import "errors"

var (
	errSessionClosed = errors.New("terminal session is closed")
	errInvalidParams = errors.New("invalid parameters")
	errSessionNotFound = errors.New("session not found")
)