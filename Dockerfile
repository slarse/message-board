FROM node:21 as build-node
COPY ./frontend /src
WORKDIR /src

RUN npm install
RUN npm run build

FROM golang:1.21 as build-go
COPY ./backend /src
WORKDIR /src

env CGO_ENABLED=0
env GOOS=linux

RUN go build -o /bin/main

FROM scratch
COPY --from=build-node /src/dist /message-board-frontend
COPY --from=build-go /bin/main /main

COPY <<EOF /etc/passwd
nobody:x:65534:65534:Unprivileged User:/:/usr/bin/nologin
EOF
USER nobody

ENV MESSAGE_BOARD_FRONTEND_PATH="/message-board-frontend"
ENV PORT=8080

CMD ["/main"]
