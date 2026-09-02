---
name: create-api-handler
description: >-
  Scaffolds and wires HTTP API handlers and routes in api/cmd/server/ that call
  application services. Use when adding or extending API endpoints, HTTP
  handlers, route registration, request/response view DTOs, or handlers for
  ListService/ItemService following Listello hexagonal architecture.
---

# Create API Handler

Guide for adding HTTP endpoints in `api/cmd/server/`. Read this skill before changing handlers or routes.

Architecture context: see [README.md](../../../README.md) (handlers call application services, map domain results to view DTOs, encode JSON). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Scope

**In scope:** handler function, handler test (`httptest`), route in `server.go`, request + response DTOs in `view-dtos/`, DTO unit tests, `make api-types` after DTO changes.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Application service methods (`api/internal/application/`) — use [create-application-service](../create-application-service/SKILL.md)
- Adapter repositories (`api/internal/adapter/`) — use [create-adapter-repository](../create-adapter-repository/SKILL.md)
- Bootstrap / `main.go` service construction (`api/internal/bootstrap/`, `api/cmd/server/main.go`)
- UI client (`ui/src/lib/api/`)
- Bruno collection (`bruno/`)

## Architecture constraints

- Handlers depend on application services — never call domain or adapters directly.
- **All JSON request and response bodies use named DTOs** in `internal/view-dtos/` — never anonymous inline structs in handlers.
- Decode into a `{Action}{Resource}Request` DTO; map domain results to `{Resource}Response` via `{Resource}FromDomain`.
- Validate required request fields in the handler (empty string, missing body) before calling the service.
- Use `response.WriteJSON` and `response.WriteError` from `cmd/server/response`.
- One handler file per endpoint (`create_list.go`, `get_list.go`, …).
- Handler factory accepts the application service: `func CreateList(lists *application.ListService) http.HandlerFunc`.

## Preconditions

1. **Application service method exists** — the handler calls an existing service API.
2. **Request and response shapes defined** — add or reuse DTOs in `internal/view-dtos/`.

If the service method is missing, stop and use `create-application-service` first.

## Decision tree

1. **New endpoint?** → Add request/response DTOs → handler + test → register route in `server.go`.
2. **POST/PUT/PATCH with body?** → Add `{Action}{Resource}Request` in `view-dtos/`; decode into it in the handler.
3. **JSON response?** → Add or reuse `{Resource}Response` + `{Resource}FromDomain`; run `make api-types`.
4. **New service dependency?** → Add parameter to `newAPIServer` in `server.go`; mention `main.go` bootstrap follow-up.
5. **Path parameter?** → Use `r.PathValue("name")` (Go 1.22+ route patterns); set `req.SetPathValue` in tests.

## Scaffold checklist

```
Task progress:
- [ ] Read application service method signature
- [ ] Decide HTTP method, path, request body fields, response fields
- [ ] Add request DTO in internal/view-dtos/ (POST/PUT/PATCH)
- [ ] Add response DTO + FromDomain mapper in internal/view-dtos/
- [ ] Add DTO unit tests in view-dtos/{resource}_test.go
- [ ] Write failing handler test with httptest
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Implement handler (decode request DTO, call service, write response DTO)
- [ ] Register route in server.go
- [ ] Update newAPIServer signature if new service dependency
- [ ] Run make api-types
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Handler file | `{action}_{resource}.go` (e.g. `create_list.go`, `get_list.go`) |
| Test file | `{action}_{resource}_test.go` |
| Test package | `handlers` (same package) |
| Handler function | `{Action}{Resource}` or `{Action}` (e.g. `CreateList`, `GetAllLists`) |
| Test name | `Test{Handler}` or `Test{Handler}_{Behavior}` |
| Route pattern | `{METHOD} /api/{resource}` or `{METHOD} /api/{resource}/{param}` |
| Request DTO | `{Action}{Resource}Request` (e.g. `CreateListRequest`, `DefineItemRequest`) |
| Response DTO | `{Resource}Response` (e.g. `ListResponse`, `ItemResponse`) |
| Response mapper | `{Resource}FromDomain`, `{Plural}FromDomain` for slices |
| DTO test file | `view-dtos/{resource}_test.go` |
| Stubs | `stubListRepository`, `stubEventPublisher` in `stubs_test.go` |

## Code templates

### Handler factory (POST with JSON body)

```go
package handlers

import (
	"encoding/json"
	"net/http"

	application "github.com/bkotos/listello/internal/application"
	viewdto "github.com/bkotos/listello/internal/view-dtos"

	"github.com/bkotos/listello/cmd/server/response"
)

func {Action}(svc *application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewdto.{Action}{Resource}Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.{Field} == "" {
			response.WriteError(w, http.StatusBadRequest, "{field} is required")
			return
		}

		result, err := svc.{Method}(req.{Field} /* map other fields */)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusCreated, viewdto.{Resource}FromDomain(result))
	}
}
```

### Handler factory (GET collection)

```go
func GetAll{Resource}(svc *application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := svc.GetAll()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteJSON(w, http.StatusOK, viewdto.{Plural}FromDomain(all))
	}
}
```

### Handler factory (GET by path param)

```go
import "strings"

func Get{Resource}(svc *application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		result, err := svc.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				response.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.{FromDomain}(result))
	}
}
```

### Route registration (`server.go`)

```go
func newAPIServer(lists *application.ListService /*, items *application.ItemService */) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(lists))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(lists))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(lists))
	// mux.HandleFunc("POST /api/lists/{id}/items", handlers.DefineItem(items))
	return mux
}
```

### View DTOs (`internal/view-dtos/{resource}.go`)

Define request and response types in the same file per resource. Tygo generates TypeScript for all exported structs.

**Request DTO** (endpoints with a JSON body):

```go
// CreateListRequest is the HTTP request body for creating a list.
type CreateListRequest struct {
	Name string `json:"name"`
}
```

**Response DTO** + domain mapper:

```go
// ListResponse is the HTTP representation of a list.
type ListResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

// ListFromDomain maps a domain list to its response DTO.
func ListFromDomain(list domain.List) ListResponse {
	return ListResponse{ID: list.ID, Name: list.Name}
}

// ListsFromDomain maps domain lists to response DTOs.
func ListsFromDomain(lists []domain.List) []ListResponse { ... }
```

GET endpoints with no body still need a response DTO. Path/query params are read in the handler and passed to the service — only JSON bodies get request DTOs.

### DTO unit tests (`view-dtos/{resource}_test.go`)

Test mappers in `viewdto_test` package:

```go
func TestListFromDomain(t *testing.T) {
	list := domain.List{ID: "LS_1", Name: "Work"}
	received := viewdto.ListFromDomain(list)
	assert.Equal(t, list.ID, received.ID)
	assert.Equal(t, list.Name, received.Name)
}
```

Request DTOs are plain structs — test them indirectly via handler tests unless they contain mapping logic.

After any DTO change: `make api-types` (regenerates `api-types/index.ts` via tygo).

## Test templates

Use `httptest` with stub repositories from `stubs_test.go`. Inject a real `application.New{Service}(stubRepo, stubPublisher)`.

### POST handler

```go
func TestCreateList(t *testing.T) {
	// Arrange
	svc := application.NewListService(&stubListRepository{}, &stubEventPublisher{})
	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()

	// Act
	CreateList(svc)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var received viewdto.ListResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, "Next actions", received.Name)
}
```

### GET with path param

```go
req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID, nil)
req.SetPathValue("id", listID)
```

### Not found

```go
func TestGetList_NotFound(t *testing.T) {
	svc := application.NewListService(&stubListRepository{
		getByIDFn: func(id string) (domain.List, error) {
			return domain.List{}, fmt.Errorf("list %q not found", id)
		},
	}, &stubEventPublisher{})
	// ...
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
```

One behavioral concern per test. Add stubs to `stubs_test.go` when a new repository port is needed.

## HTTP status conventions

| Situation | Status |
|-----------|--------|
| Invalid JSON / missing required field | 400 Bad Request |
| Domain/application validation error (writes) | 400 Bad Request |
| Resource not found (`"not found"` in error) | 404 Not Found |
| Unexpected persistence error | 500 Internal Server Error |
| Successful create | 201 Created |
| Successful read | 200 OK |

Match existing handlers; do not add new status-mapping logic the tests do not require.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing handler test first. Handler may not exist yet (compile error) or return wrong status. Confirm failure is missing behavior.
2. **Stop** — Summarize spec and failure. **Do not implement** until the user approves.
3. **Green** — Implement handler + route registration. Re-run tests.

Do not write handler implementation in the same turn as a new failing spec.

## Verification

```bash
# Handler tests only
cd api && go test ./cmd/server/handlers/...

# Full API tests
make -C api test

# Regenerate TypeScript types after DTO changes (required)
make api-types
```

## Further reading

Annotated walkthroughs of list handlers: [examples.md](examples.md)
