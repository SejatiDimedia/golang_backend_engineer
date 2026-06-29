<!--
TEMPLATE: DATABASE.md
-->

# Database Design: [Project Name]

**Status:** `Draft | Approved | Implemented`
**Engine:** PostgreSQL [version] (or other — state and justify if different)

---

## 1. Entity-relationship overview

```
[Describe or ASCII-diagram the core entities and relationships.
 Example:]

Product ──< StockMovement >── Warehouse
   │
   └──< belongs to >── Category
```

## 2. Schema

### Table: `[table_name]`

| Column | Type | Constraints | Notes |
|---|---|---|---|
| id | UUID / SERIAL | PRIMARY KEY | |
| | | | |
| created_at | TIMESTAMP | NOT NULL DEFAULT now() | |
| updated_at | TIMESTAMP | NOT NULL DEFAULT now() | |

(Repeat per table.)

## 3. Relationships

| From | To | Type | On delete |
|---|---|---|---|
| | | one-to-many / many-to-many | CASCADE / RESTRICT / SET NULL |

## 4. Indexes

| Table | Index | Reason |
|---|---|---|
| | | e.g. "frequent lookup by short_code, must be unique and fast" |

## 5. Transactions and consistency

Where in this schema does correctness depend on transactional boundaries? Be explicit about isolation level requirements where relevant (especially Projects 2, 4, 6).

| Operation | Transaction boundary | Isolation concern |
|---|---|---|
| | | |

## 6. Migrations strategy

Tool used (e.g. golang-migrate), naming convention, and how migrations are applied in each environment (see `DEPLOYMENT.md` for production application strategy).

## 7. Sample queries

A few representative queries that exercise the schema's harder relationships — useful as a sanity check that the design actually serves the access patterns it claims to.

```sql
-- Example: paginated search with filter
```

---

## Changelog

| Date | Change |
|---|---|
| | Initial schema |
