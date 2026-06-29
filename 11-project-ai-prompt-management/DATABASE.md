# Database Design: AI Prompt Management API

---

## 1. Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    WORKSPACES {
        uint id PK
        varchar name "Not Null"
        timestamp created_at
        timestamp updated_at
    }
    WORKSPACE_MEMBERS {
        uint id PK
        uint workspace_id FK "On Delete CASCADE"
        uint user_id "Not Null"
        varchar role "Default: 'member', Not Null"
        timestamp created_at
    }
    PROMPTS {
        uint id PK
        uint workspace_id FK "On Delete CASCADE"
        varchar name "Not Null"
        text description
        timestamp created_at
        timestamp updated_at
    }
    PROMPT_VERSIONS {
        uint id PK
        uint prompt_id FK "On Delete CASCADE"
        integer version_number "Not Null"
        text prompt_text "Not Null"
        varchar status "Default: 'DRAFT', Not Null"
        timestamp created_at
    }
    API_KEYS {
        uint id PK
        uint workspace_id FK "On Delete CASCADE"
        varchar name "Not Null"
        varchar key_hash UK "Not Null"
        varchar masked_key "Not Null"
        timestamp created_at
        timestamp expires_at "Not Null"
    }
    ANALYTICS_LOGS {
        uint id PK
        uint api_key_id FK "On Delete CASCADE"
        uint prompt_id FK "On Delete CASCADE"
        bigint latency_ms "Not Null"
        integer token_estimate "Not Null"
        integer response_code "Not Null"
        timestamp created_at
    }

    WORKSPACES ||--o{ WORKSPACE_MEMBERS : "has"
    WORKSPACES ||--o{ PROMPTS : "owns"
    WORKSPACES ||--o{ API_KEYS : "has"
    PROMPTS ||--o{ PROMPT_VERSIONS : "contains"
    PROMPTS ||--o{ ANALYTICS_LOGS : "monitored_by"
    API_KEYS ||--o{ ANALYTICS_LOGS : "triggers"
```

## 2. Database Indexes

Untuk memitigasi brute-force attack dan mempercepat validasi token:

```sql
CREATE UNIQUE INDEX idx_api_keys_key_hash ON api_keys (key_hash);
CREATE INDEX idx_prompt_versions_prompt_id ON prompt_versions (prompt_id);
CREATE INDEX idx_analytics_logs_prompt_id ON analytics_logs (prompt_id);
```

**Justifikasi Indeks:**
- `key_hash` dilindungi indeks unik untuk meminimalisasi overhead query DB saat Redis cache mengalami miss.
- `prompt_versions` didukung indeks pada `prompt_id` untuk mempercepat query baca snapshot versi aktif (`status = 'ACTIVE'`).
