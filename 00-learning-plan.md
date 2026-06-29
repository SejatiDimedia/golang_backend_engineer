# Learning Plan

**Status:** Active
**Last updated:** 2026-06-29
**Owner:** Solo learner (this repository's author)

---

## 1. Purpose

This document defines *how* the learning happens — the method, the weekly structure, and the standards each artifact must meet. It does not define *what* gets learned in what order; that's [`01-roadmap.md`](./01-roadmap.md).

## 2. Starting context

| Dimension | Value |
|---|---|
| Prior Go experience | None — first exposure to the language |
| Prior programming background | Yes — comfortable with at least one other language; transferable concepts (types, control flow, basic data structures) already internalized |
| Prior backend experience | Assumed minimal-to-none at repo start; this plan does not assume prior REST API, database, or deployment experience and introduces each explicitly |
| Time available | 20+ hours/week (near full-time) |
| Primary goal | Structured self-learning with durable, well-documented reference material; portfolio value is a byproduct, not the target |
| Secondary goal | Build enough production judgment to reason about real backend systems, not just pass syntax-level familiarity |

Because the goal is depth over speed-to-portfolio, this plan does **not** compress fundamentals. Time saved by prior programming experience is reinvested into Go-specific idioms that don't transfer cleanly from other languages (see §5).

## 3. Method

This repository follows a **build-document-reflect** loop for every unit of learning, whether that unit is a single language feature or an entire project:

1. **Build** — write the smallest version of the thing that works. Exercises for syntax/mechanics; full implementation for projects.
2. **Document** — capture the decision, not just the output. A PRD before code. An ADR when a nontrivial choice is made. Notes when a concept clicks (or doesn't).
3. **Reflect** — after each project, write `LESSONS-LEARNED.md` honestly, including what was misunderstood initially. Misunderstandings are signal, not failure, and are kept in the record rather than quietly fixed and forgotten.

This loop applies at two scales:
- **Daily/weekly scale** — topics in `02-notes/`, drills in `03-exercises/`.
- **Project scale** — the eight portfolio projects, each fully documented per the template in `99-templates/`.

## 4. Weekly cadence (20+ hrs/week)

This is a target shape, not a rigid schedule — adjust week to week, but don't let any category go to zero for more than one week in a row.

| Block | Hours/week (approx.) | Activity |
|---|---|---|
| Fundamentals / new concepts | 4–6 | Reading, official docs, targeted exercises in `03-exercises/` |
| Project implementation | 10–12 | Writing the actual service code for the current project |
| Documentation | 3–4 | PRD, architecture, API spec, ADRs — written *during* the project, not after |
| Review / refactor | 1–2 | Revisiting a previous project with newly learned concepts (see §6) |
| Reflection / notes consolidation | 1 | Updating `02-notes/`, writing lessons learned, journal entry |

At 20+ hrs/week, each beginner project (Projects 1–2) should take roughly **2–3 weeks** including documentation; intermediate projects (3–5) roughly **3–4 weeks**; advanced projects (6–8) roughly **4–6 weeks**, since they introduce more unfamiliar infrastructure (queues, OAuth, object storage) per project. See [`01-roadmap.md`](./01-roadmap.md) §4 for the full timeline and how it can flex.

## 5. Go-specific idioms to deliberately slow down for

Prior programming experience accelerates syntax acquisition but does not transfer for these — budget explicit time, don't assume osmosis:

- **Error handling as values, not exceptions.** No `try/catch` mental model. Errors are returned, checked, and wrapped explicitly at every call site.
- **Interfaces are implicit and small.** No `implements` keyword; satisfaction is structural. Idiomatic Go favors small, single-method interfaces defined at the point of use, not large upfront contracts.
- **Concurrency model.** Goroutines and channels are not "threads with extra syntax" — the idioms (worker pools, context cancellation, select statements) need dedicated practice, especially before Project 4 (Digital Wallet, race conditions) and Project 6 (Notification Service, workers/queues).
- **No classical OOP inheritance.** Composition via embedding, not inheritance. This affects how all eight projects structure their domain models from Project 1 onward.
- **Package-level visibility.** Exported vs. unexported (capitalization-based) instead of access modifiers — affects package design from the first project.
- **Explicit context propagation.** `context.Context` threading through call chains for cancellation/timeouts/request-scoped values — becomes load-bearing starting around Project 3 (scheduling) and critical by Project 5 (streaming/multipart).

These get dedicated entries in `02-notes/` and `13-cheatsheet/` as they're encountered, not just inline comments in project code.

## 6. The "no project is static" rule

Every project, once documented, is a candidate for revisit when a later project teaches something that would have improved it. This is not optional polish — it's the mechanism by which the roadmap simulates how real codebases evolve.

Rules for a revisit:
- The revisit is **scoped and logged**, not a silent rewrite. Open or update an ADR explaining what changed and why.
- The revisit is **recorded in that project's `CHANGELOG.md`**, including the version bump and the originating project that motivated the change (e.g. "Repository pattern refactor introduced after Project 2").
- The original decision is **not deleted from history** — old ADRs are marked `Superseded`, not removed. The record of "I used to think X, now I think Y because Z" is itself a learning artifact.

## 7. Documentation standards

Every document produced in this repo, regardless of project, must:

- Be written for a reader with **zero prior context** on this specific repo.
- State **decisions and reasoning**, not just outcomes ("we use PostgreSQL" is incomplete; "we use PostgreSQL because the domain has relational integrity constraints that matter more than horizontal write scale at this stage" is the standard).
- Avoid tutorial language. This is engineering documentation, not a walkthrough — write as a contributor documenting a system, not as a teacher explaining a concept.
- Use the matching template from `99-templates/` so structure stays consistent across all eight projects.
- Be kept current. A document describing a decision that's been superseded is either updated or explicitly marked superseded with a pointer to the new decision — never left silently stale.

## 8. Definition of done (per project)

A project is not "complete" until all of the following exist and are internally consistent with each other and with the actual code:

- [ ] `README.md`
- [ ] `PRD.md`
- [ ] `ROADMAP.md` (project-local — feature sequencing within the project)
- [ ] `ARCHITECTURE.md`
- [ ] `DATABASE.md`
- [ ] `API.md`
- [ ] `SETUP.md`
- [ ] `DEPLOYMENT.md`
- [ ] `TESTING.md`
- [ ] `CHANGELOG.md`
- [ ] `adr/` (one or more ADRs)
- [ ] `FUTURE-IMPROVEMENTS.md`
- [ ] `LESSONS-LEARNED.md`
- [ ] Working code with the test suite passing
- [ ] Status updated in root `README.md` project table

## 9. Review checkpoints

At the end of every project, before starting the next one, answer in `LESSONS-LEARNED.md`:

1. What did I get wrong in the initial architecture, and when did I notice?
2. What would I change if I rebuilt this today?
3. What concept from this project most needs reinforcement before it appears again?
4. Which earlier project should be revisited in light of what I just learned?

These four questions are non-negotiable per project — skipping them is the single most common way this kind of self-study quietly turns back into tutorial-following.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Initial learning plan created. Profile: 20+ hrs/week, new to Go with prior programming background, self-learning primary goal. |
