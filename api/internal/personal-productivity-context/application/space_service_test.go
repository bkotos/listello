package application_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

func TestSpaceService_CreateSpace_PersistsSpace(t *testing.T) {
	// Arrange
	const spaceName = "Personal"
	repo := NewMockSpaceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewSpaceService(repo, publisher)

	repo.EXPECT().
		Save(mock.MatchedBy(func(space domain.Space) bool {
			return space.Name == spaceName && strings.HasPrefix(space.ID, "SP_")
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	space, err := svc.CreateSpace(spaceName)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, spaceName, space.Name)
	assert.True(t, strings.HasPrefix(space.ID, "SP_"))
}

func TestSpaceService_CreateSpace_PublishesEvent(t *testing.T) {
	// Arrange
	const spaceName = "Personal"
	repo := NewMockSpaceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewSpaceService(repo, publisher)

	var published domain.Event
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.Space")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			meta, ok := event.Metadata.(domain.EventMetadataSpaceCreated)
			return event.Name == domain.EventSpaceCreated && ok && strings.HasPrefix(meta.ID, "SP_")
		})).
		Return(nil)

	// Act
	space, err := svc.CreateSpace(spaceName)

	// Assert
	require.NoError(t, err)
	meta, ok := published.Metadata.(domain.EventMetadataSpaceCreated)
	require.True(t, ok)
	assert.Equal(t, space.ID, meta.ID)
	assert.NotEmpty(t, published.Timestamp)
}
