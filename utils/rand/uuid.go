package urand

import "github.com/google/uuid"

func RandUuid() string {
	return uuid.New().String()
}

func IsValidUuid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
