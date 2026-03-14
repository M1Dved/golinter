# golinter

Линтер для Go, проверяющий лог-сообщения на соответствие правилам оформления.
Совместим с golangci-lint.

## Правила

| Правило | Описание |
|---------|----------|
| `lowercase` | Сообщение начинается со строчной буквы |
| `english` | Сообщение только на английском языке |
| `special` | Нет спецсимволов и эмодзи |
| `sensitive` | Нет чувствительных данных (password, token, secret...) |

## Поддерживаемые логгеры

- `log/slog`
- `go.uber.org/zap`

## Установка и запуск

```bash
go build -o golinter ./cmd/golinter

go vet -vettool=./golinter ./...
```

## Настройка через конфигурационный файл

Создай `.golinter.yml` в корне проекта:

```yaml
rules:
  lowercase: true
  english: true
  special: false      # отключить проверку спецсимволов
  sensitive: true

extra_keywords:       # свои ключевые слова для проверки чувствительных данных
  - ssn
  - credit_card
  - bank_account
```

Указать путь к конфигу:

```bash
./golinter -config=myconfig.yml ./...
```

## Настройка через флаги

```bash
# Отключить правила
./golinter -disable=special,sensitive ./...

# Добавить ключевые слова
./golinter -extra-keywords=ssn,credit_card ./...
```

Флаги имеют приоритет над конфигурационным файлом.

## Интеграция с golangci-lint

```bash
go build -buildmode=plugin -o golinter.so ./plugin
```

`.golangci.yml`:

```yaml
linters-settings:
  custom:
    golinter:
      path: ./golinter.so
      description: Log message style checker
```

## Тесты

```bash
go test ./pkg/analyzer/... -v -race
```
