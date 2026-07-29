# Foundry

Minimal Go project scaffold.

```
<<<<<<< Updated upstream
cmd/server/     entry point
internal/
  config/       environment config
  handler/      HTTP handlers
  server/       routing
web/            embedded HTML
=======
cmd/        entry point
internal/   private packages
>>>>>>> Stashed changes
```

## Run

```bash
make run
```

<<<<<<< Updated upstream
Open [http://localhost:8080](http://localhost:8080).

Health check: `GET /api/health` → `{"status":"ok"}`

## Idea submission

- Form: [http://localhost:8080/](http://localhost:8080/) — Gatewood Lab themed intake form
- Browse placeholder: [http://localhost:8080/ideas.html](http://localhost:8080/ideas.html)
- API: `POST /api/ideas` with JSON body (see `docs/superpowers/specs/2026-05-27-creative-idea-form-design.md`)
- Theme spec: `docs/superpowers/specs/2026-05-27-foundry-gatewood-theme-design.md`
- Brand assets: `web/img/` (koi mark, koi accent, calligraphy watermark)

Set `PORT` to change the listen port.
=======
## Build

```bash
make build
```
>>>>>>> Stashed changes
