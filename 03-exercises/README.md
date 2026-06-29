# Exercises

Small, throwaway drills for specific language mechanics — distinct from the eight portfolio projects, which are real systems. An exercise here is allowed to be ugly, incomplete, and never touched again once the mechanic is internalized.

## Convention

- Folder per exercise: `NN-short-name/`
- Each exercise folder needs at minimum a one-paragraph `README.md` stating what mechanic it's drilling and a `main.go` (or equivalent).
- No PRD, no architecture doc, no tests required — that rigor is reserved for the numbered projects. Exercises are deliberately lightweight.

## Phase 0 exercise set (see `01-roadmap.md`)

| Exercise | Drills |
|---|---|
| `01-cli-struct-basics/` | Structs, methods, basic I/O |
| `02-interface-swap/` | Defining a small interface, swapping implementations |
| `03-error-wrapping/` | `errors.Is`, `errors.As`, wrapping with `%w` |
| `04-goroutine-fanout/` | Basic goroutine + channel fan-out/fan-in |

Add exercises as gaps are identified — e.g. if a project reveals a shaky concept (per that project's `LESSONS-LEARNED.md` §3), the reinforcement drill belongs here, named to reference what motivated it.
