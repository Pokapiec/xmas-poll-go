FROM golang:1.27 AS build
WORKDIR /src

# Download the CSS tooling and Go modules first so they stay cached between builds.
COPY Makefile ./
RUN make tools
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make css \
    && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/xmas-poll ./cmd/xmas-poll \
    && mkdir -p /out/db

# The SQLite driver is pure Go, so the binary is static and needs no libc.
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=build /out/xmas-poll /app/xmas-poll
# The database lives in /app/db (DATABASE_PATH defaults to db/xmas_gift_poll.db).
# On Fly the volume is mounted here; without a volume this empty directory is used.
COPY --from=build --chown=65534:65534 /out/db /app/db

# Same user as the Elixir release ("nobody"), so a database file that release
# created on the volume stays writable.
USER 65534:65534

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/xmas-poll"]
