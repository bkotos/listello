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
	itemLists   map[string]string // item title -> list name
	events      []domain.Event
}

func (s *suiteState) reset() {
	s.placeholder = nil
	s.lists = make(map[string]domain.List)
	s.items = make(map[string]domain.Item)
	s.itemLists = make(map[string]string)
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

func (s *suiteState) aListNamedExists(name string) error {
	return s.theUserCreatesAListNamed(name)
}

func (s *suiteState) theUserDefinesAnItemTitledOnTheList(title, listName string) error {
	if _, ok := s.lists[listName]; !ok {
		return fmt.Errorf("list %q does not exist", listName)
	}
	item, event, err := domain.DefineItem(title)
	if err != nil {
		return err
	}
	s.items[item.Title()] = item
	s.itemLists[item.Title()] = listName
	s.record(event)
	return nil
}

func (s *suiteState) theItemShouldBeOnTheList(title, listName string) error {
	if _, ok := s.lists[listName]; !ok {
		return fmt.Errorf("expected list %q to exist", listName)
	}
	if _, ok := s.items[title]; !ok {
		return fmt.Errorf("expected item %q to exist", title)
	}
	if got := s.itemLists[title]; got != listName {
		return fmt.Errorf("expected item %q on list %q; was on %q", title, listName, got)
	}
	return nil
}

func (s *suiteState) theItemShouldBeOutstanding(title string) error {
	item, ok := s.items[title]
	if !ok {
		return fmt.Errorf("expected item %q to exist", title)
	}
	if !item.IsOutstanding() {
		return fmt.Errorf("expected item %q to be outstanding; got %q", title, item.State())
	}
	return nil
}

func (s *suiteState) anOutstandingItemTitledExistsOnTheList(title, listName string) error {
	return s.theUserDefinesAnItemTitledOnTheList(title, listName)
}

func (s *suiteState) theOwnerCompletesTheItem(title string) error {
	item, ok := s.items[title]
	if !ok {
		return fmt.Errorf("item %q does not exist", title)
	}
	item, event, err := domain.CompleteItem(item)
	if err != nil {
		return err
	}
	s.items[item.Title()] = item
	s.record(event)
	return nil
}

func (s *suiteState) theItemShouldBeComplete(title string) error {
	item, ok := s.items[title]
	if !ok {
		return fmt.Errorf("expected item %q to exist", title)
	}
	if !item.IsComplete() {
		return fmt.Errorf("expected item %q to be complete; got %q", title, item.State())
	}
	return nil
}

func (s *suiteState) theListShouldHaveAnIDPrefixedWith(name, prefix string) error {
	list, ok := s.lists[name]
	if !ok {
		return fmt.Errorf("expected list %q to exist", name)
	}
	if !strings.HasPrefix(list.ID(), prefix) {
		return fmt.Errorf("expected list %q ID to start with %q; got %q", name, prefix, list.ID())
	}
	return nil
}

func (s *suiteState) theItemShouldHaveAnIDPrefixedWith(title, prefix string) error {
	item, ok := s.items[title]
	if !ok {
		return fmt.Errorf("expected item %q to exist", title)
	}
	if !strings.HasPrefix(item.ID(), prefix) {
		return fmt.Errorf("expected item %q ID to start with %q; got %q", title, prefix, item.ID())
	}
	return nil
}

func (s *suiteState) theListsShouldHaveDifferentIDs(nameA, nameB string) error {
	listA, ok := s.lists[nameA]
	if !ok {
		return fmt.Errorf("expected list %q to exist", nameA)
	}
	listB, ok := s.lists[nameB]
	if !ok {
		return fmt.Errorf("expected list %q to exist", nameB)
	}
	if listA.ID() == "" || listB.ID() == "" {
		return fmt.Errorf("expected both lists to have IDs; got %q and %q", listA.ID(), listB.ID())
	}
	if listA.ID() == listB.ID() {
		return fmt.Errorf("expected lists %q and %q to have different IDs; both were %q", nameA, nameB, listA.ID())
	}
	return nil
}

func (s *suiteState) theItemsShouldHaveDifferentIDs(titleA, titleB string) error {
	itemA, ok := s.items[titleA]
	if !ok {
		return fmt.Errorf("expected item %q to exist", titleA)
	}
	itemB, ok := s.items[titleB]
	if !ok {
		return fmt.Errorf("expected item %q to exist", titleB)
	}
	if itemA.ID() == "" || itemB.ID() == "" {
		return fmt.Errorf("expected both items to have IDs; got %q and %q", itemA.ID(), itemB.ID())
	}
	if itemA.ID() == itemB.ID() {
		return fmt.Errorf("expected items %q and %q to have different IDs; both were %q", titleA, titleB, itemA.ID())
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
