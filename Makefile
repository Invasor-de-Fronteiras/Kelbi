LAUNCHER_FOLDER_PATH=www/kelbi
MIGRATION_FOLDER_PATH=database

install-launcher-deps:
	cd $(LAUNCHER_FOLDER_PATH) && npm install

build-launcher: install-launcher-deps
	cd $(LAUNCHER_FOLDER_PATH) && npm run build

format-launcher:
	cd $(LAUNCHER_FOLDER_PATH) && npm run format

lint-launcher:
	cd $(LAUNCHER_FOLDER_PATH) && npm run lint

format: format-launcher lint-launcher
	go fmt .

run-db:
	docker-compose up -d

install-migrations-deps:
	cd $(MIGRATION_FOLDER_PATH) && npm install

setup-migrations: install-migrations-deps
	cd $(MIGRATION_FOLDER_PATH) && cp .env.example .env

run-migrations: setup-migrations
	cd $(MIGRATION_FOLDER_PATH) && npx prisma migrate deploy


update-migrations: install-migrations-deps
	cd $(MIGRATION_FOLDER_PATH) && npx prisma db pull

setup: run-db run-migrations build-launcher
	cp config.example.json config.json

make start:
	go run .