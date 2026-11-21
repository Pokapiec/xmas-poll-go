package server

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"xmas-poll/internal/store"
)

type testApp struct {
	t      *testing.T
	dbPath string
	store  *store.Store
	srv    *httptest.Server
	client *http.Client
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })

	handler, err := New(st, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	jar, _ := cookiejar.New(nil)
	return &testApp{t: t, dbPath: dbPath, store: st, srv: srv, client: &http.Client{Jar: jar}}
}

// do sends a request, follows redirects and returns the final response's status, path and body.
func (a *testApp) do(method, path string, form url.Values) (int, string, string) {
	a.t.Helper()
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, a.srv.URL+path, body)
	if err != nil {
		a.t.Fatal(err)
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := a.client.Do(req)
	if err != nil {
		a.t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		a.t.Fatal(err)
	}
	return resp.StatusCode, resp.Request.URL.Path, string(b)
}

func (a *testApp) get(path string) (int, string, string) { return a.do(http.MethodGet, path, nil) }

func (a *testApp) post(path string, form url.Values) (int, string, string) {
	if form == nil {
		form = url.Values{}
	}
	return a.do(http.MethodPost, path, form)
}

func mustContain(t *testing.T, body string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(body, want) {
			t.Errorf("body does not contain %q", want)
		}
	}
}

func TestFullFlow(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()

	// The app defaults to Polish; this flow runs in English.
	app.get("/?locale=en")
	status, partyURL, body := app.post("/p", url.Values{"name": {"Wigilia"}})
	if status != http.StatusOK || !strings.HasPrefix(partyURL, "/p/") {
		t.Fatalf("create party: status %d, path %s", status, partyURL)
	}
	mustContain(t, body, ">Wigilia</h1>", "Add at least two people to draw.")
	partyID := strings.TrimPrefix(partyURL, "/p/")

	// Shuffling with fewer than two people is refused.
	app.post("/p/"+partyID+"/people", url.Values{"name": {"Ania"}})
	_, _, body = app.post("/p/"+partyID+"/shuffle", nil)
	mustContain(t, body, "Add at least two people before drawing.")

	for _, name := range []string{"Bartek", "Celina"} {
		_, _, body = app.post("/p/"+partyID+"/people", url.Values{"name": {name}})
		mustContain(t, body, "Added "+name+".", ">"+name+"</p>")
	}

	_, _, body = app.post("/p/"+partyID+"/shuffle", nil)
	mustContain(t, body, "Pairs drawn!", "Send everyone their link")

	party, err := app.store.PartyByPublicID(ctx, partyID)
	if err != nil {
		t.Fatal(err)
	}
	people, err := app.store.People(ctx, party.ID)
	if err != nil {
		t.Fatal(err)
	}
	receivers := map[int64]bool{}
	for _, p := range people {
		if !p.ReceiverID.Valid || p.ReceiverID.V == p.ID {
			t.Fatalf("bad receiver for %s: %+v", p.Name, p.ReceiverID)
		}
		receivers[p.ReceiverID.V] = true
	}
	if len(receivers) != len(people) {
		t.Fatalf("receivers are not unique: %v", receivers)
	}

	// Shuffling again once everyone is assigned is refused and changes nothing.
	_, _, body = app.post("/p/"+partyID+"/shuffle", nil)
	mustContain(t, body, "Everyone already has someone to buy for.")

	ania := people[0]
	giftsURL := "/g/" + ania.PublicID

	status, _, body = app.get(giftsURL)
	if status != http.StatusOK {
		t.Fatalf("gifts page: status %d", status)
	}
	mustContain(t, body, `lang="en"`, "Hi, Ania!", "Add 3 more gifts")

	// Switching language sticks via the cookie.
	app.get(giftsURL + "?locale=pl")
	_, _, body = app.get(giftsURL)
	mustContain(t, body, `lang="pl"`, "Cześć, Ania!", "Dodaj jeszcze 3 prezenty")
	app.get(giftsURL + "?locale=en")

	// The receiver is only revealed after adding three gifts.
	for i, gift := range []string{"Socks", "Book", "Chocolate"} {
		if strings.Contains(body, "You&#39;re buying a gift for") {
			t.Fatalf("receiver shown after only %d gifts", i)
		}
		_, _, body = app.post(giftsURL, url.Values{"name": {gift}})
		mustContain(t, body, "Gift added.", gift)
	}
	receiver, err := app.store.PersonByID(ctx, ania.ReceiverID.V)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, body,
		">"+receiver.Name+"</h2>",
		"You&#39;re buying a gift for this person",
		"Their list is still empty.",
		"door-swing") // the door swings open right after the third gift

	// It opens with the animation only once.
	_, _, body = app.get(giftsURL)
	if strings.Contains(body, "door-swing") {
		t.Fatal("door animation played again on reload")
	}

	// Gift counts show up on the organizer page.
	_, _, body = app.get(partyURL)
	mustContain(t, body, "3 of 3 gifts")
}

func TestValidationErrors(t *testing.T) {
	app := newTestApp(t)

	status, _, body := app.post("/p", url.Values{"name": {"   "}})
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("blank party name: status %d", status)
	}
	mustContain(t, body, "Uzupełnij to pole.")

	status, _, body = app.post("/p", url.Values{"name": {strings.Repeat("x", maxNameLen+1)}})
	if status != http.StatusUnprocessableEntity {
		t.Fatalf("long party name: status %d", status)
	}
	mustContain(t, body, "Ten tekst jest za długi.")
}

func TestLinks(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()

	party, _ := app.store.CreateParty(ctx, "A")
	person, _ := app.store.CreatePerson(ctx, party.ID, "Ania")
	if status, path, _ := app.get(giftsPath(person)); status != http.StatusOK || path != "/g/"+person.PublicID {
		t.Fatalf("person page: status %d, path %s", status, path)
	}

	// Long links from earlier versions, with the UUIDs they used, still lead to the right page.
	const legacyParty, legacyPerson = "461e316d-b1bc-4db2-a878-44e59493a12a", "b8d457ac-6a0f-421c-8280-b7f0b0c6d82b"
	db, err := sql.Open("sqlite", app.dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`UPDATE parties SET public_id = ? WHERE id = ?`, legacyParty, party.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`UPDATE people SET public_id = ? WHERE id = ?`, legacyPerson, person.ID); err != nil {
		t.Fatal(err)
	}
	if status, path, _ := app.get("/parties/" + legacyParty + "/people/new"); status != http.StatusOK || path != "/p/"+legacyParty {
		t.Errorf("legacy organizer link: status %d, path %s", status, path)
	}
	if status, path, body := app.get("/parties/" + legacyParty + "/people/" + legacyPerson + "/gifts"); status != http.StatusOK || path != "/g/"+legacyPerson {
		t.Errorf("legacy person link: status %d, path %s", status, path)
	} else {
		mustContain(t, body, "Ania")
	}

	for _, path := range []string{"/p/nope", "/g/nope", "/parties/nope/people/new"} {
		if status, _, _ := app.get(path); status != http.StatusNotFound {
			t.Errorf("%s: status %d, want 404", path, status)
		}
	}
}

func TestDeleteGiftOnlyOwn(t *testing.T) {
	app := newTestApp(t)
	ctx := context.Background()

	party, _ := app.store.CreateParty(ctx, "Party")
	ania, _ := app.store.CreatePerson(ctx, party.ID, "Ania")
	bartek, _ := app.store.CreatePerson(ctx, party.ID, "Bartek")
	if err := app.store.CreateGift(ctx, bartek.ID, "Socks"); err != nil {
		t.Fatal(err)
	}
	gifts, _ := app.store.Gifts(ctx, bartek.ID)

	// Ania can't delete Bartek's gift through her own page.
	_, _, body := app.post(giftsPath(ania)+"/"+itoa(gifts[0].ID)+"/delete", nil)
	mustContain(t, body, "Ten prezent został już usunięty.")
	if left, _ := app.store.Gifts(ctx, bartek.ID); len(left) != 1 {
		t.Fatalf("gift was deleted by another person")
	}

	_, _, body = app.post(giftsPath(bartek)+"/"+itoa(gifts[0].ID)+"/delete", nil)
	mustContain(t, body, "Usunięto prezent.")
	if left, _ := app.store.Gifts(ctx, bartek.ID); len(left) != 0 {
		t.Fatalf("gift was not deleted")
	}
}

func TestCrossOriginPostRejected(t *testing.T) {
	app := newTestApp(t)
	req, _ := http.NewRequest(http.MethodPost, app.srv.URL+"/p", strings.NewReader("name=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status %d, want 403", resp.StatusCode)
	}
}

func TestPlural(t *testing.T) {
	cases := []struct {
		lang string
		n    int
		want string
	}{
		{"pl", 1, "prezent"}, {"pl", 2, "prezenty"}, {"pl", 4, "prezenty"}, {"pl", 5, "prezentów"},
		{"pl", 12, "prezentów"}, {"pl", 22, "prezenty"}, {"pl", 0, "prezentów"},
		{"en", 1, "gift"}, {"en", 2, "gifts"},
	}
	for _, c := range cases {
		if got := plural(c.lang, c.n, "gift"); got != c.want {
			t.Errorf("plural(%s, %d) = %q, want %q", c.lang, c.n, got, c.want)
		}
	}
}

func TestLinkify(t *testing.T) {
	got := string(linkify("Scarf <b> https://www.shop.example/items/red-scarf-with-long-name?x=1, thanks"))
	want := `Scarf &lt;b&gt; <a href="https://www.shop.example/items/red-scarf-with-long-name?x=1" target="_blank" rel="noopener noreferrer" class="link">shop.example/items/red-scarf-with-lo…</a>, thanks`
	if got != want {
		t.Errorf("linkify:\n got %s\nwant %s", got, want)
	}
}

func itoa(n int64) string { return strconv.FormatInt(n, 10) }
