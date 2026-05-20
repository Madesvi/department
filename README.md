### 1. Запуск контейнера PostgreSQL

Запустите тестовый контейнер локально с помощью следующей команды:

```bash
docker run --name my-db \
  -e POSTGRES_PASSWORD='secure' \
  -e POSTGRES_USER=serverdbtest \
  -e POSTGRES_DB=restdb \
  -p 5432:5432 \
  -d postgres:latest
```

### 2. Настройка окружения (.env)

Для подключения приложения и утилиты миграций создайте файл `.env` в корне проекта:

```env
DB_USER=serverdbtest
DB_PASSWORD=secure
DB_NAME=restdb
DB_PORT=5432
DB_HOST=localhost
DB_SSLMODE=disable
```

### 3. Управление миграциями через Makefile

Используйте быстрые команды для управления схемой БД:

*   **Применить все миграции (Up):**
    ```bash
    make migrate-up
    ```
*   **Откатить последнюю миграцию (Down):**
    ```bash
    make migrate-down
    ```
*   **Проверить статус миграций:**
    ```bash
    make migrate-status
    ```