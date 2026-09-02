# API Handler Examples

Annotated references from the Listello codebase. Read when implementing a new HTTP endpoint.

## 0. Request and response DTOs (required for new endpoints)

All JSON bodies use named types in `internal/view-dtos/`. Do not decode into anonymous structs in handlers (legacy `CreateList` still uses an inline struct — migrate when touching it).

**File:** `api/internal/view-dtos/item.go`

```go
type DefineItemRequest struct {
	Title string `json:"title"`
}

type ItemDto struct {
	ID    string `json:"ID"`
	Title string `json:"Title"`
	// ...
}

func ItemFromDomain(item domain.Item) ItemDto { ... }
```

**Target handler pattern:**

```go
var req viewdto.DefineItemRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil { ... }

item, err := itemService.DefineItem(id, req.Title)
response.WriteJSON(w, http.StatusCreated, viewdto.ItemFromDomain(item))
```

Do **not** add handler-level validation for fields the domain already validates (e.g. empty `title` on `DefineItem`). Let the service return an error and map it to 400.

**DTO tests** in `view-dtos/{resource}_test.go` cover `FromDomain` mappers. Run `make api-types` after adding or changing DTOs.

## 1. POST `DefineItem` — write endpoint with path param

**File:** `api/cmd/server/handlers/define_item.go`

```go
func DefineItem(itemService application.ItemService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteError(w, http.StatusBadRequest, "id is required")
			return
		}

		var req viewdto.DefineItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		item, err := itemService.DefineItem(id, req.Title)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusCreated, viewdto.ItemFromDomain(item))
	}
}
```

Key points:

- Factory accepts `application.ItemService` **interface**, not a concrete struct or pointer.
- Parameter named `itemService` (not `items` — collides with `/items` route).
- No `title is required` check — domain owns that validation.
- Return 201 Created with `viewdto.ItemDto`.

### Test

**File:** `api/cmd/server/handlers/define_item_test.go`

```go
import appmocks "github.com/bkotos/listello/internal/application/mocks"

func TestDefineItem(t *testing.T) {
	const listID = "LS_1"
	expected := domain.Item{ID: "IT_1", ListID: listID, Title: "Buy milk", State: domain.ItemOutstanding}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem(listID, "Buy milk").Return(expected, nil)

	body := bytes.NewBufferString(`{"title":"Buy milk"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists/"+listID+"/items", body)
	req.SetPathValue("id", listID)
	rec := httptest.NewRecorder()

	DefineItem(itemService)(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	// assert response maps from expected
}
```

- Mock the **service interface** — do not wire stub repositories.
- `EXPECT()` asserts the handler called the service with the correct arguments.

## 2. POST `CreateList` — write endpoint (legacy request decode)

**File:** `api/cmd/server/handlers/create_list.go`

```go
func CreateList(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ...
		list, err := listService.CreateList(body.Name)
		response.WriteJSON(w, http.StatusCreated, viewdto.ListFromDomain(list))
	}
}
```

### Test

**File:** `api/cmd/server/handlers/create_list_test.go`

```go
listService := appmocks.NewMockListService(t)
listService.EXPECT().CreateList("Next actions").Return(domain.List{ID: "LS_1", Name: "Next actions"}, nil)
CreateList(listService)(rec, req)
```

## 3. GET `GetAllLists` — collection read

**File:** `api/cmd/server/handlers/get_all_lists.go`

```go
func GetAllLists(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all, err := listService.GetAll()
		// ...
		response.WriteJSON(w, http.StatusOK, viewdto.ListsFromDomain(all))
	}
}
```

### Test

```go
listService := appmocks.NewMockListService(t)
listService.EXPECT().GetAll().Return(expected, nil)
```

## 4. GET `GetList` — resource by ID

**File:** `api/cmd/server/handlers/get_list.go`

```go
func GetList(listService application.ListService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		list, err := listService.GetByID(id)
		// 404 when err contains "not found"
		response.WriteJSON(w, http.StatusOK, viewdto.ListFromDomain(list))
	}
}
```

### Tests

Success:

```go
listService := appmocks.NewMockListService(t)
listService.EXPECT().GetByID(listID).Return(expected, nil)
```

Not found:

```go
listService.EXPECT().GetByID(listID).Return(domain.List{}, fmt.Errorf("list %q not found", listID))
// assert 404
```

## 5. Route registration

**File:** `api/cmd/server/server.go`

```go
func newAPIServer(listService application.ListService, itemService application.ItemService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /api/lists", handlers.GetAllLists(listService))
	mux.HandleFunc("GET /api/lists/{id}", handlers.GetList(listService))
	mux.HandleFunc("POST /api/lists", handlers.CreateList(listService))
	mux.HandleFunc("POST /api/lists/{id}/items", handlers.DefineItem(itemService))
	return mux
}
```

## 6. Service mocks (handler tests)

**File:** `api/internal/application/mocks/mocks.go` (mockery-generated)

Service interfaces (`ListService`, `ItemService`) are mocked here so handler tests can import them:

```go
import appmocks "github.com/bkotos/listello/internal/application/mocks"
```

Repository port mocks remain in `mocks_test.go` for **application-layer** tests only.

When adding a new `{Aggregate}Service` interface, register it in `api/.mockery.yml` with `config` pointing to `mocks/` and run `make -C api mocks`.

## 7. Response helpers

**File:** `api/cmd/server/response/json.go`

```go
func WriteJSON(w http.ResponseWriter, status int, payload any)
func WriteError(w http.ResponseWriter, status int, message string) // {"error": "..."}
```

## 8. View DTO naming

| Type | Purpose | Example |
|------|---------|---------|
| `{Action}{Resource}Request` | JSON request body | `DefineItemRequest` |
| `{Resource}Dto` / `{Resource}Response` | JSON response body | `ItemDto`, `ListResponse` |
| `{Resource}FromDomain` | domain → response | `ItemFromDomain(item)` |

After adding or changing any DTO, run `make api-types`.
