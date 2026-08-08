package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	assert.Equal(t, 5, got)
}
