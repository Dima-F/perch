# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Perch is a personal fishing log web app. It tracks fishing sessions, catches, locations, and lures. The UI is in Ukrainian. The database is a SQLite file stored at `/mnt/c/Users/fiial/Dropbox/fishing/perch` by default (override with `DB_PATH` env var). The server listens on `:3000` by default (override with `ADDR`).

## Commands

```bash
# Development with hot-reload
air

# Run without hot-reload
go run ./cmd/server

# Build
go build -o ./bin/server ./cmd/server

# Regenerate templ files after editing .templ files
templ generate

# Run tests
go test ./...
```

> Install the templ CLI: `go install github.com/a-h/templ/cmd/templ@latest`

## Architecture

The stack is Go + Fiber (HTTP) + templ (HTML templates) + SQLite (`modernc.org/sqlite`, CGO-free) + Pico CSS + MDI icons + HTMX.

**Request flow:**
```
Fiber router (main.go) → handler/* → repository/sqlite/* → SQLite DB
                                   ↓
                         templates/pages/*.templ (rendered via templ + fiber adaptor)
```

**Layers:**
- `internal/models/` — plain Go structs; nested structs (e.g. `Lure.LureModel`, `Catch.Fish`) are populated by the repository layer, not by DB tags.
- `internal/repository/repository.go` — interfaces for each domain entity.
- `internal/repository/sqlite/` — concrete SQLite implementations; used directly (not via interface) in handlers.
- `internal/handler/render.go` — shared `render()` helper (templ → Fiber response).
- `internal/handler/handler.go` — `Handlers` struct + `New()` wiring.
- `internal/handler/*.go` — one file per domain; each handler has a `Register(r fiber.Router)` method that declares its own routes.
- `internal/templates/` — `layouts/base.templ` is the shell; `pages/*.templ` are full-page components rendered server-side.

**Route registration pattern:**
Each handler owns its routes via `Register(r fiber.Router)`. In `main.go`:
```go
app.Get("/", h.Sessions.List)
h.Sessions.Register(app.Group("/sessions"))
h.Catches.Register(app.Group("/catches"))
h.Locations.Register(app.Group("/locations"))
h.Lures.Register(app.Group("/lures"))
```

**Template workflow:** Edit `.templ` files, run `templ generate` to regenerate the corresponding `*_templ.go` files. Both files must be committed together. Never hand-edit `*_templ.go`.

## Frontend

- **Pico CSS** — classless/semantic CSS framework. Styled via HTML elements directly; custom overrides in `static/css/style.css` using Pico CSS variables (e.g. `--pico-spacing`).
- **MDI icons** — used as `<i class="mdi mdi-*"></i>`. Icon-only buttons use `class="btn-icon"` + `title` attribute.
- **Buttons** — use Pico conventions: `role="button"` on `<a>`, `.secondary.outline` for secondary actions, `.outline.danger` for destructive actions.
- All JS/CSS libraries are vendored locally in `static/vendor/` — no external CDN requests.
