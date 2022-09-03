LAUNCHER_FOLDER_PATH=www/kelbi
MIGRATION_FOLDER_PATH=database

install-launcher-deps:
	cd $(LAUNCHER_FOLDER_PATH) && npm install

build-launcher: install-launcher-deps
	cd $(LAUNCHER_FOLDER_PATH) && npm run build

run-db:
	docker-compose up -d

install-migrations-deps:
	cd $(MIGRATION_FOLDER_PATH) && npm install

setup-migrations: install-migrations-deps
	cd $(MIGRATION_FOLDER_PATH) && cp .env.example .env

run-migrations: setup-migrations
	cd $(MIGRATION_FOLDER_PATH) && npx prisma migrate deploy

setup: run-db run-migrations build-launcher
	cp config.example.json config.json

make start:
	go run .