package domain_test

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"slices"
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

type suiteState struct {
	placeholder *domain.Placeholder
	lists       map[string]domain.List
	events      []domain.Event
}

func (s *suiteState) reset() {
	s.placeholder = nil
	s.lists = make(map[string]domain.List)
	s.events = nil
}

func (s *suiteState) record(event domain.Event) {
	s.events = append(s.events, event)
}

func (s *suiteState) iCreateADomainPlaceholder() {
	s.placeholder = domain.NewPlaceholder()
}

func (s *suiteState) thePlaceholderShouldExist() error {
	if s.placeholder == nil {
		return fmt.Errorf("expected placeholder to exist, got nil")
	}
	return nil
}

func (s *suiteState) theUserCreatesAListNamed(name string) error {
	list, event, err := domain.CreateList(name)
	if err != nil {
		return err
	}
	s.lists[list.Name()] = list
	s.record(event)
	return nil
}

func (s *suiteState) aEventShouldHaveOccurred(eventName string) error {
	if slices.ContainsFunc(s.events, func(e domain.Event) bool {
		return e.Name == eventName
	}) {
		return nil
	}
	return fmt.Errorf("expected event %q to have occurred; got %v", eventName, eventNames(s.events))
}

func (s *suiteState) theListShouldExist(name string) error {
	if _, ok := s.lists[name]; !ok {
		return fmt.Errorf("expected list %q to exist", name)
	}
	return nil
}

func eventNames(events []domain.Event) []string {
	names := make([]string, len(events))
	for i, e := range events {
		names[i] = e.Name
	}
	return names
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	s := &suiteState{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.reset()
		return ctx, nil
	})

	ctx.Step(`^I create a domain placeholder$`, s.iCreateADomainPlaceholder)
	ctx.Step(`^the placeholder should exist$`, s.thePlaceholderShouldExist)
	ctx.Step(`^the user creates a list named "([^"]*)"$`, s.theUserCreatesAListNamed)
	ctx.Step(`^a "([^"]*)" event should have occurred$`, s.aEventShouldHaveOccurred)
	ctx.Step(`^the list "([^"]*)" should exist$`, s.theListShouldExist)
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
