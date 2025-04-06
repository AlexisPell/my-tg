Актуальные порты:
Сервисные:
- zookeeper:2181
- kafka:9092
- schema-registry:8081 (temporarily switched off)
- kafka-ui:4200
- traefik:80 (web), 8080 (dashboard), 8082 (metrics)
- grafana:3000
- prometheus:9090
Функциональные:
- websocket:8000
- messages:8001 // 9001 for /metrics
- users:8002

1. Architecture
- API Gateway — маршрутизация запросов
- WebSocket Service — обработка реального времени, обновления чатов и сообщений
- Auth Service — управление пользователями, JWT/OAuth, интеграция с Keycloak
- User Service — хранение информации о пользователях
- Chat Service — управление чатами, добавление/удаление пользователей в группы
- Message Service — отправка и хранение сообщений в ScyllaDB
- Notification Service

    2. Tech list
- Go (gRPC + REST)
- ScyllaDB
- Kafka
- API-gateway: Traefik
- Docker + Kubernetes
- Tracing + observability: Grafana, Prometheus, Loki, Tempo

    3. Services
- Auth Service
    Используем Keycloak для управления пользователями
    Генерация и валидация JWT-токенов
    gRPC API для проверки пользователей
- User Service
    PostgreSQL для хранения данных о пользователях
    REST/gRPC API для работы с профилями
    Подключение к Auth Service для верификации
- Chat Service
    Хранение информации о чатах группы, участники
    API для создания/удаления чатов
    gRPC для взаимодействия с Message Service
- Message Service
    ScyllaDB для хранения сообщений
    gRPC API для отправки сообщений
    Интеграция с WebSocket Service
- WebSocket Service
    Подключение пользователей через WebSocket
    Получение сообщений из Kafka
    Отправка сообщений клиентам
- Notification Service
    Генерация push-уведомлений
    Поддержка Firebase, email
    
    Очередность разработки
Настроить Auth Service (Keycloak)
Разработать User Service и API Gateway
Реализовать Chat Service (CRUD чатов)
Реализовать Message Service с ScyllaDB
Добавить WebSocket Service
Интегрировать Kafka и сервис уведомлений