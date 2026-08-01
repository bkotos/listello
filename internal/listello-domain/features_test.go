package domain_test

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"

	domain "github.com/bkotos/listello/internal/listello-domain"
)

//go:embed features/*.feature
var featuresFS embed.FS

var opts = godog.Options{
	Format: "pretty",
	Paths:  []string{"features"},
	FS:     featuresFS,
	// @wip marks scenarios that are specified but not yet implemented.
	Tags: "~@wip",
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

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
	o := opts
	o.TestingT = t
	o.FS = featuresFS
	if o.Output == nil {
		o.Output = os.Stdout
	}

	suite := godog.TestSuite{
		Name:                "listello-domain",
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}

	status := suite.Run()
	if o.ShowStepDefinitions {
		return
	}
	if status != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
