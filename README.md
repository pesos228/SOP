# Hosting Service

Сервис для управления тарифными планами и серверами.

## Общая информация

Проект представляет собой распределенную систему, организованную в стиле, близком к рекомендациям Ardan Labs. Код разделен на несколько независимых сервисов, взаимодействующих через RabbitMQ и gRPC.

Транспортные слои включают REST (Chi + oapi-codegen), GraphQL (gqlgen) и gRPC. Авторизация интегрирована с Ory Kratos.

## Технологии

- **Go 1.24** & **1.25**
- **PostgreSQL**: для хранения планов, серверов и пулов ресурсов.
- **RabbitMQ** для межсервисных команд и событий.
- **Auth**: Ory Kratos.
- **Observability**: OpenTelemetry (Tempo), Prometheus, Grafana.
- **API**: REST (HATEOAS/HAL), GraphQL, gRPC.
- **Web**: Vanilla JS SPA с поддержкой WebSocket уведомлений.

## Структура репозитория

### Основные компоненты
- `hosting-service`: Логика витрины планов и жизненного цикла серверов.
- `hosting-resources-service`: Управление квотами в пулах ресурсов (gRPC сервер).
- `hosting-provisioning-service`: Обработка очередей на выделение IP.
- `hosting-notification-service`: WebSocket-хаб для real-time обновлений клиента.
- `hosting-kit`: Shared-пакет (logger, auth, messaging, mid, otel) — инфраструктурный фундамент.
- `hosting-contracts`: Общие спецификации (Protobuf, OpenAPI, GraphQL) и структуры событий RabbitMQ.
