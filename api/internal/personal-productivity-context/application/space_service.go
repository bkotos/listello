package application

import (
	domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"
)

// SpaceRepository persists spaces.
type SpaceRepository interface {
	Save(space domain.Space) error
}

// SpaceService defines space application operations.
type SpaceService interface {
	CreateSpace(name string) (domain.Space, error)
}

type spaceService struct {
	spaceRepository SpaceRepository
	eventPublisher  EventPublisher
}

var _ SpaceService = (*spaceService)(nil)

// NewSpaceService returns a SpaceService backed by the given repository and publisher.
func NewSpaceService(spaceRepository SpaceRepository, eventPublisher EventPublisher) SpaceService {
	return &spaceService{
		spaceRepository: spaceRepository,
		eventPublisher:  eventPublisher,
	}
}

// CreateSpace creates a space via the domain and persists it.
func (s *spaceService) CreateSpace(name string) (domain.Space, error) {
	space, event, err := domain.CreateSpace(name)
	if err != nil {
		return domain.Space{}, err
	}
	if err := s.spaceRepository.Save(space); err != nil {
		return domain.Space{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.Space{}, err
	}
	return space, nil
}
