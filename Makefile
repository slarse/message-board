.PHONY: run clean

dist:
	mkdir -p dist
	cd frontend \
		&& npm i \
		&& npx vite build
	mv frontend/dist dist/message-board-frontend
	go build -C backend -o "${CURDIR}/dist/message-board-backend"

run: dist
	env MESSAGE_BOARD_FRONTEND_PATH="${CURDIR}/dist/message-board-frontend" ./dist/message-board-backend

clean:
	rm -rf dist/
