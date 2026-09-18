package sqlite

import "fmt"

// Engine is the workspace database engine.
type Engine string

const (
	EngineSQLite   Engine = "sqlite"
	EnginePostgres Engine = "postgres"
)

// ParseEngine returns a supported engine name. Empty defaults to SQLite.
func ParseEngine(s string) (Engine, error) {
	switch Engine(s) {
	case "", EngineSQLite:
		return EngineSQLite, nil
	case EnginePostgres:
		return EnginePostgres, nil
	default:
		return "", fmt.Errorf("unsupported engine %q", s)
	}
}
