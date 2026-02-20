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
