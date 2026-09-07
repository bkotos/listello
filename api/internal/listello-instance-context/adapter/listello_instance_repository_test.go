package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	adapter "github.com/bkotos/listello/internal/listello-instance-context/adapter"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

func TestListelloInstanceRepository_SaveAndGet(t *testing.T) {
	// Arrange
	repo := adapter.NewListelloInstanceRepository()
	instance, _, err := domain.CreateInstance()
	require.NoError(t, err)

	// Act
	require.NoError(t, repo.Save(instance))
	got, err := repo.Get()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, instance, got)
}
