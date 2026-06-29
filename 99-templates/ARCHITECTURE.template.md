<!--
TEMPLATE: ARCHITECTURE.md
-->

# Architecture: [Project Name]

**Status:** `Draft | Approved | Implemented`
**Last updated:** [date]

---

## 1. Architectural style

State the style and why: e.g. "Clean Architecture with handler → service → repository layering. Chosen because [reason], see ADR-00X for full tradeoff discussion."

## 2. System diagram

```
[Client]
   ↓
[HTTP Handler / Router]
   ↓
[Service Layer — business logic]
   ↓
[Repository Layer — data access]
   ↓
[PostgreSQL / Redis / etc.]
```

Adjust to the actual project. Include external dependencies (queues, object storage, third-party APIs) in the diagram, not just the internal layers.

## 3. Folder structure

```
project-name/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   └── middleware/
├── pkg/
├── migrations/
├── config/
└── docker-compose.yml
```

Adjust per project; note any deviation from previous projects' structure and why (link an ADR if the deviation is significant — folder structure drift across 8 projects should be intentional, not accidental).

## 4. Component responsibilities

| Component | Responsibility | Does NOT do |
|---|---|---|
| Handler | HTTP concerns only — parsing, status codes, response shaping | Business logic, direct DB access |
| Service | Business logic, orchestration | HTTP concerns, raw SQL |
| Repository | Data access abstraction | Business rules |

## 5. Data flow — [key operation, e.g. "Create Booking"]

Walk through one representative request end to end. Pick the operation that best exercises the system's hardest problem (e.g. for Project 4, walk through a transfer; for Project 6, walk through a queued notification with retry).

1. Request arrives at handler
2. ...

## 6. Cross-cutting concerns

| Concern | Approach |
|---|---|
| Logging | |
| Error handling | |
| Configuration | |
| Authentication (if applicable) | |
| Context propagation | |

## 7. Dependencies on other projects in this repo

If this project consumes or will later consume another project's service (e.g. Project 8 consuming Project 7's Auth Service), state it here, even if not yet implemented.

## 8. Known architectural limitations

Honest list of what's structurally weak in this design at this project's scope, with a pointer to `FUTURE-IMPROVEMENTS.md` for each.

---

## Changelog

| Date | Change |
|---|---|
| | Initial architecture |
