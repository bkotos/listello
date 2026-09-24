package commands

import (
	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
)

type testContainer struct {
	list     application.ListService
	item     application.ItemService
	instance instanceapp.ListelloInstanceService
}

func (c testContainer) ListService() application.ListService {
	return c.list
}

func (c testContainer) ItemService() application.ItemService {
	return c.item
}

func (c testContainer) InstanceService() instanceapp.ListelloInstanceService {
	return c.instance
}
