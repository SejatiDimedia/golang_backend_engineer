# Snippets

Reusable, documented code fragments — patterns that recur across projects often enough to warrant extracting, but not (yet) formalized into a shared internal package.

## Convention

Each snippet is a folder with the snippet itself plus a short `README.md` stating:
- What it does
- Which project it was first extracted from
- Any known limitation in its current generalized form

```
14-snippets/
├── http-error-response/
│   ├── README.md
│   └── error_response.go
├── graceful-shutdown/
│   ├── README.md
│   └── shutdown.go
└── pagination-helper/
    ├── README.md
    └── pagination.go
```

## When something graduates out of here

If a snippet becomes load-bearing enough across 3+ projects, that's a signal it might deserve to become a real shared internal package (relevant once Project 7's Auth Service sets the precedent for reusable services) — note that decision as an ADR in whichever project first formalizes it, and leave a pointer here.
