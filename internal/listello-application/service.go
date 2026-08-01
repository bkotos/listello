package application

// Service is the application-layer entry point.
type Service struct{}

// NewService returns a new application Service.
func NewService() *Service {
	return &Service{}
}

// Add returns the sum of a and b.
func (s *Service) Add(a, b int) int {
	return a + b
}
