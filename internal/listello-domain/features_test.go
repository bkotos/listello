package domain_test

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/require"

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
	items       map[string]domain.Item
	events      []domain.Event
}

func (s *suiteState) reset() {
	s.placeholder = nil
	s.lists = make(map[string]domain.List)
	s.items = make(map[string]domain.Item)
	s.events = nil
}

func (s *suiteState) record(event domain.Event) {
	s.events = append(s.events, event)
}

func (s *suiteState) iCreateADomainPlaceholder() {
	s.placeholder = domain.NewPlaceholder()
}

func (s *suiteState) thePlaceholderShouldExist(ctx context.Context) {
	require.NotNil(godog.T(ctx), s.placeholder)
}

func (s *suiteState) theUserCreatesAListNamed(ctx context.Context, name string) {
	list, event, err := domain.CreateList(name)
	require.NoError(godog.T(ctx), err)
	s.lists[list.Name] = list
	s.record(event)
}

func (s *suiteState) aEventShouldHaveOccurred(ctx context.Context, eventName string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			return e.Name == eventName
		}),
		"expected event %q to have occurred; got %v", eventName, eventNames(s.events),
	)
}

func (s *suiteState) theListShouldExist(ctx context.Context, name string) {
	require.Contains(godog.T(ctx), s.lists, name)
}

func (s *suiteState) aListNamedExists(ctx context.Context, name string) {
	s.theUserCreatesAListNamed(ctx, name)
}

func (s *suiteState) theUserDefinesAnItemTitledOnTheList(ctx context.Context, title, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	item, event, err := domain.DefineItem(s.lists[listName].ID, title)
	require.NoError(t, err)
	s.items[item.Title] = item
	s.record(event)
}

func (s *suiteState) theItemShouldBeOnTheList(ctx context.Context, title, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	require.Contains(t, s.items, title)
	require.Equal(t, s.lists[listName].ID, s.items[title].ListID)
}

func (s *suiteState) theItemShouldBeOutstanding(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Truef(t, s.items[title].IsOutstanding(), "expected item %q to be outstanding; got %q", title, s.items[title].State)
}

func (s *suiteState) anOutstandingItemTitledExistsOnTheList(ctx context.Context, title, listName string) {
	s.theUserDefinesAnItemTitledOnTheList(ctx, title, listName)
}

func (s *suiteState) theOwnerCompletesTheItem(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	item, event, err := domain.CompleteItem(s.items[title])
	require.NoError(t, err)
	s.items[item.Title] = item
	s.record(event)
}

func (s *suiteState) theItemShouldBeComplete(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Truef(t, s.items[title].IsComplete(), "expected item %q to be complete; got %q", title, s.items[title].State)
}

func (s *suiteState) theListShouldHaveAnIDPrefixedWith(ctx context.Context, name, prefix string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, name)
	require.Truef(t, strings.HasPrefix(s.lists[name].ID, prefix), "expected list %q ID to start with %q; got %q", name, prefix, s.lists[name].ID)
}

func (s *suiteState) theItemShouldHaveAnIDPrefixedWith(ctx context.Context, title, prefix string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Truef(t, strings.HasPrefix(s.items[title].ID, prefix), "expected item %q ID to start with %q; got %q", title, prefix, s.items[title].ID)
}

func (s *suiteState) theListsShouldHaveDifferentIDs(ctx context.Context, nameA, nameB string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, nameA)
	require.Contains(t, s.lists, nameB)
	idA := s.lists[nameA].ID
	idB := s.lists[nameB].ID
	require.NotEmpty(t, idA)
	require.NotEmpty(t, idB)
	require.NotEqual(t, idA, idB)
}

func (s *suiteState) theItemsShouldHaveDifferentIDs(ctx context.Context, titleA, titleB string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, titleA)
	require.Contains(t, s.items, titleB)
	idA := s.items[titleA].ID
	idB := s.items[titleB].ID
	require.NotEmpty(t, idA)
	require.NotEmpty(t, idB)
	require.NotEqual(t, idA, idB)
}

func (s *suiteState) aEventShouldHaveOccurredWithTheIDOfList(ctx context.Context, eventName, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	s.eventOccurredWithID(ctx, eventName, s.lists[listName].ID)
}

func (s *suiteState) aEventShouldHaveOccurredWithTheIDOfItem(ctx context.Context, eventName, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	s.eventOccurredWithID(ctx, eventName, s.items[title].ID)
}

func (s *suiteState) aEventShouldHaveOccurredWithTheListIDOfList(ctx context.Context, eventName, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	listID := s.lists[listName].ID
	require.Truef(
		t,
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			return e.Name == eventName && e.ListID == listID
		}),
		"expected event %q with list ID %q; got %v", eventName, listID, eventSummaries(s.events),
	)
}

func (s *suiteState) eventOccurredWithID(ctx context.Context, eventName, id string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			return e.Name == eventName && e.ID == id
		}),
		"expected event %q with ID %q; got %v", eventName, id, eventSummaries(s.events),
	)
}

func eventNames(events []domain.Event) []string {
	names := make([]string, len(events))
	for i, e := range events {
		names[i] = e.Name
	}
	return names
}

func eventSummaries(events []domain.Event) []string {
	summaries := make([]string, len(events))
	for i, e := range events {
		summaries[i] = fmt.Sprintf("%s(%s)", e.Name, e.ID)
	}
	return summaries
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
	ctx.Step(`^a list named "([^"]*)" exists$`, s.aListNamedExists)
	ctx.Step(`^the user defines an item titled "([^"]*)" on the list "([^"]*)"$`, s.theUserDefinesAnItemTitledOnTheList)
	ctx.Step(`^the item "([^"]*)" should be on the list "([^"]*)"$`, s.theItemShouldBeOnTheList)
	ctx.Step(`^the item "([^"]*)" should be outstanding$`, s.theItemShouldBeOutstanding)
	ctx.Step(`^an outstanding item titled "([^"]*)" exists on the list "([^"]*)"$`, s.anOutstandingItemTitledExistsOnTheList)
	ctx.Step(`^the owner completes the item "([^"]*)"$`, s.theOwnerCompletesTheItem)
	ctx.Step(`^the item "([^"]*)" should be complete$`, s.theItemShouldBeComplete)
	ctx.Step(`^the list "([^"]*)" should have an ID prefixed with "([^"]*)"$`, s.theListShouldHaveAnIDPrefixedWith)
	ctx.Step(`^the item "([^"]*)" should have an ID prefixed with "([^"]*)"$`, s.theItemShouldHaveAnIDPrefixedWith)
	ctx.Step(`^the lists "([^"]*)" and "([^"]*)" should have different IDs$`, s.theListsShouldHaveDifferentIDs)
	ctx.Step(`^the items "([^"]*)" and "([^"]*)" should have different IDs$`, s.theItemsShouldHaveDifferentIDs)
	ctx.Step(`^a "([^"]*)" event should have occurred with the ID of list "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheIDOfList)
	ctx.Step(`^a "([^"]*)" event should have occurred with the ID of item "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheIDOfItem)
	ctx.Step(`^a "([^"]*)" event should have occurred with the list ID of list "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheListIDOfList)
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
	require.Zero(t, status, "non-zero status returned, failed to run feature tests")
}
