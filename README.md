# Message Board

A simple Reddit-inspired message board application where you can post messages
and reply to existing ones.

## Makefile

There is a [Makefile](Makefile) in the root of the repository for running chores.
Everything described in this document is more precisely described there.
This README is mostly for those that don't have GNU Make.

See the [Makefile](Makefile) for more stuff you can do.

## How to use

> Makefile: `make start`

To run this application, you only need to have `Docker` and the `Docker
Compose` plugin installed. Run the following commands:

```bash
source .env
docker compose up
```

> Note: Older versions of Docker and Compose require the command to be written
> `docker-compose up`.

This will build a shockingly small Docker container for the app, spin up a
database and run some rudimentary migrations to fill it with data. The backend
serves the frontend so it's all self-contained.

The app should then be accessible from `http://127.0.0.1:8000`.

## Available users

There's no authentication or sign-in, but there are three available users.

```
John
Jane
Paul
```

By default, you are John. To impersonate another user, add the
`?author=<username>` query parameter. So, to act as Paul, you'd
go to `http://127.0.0.1:8000?author=Paul`.

It's OF COURSE case sensitive :)

## How to run the backend tests

> Makefile: `make test-backend`

To run the backend tests, you additionally need to have a local install of
`Go`. With that, first make sure that there's nothing running from the Compose
file, then start the database and finally run the tests.

```bash
source .env
docker compose down
docker compose run migrate
cd backend/
go test ./...
```

## How to run the frontend tests

You're in luck, I didn't have time to write any!

## How to shut it all down

> Makefile: `make stop`

Here's how to shut it all down:

```
docker compose down
```
