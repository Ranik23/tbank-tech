# Умный Makefile для работы с Docker Compose

# Переменные
DOCKER_COMPOSE = sudo docker-compose
DOCKER_SYSTEM_PRUNE = sudo docker system prune -a -f
DOCKER_COMPOSE_YML = docker-compose.yml

# Сборка и запуск всех контейнеров с пересборкой
all-build:
	@echo "Building and starting Docker containers..."
	$(DOCKER_COMPOSE) up --build

# Сборка и запуск всех контейнеров без пересборки
all:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up

# Очистка всех остановленных контейнеров, неиспользуемых образов и тому подобного
clear:
	@echo "Pruning unused Docker objects..."
	$(DOCKER_SYSTEM_PRUNE)

# Остановка всех контейнеров
stop:
	@echo "Stopping all Docker containers..."
	$(DOCKER_COMPOSE) down

# Проверка состояния контейнеров
status:
	@echo "Checking the status of Docker containers..."
	$(DOCKER_COMPOSE) ps

# Логи для всех контейнеров
logs:
	@echo "Fetching logs for all Docker containers..."
	$(DOCKER_COMPOSE) logs -f

# Удаление всех контейнеров и томов
clean:
	@echo "Removing all Docker containers and volumes..."
	$(DOCKER_COMPOSE) down -v

# Остановить контейнеры, если они не были остановлены в предыдущих командах
stop-if-running:
	@echo "Stopping Docker containers if running..."
	$(DOCKER_COMPOSE) ps -q | xargs -I {} sudo docker stop {}

# Сборка и запуск в фоновом режиме
all-background:
	@echo "Starting Docker containers in the background..."
	$(DOCKER_COMPOSE) up -d

.PHONY: all-build all clear stop status logs clean stop-if-running all-background
