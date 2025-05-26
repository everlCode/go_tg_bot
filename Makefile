connect:
	 ssh root@89.111.172.150
# Основной compose-файл
COMPOSE_FILE=docker-compose.yml

# Override-файл для дев-сборки
COMPOSE_FILE_DEV=docker-compose.yml:docker-compose.override.yml

.PHONY: build-dev build-prod up-dev up-prod down

# Сборка для девелопмента (используется override)
build-dev:
	COMPOSE_FILE=$(COMPOSE_FILE_DEV) docker compose build

# Сборка для продакшена (только основной файл)
build-prod:
	COMPOSE_FILE=$(COMPOSE_FILE) docker compose build

# Запуск для девелопмента
up-dev:
	docker compose -d

# Запуск для продакшена
up-prod:
	docker compose -f docker-compose.yml up -d

# Остановка и удаление контейнеров
stop:
	docker compose down
