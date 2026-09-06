package commands

import application "github.com/bkotos/listello/internal/personal-productivity-context/application"

type testContainer struct {
	list application.ListService
	item application.ItemService
}

func (c testContainer) ListService() application.ListService {
	return c.list
}

func (c testContainer) ItemService() application.ItemService {
	return c.item
}
