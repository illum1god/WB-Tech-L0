# Order Service

## Описание

Демонстрационный сервис для отображения данных о заказах. Сервис подключается к Apache Kafka, получает данные о заказах, сохраняет их в PostgreSQL и кэширует в памяти. В случае перезапуска сервиса кэш восстанавливается из базы данных. HTTP-сервер предоставляет API для получения данных о заказах по их ID, а также простой веб-интерфейс для отображения этих данных.

## Стек технологий

- **Go/Golang**
- **PostgreSQL**
- **Docker Compose**
- **Apache Kafka**
- **Gin (HTTP-сервер)**

## Запуск проекта

1. **Запуск сервиса:**
   ```bash
   make migrate-up
   make run
   ```

2. **Остановка сервиса:**
   ```bash
   Ctrl+C
   make migrate-down
   ```

## Функциональность

1. **Подключение к Apache Kafka:**
   - Подписка на топик `orders`.
   - Получение и обработка данных о заказах.

2. **Сохранение данных:**
   - Запись полученных данных в PostgreSQL.
   - Кэширование данных в памяти.

3. **Восстановление кэша:**
   - При перезапуске сервиса кэш восстанавливается из базы данных.

4. **HTTP-сервер:**
   - Предоставляет API для получения данных о заказах по ID.
   - Простой веб-интерфейс для отображения данных о заказах.

## Ссылка на видео с работой сервиса

- https://disk.yandex.ru/i/DqzROJiukX-tfg
