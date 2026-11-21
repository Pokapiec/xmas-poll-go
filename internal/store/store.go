// Package store keeps parties, people and gifts in SQLite.
package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

var (
	ErrNotFound        = errors.New("not found")
	ErrNotEnoughPeople = errors.New("at least two people are needed to shuffle")
	ErrAlreadyDrawn    = errors.New("everyone already has someone to buy for")
)

type Party struct {
	ID       int64
	PublicID string
	Name     string
}

type Person struct {
	ID         int64
	PartyID    int64
	PublicID   string
	Name       string
	ReceiverID sql.Null[int64]
	GiftCount  int // only filled in by Store.People
}

type Gift struct {
	ID       int64
	PersonID int64
	Name     string
}

type Store struct {
	db *sql.DB
}

// Open opens (or creates) the SQLite database at path and applies the schema.
// The directory holding the file is created if needed.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	dsn := "file:" + path +
		"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func timestamp() string { return time.Now().UTC().Format(time.RFC3339) }

const (
	publicIDAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	publicIDLen      = 16 // about 95 random bits: short enough to paste, impossible to guess
)

// newPublicID returns a random ID for use in links.
func newPublicID() string {
	b := make([]byte, 0, publicIDLen)
	buf := make([]byte, 32)
	for len(b) < publicIDLen {
		rand.Read(buf)
		for _, c := range buf {
			// 248 = 4*62: dropping the bytes above it keeps every character equally likely.
			if c < 248 && len(b) < publicIDLen {
				b = append(b, publicIDAlphabet[c%62])
			}
		}
	}
	return string(b)
}

func notFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *Store) CreateParty(ctx context.Context, name string) (Party, error) {
	p := Party{PublicID: newPublicID(), Name: name}
	ts := timestamp()
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO parties (name, public_id, inserted_at, updated_at) VALUES (?, ?, ?, ?) RETURNING id`,
		p.Name, p.PublicID, ts, ts,
	).Scan(&p.ID)
	return p, err
}

func (s *Store) PartyByPublicID(ctx context.Context, publicID string) (Party, error) {
	var p Party
	err := s.db.QueryRowContext(ctx,
		`SELECT id, public_id, name FROM parties WHERE public_id = ?`, publicID,
	).Scan(&p.ID, &p.PublicID, &p.Name)
	return p, notFound(err)
}

func (s *Store) PartyByID(ctx context.Context, id int64) (Party, error) {
	var p Party
	err := s.db.QueryRowContext(ctx,
		`SELECT id, public_id, name FROM parties WHERE id = ?`, id,
	).Scan(&p.ID, &p.PublicID, &p.Name)
	return p, notFound(err)
}

func (s *Store) RenameParty(ctx context.Context, partyID int64, name string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE parties SET name = ?, updated_at = ? WHERE id = ?`, name, timestamp(), partyID)
	return err
}

// People returns everyone in the party along with how many gifts each has on their wishlist.
func (s *Store) People(ctx context.Context, partyID int64) ([]Person, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id, p.party_id, p.public_id, p.name, p.receiver_id, COUNT(g.id)
		FROM people p
		LEFT JOIN gifts g ON g.person_id = p.id
		WHERE p.party_id = ?
		GROUP BY p.id
		ORDER BY p.id`, partyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.PartyID, &p.PublicID, &p.Name, &p.ReceiverID, &p.GiftCount); err != nil {
			return nil, err
		}
		people = append(people, p)
	}
	return people, rows.Err()
}

func (s *Store) CreatePerson(ctx context.Context, partyID int64, name string) (Person, error) {
	p := Person{PartyID: partyID, PublicID: newPublicID(), Name: name}
	ts := timestamp()
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO people (name, public_id, party_id, inserted_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id`,
		p.Name, p.PublicID, p.PartyID, ts, ts,
	).Scan(&p.ID)
	return p, err
}

const personColumns = `id, party_id, public_id, name, receiver_id`

func scanPerson(row *sql.Row) (Person, error) {
	var p Person
	err := row.Scan(&p.ID, &p.PartyID, &p.PublicID, &p.Name, &p.ReceiverID)
	return p, notFound(err)
}

func (s *Store) PersonByPublicID(ctx context.Context, publicID string) (Person, error) {
	return scanPerson(s.db.QueryRowContext(ctx,
		`SELECT `+personColumns+` FROM people WHERE public_id = ?`, publicID))
}

// PersonNumber returns the person's position in their party (1-based, in the
// order people were added), the same number the organizer page shows.
func (s *Store) PersonNumber(ctx context.Context, p Person) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM people WHERE party_id = ? AND id <= ?`, p.PartyID, p.ID).Scan(&n)
	return n, err
}

func (s *Store) PersonByID(ctx context.Context, id int64) (Person, error) {
	return scanPerson(s.db.QueryRowContext(ctx,
		`SELECT `+personColumns+` FROM people WHERE id = ?`, id))
}

// ShufflePeople assigns everyone in the party a random person to buy a gift for,
// making sure nobody gets themselves. Earlier assignments are replaced.
//
// It refuses with ErrAlreadyDrawn when everyone already has someone. That check
// runs inside the write transaction (opened with BEGIN IMMEDIATE, see the DSN in
// Open), so two overlapping requests can't both draw: the second waits for the
// first to commit and then sees its result.
func (s *Store) ShufflePeople(ctx context.Context, partyID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		`SELECT id, receiver_id IS NOT NULL FROM people WHERE party_id = ? ORDER BY id`, partyID)
	if err != nil {
		return err
	}
	var ids []int64
	paired := 0
	for rows.Next() {
		var id int64
		var hasReceiver bool
		if err := rows.Scan(&id, &hasReceiver); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
		if hasReceiver {
			paired++
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	if len(ids) < 2 {
		return ErrNotEnoughPeople
	}
	if paired == len(ids) {
		return ErrAlreadyDrawn
	}

	receivers := derange(ids)
	ts := timestamp()
	for i, giver := range ids {
		if _, err := tx.ExecContext(ctx,
			`UPDATE people SET receiver_id = ?, updated_at = ? WHERE id = ?`, receivers[i], ts, giver); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) Gifts(ctx context.Context, personID int64) ([]Gift, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, person_id, name FROM gifts WHERE person_id = ? ORDER BY id`, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gifts []Gift
	for rows.Next() {
		var g Gift
		if err := rows.Scan(&g.ID, &g.PersonID, &g.Name); err != nil {
			return nil, err
		}
		gifts = append(gifts, g)
	}
	return gifts, rows.Err()
}

func (s *Store) CreateGift(ctx context.Context, personID int64, name string) error {
	ts := timestamp()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO gifts (name, person_id, inserted_at, updated_at) VALUES (?, ?, ?, ?)`,
		name, personID, ts, ts)
	return err
}

// DeleteGift deletes a gift only if it belongs to the given person.
func (s *Store) DeleteGift(ctx context.Context, personID, giftID int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM gifts WHERE id = ? AND person_id = ?`, giftID, personID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
