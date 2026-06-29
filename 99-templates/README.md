# Document Templates

This folder holds the canonical skeleton for every document type used across all eight projects in this repository. Projects do not redefine document structure — they instantiate these templates and fill them in.

## Why this folder exists

With 8 projects × 12 documents each, structural drift is the default outcome unless there's a single source of truth for "what does a PRD in this repo look like." This folder is that source of truth.

## Usage

When starting a new project (or backfilling docs for an existing one):

```bash
cp 99-templates/README.template.md 04-project-url-shortener/README.md
cp 99-templates/PRD.template.md 04-project-url-shortener/PRD.md
cp 99-templates/ROADMAP.template.md 04-project-url-shortener/ROADMAP.md
cp 99-templates/ARCHITECTURE.template.md 04-project-url-shortener/ARCHITECTURE.md
cp 99-templates/DATABASE.template.md 04-project-url-shortener/DATABASE.md
cp 99-templates/API.template.md 04-project-url-shortener/API.md
cp 99-templates/SETUP.template.md 04-project-url-shortener/SETUP.md
cp 99-templates/DEPLOYMENT.template.md 04-project-url-shortener/DEPLOYMENT.md
cp 99-templates/TESTING.template.md 04-project-url-shortener/TESTING.md
cp 99-templates/CHANGELOG.template.md 04-project-url-shortener/CHANGELOG.md
cp 99-templates/FUTURE-IMPROVEMENTS.template.md 04-project-url-shortener/FUTURE-IMPROVEMENTS.md
cp 99-templates/LESSONS-LEARNED.template.md 04-project-url-shortener/LESSONS-LEARNED.md
mkdir -p 04-project-url-shortener/adr
cp 99-templates/ADR.template.md 04-project-url-shortener/adr/001-<short-title>.md
```

Delete the HTML comment block at the top of each file once instantiated — it's instructional scaffolding for the template itself, not meant to ship in the real document.

## Template index

| Template | Maps to |
|---|---|
| `README.template.md` | Project entry point |
| `PRD.template.md` | Product requirements |
| `ROADMAP.template.md` | Project-local feature sequencing (not the repo-level roadmap) |
| `ARCHITECTURE.template.md` | System design |
| `DATABASE.template.md` | Schema and data design |
| `API.template.md` | Endpoint specification |
| `SETUP.template.md` | Local dev setup |
| `DEPLOYMENT.template.md` | Deployment story |
| `TESTING.template.md` | Test strategy |
| `ADR.template.md` | One per significant decision, lives in each project's `adr/` subfolder |
| `CHANGELOG.template.md` | Version history including revisits |
| `FUTURE-IMPROVEMENTS.template.md` | Deferred work tracker |
| `LESSONS-LEARNED.template.md` | Retrospective |

## Changing a template

If a template itself needs to improve (e.g. you realize halfway through Project 3 that the ARCHITECTURE template is missing a section every project needs), update it here — but note in this README's changelog when and why, and consider whether already-completed projects should backfill the new section.

---

## Changelog

| Date | Change |
|---|---|
| 2026-06-29 | Initial template set created: README, PRD, ROADMAP, ARCHITECTURE, DATABASE, API, SETUP, DEPLOYMENT, TESTING, ADR, CHANGELOG, FUTURE-IMPROVEMENTS, LESSONS-LEARNED. |
