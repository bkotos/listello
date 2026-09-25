package viewmodel

// ItemComment is a comment on a list item for query views.
type ItemComment struct {
	ID        string
	ItemID    string
	UserID    string
	UserName  string
	Body      string
	CreatedAt string
}
