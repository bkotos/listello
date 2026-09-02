# Application Service Examples

Annotated references from the Listello codebase. Read when implementing a new or extended service.

## 1. Write command: `ListService.CreateList`

**File:** `api/internal/application/list_service.go`

Full write path: domain → save → publish.

```go
func (s *ListService) CreateList(name string) (domain.List, error) {
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

Key points:

- Domain owns validation and event construction; application only orchestrates.
- Errors from domain, repository, or publisher propagate unchanged.
- Return the domain aggregate on success.

### Test: persistence (`TestListService_CreateList_PersistsList`)

**File:** `api/internal/application/list_service_test.go`

- Mocks: `NewMockListRepository(t)`, `NewMockEventPublisher(t)`.
- `Save` expectation uses `mock.MatchedBy` to assert aggregate fields (name, ID prefix).
- `Publish` expectation accepts any `domain.Event` (persistence is the focus).
- Asserts returned list name and ID prefix.

### Test: event publish (`TestListService_CreateList_PublishesEvent`)

- Captures published event via closure in `mock.MatchedBy`.
- Asserts `event.Name == domain.EventListCreated`.
- Asserts metadata type `domain.EventMetadataListCreated` and ID matches returned list.
- Asserts `published.Timestamp` is non-empty.

**Pattern:** Split persistence and event assertions into separate tests.

## 2. Read methods: `ListService.GetAll` and `GetByID`

**File:** `api/internal/application/list_service.go`

```go
func (s *ListService) GetAll() ([]domain.List, error) {
	return s.listRepository.GetAll()
}

func (s *ListService) GetByID(id string) (domain.List, error) {
	return s.listRepository.GetByID(id)
}
```

Key points:

- No domain call, no event publish.
- One-line delegation to repository.

### Test: `TestListService_GetAll_ReturnsListsFromRepository`

- Sets up `expected` slice of domain values.
- `repo.EXPECT().GetAll().Return(expected, nil)`.
- Asserts `received` equals `expected`.

### Test: `TestListService_GetByID_ReturnsListFromRepository`

- Same pattern with `GetByID(listID)` and a single `domain.List`.

## 3. Partial scaffold: `ItemService`

**File:** `api/internal/application/item_service.go`

Shows a new aggregate service with repository port and red-phase stub:

```go
// ItemRepository persists items and their list membership.
type ItemRepository interface {
	Save(listID string, item domain.Item) error
}

// ItemService coordinates item aggregate commands and persistence.
type ItemService struct {
	itemRepository ItemRepository
	eventPublisher EventPublisher
}

func NewItemService(itemRepository ItemRepository, eventPublisher EventPublisher) *ItemService {
	return &ItemService{
		itemRepository: itemRepository,
		eventPublisher: eventPublisher,
	}
}

// DefineItem defines an item on a list via the domain and persists it.
func (s *ItemService) DefineItem(listID, title string) (domain.Item, error) {
	return domain.Item{}, fmt.Errorf("not implemented")
}
```

Key points:

- Repository port colocated with service in the same file.
- Constructor takes repository + shared `EventPublisher`.
- Stub returns `not implemented` so tests fail for the right reason during red phase.
- `ItemRepository` is registered in `api/.mockery.yml` for mock generation.

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
```

After adding a new port interface, add it here and run `make -C api mocks`.
