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

func TestUserService_CreateUser_PersistsUser(t *testing.T) {
	// Arrange
	const userName = "Alex"
	repo := NewMockUserRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewUserService(repo, publisher)

	repo.EXPECT().
		Save(mock.MatchedBy(func(user domain.User) bool {
			return user.Name == userName && strings.HasPrefix(user.ID, "US_")
		})).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.AnythingOfType("event.Event")).
		Return(nil)

	// Act
	user, err := svc.CreateUser(userName)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, userName, user.Name)
	assert.True(t, strings.HasPrefix(user.ID, "US_"))
}

func TestUserService_CreateUser_PublishesEvent(t *testing.T) {
	// Arrange
	const userName = "Alex"
	repo := NewMockUserRepository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.NewUserService(repo, publisher)

	var published domain.Event
	repo.EXPECT().
		Save(mock.AnythingOfType("domain.User")).
		Return(nil)
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			meta, ok := event.Metadata.(domain.EventMetadataUserCreated)
			return event.Name == domain.EventUserCreated && ok && strings.HasPrefix(meta.ID, "US_")
		})).
		Return(nil)

	// Act
	user, err := svc.CreateUser(userName)

	// Assert
	require.NoError(t, err)
	meta, ok := published.Metadata.(domain.EventMetadataUserCreated)
	require.True(t, ok)
	assert.Equal(t, user.ID, meta.ID)
	assert.NotEmpty(t, published.Timestamp)
}
