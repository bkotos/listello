package domain_test

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"reflect"
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
	items       map[string]*domain.Item
	events      []domain.Event
	lastErr     error
}

func (s *suiteState) reset() {
	s.placeholder = nil
	s.lists = make(map[string]domain.List)
	s.items = make(map[string]*domain.Item)
	s.events = nil
	s.lastErr = nil
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
	s.lastErr = err
	if err != nil {
		return
	}
	s.lists[list.Name] = list
	s.record(event)
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

func (s *suiteState) theListShouldExist(ctx context.Context, name string) {
	require.Contains(godog.T(ctx), s.lists, name)
}

func (s *suiteState) aListNamedExists(ctx context.Context, name string) {
	s.theUserCreatesAListNamed(ctx, name)
}

func (s *suiteState) theUserDefinesAnItemTitledOnTheList(ctx context.Context, title, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	item, event, err := domain.DefineItem(s.lists[listName], title)
	s.lastErr = err
	if err != nil {
		return
	}
	s.items[item.Title] = &item
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
	item, event, err := domain.CompleteItem(*s.items[title])
	require.NoError(t, err)
	s.items[item.Title] = &item
	s.record(event)
}

func (s *suiteState) theItemShouldBeComplete(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Truef(t, s.items[title].IsComplete(), "expected item %q to be complete; got %q", title, s.items[title].State)
}

func (s *suiteState) anInboxListExists(ctx context.Context) {
	s.inboxList(ctx)
}

func (s *suiteState) inboxList(ctx context.Context) domain.List {
	const inboxName = "Inbox"
	if list, ok := s.lists[inboxName]; ok {
		return list
	}
	// Inbox is not created via CreateList; seed it for capture/define scenarios.
	list := domain.List{ID: "LS_inbox", Name: inboxName}
	s.lists[inboxName] = list
	return list
}

func (s *suiteState) theUserCapturesAnInboxItem(ctx context.Context, title string) {
	item, event, err := domain.CaptureItem(s.inboxList(ctx), title)
	require.NoError(godog.T(ctx), err)
	s.items[item.Title] = &item
	s.record(event)
}

func (s *suiteState) theUserCapturesAnInboxItemOnTheList(ctx context.Context, title, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.lists, listName)
	_, _, err := domain.CaptureItem(s.lists[listName], title)
	s.lastErr = err
}

func (s *suiteState) theCaptureShouldFailWithError(ctx context.Context, message string) {
	s.theOperationShouldFailWithError(ctx, message)
}

func (s *suiteState) theOperationShouldFailWithError(ctx context.Context, message string) {
	require.EqualError(godog.T(ctx), s.lastErr, message)
}

func (s *suiteState) aCapturedItemExists(ctx context.Context, title string) {
	s.theUserCapturesAnInboxItem(ctx, title)
}

func (s *suiteState) theOwnerModifiesTheTitleOfTheItemTo(ctx context.Context, title, newTitle string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].ModifyTitle(newTitle)
	require.NoError(t, err)
	item := s.items[title]
	delete(s.items, title)
	s.items[item.Title] = item
	s.record(event)
}

func (s *suiteState) theItemShouldExist(ctx context.Context, title string) {
	require.Contains(godog.T(ctx), s.items, title)
}

func (s *suiteState) theOwnerModifiesTheDescriptionOfTheItemTo(ctx context.Context, title, description string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].ModifyDescription(description)
	require.NoError(t, err)
	s.record(event)
}

func (s *suiteState) theItemShouldHaveDescription(ctx context.Context, title, description string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Equal(t, description, s.items[title].Description)
}

func (s *suiteState) theOwnerModifiesTheDueDateOfTheItemTo(ctx context.Context, title, dueDate string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].ModifyDueDate(dueDate)
	s.lastErr = err
	if err != nil {
		return
	}
	s.record(event)
}

func (s *suiteState) theItemShouldBeDueOn(ctx context.Context, title, dueDate string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Equal(t, dueDate, s.items[title].DueDate)
}

func (s *suiteState) theOwnerRemovesTheDueDateOfTheItem(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].RemoveDueDate()
	require.NoError(t, err)
	s.record(event)
}

func (s *suiteState) theItemShouldHaveNoDueDate(ctx context.Context, title string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Empty(t, s.items[title].DueDate)
}

func (s *suiteState) theOwnerTagsTheItemWith(ctx context.Context, title, tag string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].Tag(tag)
	require.NoError(t, err)
	s.record(event)
}

func (s *suiteState) theItemShouldBeTaggedWith(ctx context.Context, title, tag string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Contains(t, s.items[title].Tags, tag)
}

func (s *suiteState) theOwnerUntagsTheItemWith(ctx context.Context, title, tag string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].Untag(tag)
	require.NoError(t, err)
	s.record(event)
}

func (s *suiteState) theItemShouldNotBeTaggedWith(ctx context.Context, title, tag string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.NotContains(t, s.items[title].Tags, tag)
}

func (s *suiteState) theOwnerChangesThePriorityOfTheItemTo(ctx context.Context, title, priority string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	event, err := s.items[title].ChangePriority(domain.ItemPriority(priority))
	require.NoError(t, err)
	s.record(event)
}

func (s *suiteState) theItemShouldHavePriority(ctx context.Context, title, priority string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Equal(t, domain.ItemPriority(priority), s.items[title].Priority)
}

func (s *suiteState) theOwnerMovesTheItemToTheList(ctx context.Context, title, listName string) {
	t := godog.T(ctx)
	require.Contains(t, s.items, title)
	require.Contains(t, s.lists, listName)
	event, err := s.items[title].Move(s.lists[listName])
	require.NoError(t, err)
	s.record(event)
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
			return e.Name == domain.EventName(eventName) && eventListID(e) == listID
		}),
		"expected event %q with list ID %q; got %v", eventName, listID, eventSummaries(s.events),
	)
}

func (s *suiteState) aEventShouldHaveOccurredWithTitle(ctx context.Context, eventName, title string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			meta, ok := e.Metadata.(domain.EventMetadataItemTitleChanged)
			return e.Name == domain.EventName(eventName) && ok && meta.Title == title
		}),
		"expected event %q with title %q; got %v", eventName, title, eventSummaries(s.events),
	)
}

func (s *suiteState) aEventShouldHaveOccurredWithDescription(ctx context.Context, eventName, description string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			meta, ok := e.Metadata.(domain.EventMetadataItemDescriptionChanged)
			return e.Name == domain.EventName(eventName) && ok && meta.Description == description
		}),
		"expected event %q with description %q; got %v", eventName, description, eventSummaries(s.events),
	)
}

func (s *suiteState) aEventShouldHaveOccurredWithDueDate(ctx context.Context, eventName, dueDate string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			meta, ok := e.Metadata.(domain.EventMetadataDueDateAddedToItem)
			return e.Name == domain.EventName(eventName) && ok && meta.DueDate == dueDate
		}),
		"expected event %q with due date %q; got %v", eventName, dueDate, eventSummaries(s.events),
	)
}

func (s *suiteState) aEventShouldHaveOccurredWithTag(ctx context.Context, eventName, tag string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			meta, ok := e.Metadata.(domain.EventMetadataTagAddedToItem)
			return e.Name == domain.EventName(eventName) && ok && meta.Tag == tag
		}),
		"expected event %q with tag %q; got %v", eventName, tag, eventSummaries(s.events),
	)
}

func (s *suiteState) aEventShouldHaveOccurredWithPriority(ctx context.Context, eventName, priority string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			meta, ok := e.Metadata.(domain.EventMetadataSubtaskPriorityChanged)
			return e.Name == domain.EventName(eventName) && ok && meta.Priority == domain.ItemPriority(priority)
		}),
		"expected event %q with priority %q; got %v", eventName, priority, eventSummaries(s.events),
	)
}

func (s *suiteState) eventOccurredWithID(ctx context.Context, eventName, id string) {
	require.Truef(
		godog.T(ctx),
		slices.ContainsFunc(s.events, func(e domain.Event) bool {
			return e.Name == domain.EventName(eventName) && eventEntityID(e) == id
		}),
		"expected event %q with ID %q; got %v", eventName, id, eventSummaries(s.events),
	)
}

func eventEntityID(e domain.Event) string {
	switch e.Metadata.(type) {
	case domain.EventMetadataListCreated,
		domain.EventMetadataItemDefined,
		domain.EventMetadataItemCompleted,
		domain.EventMetadataItemCaptured,
		domain.EventMetadataItemTitleChanged,
		domain.EventMetadataItemDescriptionChanged,
		domain.EventMetadataDueDateAddedToItem,
		domain.EventMetadataDueDateRemovedFromItem,
		domain.EventMetadataTagAddedToItem,
		domain.EventMetadataTagRemovedFromItem,
		domain.EventMetadataSubtaskPriorityChanged,
		domain.EventMetadataItemMovedToOtherList:
		return reflect.ValueOf(e.Metadata).FieldByName("ID").String()
	default:
		return ""
	}
}

func eventListID(e domain.Event) string {
	switch meta := e.Metadata.(type) {
	case domain.EventMetadataItemDefined:
		return meta.ListID
	case domain.EventMetadataItemCaptured:
		return meta.ListID
	case domain.EventMetadataItemMovedToOtherList:
		return meta.ListID
	default:
		return ""
	}
}

func eventNames(events []domain.Event) []string {
	names := make([]string, len(events))
	for i, e := range events {
		names[i] = string(e.Name)
	}
	return names
}

func eventSummaries(events []domain.Event) []string {
	summaries := make([]string, len(events))
	for i, e := range events {
		summaries[i] = fmt.Sprintf("%s(%+v)", e.Name, e.Metadata)
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
	ctx.Step(`^an inbox list exists$`, s.anInboxListExists)
	ctx.Step(`^the user defines an item titled "([^"]*)" on the list "([^"]*)"$`, s.theUserDefinesAnItemTitledOnTheList)
	ctx.Step(`^the item "([^"]*)" should be on the list "([^"]*)"$`, s.theItemShouldBeOnTheList)
	ctx.Step(`^the item "([^"]*)" should be outstanding$`, s.theItemShouldBeOutstanding)
	ctx.Step(`^an outstanding item titled "([^"]*)" exists on the list "([^"]*)"$`, s.anOutstandingItemTitledExistsOnTheList)
	ctx.Step(`^the owner completes the item "([^"]*)"$`, s.theOwnerCompletesTheItem)
	ctx.Step(`^the item "([^"]*)" should be complete$`, s.theItemShouldBeComplete)
	ctx.Step(`^creating the list should fail with error "([^"]*)"$`, s.theOperationShouldFailWithError)
	ctx.Step(`^defining the item should fail with error "([^"]*)"$`, s.theOperationShouldFailWithError)
	ctx.Step(`^modifying the due date should fail with error "([^"]*)"$`, s.theOperationShouldFailWithError)
	ctx.Step(`^the user captures an inbox item "([^"]*)"$`, s.theUserCapturesAnInboxItem)
	ctx.Step(`^the user captures an inbox item "([^"]*)" on the list "([^"]*)"$`, s.theUserCapturesAnInboxItemOnTheList)
	ctx.Step(`^the capture should fail with error "([^"]*)"$`, s.theCaptureShouldFailWithError)
	ctx.Step(`^a captured item "([^"]*)" exists$`, s.aCapturedItemExists)
	ctx.Step(`^the owner modifies the title of the item "([^"]*)" to "([^"]*)"$`, s.theOwnerModifiesTheTitleOfTheItemTo)
	ctx.Step(`^the item "([^"]*)" should exist$`, s.theItemShouldExist)
	ctx.Step(`^the owner modifies the description of the item "([^"]*)" to "([^"]*)"$`, s.theOwnerModifiesTheDescriptionOfTheItemTo)
	ctx.Step(`^the item "([^"]*)" should have description "([^"]*)"$`, s.theItemShouldHaveDescription)
	ctx.Step(`^the owner modifies the due date of the item "([^"]*)" to "([^"]*)"$`, s.theOwnerModifiesTheDueDateOfTheItemTo)
	ctx.Step(`^the item "([^"]*)" should be due on "([^"]*)"$`, s.theItemShouldBeDueOn)
	ctx.Step(`^the owner removes the due date of the item "([^"]*)"$`, s.theOwnerRemovesTheDueDateOfTheItem)
	ctx.Step(`^the item "([^"]*)" should have no due date$`, s.theItemShouldHaveNoDueDate)
	ctx.Step(`^the owner tags the item "([^"]*)" with "([^"]*)"$`, s.theOwnerTagsTheItemWith)
	ctx.Step(`^the item "([^"]*)" should be tagged with "([^"]*)"$`, s.theItemShouldBeTaggedWith)
	ctx.Step(`^the owner untags the item "([^"]*)" with "([^"]*)"$`, s.theOwnerUntagsTheItemWith)
	ctx.Step(`^the item "([^"]*)" should not be tagged with "([^"]*)"$`, s.theItemShouldNotBeTaggedWith)
	ctx.Step(`^the owner changes the priority of the item "([^"]*)" to "([^"]*)"$`, s.theOwnerChangesThePriorityOfTheItemTo)
	ctx.Step(`^the item "([^"]*)" should have priority "([^"]*)"$`, s.theItemShouldHavePriority)
	ctx.Step(`^the owner moves the item "([^"]*)" to the list "([^"]*)"$`, s.theOwnerMovesTheItemToTheList)
	ctx.Step(`^the list "([^"]*)" should have an ID prefixed with "([^"]*)"$`, s.theListShouldHaveAnIDPrefixedWith)
	ctx.Step(`^the item "([^"]*)" should have an ID prefixed with "([^"]*)"$`, s.theItemShouldHaveAnIDPrefixedWith)
	ctx.Step(`^the lists "([^"]*)" and "([^"]*)" should have different IDs$`, s.theListsShouldHaveDifferentIDs)
	ctx.Step(`^the items "([^"]*)" and "([^"]*)" should have different IDs$`, s.theItemsShouldHaveDifferentIDs)
	ctx.Step(`^a "([^"]*)" event should have occurred with the ID of list "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheIDOfList)
	ctx.Step(`^a "([^"]*)" event should have occurred with the ID of item "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheIDOfItem)
	ctx.Step(`^a "([^"]*)" event should have occurred with the list ID of list "([^"]*)"$`, s.aEventShouldHaveOccurredWithTheListIDOfList)
	ctx.Step(`^a "([^"]*)" event should have occurred with title "([^"]*)"$`, s.aEventShouldHaveOccurredWithTitle)
	ctx.Step(`^a "([^"]*)" event should have occurred with description "([^"]*)"$`, s.aEventShouldHaveOccurredWithDescription)
	ctx.Step(`^a "([^"]*)" event should have occurred with due date "([^"]*)"$`, s.aEventShouldHaveOccurredWithDueDate)
	ctx.Step(`^a "([^"]*)" event should have occurred with tag "([^"]*)"$`, s.aEventShouldHaveOccurredWithTag)
	ctx.Step(`^a "([^"]*)" event should have occurred with priority "([^"]*)"$`, s.aEventShouldHaveOccurredWithPriority)
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
