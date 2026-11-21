-- Table and column names match the original Ecto schema, so a database
-- created by the Elixir version of the app can be opened as-is.
-- Every statement is idempotent; it runs on each startup.

CREATE TABLE IF NOT EXISTS parties (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    public_id   TEXT NOT NULL,
    inserted_at TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS parties_public_id_index ON parties (public_id);

CREATE TABLE IF NOT EXISTS people (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    public_id   TEXT NOT NULL,
    party_id    INTEGER NOT NULL REFERENCES parties (id),
    -- The person this person buys a gift for. NULL until names are shuffled.
    receiver_id INTEGER REFERENCES people (id),
    inserted_at TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS people_party_id_index ON people (party_id);
CREATE UNIQUE INDEX IF NOT EXISTS people_public_id_index ON people (public_id);

CREATE TABLE IF NOT EXISTS gifts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    person_id   INTEGER NOT NULL REFERENCES people (id),
    inserted_at TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS gifts_person_id_index ON gifts (person_id);
