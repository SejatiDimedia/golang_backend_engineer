<!--
TEMPLATE: TESTING.md
-->

# Testing Strategy: [Project Name]

---

## 1. Scope of testing for this project

State explicitly what level of testing this project targets and why — this should scale with project difficulty per `00-learning-plan.md`. A Project 1 testing strategy that's intentionally lighter than Project 7's is correct, not a shortcut, *if stated as a deliberate scoping decision here*.

## 2. Test types in use

| Type | Used? | Tooling | Scope |
|---|---|---|---|
| Unit tests | Yes/No | `testing` + `testify`/`mockery` etc. | Service layer logic, isolated from DB |
| Integration tests | Yes/No | | Repository layer against a real (test) DB |
| End-to-end / API tests | Yes/No | | Full request/response cycle |
| Load/performance tests | Yes/No | | Only if relevant to this project's learning goals |

## 3. What is covered

| Component | Coverage approach |
|---|---|
| Handlers | |
| Services | |
| Repositories | |

## 4. What is explicitly NOT covered, and why

Be honest. "Not tested because out of scope for this project's learning goals" is a legitimate answer — but write it down so it reads as a decision, not an omission.

## 5. Test data strategy

How test data is set up/torn down (fixtures, factories, a dedicated test database, transactions rolled back per test, etc.)

## 6. Running tests

```bash
go test ./... -v
go test ./... -cover
```

## 7. CI integration

State whether tests run in CI for this project, and if not yet, when that's planned (likely tied to a later project's deployment maturity).

---

## Changelog

| Date | Change |
|---|---|
| | Initial testing strategy |
