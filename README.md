# **Music Max API**

## **Описание**

Music Max API — это серверное приложение, предоставляющее API для работы с музыкальными данными. Оно включает функциональность для работы с песнями и текстами песен (lyrics). API поддерживает базовые CRUD операции с базой данных PostgreSQL.

## **Требования**

Перед запуском убедитесь, что у вас установлены:

- [Go 1.20+](https://golang.org/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Docker](https://www.docker.com/products/docker-desktop) (для работы с базой данных в контейнере)
- [Make](https://www.gnu.org/software/make/)

## **Настройка проекта**

### **1. Клонирование репозитория**

```bash
git clone https://github.com/MyrzakhmetSmagul/music_max.git
cd music_max
```

### **2. Настройка переменных окружения**
Создайте файл .env по пути `./config/.env` проекта и укажите следующие переменные:

```env
# MODES ['DEV', 'TEST', 'PROD'] для уровня логировании
MODE=DEV
# database configs
DB_HOST=localhost
DB_PORT=5435
DB_USER=user
DB_PASSWORD=password
DB_NAME=db_name
DB_SSLMODE=disable
# http server configs
HTTP_PORT=:3030
HTTP_READ_TIMEOUT=5
HTTP_WRITE_TIMEOUT=10
# API configs
API_ADDR=http://path/to/music/info/info
```

### **3. Настройка `compose.yml`**
Создайте файл `compose.yml` в корне проекта и укажите следующие переменные:

```yml
services:
  db:
    image: postgres
    restart: always
    user: postgres
    volumes:
      - db-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=db_name
      - POSTGRES_PASSWORD=password
      - POSTGRES_USER=user
    ports:
      - 5435:5432
    healthcheck:
      test: [ "CMD", "pg_isready" ]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
    db-data:
```

### **4. Запуск Docker контейнера с PostgreSQL**
Чтобы запустить базу данных PostgreSQL в контейнере, выполните команду:

```bash
make upDB
```

### **5. Инициализация базы данных**
После запуска контейнера выполните миграции для инициализации структуры базы данных:

```bash
make initDB
```

### **6. Запуск сервера**
Запустите приложение с помощью команды:

```bash
make run
```

API будет доступно по адресу: `http://localhost:3030/api/v1`.

### **6. Остановка контейнера**
Чтобы остановить и удалить контейнеры:

```bash
make downDB
```

**Swagger документация**

Для автоматической генерации документации API используется Swagger. Swagger-документация доступна по следующему пути:
`http://localhost:8080/swagger/index.html`

**Команды Makefile**
* `make upDB` — запуск Docker контейнера с базой данных.
* `make initDB` — инициализация базы данных.
* `make run` — запуск приложения.
* `make downDB` — остановка и удаление контейнеров.

**Автор**

[Myrzakhmet] — разработчик этого проекта.