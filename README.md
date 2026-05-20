### 1. Настройка окружения (.env)

Для подключения приложения и утилиты миграций создайте файл `.env` в корне проекта:

```env
DB_USER=serverdbtest
DB_PASSWORD=secure
DB_NAME=restdb
DB_PORT=5455
DB_HOST=localhost
DB_SSLMODE=disable
# Application server configuration
SERVER_PORT=3000
APP_ENV=development
# APP_ENV=production
LOG_LEVEL=DEBUG
```

Или выполните комманду:
```bash
cp .env.example .env
```

### 2. Запуск контейнеров
Запустите сборку приложения и развертывание сети контейнеров:
```bash
docker compose up --build
```

### 3. Применение Миграций:
Откройте новое окно терминала на вашем хосте и примените SQL-миграции через `goose`:
```bash
make migrate_up
```
* API-сервер станет доступен локально по адресу: `http://localhost:3000`
* База данных PostgreSQL будет доступна для внешних клиентов (DBeaver/Postman) по порту `5455`


### 4. Управление миграциями через Makefile

Используйте быстрые команды для управления схемой БД:

*   **Создание миграции:**
    ```bash
    make migrate_create
    ```
*   **Применить все миграции (Up):**
    ```bash
    make migrate_up
    ```
*   **Откатить последнюю миграцию (Down):**
    ```bash
    make migrate_down
    ```
*   **Проверить статус миграций:**
    ```bash
    make migrate_status
    ```