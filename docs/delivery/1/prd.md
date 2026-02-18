# PBI-1: Project Foundation & Dev Environment

[View in Backlog](../backlog.md)

## Overview

Scaffold the monorepo with a Next.js frontend, Go backend, and Docker-based dev environment. This PBI establishes the foundational project structure, build tooling, and local development workflow that all subsequent PBIs depend on.

## Problem Statement

No project structure exists yet. Every subsequent PBI needs a working frontend app, backend server, and containerised dev environment to build against. Without this foundation, no feature work can begin.

## User Stories

- As a developer, I want a monorepo with clearly separated frontend and backend directories so that I can work on either independently
- As a developer, I want a single command to start the full dev environment so that I can begin developing immediately
- As a developer, I want basic CI checks (lint, build) so that code quality is enforced from the start

## Technical Approach

- Monorepo with `frontend/` (Next.js + TypeScript) and `backend/` (Go module) at the root
- `pnpm` as the package manager for the frontend
- Docker Compose for local dev environment (frontend dev server, Go backend, potentially a Postgres for diagram storage)
- Basic linting: ESLint + Prettier for frontend, `golangci-lint` for backend
- Makefile or similar for common dev commands

## UX/UI Considerations

N/A — infrastructure PBI.

## Acceptance Criteria

1. Monorepo structure matches PRD repository layout (`frontend/`, `backend/`, `docker/`, `docs/`)
2. `pnpm install` and `pnpm dev` starts the Next.js frontend on a dev server
3. `go build ./...` compiles the Go backend without errors
4. Docker Compose brings up all dev services with a single command
5. ESLint and Prettier configured for frontend; `golangci-lint` configured for backend
6. Basic CI pipeline runs lint and build checks

## Dependencies

- **Depends on**: None
- **External**: Node.js, Go, Docker, pnpm

## Open Questions

- Should we use a task runner (Makefile, Taskfile, or Turborepo) for cross-project commands?
- What Go version to target? (1.22+ recommended)
- Do we need a database for diagram persistence at this stage, or start with in-memory/file storage?

## Related Tasks

[View Tasks](./tasks.md)
