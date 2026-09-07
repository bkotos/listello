package domain_test

import (
	"context"
	"embed"
	"flag"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domain "github.com/bkotos/listello/internal/listello-instance-context/domain"
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
	instance *domain.ListelloInstance
	events   []domain.Event
	lastErr  error
}

func (s *suiteState) reset() {
	s.instance = nil
	s.events = nil
	s.lastErr = nil
}

func (s *suiteState) record(event domain.Event) {
	s.events = append(s.events, event)
}

func (s *suiteState) theUserCreatesAnInstance() {
	instance, ev, err := domain.CreateInstance()
	s.lastErr = err
	if err != nil {
		return
	}
	s.instance = &instance
	s.record(ev)
}

func (s *suiteState) aEventShouldHaveOccurred(ctx context.Context, eventName string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			return e.Name == domain.EventName(eventName)
		}),
		"expected event %q to have occurred; got %v", eventName, eventNames(s.events),
	)
}

func (s *suiteState) theInstanceShouldExist(ctx context.Context) {
	require.NotNil(godog.T(ctx), s.instance)
}

func (s *suiteState) theInstanceShouldHaveAnIDPrefixedWith(ctx context.Context, prefix string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	require.Truef(t, strings.HasPrefix(s.instance.ID, prefix), "expected instance ID to start with %q; got %q", prefix, s.instance.ID)
}

func (s *suiteState) theInstanceIDAfterThePrefixShouldBeAUUID(ctx context.Context, prefix string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	remainder := strings.TrimPrefix(s.instance.ID, prefix)
	_, err := uuid.Parse(remainder)
	require.NoErrorf(t, err, "expected instance ID after prefix %q to be a UUID; got %q", prefix, remainder)
}

func (s *suiteState) anInstanceExists() {
	s.theUserCreatesAnInstance()
}

func (s *suiteState) theUserSelectsHostingMode(ctx context.Context, mode string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	ev, err := s.instance.SelectHostingMode(domain.HostingMode(mode))
	s.lastErr = err
	if err != nil {
		return
	}
	s.record(ev)
}

func (s *suiteState) theInstanceShouldHaveHostingMode(ctx context.Context, mode string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	require.Equal(t, domain.HostingMode(mode), s.instance.HostingMode)
}

func (s *suiteState) theUserSelectsPersistenceLocation(ctx context.Context, location string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	ev, err := s.instance.SelectPersistenceLocation(location)
	s.lastErr = err
	if err != nil {
		return
	}
	s.record(ev)
}

func (s *suiteState) theInstanceShouldHavePersistenceLocation(ctx context.Context, location string) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	require.Equal(t, location, s.instance.PersistenceLocation)
}

func (s *suiteState) selectingThePersistenceLocationShouldFailWithError(ctx context.Context, message string) {
	require.EqualError(godog.T(ctx), s.lastErr, message)
}

func (s *suiteState) theUserInitializesPersistence(ctx context.Context) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	ev, err := s.instance.InitializePersistence()
	s.lastErr = err
	if err != nil {
		return
	}
	s.record(ev)
}

func (s *suiteState) theInstanceShouldHavePersistenceInitialized(ctx context.Context) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	require.True(t, s.instance.IsPersistenceInitialized())
}

func (s *suiteState) theUserCompletesSetup(ctx context.Context) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	ev, err := s.instance.CompleteSetup()
	s.lastErr = err
	if err != nil {
		return
	}
	s.record(ev)
}

func (s *suiteState) theInstanceShouldHaveSetupCompleted(ctx context.Context) {
	t := godog.T(ctx)
	require.NotNil(t, s.instance)
	require.True(t, s.instance.IsSetupCompleted())
}

func eventNames(events []domain.Event) []string {
	names := make([]string, len(events))
	for i, e := range events {
		names[i] = string(e.Name)
	}
	return names
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	s := &suiteState{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		s.reset()
		return ctx, nil
	})

	ctx.Step(`^the user creates an instance$`, s.theUserCreatesAnInstance)
	ctx.Step(`^an instance exists$`, s.anInstanceExists)
	ctx.Step(`^the user selects hosting mode "([^"]*)"$`, s.theUserSelectsHostingMode)
	ctx.Step(`^a "([^"]*)" event should have occurred$`, s.aEventShouldHaveOccurred)
	ctx.Step(`^the instance should exist$`, s.theInstanceShouldExist)
	ctx.Step(`^the instance should have an ID prefixed with "([^"]*)"$`, s.theInstanceShouldHaveAnIDPrefixedWith)
	ctx.Step(`^the instance ID after the prefix "([^"]*)" should be a UUID$`, s.theInstanceIDAfterThePrefixShouldBeAUUID)
	ctx.Step(`^the instance should have hosting mode "([^"]*)"$`, s.theInstanceShouldHaveHostingMode)
	ctx.Step(`^the user selects persistence location "([^"]*)"$`, s.theUserSelectsPersistenceLocation)
	ctx.Step(`^the instance should have persistence location "([^"]*)"$`, s.theInstanceShouldHavePersistenceLocation)
	ctx.Step(`^selecting the persistence location should fail with error "([^"]*)"$`, s.selectingThePersistenceLocationShouldFailWithError)
	ctx.Step(`^the user initializes persistence$`, s.theUserInitializesPersistence)
	ctx.Step(`^the instance should have persistence initialized$`, s.theInstanceShouldHavePersistenceInitialized)
	ctx.Step(`^the user completes setup$`, s.theUserCompletesSetup)
	ctx.Step(`^the instance should have setup completed$`, s.theInstanceShouldHaveSetupCompleted)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t
	o.FS = featuresFS
	if o.Output == nil {
		o.Output = os.Stdout
	}

	suite := godog.TestSuite{
		Name:                "domain",
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}

	status := suite.Run()
	if o.ShowStepDefinitions {
		return
	}
	require.Zero(t, status, "non-zero status returned, failed to run feature tests")
}
