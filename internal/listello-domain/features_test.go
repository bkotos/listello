package domain_test

import (
	"context"
	"embed"
	"fmt"
	"testing"

	"github.com/cucumber/godog"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

//go:embed features/*.feature
var featuresFS embed.FS

type placeholderSuite struct {
	placeholder *domain.Placeholder
}

func (s *placeholderSuite) iCreateADomainPlaceholder() {
	s.placeholder = domain.NewPlaceholder()
}

func (s *placeholderSuite) thePlaceholderShouldExist() error {
	if s.placeholder == nil {
		return fmt.Errorf("expected placeholder to exist, got nil")
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	s := &placeholderSuite{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.placeholder = nil
		return ctx, nil
	})

	ctx.Step(`^I create a domain placeholder$`, s.iCreateADomainPlaceholder)
	ctx.Step(`^the placeholder should exist$`, s.thePlaceholderShouldExist)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			FS:       featuresFS,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
