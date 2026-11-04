# Shortener

Сервис коротких ссылок. Пользователь входит через Keycloak и в веб-интерфейсе создает ссылку, смотрит список, меняет исходный адрес и удаляет ссылку. Посетитель открывает адрес из восьми букв и цифр и получает перенаправление на исходный сайт. MCP-клиент вызывает те же операции инструментами `shorten_url`, `list_urls`, `resolve_url`, `update_url`, `delete_url` и читает ресурсы ссылок.

Идентификатор генерируется случайно. Исходный адрес хранится в Redis по этому идентификатору, а порядок списка задает sorted set `shortener:ids`. Запросы `/api` требуют access token клиента `shortener-web`. Переход по короткой ссылке и `/mcp` открыты без токена.

## Контекст

![Контекстная диаграмма](schemas/images/system-context.png)

[Исходный код диаграммы](schemas/system-context.puml)

## Контейнеры

![Контейнерная диаграмма](schemas/images/system-containers.png)

[Исходный код диаграммы](schemas/system-containers.puml)

Веб-интерфейс на nginx отдает SPA и проксирует `/api`, идентификатор короткой ссылки и страницы Keycloak. Сервис на Go слушает HTTP API, MCP и редирект. Keycloak выдает токен, сервис проверяет его по JWKS.

## Компоненты сервиса

![Компоненты сервиса shortener](schemas/images/shortener-components.png)

[Исходный код диаграммы](schemas/shortener-components.puml)

## Создание ссылки

![Создание короткой ссылки](schemas/images/shorten-sequence.png)

[Исходный код диаграммы](schemas/shorten-sequence.puml)

## Переход по ссылке

![Переход по короткой ссылке](schemas/images/redirect-sequence.png)

[Исходный код диаграммы](schemas/redirect-sequence.puml)

## Вызов MCP

![Вызов инструмента MCP](schemas/images/mcp-sequence.png)

[Исходный код диаграммы](schemas/mcp-sequence.puml)

## Модель хранения

![Модель хранения ссылок](schemas/images/storage-model.png)

[Исходный код диаграммы](schemas/storage-model.puml)
