package utils

import "github.com/google/uuid"

func IDGenerator() string {
	id := uuid.New()
	return id.String()
}
