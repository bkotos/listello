---
name: create-application-service
description: >-
  Scaffolds and implements application-layer services in api/internal/application/
  (repository ports, service struct, unit tests, mockery). Use when adding or
  extending use cases, application services, repository interfaces, or
  ListService/ItemService methods following Listello hexagonal architecture.
---

# Create Application Service

Guide for adding or extending use cases in `api/internal/application/`. Read this skill before changing application-layer code.

Architecture context: see [README.md](../../../README.md) (Application layer coordinates domain commands, repository ports, and `EventPublisher`). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Scope

**In scope:** repository ports, service struct, constructor, methods, unit tests, mockery config/regeneration.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Adapter implementations (`api/internal/adapter/`)
- Bootstrap wiring (`api/internal/bootstrap/bootstrap.go`)
- HTTP handlers (`api/cmd/server/handlers/`)
- CLI commands (`api/cmd/cli/commands/`)

## Architecture constraints

- Application depends only on `internal/domain` and its own interfaces.
- Never import adapters, HTTP, CLI, or database packages.
- **Write commands:** call domain → `repository.Save` → `eventPublisher.Publish(event)`.
- **Reads:** delegate directly to repository; no domain call, no event publish.

## Decision tree

1. **New aggregate?** → Create `{aggregate}_service.go` + `{aggregate}_service_test.go`.
2. **Existing service, new command?** → Add method to existing service file; add tests.
3. **Existing service, new read?** → Add repository method (if needed), service method, test.
4. **Domain command missing?** → Stop. Domain must exist first (Gherkin + godog in `internal/domain`).

## Scaffold checklist

```
Task progress:
- [ ] Read domain API (command signature, event type, metadata)
- [ ] Decide: new service vs extend existing
- [ ] Add/update repository port interface (same file as service)
- [ ] Add service struct + constructor (if new aggregate)
- [ ] Add method stub returning `fmt.Errorf("not implemented")` for red phase
- [ ] Write failing unit test(s) in application_test package
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Implement minimum production code
- [ ] Update api/.mockery.yml if port interface changed; run make -C api mocks
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Service file | `{aggregate}_service.go` |
| Test file | `{aggregate}_service_test.go` |
| Test package | `application_test` |
| Repository port | `{Aggregate}Repository` (same file as service) |
| Service type | `{Aggregate}Service` |
| Constructor | `New{Aggregate}Service(repo, eventPublisher)` |
| Test name | `Test{Service}_{Method}_{Behavior}` |
| Domain import | `domain "github.com/bkotos/listello/internal/domain"` |

Shared port: `EventPublisher` lives in `event_publisher.go`.

## Code templates

### New aggregate service

```go
package application

import (
	domain "github.com/bkotos/listello/internal/domain"
)

// {Aggregate}Repository persists {aggregates}.
type {Aggregate}Repository interface {
	Save(/* args */) error
}

// {Aggregate}Service coordinates {aggregate} commands and persistence.
type {Aggregate}Service struct {
	{aggregate}Repository {Aggregate}Repository
	eventPublisher        EventPublisher
}

// New{Aggregate}Service returns a {Aggregate}Service backed by the given repository and publisher.
func New{Aggregate}Service({aggregate}Repository {Aggregate}Repository, eventPublisher EventPublisher) *{Aggregate}Service {
	return &{Aggregate}Service{
		{aggregate}Repository: {aggregate}Repository,
		eventPublisher:        eventPublisher,
	}
}
```

### Write command method

```go
// {MethodName} {description} via the domain and persists it.
func (s *{Aggregate}Service) {MethodName}(/* args */) (domain.{Aggregate}, error) {
	aggregate, event, err := domain.{Command}(/* args */)
	if err != nil {
		return domain.{Aggregate}{}, err
	}
	if err := s.{aggregate}Repository.Save(/* args */); err != nil {
		return domain.{Aggregate}{}, err
	}
	if err := s.eventPublisher.Publish(event); err != nil {
		return domain.{Aggregate}{}, err
	}
	return aggregate, nil
}
```

### Read method

```go
// {MethodName} returns {what} from persistence.
func (s *{Aggregate}Service) {MethodName}(/* args */) (/* return type */, error) {
	return s.{aggregate}Repository.{MethodName}(/* args */)
}
```

### Red-phase stub

When the domain exists but implementation is deferred:

```go
return domain.{Aggregate}{}, fmt.Errorf("not implemented")
```

## Test templates

Use `application_test` package with testify `require` + `assert`. Mocks come from `mocks_test.go` (mockery-generated).

### Write command — persistence

```go
func Test{Service}_{Method}_Persists{Aggregate}(t *testing.T) {
	// Arrange
	repo := NewMock{Aggregate}Repository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.New{Aggregate}Service(repo, publisher)

	repo.EXPECT().Save(/* matcher */).Return(nil)
	publisher.EXPECT().Publish(mock.AnythingOfType("domain.Event")).Return(nil)

	// Act
	result, err := svc.{Method}(/* args */)

	// Assert
	require.NoError(t, err)
	// assert on result fields
}
```

### Write command — event publish

Separate test asserting event name, metadata type, and key fields:

```go
func Test{Service}_{Method}_PublishesEvent(t *testing.T) {
	// Arrange
	var published domain.Event
	// ...
	publisher.EXPECT().
		Publish(mock.MatchedBy(func(event domain.Event) bool {
			published = event
			// check event.Name and metadata type
			return true
		})).
		Return(nil)

	// Act + Assert on published.Metadata
}
```

### Read method

```go
func Test{Service}_{Method}_Returns{Aggregate}FromRepository(t *testing.T) {
	// Arrange
	expected := /* domain value */
	repo := NewMock{Aggregate}Repository(t)
	publisher := NewMockEventPublisher(t)
	svc := application.New{Aggregate}Service(repo, publisher)

	repo.EXPECT().{Method}(/* args */).Return(expected, nil)

	// Act
	received, err := svc.{Method}(/* args */)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, received)
}
```

One behavioral concern per test. For write commands, test persistence **and** event publish in separate tests.

## TDD workflow

Follow `.cursor/rules/tdd.mdc` for application layer:

1. **Red** — Write failing test first. Production may be a `not implemented` stub. Run tests; confirm failure is missing behavior, not a broken harness.
2. **Stop** — Summarize the new/changed spec and the failure. **Do not implement** until the user explicitly approves.
3. **Green** — Implement the bare minimum to pass. No extra validation or edge cases the test does not assert. Re-run tests.

Do not write production implementation in the same turn as a new failing spec.

## Mockery

When a port interface is added or changed:

1. Add the interface under `packages.github.com/bkotos/listello/internal/application.interfaces` in `api/.mockery.yml`.
2. Run `make -C api mocks` (requires mockery v3).
3. Use generated mocks from `mocks_test.go` — **never hand-edit mocks**.

Existing ports: `ListRepository`, `ItemRepository`, `EventPublisher`.

## Verification

```bash
# All API tests
make -C api test

# Application layer only
cd api && go test ./internal/application/...
```

## Further reading

Annotated walkthroughs of `ListService` and `ItemService`: [examples.md](examples.md)
