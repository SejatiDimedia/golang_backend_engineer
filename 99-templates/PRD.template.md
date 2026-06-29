<!--
TEMPLATE: PRD.md (Product Requirements Document)
Write this BEFORE writing implementation code. If code already exists when this
is written, that's a signal the build-document-reflect loop was skipped — note
it honestly in LESSONS-LEARNED.md rather than backdating this document.
-->

# PRD: [Project Name]

**Status:** `Draft | Approved | Implemented`
**Author:** [you]
**Last updated:** [date]

---

## 1. Problem statement

What real-world problem does this system solve? Write this as if pitching to a stakeholder who doesn't know or care what language it's built in. One paragraph, no jargon.

## 2. Goals

- Goal 1 (business-framed, not technical — e.g. "allow staff to track stock levels in real time" not "build a CRUD API")
- Goal 2
- Goal 3

## 3. Non-goals

Explicitly out of scope for this version. This section matters as much as Goals — it's what prevents scope creep and documents intentional simplification. Every non-goal here should map to an entry in `FUTURE-IMPROVEMENTS.md`.

- Non-goal 1
- Non-goal 2

## 4. Target users / personas

Who uses this system and how. Even for a portfolio project, define 1–2 concrete personas — it disciplines every later design decision.

| Persona | Need | Frequency of use |
|---|---|---|
| | | |

## 5. Functional requirements

Numbered, testable requirements. Each should be specific enough that a test case could be written directly from it.

| ID | Requirement | Priority |
|---|---|---|
| FR-1 | | Must |
| FR-2 | | Must |
| FR-3 | | Should |

## 6. Non-functional requirements

| Category | Requirement |
|---|---|
| Performance | e.g. expected request volume, acceptable latency |
| Security | e.g. auth requirements, data sensitivity |
| Availability | e.g. acceptable downtime for a learning project vs. what production would actually require — state both |
| Scalability | What scale this is designed for *now*, and what would break first if scale increased 10x — answering this is the exercise, not actually building for 10x |
| Data consistency | Especially relevant for Projects 2, 4, 6 — state isolation/consistency requirements explicitly |

## 7. Constraints

Technical, time, or learning-objective constraints that shaped requirements. E.g. "must use PostgreSQL, not because it's objectively correct for this domain, but because relational modeling is the learning objective for this project."

## 8. Success criteria

How do you know this project achieved what it set out to. Tie back to the learning goals in `01-roadmap.md` as well as the functional requirements above.

## 9. Open questions

Anything unresolved at the time of writing. Resolve and move to an ADR once decided — don't let this section become a graveyard of unanswered questions by the time the project ships.

---

## Revision history

| Date | Change |
|---|---|
| | Initial draft |
