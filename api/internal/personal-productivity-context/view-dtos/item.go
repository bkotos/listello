package viewdto

import domain "github.com/bkotos/listello/internal/personal-productivity-context/domain"

// DefineItemRequest is the HTTP request body for defining an item on a list.
type DefineItemRequest struct {
	Title string `json:"title"`
}

// ModifyItemTitleRequest is the HTTP request body for changing an item's title.
type ModifyItemTitleRequest struct {
	Title string `json:"title"`
}

// MoveItemRequest is the HTTP request body for moving an item to another list.
type MoveItemRequest struct {
	ListID string `json:"listID"`
}

// CommentItemRequest is the HTTP request body for commenting on an item.
type CommentItemRequest struct {
	Body string `json:"body"`
}

// AddSubtaskRequest is the HTTP request body for adding a subtask to an item.
type AddSubtaskRequest struct {
	Title string `json:"title"`
}

// PrioritizeSubtaskRequest is the HTTP request body for changing a subtask's priority.
type PrioritizeSubtaskRequest struct {
	Priority string `json:"priority"`
}

// SubtaskDto is the HTTP representation of a subtask on an item.
type SubtaskDto struct {
	ID        string `json:"ID"`
	Title     string `json:"Title"`
	Completed bool   `json:"Completed"`
	Priority  string `json:"Priority,omitempty"`
}

// ItemDto is the HTTP representation of an item.
type ItemDto struct {
	ID          string       `json:"ID"`
	ListID      string       `json:"ListID"`
	ParentID    string       `json:"ParentID"`
	Title       string       `json:"Title"`
	Description string       `json:"Description"`
	DueDate     string       `json:"DueDate"`
	Tags        []string     `json:"Tags"`
	Priority    string       `json:"Priority"`
	State       string       `json:"State"`
	Subtasks    []SubtaskDto `json:"Subtasks,omitempty"`
}

// ItemFromDomain maps a domain item to its response DTO.
func ItemFromDomain(item domain.Item) ItemDto {
	subtasks := make([]SubtaskDto, len(item.Subtasks))
	for i, subtask := range item.Subtasks {
		subtasks[i] = SubtaskDto{
			ID:        subtask.ID,
			Title:     subtask.Title,
			Completed: subtask.Completed,
			Priority:  string(subtask.Priority),
		}
	}
	return ItemDto{
		ID:          item.ID,
		ListID:      item.ListID,
		ParentID:    item.ParentID,
		Title:       item.Title,
		Description: item.Description,
		DueDate:     item.DueDate,
		Tags:        item.Tags,
		Priority:    string(item.Priority),
		State:       string(item.State),
		Subtasks:    subtasks,
	}
}

// ItemsFromDomain maps domain items to response DTOs.
func ItemsFromDomain(items []domain.Item) []ItemDto {
	response := make([]ItemDto, len(items))
	for i, item := range items {
		response[i] = ItemFromDomain(item)
	}
	return response
}
