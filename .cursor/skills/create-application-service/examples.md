# Application Service Examples

Annotated references from the Listello codebase. Read when implementing a new or extended service.

## 1. Service interface + implementation: `ListService`

**File:** `api/internal/application/list_service.go`

Each aggregate exposes an exported **interface** and an unexported **struct**. The constructor returns the interface; a compile-time check ensures the struct satisfies it.

```go
// ListService defines list application operations.
type ListService interface {
	CreateList(name string) (domain.List, error)
	GetAll() ([]domain.List, error)
	GetByID(id string) (domain.List, error)
}

type listService struct {
	listRepository ListRepository
	eventPublisher EventPublisher
}

var _ ListService = (*listService)(nil)

func NewListService(listRepository ListRepository, eventPublisher EventPublisher) ListService {
	return &listService{
		listRepository: listRepository,
		eventPublisher: eventPublisher,
	}
}
```

Key points:

- Handlers and CLI depend on `ListService` (interface), not `*listService`.
- `var _ ListService = (*listService)(nil)` catches drift if methods are added to the interface but not the struct.
- Repository ports stay in the same file as the service.

### Write command: `CreateList`

```go
func (s *listService) CreateList(name string) (domain.List, error) {
	list, event, err := domain.CreateList(name)
	if err != nil {
		return domain.List{}, err
	}
	if err := s.listRepository.Save(list); err != nil {
		return domain.List{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.List{}, err
	}
	return list, nil
}
```

- Domain owns validation and event construction; application only orchestrates.
- Errors from domain, repository, or publisher propagate unchanged.

### Test: persistence (`TestListService_CreateList_PersistsList`)

**File:** `api/internal/application/list_service_test.go`

- Mocks: `NewMockListRepository(t)`, `NewMockEventPublisher(t)` from `mocks_test.go`.
- `Save` expectation uses `mock.MatchedBy` to verify a non-empty item was passed (do not re-assert domain field values).
- `Publish` expectation accepts any `domain.Event` (persistence is the focus).

### Test: event publish (`TestListService_CreateList_PublishesEvent`)

- Captures published event via closure in `mock.MatchedBy`.
- Asserts `event.Name`, metadata type, and non-empty timestamp.
- Does not re-test domain field values on the returned aggregate.

**Pattern:** Split persistence and event assertions into separate tests. Application tests mock **repository ports**, not service interfaces.

## 2. Read methods: `GetAll` and `GetByID`

```go
func (s *listService) GetAll() ([]domain.List, error) {
	return s.listRepository.GetAll()
}

func (s *listService) GetByID(id string) (domain.List, error) {
	return s.listRepository.GetByID(id)
}
```

- No domain call, no event publish.
- One-line delegation to repository.

## 3. `ItemService` with cross-aggregate dependency

**File:** `api/internal/application/item_service.go`

```go
// ItemService defines item application operations.
type ItemService interface {
	DefineItem(listID, title string) (domain.Item, error)
}

type itemService struct {
	listRepository ListRepository
	itemRepository ItemRepository
	eventPublisher EventPublisher
}

var _ ItemService = (*itemService)(nil)

func NewItemService(listRepository ListRepository, itemRepository ItemRepository, eventPublisher EventPublisher) ItemService {
	return &itemService{ ... }
}

func (s *itemService) DefineItem(listID, title string) (domain.Item, error) {
	list, err := s.listRepository.GetByID(listID)
	// domain.DefineItem → Save(item) → Publish(event)
}
```

Key points:

- `ItemRepository.Save(item domain.Item)` — no separate `listID` param; `listID` is on `domain.Item`.
- `itemService` depends on `ListRepository` to load the list before calling domain.
- Application tests mock repository ports; handler tests mock `ItemService` (see `create-api-handler` skill).

## 4. Shared `EventPublisher` port

**File:** `api/internal/application/event_publisher.go`

```go
type EventPublisher interface {
	Publish(event domain.Event) error
}
```

All command services take `EventPublisher` as a constructor dependency. Do not create per-service publisher interfaces.

## 5. Mockery config

**File:** `api/.mockery.yml`

```yaml
packages:
  github.com/bkotos/listello/internal/application:
    interfaces:
      ListRepository:
      EventPublisher:
      ItemRepository:
      ListService:
        config:
          dir: '{{.InterfaceDir}}/mocks'
          filename: mocks.go
          pkgname: mocks
      ItemService:
        config:
          dir: '{{.InterfaceDir}}/mocks'
          filename: mocks.go
          pkgname: mocks
```

- **Repository mocks** → `mocks_test.go` (default config) for `application_test` package.
- **Service mocks** → `mocks/mocks.go` importable by handler tests as `appmocks "github.com/bkotos/listello/internal/application/mocks"`.

After adding a new interface, register it in `.mockery.yml` and run `make -C api mocks`.
