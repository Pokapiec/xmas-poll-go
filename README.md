# Xmas Gift Poll

A small Christmas gift poll. The organizer creates a party, adds people and
shuffles who buys a present for whom. Each person gets their own link where they
write down at least three gifts they would like; after that they see who they
are buying for and that person's wishlist.

Go rewrite of the Phoenix LiveView app: standard library `net/http` and
`html/template`, SQLite through the pure-Go `modernc.org/sqlite` driver,
Tailwind + daisyUI for styling and htmx (`hx-boost`) for page-reload-free forms.

The UI is styled as a printed advent calendar: your draw sits behind a door
that opens once you've written three gifts. Polish (default) and English on
every page. Icons are [Lucide](https://lucide.dev) (ISC), inlined from
`web/templates/icons/`; the heading font is Bricolage Grotesque (OFL), self-hosted
from `web/static/fonts/`. Product notes live in `PRODUCT.md`.

## Layout

```
cmd/xmas-poll/     main: config, wiring, graceful shutdown
internal/store/    SQLite: schema.sql, queries, link IDs, the shuffle
internal/server/   routes, handlers, translations, template helpers
web/               templates/, static/ (served), assets/ (Tailwind source); embedded into the binary
db/                SQLite files: the local default and the Fly volume mount point
```

## Running locally

```sh
make run            # or: go run ./cmd/xmas-poll  ->  http://localhost:4000, database in db/
make test
```

`web/static/app.css` is generated and checked in, so `go run ./cmd/xmas-poll`
works without any extra tools. After changing classes in `web/templates/`, rebuild it:

```sh
make css            # or `make css-watch` while editing
```

This downloads the Tailwind standalone CLI and daisyUI into `bin/` and
`web/assets/vendor/` (both git-ignored). No Node.js needed.

## Configuration

| Variable        | Default                  | Description          |
| --------------- | ------------------------ | -------------------- |
| `PORT`          | `4000`                   | HTTP port            |
| `DATABASE_PATH` | `db/xmas_gift_poll.db`   | SQLite database file |

The schema is created on startup (`internal/store/schema.sql`). Table and column
names match the Elixir version, so an existing database from it can be reused as-is.

## Deploying to Fly.io

```sh
fly launch --no-deploy   # first time only; or edit `app` in fly.toml
fly volumes create db --region ams --size 1
fly deploy
```

The image runs from `/app` and the `db` volume is mounted at `/app/db`, so the
database is `/app/db/xmas_gift_poll.db`. It is a static binary on distroless
(~27 MB) running as uid 65534 (`nobody`), same as the Elixir release, so a
database file that release left on a volume stays writable. To reuse the old
Elixir volume, mount it at `/app/db` and set `DATABASE_PATH=/app/db/name.db`.

## Pages

| Path           | Who        |
| -------------- | ---------- |
| `/`            | organizer: create a party |
| `/p/{party}`   | organizer: add people, draw pairs, copy/share links, rename |
| `/g/{person}`  | person: wishlist, then who they buy for (PL by default, EN switch) |

IDs are 16 random alphanumeric characters. Parties and people created before
short links keep their UUIDs, and the old long URLs
(`/parties/{party}/people/new`, `/parties/{party}/people/{person}/gifts`)
redirect to the short ones.
