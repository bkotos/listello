package main

import (
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
)

type container struct {
	list application.ListService
	item application.ItemService
}

func (c *container) ListService() application.ListService {
	return c.list
}

func (c *container) ItemService() application.ItemService {
	return c.item
}
