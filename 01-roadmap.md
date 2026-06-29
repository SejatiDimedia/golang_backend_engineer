# Roadmap

**Status:** Active
**Last updated:** 2026-06-29

This is the master index for the entire repository: the topic sequence, the project map, the milestones, and the timeline. If `00-learning-plan.md` is the *how*, this document is the *what* and *when*.

---

## 1. Phase overview

| Phase | Projects | Theme |
|---|---|---|
| Phase 0 | — | Go fundamentals refresh (no project yet) |
| Phase 1 — Foundations | 1, 2 | CRUD, REST, relational data, repository pattern |
| Phase 2 — Core Backend Skills | 3, 4 | Auth, scheduling, concurrency, financial consistency |
| Phase 3 — Infrastructure Depth | 5, 6 | Object storage, streaming, queues, background processing |
| Phase 4 — Platform Engineering | 7, 8 | Reusable auth service, RBAC, modular multi-tenant architecture |
| Phase 5 — Capstone | — | Composed ecosystem of all services behind a gateway |

Each phase assumes the previous phase's concepts are load-bearing, not optional review. Project 4 assumes you can already build a REST CRUD service without re-deriving it (Projects 1–2); Project 7's auth service assumes you've already built ad hoc auth once in Project 3 and felt its limitations firsthand.

---

## 2. Phase 0 — Fundamentals refresh (pre-project)

Not a project; a deliberate runway before Project 1 starts. Tracked in `02-notes/` and `03-exercises/`.

| Topic | Notes location | Why it's gated before Project 1 |
|---|---|---|
| Syntax, types, control flow | `02-notes/01-syntax-and-types.md` | Baseline fluency |
| Structs and methods | `02-notes/02-structs-methods.md` | Domain modeling starts immediately in Project 1 |
| Interfaces (implicit satisfaction) | `02-notes/03-interfaces.md` | Repository pattern in Project 1 depends on this |
| Error handling idioms | `02-notes/04-error-handling.md` | Every layer of every project returns errors explicitly |
| Packages, modules, visibility | `02-notes/05-packages-modules.md` | Clean Architecture folder structure depends on this |
| Slices, maps, basic collections | `02-notes/06-collections.md` | Used everywhere |
| Goroutines/channels (intro only — depth comes later) | `02-notes/07-concurrency-intro.md` | Just enough to not be surprised; depth revisited before Project 4 |

**Exit criterion for Phase 0:** can write a small CLI tool with structs, an interface-based abstraction, and explicit error propagation, without referring back to docs for syntax. Tracked as an exercise in `03-exercises/`.

---

## 3. Project map (full detail)

### Project 1 — URL Shortener Service
**Folder:** [`04-project-url-shortener/`](./04-project-url-shortener/)
**Difficulty:** Beginner
**Phase:** 1 — Foundations

| Aspect | Detail |
|---|---|
| New concepts | Struct/interface design, error handling in practice, CRUD, REST routing (Gin), PostgreSQL via GORM or sqlx, env config, Docker basics |
| Carried forward from | Phase 0 fundamentals only — first project |
| Features | Shorten URL, custom alias, expiration, click counter, basic analytics, health check |
| Stack | Gin, PostgreSQL, GORM or sqlx, Docker |
| Key design question | GORM (faster to ship, more magic) vs sqlx (more explicit, closer to SQL) — decided and recorded as ADR-001 in this project |

### Project 2 — Inventory Management API
**Folder:** [`05-project-inventory-management/`](./05-project-inventory-management/)
**Difficulty:** Beginner → Intermediate
**Phase:** 1 — Foundations

| Aspect | Detail |
|---|---|
| New concepts | SQL transactions, complex queries (joins, aggregates), repository pattern (formalized), service layer separation, relational design with real foreign-key relationships |
| Carried forward from | Project 1's REST/CRUD/Docker baseline |
| Features | Products, categories, suppliers, stock in/out, inventory history, search, pagination, CSV export |
| Key design question | How stock in/out updates maintain consistency under concurrent requests — first real transactional integrity problem in the repo |

### Project 3 — Booking Management System
**Folder:** [`06-project-booking-system/`](./06-project-booking-system/)
**Difficulty:** Intermediate
**Phase:** 2 — Core Backend Skills
**Domain (pick one, document the choice in PRD):** Barbershop / Clinic / Coworking Space / Photography Studio

| Aspect | Detail |
|---|---|
| New concepts | JWT auth (first auth implementation in the repo), middleware chains, time/timezone handling, availability-checking logic, background jobs, basic scheduling |
| Carried forward from | Project 2's repository/service layers, transactional thinking |
| Features | Auth, booking, schedule, availability checking, cancellation, notification (basic), admin dashboard API |
| Key design question | How double-booking is prevented under concurrent requests — connects forward to Project 4's concurrency-safety theme |

### Project 4 — Digital Wallet API
**Folder:** [`07-project-digital-wallet/`](./07-project-digital-wallet/)
**Difficulty:** Intermediate
**Phase:** 2 — Core Backend Skills

| Aspect | Detail |
|---|---|
| New concepts | Redis, mutexes, race conditions (now load-bearing, not introductory), idempotency keys, financial-grade data consistency, isolation levels |
| Carried forward from | Project 3's auth (reused, not rebuilt — JWT middleware ported forward) |
| Features | Register, login, wallet, balance, top-up, transfer, transaction history, audit logs |
| Key design question | How a transfer is made atomic and idempotent under concurrent retries — the most consistency-critical project in the repo before the Auth Service |

### Project 5 — File Management Service
**Folder:** [`08-project-file-management/`](./08-project-file-management/)
**Difficulty:** Intermediate → Advanced
**Phase:** 3 — Infrastructure Depth

| Aspect | Detail |
|---|---|
| New concepts | Object storage (MinIO), streaming I/O, multipart upload, `context.Context` propagation under real load (cancellation on large uploads), file processing, thumbnailing |
| Carried forward from | Project 4's auth + Redis patterns |
| Features | Upload, download, folder management, search, share links, versioning, thumbnail generation |
| Key design question | Streaming vs. buffering for large files — directly shapes memory/performance characteristics, recorded as an ADR |

### Project 6 — Notification Service
**Folder:** [`09-project-notification-service/`](./09-project-notification-service/)
**Difficulty:** Intermediate → Advanced
**Phase:** 3 — Infrastructure Depth

| Aspect | Detail |
|---|---|
| New concepts | Message queues, worker pools, retry/backoff strategy, cron scheduling, Redis as a queue backend, background processing patterns |
| Carried forward from | Project 5's async/context patterns; Project 3's "basic notification" stub gets replaced by this real service (see §5, dependency notes) |
| Features | Email, push, webhook, scheduled notifications |
| Key design question | At-least-once vs. exactly-once delivery guarantees, and how retries avoid duplicate sends — recorded as an ADR |

### Project 7 — Authentication Service
**Folder:** [`10-project-auth-service/`](./10-project-auth-service/)
**Difficulty:** Advanced
**Phase:** 4 — Platform Engineering

| Aspect | Detail |
|---|---|
| New concepts | OAuth flows, refresh token rotation, email verification, forgot-password flow, RBAC, security hardening practices |
| Carried forward from | Project 3's JWT basics and Project 4's reused-auth experience — this project exists because that ad hoc auth doesn't scale across services |
| Features | Register, login, refresh token, OAuth, email verification, forgot password, RBAC |
| Key design question | This service is explicitly designed to be **reusable by every other project** — interface boundaries are chosen for that, and the decision is recorded as an ADR with the migration plan for retrofitting Projects 3–6 |

### Project 8 — AI Prompt Management API
**Folder:** [`11-project-ai-prompt-management/`](./11-project-ai-prompt-management/)
**Difficulty:** Advanced
**Phase:** 4 — Platform Engineering

| Aspect | Detail |
|---|---|
| New concepts | Modular/multi-tenant architecture, API-key authentication (distinct from user auth), version-control concepts applied to a domain model (prompt versioning), team-workspace data modeling, usage analytics |
| Carried forward from | Project 7's Auth Service (consumed, not reimplemented) |
| Features | Prompt collections, versioning, variables, prompt testing, team workspaces, API keys, usage analytics, history |
| Key design question | How prompt versioning is modeled — full snapshots vs. diffs — recorded as an ADR with the tradeoff analysis |

---

## 4. Timeline (indicative, 20+ hrs/week)

| Project | Estimated duration | Cumulative |
|---|---|---|
| Phase 0 (fundamentals) | 1–2 weeks | Week 2 |
| Project 1 — URL Shortener | 2–3 weeks | Week 5 |
| Project 2 — Inventory Management | 2–3 weeks | Week 8 |
| Project 3 — Booking System | 3–4 weeks | Week 12 |
| Project 4 — Digital Wallet | 3–4 weeks | Week 16 |
| Project 5 — File Management | 4–5 weeks | Week 21 |
| Project 6 — Notification Service | 4–5 weeks | Week 26 |
| Project 7 — Auth Service | 4–6 weeks | Week 32 |
| Project 8 — AI Prompt Management | 4–6 weeks | Week 38 |
| Capstone design + integration | 3–5 weeks | Week 43 |

This is a planning instrument, not a deadline. It exists so that drift is visible and discussable, not so any week is "behind." Update the **actual** column (once added, after Project 1 completes) rather than the estimate — estimates are not revised retroactively.

## 5. Cross-project dependency notes

A few dependencies are easy to miss because they're not adjacent project numbers:

- **Project 3's notification feature is intentionally a stub.** It exists to make the booking system functional, but it is explicitly *not* the real notification architecture — that's Project 6. When Project 6 is done, Project 3 should be revisited to swap the stub for the real service (logged as a revisit in Project 3's `CHANGELOG.md`).
- **Project 3 and 4's auth is intentionally ad hoc.** Both hand-roll JWT middleware. This is deliberate — the pain of duplicating auth logic across two services is what motivates Project 7's design. Do not prematurely build a shared auth service before Project 7; the friction is the lesson.
- **Once Project 7 ships**, Projects 3, 4, and 6 are candidates for an auth-retrofit pass. This is tracked as a milestone (§6) and documented as a migration ADR in Project 7, not silently done.
- **Project 8 is the first project that should consume Project 7 from day one** rather than retrofit it — by this point the reusable service exists and using it natively is itself the lesson (vs. Projects 3/4/6 which show the "before" state).

## 6. Milestones

| Milestone | Trigger | What it represents |
|---|---|---|
| M1 — Fundamentals exit | Phase 0 exit criterion met | Ready to build, not just read |
| M2 — First working production-shaped service | Project 1 complete | Can take a service from PRD to deployed Docker container |
| M3 — Relational thinking | Project 2 complete | Can design and query a real relational schema under transactional constraints |
| M4 — Auth + time-based logic | Project 3 complete | Can build stateful, time-aware backend logic with naive auth |
| M5 — Consistency-critical systems | Project 4 complete | Can reason about race conditions and idempotency in financial-grade logic |
| M6 — Infra-aware backend engineering | Projects 5–6 complete | Can build around object storage and async/queue-based systems |
| M7 — Platform thinking | Project 7 complete | Can design a service meant to be *consumed* by other services, not just stand alone |
| M8 — Full portfolio | Project 8 complete | Eight documented, production-shaped services |
| M9 — Capstone | Capstone integration complete | Can compose independently-built services into one coherent system with a gateway, shared infra, and monitoring |

## 7. Capstone preview

Full capstone architecture is designed in detail once Project 7 is complete (it's the dependency every other service routes through). Target shape:

```
Client
  ↓
API Gateway
  ↓
Authentication Service (Project 7)
  ↓
┌─────────────┬──────────────────┬───────────────┬─────────────────────┐
Inventory Svc   Notification Svc    File Svc        AI Prompt Svc
(Project 2)      (Project 6)       (Project 5)       (Project 8)
  ↓                  ↓                 ↓                  ↓
                  Shared PostgreSQL · Redis · Message Queue
                                  ↓
                              Monitoring
```

Booking (Project 3) and Wallet (Project 4) are evaluated for capstone inclusion once the gateway design is in place — they may be included as additional services behind the gateway, or kept as standalone portfolio pieces, depending on what the capstone is trying to demonstrate. That decision gets its own ADR when the time comes.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Initial roadmap created: phase structure, full 8-project map, indicative timeline, milestones, capstone preview. |
