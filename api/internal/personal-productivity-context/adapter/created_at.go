package adapter

import (
	"time"

	"github.com/bkotos/listello/internal/sqlite"
)

func newCreatedAt() string {
	return sqlite.FormatCreatedAt(time.Now())
}
