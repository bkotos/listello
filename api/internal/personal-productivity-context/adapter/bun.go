package adapter

import (
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	"github.com/bkotos/listello/internal/sqlite"
)

func openBun(workspace *sqlite.WorkspaceDB) (*bun.DB, error) {
	engine, err := workspace.Engine()
	if err != nil {
		return nil, err
	}
	sqldb, err := workspace.DB()
	if err != nil {
		return nil, err
	}
	switch engine {
	case sqlite.EngineSQLite:
		return bun.NewDB(sqldb, sqlitedialect.New()), nil
	case sqlite.EnginePostgres:
		return bun.NewDB(sqldb, pgdialect.New()), nil
	default:
		return nil, fmt.Errorf("unsupported engine %q", engine)
	}
}
