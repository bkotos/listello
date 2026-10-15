package viewdto

import viewmodel "github.com/bkotos/listello/internal/personal-productivity-context/view-models"

// ItemCommentDto is the HTTP representation of a comment on an item.
type ItemCommentDto struct {
	ID        string `json:"ID"`
	ItemID    string `json:"ItemID"`
	UserID    string `json:"UserID"`
	UserName  string `json:"UserName"`
	Body      string `json:"Body"`
	CreatedAt string `json:"CreatedAt"`
}

// ItemCommentFromView maps an item comment view model to its response DTO.
func ItemCommentFromView(comment viewmodel.ItemComment) ItemCommentDto {
	return ItemCommentDto{
		ID:        comment.ID,
		ItemID:    comment.ItemID,
		UserID:    comment.UserID,
		UserName:  comment.UserName,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt,
	}
}

// ItemCommentsFromView maps item comment view models to response DTOs.
func ItemCommentsFromView(comments []viewmodel.ItemComment) []ItemCommentDto {
	response := make([]ItemCommentDto, len(comments))
	for i, comment := range comments {
		response[i] = ItemCommentFromView(comment)
	}
	return response
}
