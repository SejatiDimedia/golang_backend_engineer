<!--
TEMPLATE: API.md
For larger projects, consider also maintaining an OpenAPI/Swagger spec file
alongside this human-readable version — note that decision in ARCHITECTURE.md
if adopted, and keep both in sync or this doc becomes the stale one.
-->

# API Specification: [Project Name]

**Base URL:** `http://localhost:PORT/api/v1`
**Auth:** `None | Bearer JWT | API Key` — describe scheme briefly, full detail in ARCHITECTURE.md or the Auth Service's own API.md if delegated

---

## Conventions

- All responses are JSON.
- Errors follow this shape:
```json
{
  "error": {
    "code": "string",
    "message": "human readable"
  }
}
```
- Pagination (where applicable): `?page=1&limit=20`, response includes `meta.total`, `meta.page`, `meta.limit`.

## Endpoints

### `[METHOD] /resource`

**Description:** What this does.

**Auth required:** Yes/No

**Request:**
```json
{}
```

**Response — 200 OK:**
```json
{}
```

**Response — error cases:**

| Status | Condition |
|---|---|
| 400 | |
| 404 | |
| 409 | |

---

(Repeat per endpoint. Group related endpoints under `##` headers, e.g. `## URLs`, `## Analytics`.)

## Rate limiting

State if applicable, and what's deferred to `FUTURE-IMPROVEMENTS.md` if not yet implemented.

## Versioning strategy

How this API will version if breaking changes are introduced later (even if v1 is the only version that ever ships for a given project, state the intended strategy).

---

## Changelog

| Date | Change |
|---|---|
| | Initial API spec |
