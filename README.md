# Golang Backend Engineering — Learning Repository

A structured, documentation-first journey from Go fundamentals to production-grade backend engineering, organized the way a real engineering team organizes its work: plans, decisions, architecture, and projects, all written down as they happen — not reconstructed afterward.

This repository is **not** a tutorial archive. It is a working engineering record. Every project here is treated as if it were going into production: it has a PRD, an architecture doc, a database design, an API spec, a testing strategy, and a changelog. The code matters, but the documentation is what makes the learning durable.

---

## Why this exists

Learning a language by following tutorials produces syntax familiarity. It does not produce the judgment to design a system, justify a tradeoff, or explain a decision to someone else six months later. This repository closes that gap by forcing every concept learned to land in a real artifact: a document, a decision record, a working service.

The standard this repo holds itself to: **if a backend engineer with no prior context opened this repo, they should be able to understand what was built, why it was built that way, and what I'd do differently next time — without reading the code first.**

## How this repository is organized

```
golang-backend-roadmap/
├── 00-learning-plan.md              # Time budget, learning method, weekly cadence
├── 01-roadmap.md                    # Full sequence of topics + project map + milestones
├── 02-notes/                        # Topic-by-topic engineering notes (not tutorials — notes)
├── 03-exercises/                    # Small, throwaway drills for specific language mechanics
├── 04-project-url-shortener/        # Project 1 — Beginner
├── 05-project-inventory-management/ # Project 2 — Beginner → Intermediate
├── 06-project-booking-system/       # Project 3 — Intermediate
├── 07-project-digital-wallet/       # Project 4 — Intermediate
├── 08-project-file-management/      # Project 5 — Intermediate → Advanced
├── 09-project-notification-service/ # Project 6 — Intermediate → Advanced
├── 10-project-auth-service/         # Project 7 — Advanced
├── 11-project-ai-prompt-management/ # Project 8 — Advanced
├── 12-interview/                    # Interview prep distilled from real project decisions
├── 13-cheatsheet/                   # Quick-reference sheets per topic
├── 14-snippets/                     # Reusable, documented code fragments
├── 15-resources/                    # Curated external references, books, articles
├── 99-templates/                    # Document templates used by every project
└── README.md                        # You are here
```

## How to read this repo

- Start with [`00-learning-plan.md`](./00-learning-plan.md) to understand the method and cadence.
- [`01-roadmap.md`](./01-roadmap.md) is the master index — every topic, every project, every milestone, in order.
- Each numbered project folder (`04-` through `11-`) is self-contained: README, PRD, architecture, database, API spec, setup, deployment, testing, ADRs, changelog, lessons learned. None of them assume you've read the code.
- `99-templates/` holds the blank skeleton for every document type used across this repo. Projects don't redefine document structure — they fill in the template.

## Project sequence

| # | Project | Focus | Status |
|---|---------|-------|--------|
| 1 | [URL Shortener](./04-project-url-shortener/) | Fundamentals, REST, PostgreSQL, Docker | Documented |
| 2 | [Inventory Management](./05-project-inventory-management/) | Transactions, repository pattern, relational design | Documented |
| 3 | [Booking System](./06-project-booking-system/) | Auth, scheduling, background jobs | Documented |
| 4 | [Digital Wallet](./07-project-digital-wallet/) | Financial consistency, Redis, concurrency | Documented |
| 5 | [File Management](./08-project-file-management/) | Object storage, streaming, multipart upload | Documented |
| 6 | [Notification Service](./09-project-notification-service/) | Queues, workers, retries, cron | Documented |
| 7 | [Auth Service](./10-project-auth-service/) | OAuth, RBAC, token rotation, reusable service | Documented |
| 8 | [AI Prompt Management](./11-project-ai-prompt-management/) | Modular architecture, API keys, versioning | Not started |

Status is updated as each project moves through `Not started → In progress → Documented → Refactored`.

## Engineering principles followed in this repo

1. **No project is static.** Every project gets revisited after later projects teach better patterns. A revisit is logged in that project's `CHANGELOG.md` and, if the decision is significant, a new ADR.
2. **Decisions are written down before they're acted on**, not rationalized afterward. See any project's `adr/` folder.
3. **Documentation describes the system as it is**, not as it was planned. Docs are updated alongside code, not at the end.
4. **Nothing ships without a testing strategy**, even if the strategy is intentionally minimal for a given project's scope — and that scoping decision is itself documented.
5. **Security, scalability, and deployment are considered from project 1**, even when the answer for a beginner project is "out of scope, see Future Improvements" — the muscle of asking the question matters more than the answer at that stage.

## Final capstone

After all eight projects are complete, this repo concludes with a capstone: a composed backend ecosystem connecting the Auth Service, Inventory Service, Notification Service, File Service, and AI Prompt Service behind a shared gateway, with shared PostgreSQL/Redis infrastructure, async messaging, and monitoring. Design work for the capstone begins once Project 7 (Auth Service) is documented, since it's the dependency every other service routes through.

---

*This repository is a living document. Last structural update: see [`01-roadmap.md`](./01-roadmap.md) changelog section.*
