# Follow-Ups

Ideas, improvements, and deferred work captured during planning and implementation.
Review periodically — good candidates become new PBIs via `/new-pbi`.

## Open

| # | Type | Summary | Source | Date | Notes |
|---|------|---------|--------|------|-------|
| 1 | enhancement | Support request validation in Prism (validate incoming requests against OpenAPI spec) | PBI-8 | 2026-02-19 | Prism supports request validation out of the box; currently only response mocking is configured. Enabling `--errors` flag would reject malformed requests — useful for testing client code. |
| 2 | enhancement | Support complex response schemas (nested objects, arrays, $ref) in OpenAPI generator | PBI-8 | 2026-02-19 | MVP treats responseSchema as a flat JSON Schema string. Users may want richer schemas with references, allOf/oneOf, and nested arrays. Worth revisiting after MVP feedback. |
| 3 | tech-debt | Spec file cleanup on container teardown | PBI-8 | 2026-02-19 | Spec files are written to OS temp dir. They'll be cleaned up by the OS eventually, but explicit cleanup in TeardownAll would be cleaner. Low priority since temp files are small. |
| 4 | feature | Dynamic response examples using JSON Schema Faker | PBI-8 | 2026-02-19 | PRD mentions JSON Schema Faker concepts. Prism's built-in dynamic mocking generates random data matching schemas. A dedicated faker could produce more realistic/deterministic examples. Defer until user feedback indicates need. |
| 5 | enhancement | Auto-deploy mode (deploy on diagram change without button click) | PBI-9 | 2026-02-20 | PBI open question: should deploy be automatic? Current plan uses explicit buttons for user control. Auto-deploy could be a toggle in settings for power users who want live-as-you-design feedback. |
| 6 | enhancement | Partial deploy failure recovery (retry individual containers) | PBI-9 | 2026-02-20 | PBI open question: how to handle partial failures. Current plan surfaces errors but doesn't offer per-container retry. A "Retry Failed" button per node would improve UX for diagrams where only one service fails. |
| 7 | tech-debt | Centralized frontend API base URL configuration | PBI-9 | 2026-02-20 | Frontend currently has no shared API client config. Deploy API will hardcode localhost:8080. Should extract to an env-based config shared across all API calls (diagram CRUD + deploy). |
| 8 | enhancement | WebSocket reconnection with exponential backoff and status indicator | PBI-9 | 2026-02-20 | Basic reconnection planned for 9-4. A connection status indicator (green/red dot in toolbar) and exponential backoff with jitter would improve robustness for long-running deploy sessions. |
| 9 | tech-debt | ~~E2E deploy tests with real Docker daemon~~ | PBI-9 | 2026-02-20 | **Subsumed by PBI-12.** Current E2E tests (9-7) use a mock orchestrator for speed and CI compatibility. Integration tests against a real Docker daemon would catch orchestration bugs that mocks miss (image pulls, network creation, port conflicts). Gate behind a build tag like `//go:build docker_integration`. |
| 10 | feature | Container log streaming to frontend | PBI-9 | 2026-02-20 | Users may want to see container stdout/stderr in the UI for debugging. Could extend the WebSocket to stream logs per container, or add a `GET /api/deploy/{nodeId}/logs` endpoint with SSE. Useful for diagnosing startup failures beyond the error status badge. |
| 11 | feature | Server session persistence for container tracking (SQLite) | PBI-9 | 2026-02-20 | Currently container tracking is in-memory only — if the backend restarts, orphaned containers are left running with no way to recover. Persist container mappings (node ID → container ID, ports) to SQLite so on startup the server can detect and re-adopt or clean up previously managed containers. |
