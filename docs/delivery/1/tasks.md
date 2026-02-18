# Tasks for PBI 1: Project Foundation & Dev Environment

This document lists all tasks associated with PBI 1.

**Parent PBI**: [PBI 1: Project Foundation & Dev Environment](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 1-1 | [Monorepo Structure & Go Backend Scaffold](./1-1.md) | Proposed | Create monorepo directory layout and scaffold Go backend with module init, main entrypoint, and basic HTTP server |
| 1-2 | [Next.js Frontend Scaffold](./1-2.md) | Proposed | Scaffold Next.js app with TypeScript and pnpm in the frontend directory |
| 1-3 | [Docker Compose Dev Environment](./1-3.md) | Proposed | Create Docker Compose configuration to run frontend and backend in containers with a shared network |
| 1-4 | [Linting & Formatting Configuration](./1-4.md) | Proposed | Configure ESLint + Prettier for frontend and golangci-lint for backend |
| 1-5 | [Makefile for Dev Commands](./1-5.md) | Proposed | Create a Makefile with common dev commands (dev, build, lint, test) |
| 1-6 | [E2E CoS Test](./1-6.md) | Proposed | Verify all PBI-1 acceptance criteria are met end-to-end |

## Dependency Graph

```
1-1 (Go backend scaffold)
 └──► 1-3 (Docker Compose)
1-2 (Next.js frontend scaffold)
 └──► 1-3 (Docker Compose)
1-1 ──► 1-4 (Linting — backend)
1-2 ──► 1-4 (Linting — frontend)
1-3 ──► 1-5 (Makefile — wraps docker/build/lint commands)
1-4 ──► 1-5 (Makefile — wraps lint commands)
1-5 ──► 1-6 (E2E CoS Test — verifies everything)
```

## Implementation Order

1. **1-1** — Go backend scaffold (no dependencies, foundational)
2. **1-2** — Next.js frontend scaffold (no dependencies, can parallel with 1-1)
3. **1-3** — Docker Compose (depends on 1-1, 1-2 — needs both projects to containerise)
4. **1-4** — Linting & formatting (depends on 1-1, 1-2 — needs source files to lint)
5. **1-5** — Makefile (depends on 1-3, 1-4 — wraps all commands)
6. **1-6** — E2E CoS Test (depends on all above — final validation)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|-----------|-------------------|
| 1-1 | Simple | None |
| 1-2 | Simple | create-next-app (well-known, no guide needed) |
| 1-3 | Medium | None (Docker Compose is config, not a package) |
| 1-4 | Simple | ESLint, Prettier, golangci-lint (well-known, no guides needed) |
| 1-5 | Simple | None |
| 1-6 | Simple | None |

## External Package Research Required

None — all tools used (Next.js, Go, Docker Compose, ESLint, Prettier, golangci-lint) are well-established with stable, well-documented APIs.
