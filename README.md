Product Accounting Service
REST API сервис для учёта товаров и категорий, написанный на Go.
Стек
Go 1.24 · Echo · PostgreSQL · Docker · Docker Compose · Goose · golangci-lint

Возможности:

CRUD операции для товаров и категорий  
Чистая многослойная архитектура (handler → service → repository)  
Миграции базы данных через Goose  
Middleware с логированием каждого запроса (request_id, метод, путь, статус, время)  
Graceful shutdown  
Валидация входящих данных  
Unit-тесты с моками (testify + mockery)  
Конфигурация через переменные окружения  
Dockerfile с multi-stage сборкой + Docker Compose  

