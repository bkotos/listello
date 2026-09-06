package commands

import application "github.com/bkotos/listello/internal/personal-productivity-context/application"

type Container interface {
	ListService() application.ListService
	ItemService() application.ItemService
}
