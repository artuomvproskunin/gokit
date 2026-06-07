# gokit — переиспользуемые Go-библиотеки

Module: github.com/artuomvproskunin/gokit. ОДИН модуль, подпакеты. Без доменной логики автомоек.

## Подпакеты
- config   — загрузка/валидация конфигурации из env
- log      — slog-логгер (JSON/текст, request-scoped поля)
- postgres — pgxpool, health, TxManager (WithTx), раннер goose-миграций
- httpx    — mux, middleware (request-id/recovery/log/CORS/timeout), единый формат ошибки,
             decode/encode + валидация
- auth     — session store, OTP/magic link, session-middleware, API-key middleware

## Правила (это библиотека)
- Публичный API минимальный и стабильный; ломающие изменения → новый minor/major тег + CHANGELOG.
- Без завязки на проект; конфигурация через параметры/опции; без скрытых глобалов.
- context.Context первым аргументом; явное управление ресурсами (Close/таймауты).
- Тесты обязательны; postgres/auth — интеграционные через testcontainers.
- Семантическое версионирование тегами vX.Y.Z; потребители подключают по версии.

## Команды
- go test ./... ; golangci-lint run ; релиз: git tag vX.Y.Z && git push --tags
