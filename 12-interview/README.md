# Interview Preparation

Interview prep distilled from real decisions made in this repo's projects — not generic interview-question lists copied from elsewhere. The premise: the strongest interview answers come from having actually made and justified a decision, not from memorizing a "correct" answer.

## Convention

Each file should, where possible, link back to the ADR or section of a project's documentation that the answer is grounded in.

## Planned structure

| File | Covers |
|---|---|
| `01-go-fundamentals-qa.md` | Language mechanics: interfaces, error handling, concurrency basics |
| `02-system-design-qa.md` | Architecture tradeoffs actually made across the 8 projects |
| `03-database-qa.md` | Schema design, transactions, indexing — grounded in Projects 2, 4 |
| `04-concurrency-qa.md` | Goroutines, channels, race conditions — grounded in Project 4 |
| `05-behavioral-storytelling.md` | "Tell me about a time..." answers built from real `LESSONS-LEARNED.md` entries across projects |

This folder is populated progressively as projects complete — it's most credible when written from real friction encountered, not anticipated in advance. Avoid front-loading generic content here before the projects that would ground it exist.
