package util

import "github.com/google/uuid"

var UUID = func() uuid.UUID {
	return uuid.New()
}
