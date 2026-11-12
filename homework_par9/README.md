# Пара 9 - API-first проектирование сервиса ozonroute

API нужен сотруднику ПВЗ, чтобы принять заказ, найти его, выдать клиенту или оформить возврат.

## Контракт

`api.yaml` описывает пять операций в OpenAPI 3.0.3. Проект опирается на [исходный код ozonroute](https://github.com/4s4ken/ozonroute/tree/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0).

| HTTP | Путь | Сценарий и источник |
| --- | --- | --- |
| GET | `/v1/orders?userId=42` | Список заказов клиента. [ListOrders](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/app/list_orders.go) |
| GET | `/v1/orders/{orderId}` | Карточка одного заказа. [GetOrder](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/app/service.go) |
| POST | `/v1/orders` | Приём заказа от курьера. [Accept](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/app/accept_order.go) |
| PATCH | `/v1/orders/{orderId}` | Выдача или приём возврата от клиента через `ISSUE` или `RETURN`. [Process](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/app/process_order.go) |
| DELETE | `/v1/orders/{orderId}/custody` | Снятие заказа с учёта ПВЗ при возврате курьеру. [Return](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/app/return_order.go) |

Это новый REST-контракт для учебного задания. Действующие маршруты gRPC-Gateway описаны в [pvz.proto](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/api/pvz/pvz.proto). PATCH берёт один заказ из пакетного сценария `Process`. DELETE относится к учёту заказа в ПВЗ. После возврата карточка по-прежнему доступна через GET со статусом `RETURNED`.

Список поддерживает фильтр `inPvz` и нумерацию страниц с 1. При подключении к SQL-хранилищу нужно передавать `page - 1`, так как оно вычисляет смещение от нуля. [Реализация пагинации](https://github.com/4s4ken/ozonroute/blob/96751e3f4cd2d0fb04cfc2a8ea3cc036879f42d0/internal/storage/db/list.go)

В `components` вынесены отдельные DTO для приёмки, обработки, карточки, списка, результата смены статуса и ошибки. Идентификаторы передаются десятичными строками, чтобы сохранить точность `uint64`. Код адаптера должен проверять верхнюю границу при преобразовании строки в число. Данные хранилища и события outbox остаются внутри приложения.

Все пять операций описывают успешный ответ `200` по условию задания. `400` означает ошибку входных данных, `404` - отсутствие заказа, `409` - дубликат или конфликт с текущим состоянием заказа. `401` означал бы отсутствие действительных учётных данных, а `403` - отсутствие доступа у опознанного клиента. Оба кода относятся к проверке доступа, которой в этом учебном контракте нет.

## Как делали

Спецификацию и каркас подготовил OpenAI Codex на базе GPT-6. Точный идентификатор модели в сессии недоступен.

Промпты пользователя в сокращении:

1. «Спроектируй API в Swagger-стиле. Подготовь OpenAPI 3.0 с GET списка, GET одной записи, POST, PUT или PATCH и DELETE. Для каждой операции укажи summary, параметры и ответы. Вынеси DTO в components. Сгенерируй сервис и опиши процесс в README».
2. «Изучи GitHub-репозиторий 4s4ken/ozonroute и выполни задание на основе реального функционала. Создай homework_par9 с api.yaml, запусти openapi-generator для Go, проверь валидность и опиши, что редактировалось вручную».

До генерации агент сверил операции с protobuf и Go-кодом. В `api.yaml` вручную описаны маршруты, публичные DTO, причины ошибок и примеры тел запросов. При сверке с реализацией уточнены сроки для PATCH и причины конфликта для DELETE. Пользователь ручных правок не вносил. Файлы генератора остались без ручных изменений. Команда `go mod tidy` создала `go.sum`.

## Валидация и генерация

Из папки `homework_par9`, с OpenAPI Generator CLI 7.20.0:

```sh
openapi-generator-cli validate -i api.yaml
openapi-generator-cli generate -i api.yaml -g go-server -o service --git-user-id 4s4ken --git-repo-id ozonroute/homework_par9/service --additional-properties=packageName=openapi,outputAsLibrary=false,hideGenerationTimestamp=true
cd service
go mod tidy
go test ./...
go run .
```

Проверка OpenAPI Generator завершилась без замечаний. `go test ./...` проверяет, что пакеты каркаса компилируются. При локальном запуске на порту 8080 проверены все пять маршрутов. Каждый принял корректный запрос и вернул `501 Not Implemented`. Для GET списка без `userId` стандартный обработчик генератора вернул `422`, хотя контракт задаёт `400`. При подключении к ozonroute понадобится реализовать бизнес-обработчики и привести обработку ошибок к контракту.

## Импорт в Postman

Выберите **Import**, откройте `api.yaml` и создайте коллекцию из определения OpenAPI. Для локального запуска используйте `http://localhost:8080`. Примеры тел POST и PATCH находятся в спецификации.
