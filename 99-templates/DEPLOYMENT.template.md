<!--
TEMPLATE: DEPLOYMENT.md
For early projects this may be intentionally minimal (e.g. "this project is
not deployed beyond local Docker; here's what production deployment would
require"). State that explicitly rather than leaving the document thin
without explanation.
-->

# Deployment: [Project Name]

**Deployment status:** `Local only | Deployed to [environment]`

---

## 1. Deployment target

Where this runs (or would run). For early portfolio projects, it's acceptable for this to be "not deployed — documented as if it would be" — state that plainly rather than implying a deployment that doesn't exist.

## 2. Build process

```bash
docker build -t [image-name] .
```

Reference the `Dockerfile` and explain any non-obvious choices (multi-stage build, base image choice, etc.) — link an ADR if the choice was non-trivial.

## 3. Configuration management

How environment-specific config (secrets, DB URLs) is handled per environment. Even if this is just ".env files locally, would be [secrets manager] in real production" — state the gap honestly.

## 4. Database migrations in deployment

How/when migrations run relative to deploys (before, automatically on boot, manual gate).

## 5. Health checks and readiness

What endpoint(s) an orchestrator would use to know this service is healthy.

## 6. Rollback strategy

What happens if a deploy is bad. Even a simple answer ("redeploy previous image tag") is worth stating explicitly.

## 7. What real production would add

Be honest about the gap between this project's deployment story and what an actual production system would need at this project's scale: e.g. CI/CD pipeline, blue-green deploy, autoscaling, secrets rotation, multi-region. List these as deferred, not as oversights — and cross-reference `FUTURE-IMPROVEMENTS.md`.

---

## Changelog

| Date | Change |
|---|---|
| | Initial deployment doc |
