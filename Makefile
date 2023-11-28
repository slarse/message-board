.PHONY: run-local clean create-migration start start-db-only stop build test-backend dist-watch

export POSTGRES_USER ?= postgres
export POSTGRES_PASSWORD ?= postgres
export POSTGRES_DB ?= message_board

export DB_USER ?= ${POSTGRES_USER}
export DB_PASSWORD ?= ${POSTGRES_PASSWORD}
export DB_HOST ?= 127.0.0.1
export DB_PORT ?= 5432
export DB_NAME ?= ${POSTGRES_DB}

start: build
	docker compose up --force-recreate -d

stop:
	docker compose down

build:
	docker compose build

dist:
	mkdir -p dist
	cd frontend \
		&& npm i \
		&& npm run build
	mv frontend/dist dist/message-board-frontend
	go build -C backend -o "${CURDIR}/dist/message-board-backend"

dist-watch:
	find ./backend/app ./frontend/src -type f | entr -s 'make clean dist'

run-local: dist start-db-only
	env MESSAGE_BOARD_FRONTEND_PATH="${CURDIR}/dist/message-board-frontend" PORT=8000 \
		./dist/message-board-backend

start-db-only:
	docker compose up -d migrate

create-migration:
	docker run --rm \
		-v ${CURDIR}/backend/migrations:/migrations \
		migrate/migrate create -ext sql -dir /migrations -seq $(name)

test-backend: start-db-only
	cd backend && go test -v ./...

clean:
	rm -rf dist/
