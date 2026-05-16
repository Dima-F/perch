# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Perch is a personal fishing log web app. It tracks fishing sessions, catches, locations, and lures. The UI is in Ukrainian. The database is a SQLite file stored at `/mnt/c/Users/fiial/Dropbox/fishing/perch` by default (override with `DB_PATH` env var). The server listens on `:3000` by default (override with `ADDR`).

## Commands

```bash
# Run the server
go run ./cmd/server

# Build
go build ./cmd/server

# Regenerate templ files after editing .templ files
templ generate

# Run tests
go test ./...
```

> Install the templ CLI: `go install github.com/a-h/templ/cmd/templ@latest`

## Architecture

The stack is Go + Fiber (HTTP) + templ (HTML templates) + SQLite (`modernc.org/sqlite`, CGO-free).

**Request flow:**
```
Fiber router (main.go) → handler/* → repository/sqlite/* → SQLite DB
                                   ↓
                         templates/pages/*.templ (rendered via templ + fiber adaptor)
```

**Layers:**
- `internal/models/` — plain Go structs; nested structs (e.g. `Lure.LureModel`, `Catch.Fish`) are populated by the repository layer, not by DB tags.
- `internal/repository/repository.go` — interfaces for each domain entity.
- `internal/repository/sqlite/` — concrete SQLite implementations; these are used directly (not via interface) in handlers.
- `internal/handler/` — one file per domain; `handler.go` wires everything into `Handlers`.
- `internal/templates/` — `layouts/base.templ` is the shell; `pages/*.templ` are full-page components rendered server-side. HTMX is loaded from CDN for any future partial updates.

**Template workflow:** Edit `.templ` files, run `templ generate` to regenerate the corresponding `*_templ.go` files. Both files must be committed together. Never hand-edit `*_templ.go`.
