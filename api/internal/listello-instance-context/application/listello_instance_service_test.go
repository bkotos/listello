package application_test

import (
	"errors"
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
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

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
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

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
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

	repo.EXPECT().
		Get().
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
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

	repo.EXPECT().
		Get().
		Return(nil, nil)

	// Act
	received, err := svc.GetInstance()

	// Assert
	require.NoError(t, err)
	assert.Nil(t, received)
}

func TestListelloInstanceService_SelectHostingMode_PersistsInstance(t *testing.T) {
	// Arrange
	instance := domain.ListelloInstance{ID: "LI_1"}
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

	repo.EXPECT().
		Get().
		Return(&instance, nil)
	repo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.ListelloInstance) bool {
			return saved.HostingMode == domain.HostingModeLocal
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	_, err := svc.SelectHostingMode(domain.HostingModeLocal)

	// Assert
	require.NoError(t, err)
}

func TestListelloInstanceService_SelectHostingMode_PublishesEvent(t *testing.T) {
	// Arrange
	instance := domain.ListelloInstance{ID: "LI_1"}
	repo := NewMockListelloInstanceRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, NewMockPersistenceAdapter(t), publisher)

	var published domain.Event
	repo.EXPECT().
		Get().
		Return(&instance, nil)
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.ListelloInstance")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			_, ok := event.Metadata.(domain.EventMetadataHostingModeSelected)
			return event.Name == domain.EventHostingModeSelected && ok
		})).
		Return(nil)

	// Act
	_, err := svc.SelectHostingMode(domain.HostingModeLocal)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataHostingModeSelected)
	require.True(t, ok)
	assert.Equal(t, domain.HostingModeLocal, metadata.Mode)
	assert.NotEmpty(t, published.Timestamp)
}

func TestListelloInstanceService_SelectPersistenceLocation_PersistsInstance(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	instance := domain.ListelloInstance{ID: "LI_1", HostingMode: domain.HostingModeLocal}
	observation := usablePersistenceLocationObservation()
	repo := NewMockListelloInstanceRepository(t)
	observer := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, observer, publisher)

	repo.EXPECT().
		Get().
		Return(&instance, nil)
	observer.EXPECT().
		ObservePersistenceLocation(location).
		Return(observation, nil)
	repo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.ListelloInstance) bool {
			return saved.Persistence.Location == location
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	_, err := svc.SelectPersistenceLocation(location)

	// Assert
	require.NoError(t, err)
}

func TestListelloInstanceService_SelectPersistenceLocation_PublishesEvent(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	instance := domain.ListelloInstance{ID: "LI_1", HostingMode: domain.HostingModeLocal}
	observation := usablePersistenceLocationObservation()
	repo := NewMockListelloInstanceRepository(t)
	observer := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, observer, publisher)

	var published domain.Event
	repo.EXPECT().
		Get().
		Return(&instance, nil)
	observer.EXPECT().
		ObservePersistenceLocation(location).
		Return(observation, nil)
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.ListelloInstance")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			_, ok := event.Metadata.(domain.EventMetadataLocalPersistenceLocationSelected)
			return event.Name == domain.EventLocalPersistenceLocationSelected && ok
		})).
		Return(nil)

	// Act
	_, err := svc.SelectPersistenceLocation(location)

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataLocalPersistenceLocationSelected)
	require.True(t, ok)
	assert.Equal(t, location, metadata.Location)
	assert.NotEmpty(t, published.Timestamp)
}

func TestListelloInstanceService_InitializePersistence_PersistsInstance(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	instance := domain.ListelloInstance{
		ID:          "LI_1",
		Persistence: domain.Persistence{Location: location},
	}
	repo := NewMockListelloInstanceRepository(t)
	persistence := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, persistence, publisher)

	repo.EXPECT().
		Get().
		Return(&instance, nil)
	persistence.EXPECT().
		ProvisionStorage(location).
		Return(nil)
	repo.EXPECT().
		Save(mock.MatchedBy(func(saved domain.ListelloInstance) bool {
			return saved.Persistence.IsInitialized()
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	_, err := svc.InitializePersistence()

	// Assert
	require.NoError(t, err)
}

func TestListelloInstanceService_InitializePersistence_PublishesEvent(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	instance := domain.ListelloInstance{
		ID:          "LI_1",
		HostingMode: domain.HostingModeLocal,
		Persistence: domain.Persistence{Location: location},
	}
	repo := NewMockListelloInstanceRepository(t)
	persistence := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, persistence, publisher)

	var published domain.Event
	repo.EXPECT().
		Get().
		Return(&instance, nil)
	persistence.EXPECT().
		ProvisionStorage(location).
		Return(nil)
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.ListelloInstance")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			_, ok := event.Metadata.(domain.EventMetadataPersistenceInitialized)
			return event.Name == domain.EventPersistenceInitialized && ok
		})).
		Return(nil)

	// Act
	_, err := svc.InitializePersistence()

	// Assert
	require.NoError(t, err)
	metadata, ok := published.Metadata.(domain.EventMetadataPersistenceInitialized)
	require.True(t, ok)
	assert.Equal(t, instance.ID, metadata.ID)
	assert.Equal(t, domain.HostingModeLocal, metadata.Mode)
	assert.Equal(t, location, metadata.Location)
	assert.NotEmpty(t, published.Timestamp)
}

func TestListelloInstanceService_InitializePersistence_DoesNotInitializeWhenProvisionStorageFails(t *testing.T) {
	// Arrange
	const location = "/var/listello"
	instance := domain.ListelloInstance{
		ID:          "LI_1",
		HostingMode: domain.HostingModeLocal,
		Persistence: domain.Persistence{Location: location},
	}
	repo := NewMockListelloInstanceRepository(t)
	persistence := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, persistence, publisher)

	repo.EXPECT().
		Get().
		Return(&instance, nil)
	persistence.EXPECT().
		ProvisionStorage(location).
		Return(errors.New("provision failed"))

	// Act
	_, err := svc.InitializePersistence()

	// Assert
	require.Error(t, err)
	assert.False(t, instance.Persistence.IsInitialized())
}

func usablePersistenceLocationObservation() domain.PersistenceLocationObservation {
	var observation domain.PersistenceLocationObservation
	observation.SetParentExists(true)
	observation.SetParentWritable(true)
	observation.SetLocationExists(false)
	return observation
}

func TestListelloInstanceService_GetDefaultPersistenceLocation_ReturnsLocationFromAdapter(t *testing.T) {
	// Arrange
	const location = "/Users/me/Library/Application Support/listello"
	repo := NewMockListelloInstanceRepository(t)
	persistence := NewMockPersistenceAdapter(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewListelloInstanceService(repo, persistence, publisher)

	persistence.EXPECT().
		GetDefaultPersistenceLocation().
		Return(location, nil)

	// Act
	got, err := svc.GetDefaultPersistenceLocation()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, location, got)
}
