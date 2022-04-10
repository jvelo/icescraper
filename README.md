![Badly hand-drawn ice scraper used as a logo](./icescraper.png?raw=true|width=100px)

# Prerequisites

- Golang 1.16+
- NodeJS

# Getting started

# Development

Running migrations:

    make migrate

To re-generate the DB client stubs:

    make generate

# Running in production

Apply migrations:

    DATABASE_URL="postgres://<connection_string>/<db_name>" npx prisma migrate deploy