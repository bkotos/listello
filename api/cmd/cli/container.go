package main

import (
	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
)

type container struct {
	list     application.ListService
	item     application.ItemService
	instance instanceapp.ListelloInstanceService
}

func (c *container) ListService() application.ListService {
	return c.list
}

func (c *container) ItemService() application.ItemService {
	return c.item
}

func (c *container) InstanceService() instanceapp.ListelloInstanceService {
	return c.instance
}
