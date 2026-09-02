package commands

import application "github.com/bkotos/listello/internal/application"

type Container interface {
	ListService() application.ListService
	ItemService() application.ItemService
}
