package application

import (
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

// UserRepository persists users.
type UserRepository interface {
	Save(user domain.User) error
}

// UserService defines user application operations.
type UserService interface {
	CreateUser(name string) (domain.User, error)
}

type userService struct {
	userRepository UserRepository
	eventPublisher EventPublisher
}

var _ UserService = (*userService)(nil)

// NewUserService returns a UserService backed by the given repository and publisher.
func NewUserService(userRepository UserRepository, eventPublisher EventPublisher) UserService {
	return &userService{
		userRepository: userRepository,
		eventPublisher: eventPublisher,
	}
}

// CreateUser creates a user via the domain and persists it.
func (s *userService) CreateUser(name string) (domain.User, error) {
	user, event, err := domain.CreateUser(name)
	if err != nil {
		return domain.User{}, err
	}
	if err := s.userRepository.Save(user); err != nil {
		return domain.User{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.User{}, err
	}
	return user, nil
}
