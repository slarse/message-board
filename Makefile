.PHONY: run-local clean create-migration start

export POSTGRES_USER ?= postgres
export POSTGRES_PASSWORD ?= postgres
export POSTGRES_DB ?= message_board

dist:
	mkdir -p dist
	cd frontend \
		&& npm i \
		&& npm run build
	mv frontend/dist dist/message-board-frontend
	go build -C backend -o "${CURDIR}/dist/message-board-backend"

run-local: dist
	env MESSAGE_BOARD_FRONTEND_PATH="${CURDIR}/dist/message-board-frontend" PORT=8000 \
		./dist/message-board-backend

start:
	docker compose up --force-recreate -d

create-migration:
	docker run --rm \
		-v ${CURDIR}/backend/migrations:/migrations \
		migrate/migrate create -ext sql -dir /migrations -seq $(name)

clean:
	rm -rf dist/
