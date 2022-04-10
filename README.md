<img src="./icescraper.png" width="100" alt="Badly hand-drawn ice scraper used as a logo">

# Prerequisites

- Golang 1.16+
- NodeJS

Optionally:

- docker to build the docker image
- flyctl to run on fly.io

# Getting started

To build and run:

    $ make build
    $ cp config.yml.example config.yml
    $ # edit config, then:
    $ ./icescraper

# Development

Running migrations:

    make migrate

To re-generate the DB client stubs:

    make generate

# Running in production

To apply migrations on the production database:

    DATABASE_URL="postgres://<connection_string>/<db_name>" npx prisma migrate deploy

A sample `fly.toml` is provided for convenience, and with `flyctl` available, fly.io deploments are a matter of:

    make deploy
