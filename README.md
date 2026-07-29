# Foundry

Idea-intake site for capturing project ideas. Paper-themed UI served by Go templates.

```
cmd/           HTTP server entry point
ui/
  html/        Go html/template files
  static/      CSS and other static assets
```

## Run

```bash
make run
```

Open [http://localhost:4000](http://localhost:4000).

## Build

```bash
make build
```

## Endpoints

- `GET /` — homepage with idea form
- `POST /api/ideas` — submit an idea (`title`, `description`, `category`, `tags`, `contact`)
- `GET /static/...` — static assets
