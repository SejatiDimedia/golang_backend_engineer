<!--
TEMPLATE: ROADMAP.md (project-local)
This is NOT the same as the repo-level 01-roadmap.md. This is the feature
build sequence within a single project — closer to a sprint plan.
-->

# Roadmap: [Project Name]

**Status:** `Planning | In progress | Complete`

This document sequences the build order *within* this project. It exists so implementation order is a decision, not an accident — and so partial progress is legible to a reader who opens this mid-build.

---

## 1. Build phases

| Phase | Scope | Depends on | Status |
|---|---|---|---|
| 1 — Foundation | Project scaffold, config, DB connection, health check | — | |
| 2 — Core domain | [primary entity] CRUD | Phase 1 | |
| 3 — [feature] | | Phase 2 | |
| 4 — [feature] | | Phase 3 | |
| 5 — Hardening | Tests, error handling polish, logging | All above | |
| 6 — Deployment | Docker, deployment docs | Phase 5 | |

Adjust phase count and content per project — this is a starting skeleton, not a fixed number of phases.

## 2. Feature breakdown

For each feature listed in the PRD's functional requirements, note implementation order and any sequencing reason:

| Feature | PRD ref | Build order reason |
|---|---|---|
| | FR-1 | e.g. "must exist before stock-out logic can be tested" |

## 3. Concepts this project is exercising

Pulled from `01-roadmap.md`'s project map — restated here so this document is self-contained for a reader who hasn't read the repo-level roadmap.

- Concept 1
- Concept 2

## 4. Known risks / unknowns at planning time

What part of this project is least understood before starting. This section is most valuable when written honestly *before* the unknown is resolved — compare against `LESSONS-LEARNED.md` at the end to see how the estimate held up.

---

## Changelog

| Date | Change |
|---|---|
| | Initial roadmap |
