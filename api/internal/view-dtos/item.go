package viewdto

import domain "github.com/bkotos/listello/internal/domain"

// DefineItemRequest is the HTTP request body for defining an item on a list.
type DefineItemRequest struct {
	Title string `json:"title"`
}

// ModifyItemTitleRequest is the HTTP request body for changing an item's title.
type ModifyItemTitleRequest struct {
	Title string `json:"title"`
}

// ItemDto is the HTTP representation of an item.
type ItemDto struct {
	ID          string   `json:"ID"`
	ListID      string   `json:"ListID"`
	ParentID    string   `json:"ParentID"`
	Title       string   `json:"Title"`
	Description string   `json:"Description"`
	DueDate     string   `json:"DueDate"`
	Tags        []string `json:"Tags"`
	Priority    string   `json:"Priority"`
	State       string   `json:"State"`
}

// ItemFromDomain maps a domain item to its response DTO.
func ItemFromDomain(item domain.Item) ItemDto {
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
