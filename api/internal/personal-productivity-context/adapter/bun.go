package adapter

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/bkotos/listello/internal/sqlite"
)

func openBun(workspace *sqlite.WorkspaceDB) (*bun.DB, error) {
	sqldb, err := workspace.DB()
	if err != nil {
		return nil, err
	}
	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}
