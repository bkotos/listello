//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	application "github.com/bkotos/listello/internal/listello-application"
	domain "github.com/bkotos/listello/internal/listello-domain"
)

func main() {
	lists := &memoryListRepository{lists: make(map[string]domain.List)}
	events := &memoryEventPublisher{}
	listService := application.NewListService(lists, events)
	itemService := application.NewItemService(&memoryItemRepository{}, events)

	js.Global().Set("listello", map[string]any{
		"list": map[string]any{
			"createList": js.FuncOf(func(_ js.Value, args []js.Value) any {
				return createList(listService, args)
			}),
		},
		"item": map[string]any{
			"defineItem": js.FuncOf(func(_ js.Value, args []js.Value) any {
				return defineItem(itemService, args)
			}),
		},
	})

	select {}
}

func createList(svc *application.ListService, args []js.Value) any {
	if len(args) < 1 {
		panic("createList: name is required")
	}
	list, err := svc.CreateList(args[0].String())
	if err != nil {
		panic(err.Error())
	}
	return listToJS(list)
}

func defineItem(svc *application.ItemService, args []js.Value) any {
	if len(args) < 2 {
		panic("defineItem: listId and title are required")
	}
	item, err := svc.DefineItem(args[0].String(), args[1].String())
	if err != nil {
		panic(err.Error())
	}
	return itemToJS(item)
}

func listToJS(list domain.List) map[string]any {
	return map[string]any{
		"ID":   list.ID,
		"Name": list.Name,
	}
}

func itemToJS(item domain.Item) map[string]any {
	tags := make([]any, len(item.Tags))
	for i, t := range item.Tags {
		tags[i] = t
	}
	return map[string]any{
		"ID":          item.ID,
		"ListID":      item.ListID,
		"ParentID":    item.ParentID,
		"Title":       item.Title,
		"Description": item.Description,
		"DueDate":     item.DueDate,
		"Tags":        tags,
		"Priority":    string(item.Priority),
		"State":       string(item.State),
	}
}

type memoryListRepository struct {
	lists map[string]domain.List
}

func (r *memoryListRepository) Save(list domain.List) error {
	r.lists[list.ID] = list
	return nil
}

type memoryItemRepository struct{}

func (r *memoryItemRepository) Save(listID string, item domain.Item) error {
	return fmt.Errorf("not implemented")
}

type memoryEventPublisher struct{}

func (p *memoryEventPublisher) Publish(event domain.Event) error {
	return nil
}
