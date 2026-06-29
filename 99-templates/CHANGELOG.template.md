<!--
TEMPLATE: CHANGELOG.md
Follows Keep a Changelog conventions, adapted for a learning-project context
where "releases" are project milestones, not version numbers tied to deploys.
-->

# Changelog: [Project Name]

All notable changes to this project are documented here, including post-hoc revisits driven by lessons from later projects (per the "no project is static" rule in `00-learning-plan.md`).

Format loosely follows [Keep a Changelog](https://keepachangelog.com/).

---

## [Unreleased]

### Added
###  Changed
### Fixed

## [0.1.0] — [date] — Initial implementation

### Added
- Initial PRD, architecture, and database design
- Core CRUD for [entity]
- [feature]

### Notes
Any context about this initial release worth preserving.

---

<!--
When a later project motivates a revisit of THIS project, add an entry like:

## [0.2.0] — [date] — Repository pattern refactor (motivated by Project 2)

### Changed
- Refactored data access from direct GORM calls in handlers to a formal
  repository interface, per the pattern learned in Project 2.
- See adr/004-repository-pattern-retrofit.md for full reasoning.

### Why this isn't a "fix"
This wasn't broken — it worked. It's documented as a changed entry because
later learning produced a better pattern, and the revisit itself is the
point of this repository's structure.
-->
