package adapter

import "time"

func newCreatedAt() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
