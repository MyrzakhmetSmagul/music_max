# Название вашего Docker контейнера и базы данных
DB_NAME=music_max-db
DB_USER=musician

# Путь к SQL-скрипту
SQL_SCRIPT=./schema/up.sql

# Имя сервиса в docker-compose
SERVICE_NAME=db

# Имя контейнера
CONTAINER_NAME=music_max-db-1

.PHONY: upDB initDB run downDB

# Запускает Docker контейнер с PostgreSQL
upDB:
	docker-compose up -d $(SERVICE_NAME)

# Инициализирует базу данных, выполняя SQL-скрипт
initDB:
	docker exec -i $(CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) < $(SQL_SCRIPT)

# Запускает ваше приложение
run:
	go run ./cmd/main.go

# Останавливает и удаляет контейнер
downDB:
	docker-compose down