package adapter

import "github.com/uptrace/bun"

type listRow struct {
	bun.BaseModel `bun:"table:lists"`
	ID            string `bun:"id,pk"`
	Name          string `bun:"name"`
	CreatedAt     string `bun:"created_at"`
}

type itemRow struct {
	bun.BaseModel `bun:"table:items"`
	ID            string `bun:"id,pk"`
	ListID        string `bun:"list_id"`
	Title         string `bun:"title"`
	State         string `bun:"state"`
	CreatedAt     string `bun:"created_at"`
}

type userRow struct {
	bun.BaseModel `bun:"table:users"`
	ID            string `bun:"id,pk"`
	Name          string `bun:"name"`
}

type spaceRow struct {
	bun.BaseModel `bun:"table:spaces"`
	ID            string `bun:"id,pk"`
	Name          string `bun:"name"`
	UserID        string `bun:"user_id"`
}

type commentRow struct {
	bun.BaseModel `bun:"table:comments"`
	ID            string `bun:"id,pk"`
	ItemID        string `bun:"item_id"`
	UserID        string `bun:"user_id"`
	Body          string `bun:"body"`
	CreatedAt     string `bun:"created_at"`
}
