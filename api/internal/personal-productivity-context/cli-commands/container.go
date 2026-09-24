package commands

import (
	instanceapp "github.com/bkotos/listello/internal/listello-instance-context/application"
	application "github.com/bkotos/listello/internal/personal-productivity-context/application"
)

type Container interface {
	ListService() application.ListService
	ItemService() application.ItemService
	InstanceService() instanceapp.ListelloInstanceService
}
