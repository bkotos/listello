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

Architecture context: see [README.md](../../../README.md) (handlers call application services, map domain results to view DTOs, encode JSON). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Upstream dependencies

**Stop and do not proceed** until the application service method exists and is implemented. See [LAYER-ORDER.md](../LAYER-ORDER.md).

Before starting, verify:

| Check | How |
|-------|-----|
| Service method exists | `{Method}(...)` on `{Aggregate}Service` interface in `internal/application/` |
| Method is implemented | Not `return ..., fmt.Errorf("not implemented")` |
| Application tests pass | `go test ./internal/application/...` green for that method |

If any check fails → **stop**. Tell the user to use `create-application-service` first and get green tests. Do not write handler tests, DTOs, or routes.

Adapter repository is **not** required for this skill (handler tests use mock service interfaces).

## Scope

**In scope:** handler function, handler test (`httptest`), route in `server.go`, request + response DTOs in `view-dtos/`, DTO unit tests, `make api-types` after DTO changes, Bruno request in `bruno/api/`.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Application service methods (`api/internal/application/`) — use [create-application-service](../create-application-service/SKILL.md)
- Adapter repositories (`api/internal/adapter/`) — use [create-adapter-repository](../create-adapter-repository/SKILL.md)
- Bootstrap / `main.go` service construction (`api/internal/bootstrap/`, `api/cmd/server/main.go`)
- UI client (`ui/src/lib/api/`)

## Architecture constraints

- Handlers depend on application **service interfaces** — never call domain, adapters, or concrete service structs directly.
- **All JSON request and response bodies use named DTOs** in `internal/view-dtos/` — never anonymous inline structs in handlers (legacy handlers may still use inline structs; migrate when touching them).
- Decode into a `{Action}{Resource}Request` DTO; map domain results to response DTOs via `{Resource}FromDomain`.
- **Do not validate domain-owned fields in handlers** (e.g. empty `title` on `DefineItem`) — domain/application return errors; map those to 400.
- Validate only HTTP/transport concerns in handlers (invalid JSON, missing path param).
- Use `response.WriteJSON` and `response.WriteError` from `cmd/server/response`.
- One handler file per endpoint (`create_list.go`, `get_list.go`, …).
- Handler factory accepts the service **interface**: `func CreateList(listService application.ListService) http.HandlerFunc`.

## Do not duplicate domain logic

Domain rules live in `internal/domain` (Gherkin + godog). Handlers are thin transport adapters — do not re-implement or re-test domain behavior here.

**In handler code:**

- Do not validate fields the domain already validates (e.g. no `if req.Title == ""` when domain rejects empty titles).
- Decode JSON and read path params; pass values to the service; map errors to HTTP status.
- Propagate service/domain errors to the client — do not add handler-only business rules.

**In handler tests:**

- Mock the **service interface**; assert `EXPECT().{Method}(...)` was called with the correct arguments.
- Assert HTTP concerns: status code, response shape/mapping from the value the mock returns.
- Do **not** add tests for domain validation failures (e.g. define on inbox, empty title) — those belong in `internal/domain`.
- Do **not** re-assert domain field semantics on the returned aggregate (outstanding state, ID prefixes, etc.) — only verify the handler forwards the service result into the response DTO.

HTTP-specific behavior (e.g. mapping `"not found"` in an error to 404) is fair game — that is transport logic, not domain logic.

## Preconditions

These duplicate the upstream gate — all must pass:

1. **Application service method exists and is implemented** — the handler calls a real service API, not a stub.
2. **Request and response shapes** — you will add DTOs in `internal/view-dtos/` as part of this skill.

If the service method is missing or still `not implemented` → **stop** and use `create-application-service` first.

## Decision tree

1. **New endpoint?** → Add request/response DTOs → handler + test → register route in `server.go`.
2. **POST/PUT/PATCH with body?** → Add `{Action}{Resource}Request` in `view-dtos/`; decode into it in the handler.
3. **JSON response?** → Add or reuse response DTO (`{Resource}Dto` or `{Resource}Response`) + `{Resource}FromDomain`; run `make api-types`.
4. **New service dependency?** → Add parameter to `newAPIServer` in `server.go`; mention `main.go` bootstrap follow-up.
5. **Path parameter?** → Use `r.PathValue("name")` (Go 1.22+ route patterns); set `req.SetPathValue` in tests.
6. **Handler green?** → Add matching Bruno request in `bruno/api/` for manual testing.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream: service method implemented, application tests green (stop if not — see LAYER-ORDER.md)
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
- [ ] Register new `{Aggregate}Service` in api/.mockery.yml (mocks/ output) if new aggregate; run make -C api mocks
- [ ] Run make api-types
- [ ] Add Bruno request in bruno/api/{Action} {Resource}.bru (on green — manual smoke test)
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
| Response DTO | `{Resource}Dto` or `{Resource}Response` (e.g. `ListResponse`, `ItemDto`) |
| Response mapper | `{Resource}FromDomain`, `{Plural}FromDomain` for slices |
| DTO test file | `view-dtos/{resource}_test.go` |
| Service mocks | `appmocks "github.com/bkotos/listello/internal/application/mocks"` (mockery-generated) |
| Service handler param | `{aggregate}Service` (e.g. `listService`, `itemService`) — type is `application.{Aggregate}Service` interface |
| Bruno request | `bruno/api/{Action} {Resource}.bru` (e.g. `Define Item.bru`) |

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

func {Action}({aggregate}Service application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req viewdto.{Action}{Resource}Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		result, err := {aggregate}Service.{Method}(/* map fields from req */)
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
func GetAll{Resource}({aggregate}Service application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := {aggregate}Service.GetAll()
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

func Get{Resource}({aggregate}Service application.{Service}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		result, err := {aggregate}Service.GetByID(id)
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
func newAPIServer(listService application.ListService, itemService application.ItemService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(listService))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(listService))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(listService))
	// mux.HandleFunc("POST /api/lists/{id}/items", handlers.DefineItem(itemService))
	return mux
}
```

### Bruno request (`bruno/api/{Action} {Resource}.bru`)

Add a Bruno request when the handler goes green so the endpoint can be exercised manually against a running server. Use the `Local` environment variable `apiBase` (`http://localhost:8080`).

**GET (no body):**

```bru
meta {
  name: Get {Resource}
  type: http
  seq: 4
}

get {
  url: {{apiBase}}/api/{resources}/LS_1
}
```

**POST with JSON body:**

```bru
meta {
  name: Define Item
  type: http
  seq: 5
}

post {
  url: {{apiBase}}/api/lists/LS_1/items
  body: json
}

headers {
  content-type: application/json
}

body:json {
  {
    "title": "Buy milk"
  }
}
```

Key points:

- File lives in `bruno/api/` (same collection as other API requests).
- `meta.name` matches the action in plain language (e.g. `Define Item`, `Create List`).
- Use `{{apiBase}}` — do not hardcode the host.
- JSON field names match the Go request DTO `json` tags (e.g. `"title"` for `DefineItemRequest`).
- Use a sample list/item ID in the path (e.g. `LS_1`); create prerequisite data first if needed.
- Increment `seq` so requests stay ordered in the Bruno UI.

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

Use `httptest` with **mock service interfaces** from `internal/application/mocks`. Assert the service method was called with the correct arguments via `EXPECT()`; optionally assert HTTP status and response mapping from the value returned by the mock.

### POST handler

```go
func TestCreateList(t *testing.T) {
	// Arrange
	const listName = "Next actions"
	expected := domain.List{ID: "LS_1", Name: listName}
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().CreateList(listName).Return(expected, nil)

	body := bytes.NewBufferString(`{"name":"Next actions"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists", body)
	rec := httptest.NewRecorder()

	// Act
	CreateList(listService)(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var received viewdto.ListResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&received))
	assert.Equal(t, expected.Name, received.Name)
	assert.Equal(t, expected.ID, received.ID)
}
```

### GET with path param

```go
listService := appmocks.NewMockListService(t)
listService.EXPECT().GetByID(listID).Return(expected, nil)

req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID, nil)
req.SetPathValue("id", listID)
```

### Not found

```go
func TestGetList_NotFound(t *testing.T) {
	listService := appmocks.NewMockListService(t)
	listService.EXPECT().GetByID(listID).Return(domain.List{}, fmt.Errorf("list %q not found", listID))
	// ...
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
```

One behavioral concern per test. Do **not** stub repository ports in handler tests — mock the service interface instead. Do **not** add handler tests for domain validation failures.

## HTTP status conventions

| Situation | Status |
|-----------|--------|
| Invalid JSON / missing path param | 400 Bad Request |
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
