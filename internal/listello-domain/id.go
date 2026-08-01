package domain

import "github.com/google/uuid"

func newID(prefix string) string {
	return prefix + uuid.NewString()
}
