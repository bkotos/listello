# API Handler Examples

Annotated references from the Listello codebase. Read when implementing a new HTTP endpoint.

## 0. Request and response DTOs (required for new endpoints)

All JSON bodies use named types in `internal/view-dtos/`. Do not decode into anonymous structs in handlers.

**File:** `api/internal/view-dtos/list.go`

```go
// CreateListRequest is the HTTP request body for creating a list.
type CreateListRequest struct {
	Name string `json:"name"`
}

// ListResponse is the HTTP representation of a list.
type ListResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

func ListFromDomain(list domain.List) ListResponse { ... }
func ListsFromDomain(lists []domain.List) []ListResponse { ... }
```

**Target handler pattern:**

```go
var req viewdto.CreateListRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil { ... }
if req.Name == "" { ... }

list, err := lists.CreateList(req.Name)
response.WriteJSON(w, http.StatusCreated, viewdto.ListFromDomain(list))
```

**DTO tests** in `view-dtos/list_test.go` cover `FromDomain` mappers. Run `make api-types` after adding or changing DTOs so `api-types/index.ts` stays in sync for the UI.

> **Note:** `CreateList` still uses an inline request struct — legacy. New endpoints must add `{Action}{Resource}Request` in `view-dtos/`.

## 1. POST `CreateList` — write endpoint

**File:** `api/cmd/server/handlers/create_list.go`

```go
func CreateList(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if body.Name == "" {
			response.WriteError(w, http.StatusBadRequest, "name is required")
			return
		}

		list, err := lists.CreateList(body.Name)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusCreated, viewdto.ListFromDomain(list))
	}
}
```

Key points:

- **Legacy:** decodes into an anonymous struct — replace with `viewdto.CreateListRequest` when touching this handler.
- Application errors on writes map to 400.
- Return 201 Created with `viewdto.ListResponse` (not raw domain type).
- Factory closes over `*application.ListService`.

### Test

**File:** `api/cmd/server/handlers/create_list_test.go`

- `application.NewListService(&stubListRepository{}, &stubEventPublisher{})` — real service, stubbed ports.
- POST body as `bytes.NewBufferString(`{"name":"..."}`)`.
- Assert status 201 and decode into `viewdto.ListResponse` (or `map[string]string` until `CreateListRequest`/`ListResponse` are wired).

## 2. GET `GetAllLists` — collection read

**File:** `api/cmd/server/handlers/get_all_lists.go`

```go
func GetAllLists(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := lists.GetAll()
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		response.WriteJSON(w, http.StatusOK, viewdto.ListsFromDomain(all))
	}
}
```

Key points:

- Read endpoints: persistence errors → 500.
- Plural mapper `ListsFromDomain` for slices.

### Test

**File:** `api/cmd/server/handlers/get_all_lists_test.go`

- Stub `getAllFn` on `stubListRepository` to return expected slice.
- Assert 200 and array of objects with matching `ID` / `Name`.

## 3. GET `GetList` — resource by ID

**File:** `api/cmd/server/handlers/get_list.go`

```go
func GetList(lists *application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		list, err := lists.GetByID(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				response.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, viewdto.ListFromDomain(list))
	}
}
```

Key points:

- Path param via `r.PathValue("id")` — matches route `GET /api/lists/{id}`.
- `"not found"` in error message → 404; other errors → 500.

### Tests

**File:** `api/cmd/server/handlers/get_list_test.go`

Success:

```go
req := httptest.NewRequest(http.MethodGet, "/api/lists/"+listID, nil)
req.SetPathValue("id", listID)
```

Not found:

```go
getByIDFn: func(id string) (domain.List, error) {
	return domain.List{}, fmt.Errorf("list %q not found", id)
}
// assert 404 and {"error": "list \"LS_missing\" not found"}
```

## 4. Route registration

**File:** `api/cmd/server/server.go`

```go
func newAPIServer(lists *application.ListService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(lists))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(lists))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(lists))
	return mux
}
```

Key points:

- Go 1.22+ method-aware patterns: `"GET /api/lists/{id}"`.
- Pass application service into handler factory.
- Adding a new service (e.g. `ItemService`) requires a new `newAPIServer` parameter and routes — update `main.go` separately to construct and pass the service.

## 5. Response helpers

**File:** `api/cmd/server/response/json.go`

```go
func WriteJSON(w http.ResponseWriter, status int, payload any)
func WriteError(w http.ResponseWriter, status int, message string) // {"error": "..."}
```

Always use these — do not write JSON directly to `ResponseWriter`.

## 6. View DTOs — request and response

**File:** `api/internal/view-dtos/list.go`

| Type | Purpose | Example |
|------|---------|---------|
| `{Action}{Resource}Request` | JSON request body | `CreateListRequest`, `DefineItemRequest` |
| `{Resource}Response` | JSON response body | `ListResponse`, `ItemResponse` |
| `{Resource}FromDomain` | domain → response | `ListFromDomain(list)` |
| `{Plural}FromDomain` | slice mapper | `ListsFromDomain(lists)` |

```go
type ListResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

func ListFromDomain(list domain.List) ListResponse { ... }
func ListsFromDomain(lists []domain.List) []ListResponse { ... }
```

After adding or changing any DTO, run `make api-types` to regenerate TypeScript types in `api-types/index.ts` (tygo config in `api/tygo.yaml`).

## 7. Test stubs

**File:** `api/cmd/server/handlers/stubs_test.go`

Hand-rolled stubs implement application ports for handler tests (not mockery):

```go
type stubListRepository struct {
	saveFn    func(list domain.List) error
	getAllFn  func() ([]domain.List, error)
	getByIDFn func(id string) (domain.List, error)
}
```

Add `stubItemRepository` here when wiring item endpoints.

## 8. Next endpoint example: define item

Given `ItemService.DefineItem(listID, title string)`:

1. **DTOs** in `view-dtos/item.go`:
   - `DefineItemRequest` with `Title string \`json:"title"\``
   - `ItemResponse` with fields from `domain.Item`
   - `ItemFromDomain(item domain.Item) ItemResponse`
   - Tests in `item_test.go` for `ItemFromDomain`
2. **Route:** `POST /api/lists/{id}/items`
3. **Handler:** `define_item.go` — `r.PathValue("id")`, decode `viewdto.DefineItemRequest`, call `items.DefineItem(id, req.Title)`, return `viewdto.ItemFromDomain(item)`
4. **server.go:** add `items *application.ItemService` param and route
5. **`make api-types`** — exports new request/response types to the UI
6. **Follow-up:** `main.go` bootstrap to construct `ItemService` (out of skill scope unless asked)
