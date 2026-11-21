package store

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"testing"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

// Overlapping draw requests (a double click, two open tabs) must result in
// exactly one draw; the others see that everyone already has a pair.
func TestShufflePeopleConcurrent(t *testing.T) {
	st := openTestStore(t)
	ctx := context.Background()

	for round := range 10 {
		party, err := st.CreateParty(ctx, "Party")
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"A", "B", "C", "D"} {
			if _, err := st.CreatePerson(ctx, party.ID, name); err != nil {
				t.Fatal(err)
			}
		}

		const requests = 8
		var (
			wg       sync.WaitGroup
			start    = make(chan struct{})
			errs     = make(chan error, requests)
			drawn    int
			refusals int
		)
		for range requests {
			wg.Go(func() {
				<-start
				errs <- st.ShufflePeople(ctx, party.ID)
			})
		}
		close(start)
		wg.Wait()
		close(errs)

		for err := range errs {
			switch {
			case err == nil:
				drawn++
			case errors.Is(err, ErrAlreadyDrawn):
				refusals++
			default:
				t.Fatalf("round %d: unexpected error: %v", round, err)
			}
		}
		if drawn != 1 || refusals != requests-1 {
			t.Fatalf("round %d: %d draws and %d refusals, want 1 and %d", round, drawn, refusals, requests-1)
		}
	}
}

func TestShufflePeopleRefusals(t *testing.T) {
	st := openTestStore(t)
	ctx := context.Background()

	party, _ := st.CreateParty(ctx, "Party")
	st.CreatePerson(ctx, party.ID, "A")
	if err := st.ShufflePeople(ctx, party.ID); !errors.Is(err, ErrNotEnoughPeople) {
		t.Fatalf("one person: got %v, want ErrNotEnoughPeople", err)
	}

	st.CreatePerson(ctx, party.ID, "B")
	if err := st.ShufflePeople(ctx, party.ID); err != nil {
		t.Fatalf("first draw: %v", err)
	}
	if err := st.ShufflePeople(ctx, party.ID); !errors.Is(err, ErrAlreadyDrawn) {
		t.Fatalf("second draw: got %v, want ErrAlreadyDrawn", err)
	}

	// Someone joining after the draw allows drawing again.
	st.CreatePerson(ctx, party.ID, "C")
	if err := st.ShufflePeople(ctx, party.ID); err != nil {
		t.Fatalf("draw after a new person: %v", err)
	}
}
