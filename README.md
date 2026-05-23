# Perch

Особистий щоденник рибалки. Відстежує сесії, улови, локації та приманки. UI українською мовою.

## Стек

- Go + [Fiber v2](https://github.com/gofiber/fiber) — HTTP
- [templ](https://templ.guide/) — HTML-шаблони (server-side rendering)
- SQLite ([modernc.org/sqlite](https://gitlab.com/cznic/sqlite), CGO-free)
- [Pico CSS](https://picocss.com/) — стилі
- [MDI](https://pictogrammers.com/library/mdi/) — іконки
- HTMX — для майбутніх часткових оновлень

## Налаштування

| Змінна | За замовчуванням | Опис |
|--------|-----------------|------|
| `DB_PATH` | `/mnt/c/Users/fiial/Dropbox/fishing/perch` | Шлях до SQLite-файлу |
| `ADDR` | `:3000` | Адреса сервера |

## Команди

```bash
# Розробка з hot-reload (air + templ generate)
air

# Запуск без hot-reload
go run ./cmd/server

# Збірка бінарника
go build -o ./bin/server ./cmd/server

# Регенерація templ-шаблонів (після редагування .templ файлів)
templ generate

# Тести
go test ./...
```

## Структура

```
cmd/server/         — точка входу, wiring + реєстрація маршрутів
docs/schema.md      — ER-схема бази даних
internal/
  db/               — підключення до SQLite
  models/           — Go-структури
  repository/       — інтерфейси та SQLite-реалізації
  handler/
    render.go       — спільний templ-рендерер
    handler.go      — Handlers struct + New()
    *.go            — один файл на домен, кожен має Register(fiber.Router)
  templates/        — templ-шаблони (layouts + pages)
static/
  css/style.css     — кастомні стилі (доповнення до Pico)
  vendor/           — офлайн-копії JS/CSS бібліотек
```

## Маршрути

Кожен хендлер реєструє свої маршрути через `Register(r fiber.Router)`. У `main.go`:

```go
app.Get("/", h.Sessions.List)
h.Sessions.Register(app.Group("/sessions"))
h.Catches.Register(app.Group("/catches"))
h.Locations.Register(app.Group("/locations"))
h.Lures.Register(app.Group("/lures"))
```

## Шаблони

Редагувати тільки `.templ` файли, після чого запустити `templ generate`. Файли `*_templ.go` — згенеровані, не редагувати вручну. Комітити обидва файли разом.

## Статичні бібліотеки

JS/CSS бібліотеки зберігаються локально в `static/vendor/` (без зовнішніх CDN):

| Файл | Версія |
|------|--------|
| `vendor/pico.min.css` | Pico CSS v2 |
| `vendor/htmx.min.js` | HTMX v2.0.3 |
| `vendor/mdi/` | MDI Font v7.4.47 |
