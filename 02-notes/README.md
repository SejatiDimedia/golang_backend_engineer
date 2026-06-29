# Engineering Notes

Topic-by-topic notes written as concepts are encountered — **not** tutorials, and not a copy of documentation. The standard for a note here: it should capture the *decision-relevant* understanding of a concept — what it is, when to reach for it, what it costs, and what idiom in another language it does or doesn't map to.

## Convention

- One file per topic, numbered in rough order of first encounter (see `01-roadmap.md` Phase 0 for the initial set).
- Filename: `NN-topic-name.md`.
- Each note should be short enough to re-read in under 5 minutes — if a topic needs more than that, split it.

## Suggested note structure

```markdown
# [Topic]

## What it is
## When to use it
## Go-specific idiom (especially vs. [your prior language])
## Common mistake / gotcha
## Where it showed up in a project
```

That last section — "where it showed up in a project" — is what separates this from generic notes. Link back to the specific project and file/decision where this concept became real, once it does.

## Index

| File | Topic | First needed in |
|---|---|---|
| `01-syntax-and-types.md` | Core syntax, types | Phase 0 |
| `02-structs-methods.md` | Structs and methods | Phase 0 |
| `03-interfaces.md` | Interfaces, implicit satisfaction | Phase 0 / Project 1 |
| `04-error-handling.md` | Error handling idioms | Phase 0 / Project 1 |
| `05-packages-modules.md` | Packages, modules, visibility | Phase 0 / Project 1 |
| `06-collections.md` | Slices, maps | Phase 0 |
| `07-concurrency-intro.md` | Goroutines/channels (intro) | Phase 0, deepened before Project 4 |

(Add rows as new topics are encountered — this index should always reflect what's actually in this folder.)
