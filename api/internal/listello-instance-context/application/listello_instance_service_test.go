package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	application "github.com/bkotos/listello/internal/listello-instance-context/application"
	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
)

func TestListelloInstanceService_CreateInstance_PersistsInstance(t *testing.T) {
	// Arrange
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, publisher)

	repo.EXPECT().
		Save(mock.AnythingOfType("domain.ListelloInstance")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	_, err := svc.CreateInstance()

	// Assert
	require.NoError(t, err)
}

func TestListelloInstanceService_CreateInstance_PublishesEvent(t *testing.T) {
	// Arrange
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, publisher)

	var published domain.Event
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.ListelloInstance")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			_, ok := event.Metadata.(domain.EventMetadataInstanceCreated)
			return event.Name == domain.EventInstanceCreated && ok
		})).
		Return(nil)

	// Act
	_, err := svc.CreateInstance()

	// Assert
	require.NoError(t, err)
	_, ok := published.Metadata.(domain.EventMetadataInstanceCreated)
	require.True(t, ok)
	assert.NotEmpty(t, published.Timestamp)
}

func TestListelloInstanceService_GetInstance_ReturnsInstanceFromRepository(t *testing.T) {
	// Arrange
	expected := &domain.ListelloInstance{ID: "LI_1"}
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, publisher)

	repo.EXPECT().
		GetInstance().
		Return(expected, nil)

	// Act
	received, err := svc.GetInstance()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}

func TestListelloInstanceService_GetInstance_ReturnsNullWhenNotExists(t *testing.T) {
	// Arrange
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, publisher)

	repo.EXPECT().
		GetInstance().
		Return(nil, nil)

	// Act
	received, err := svc.GetInstance()

	// Assert
	require.NoError(t, err)
	assert.Nil(t, received)
}
