package server

import (
	"bytes"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"xmas-poll/internal/store"
	"xmas-poll/web"
)

const (
	maxNameLen   = 100
	maxGiftLen   = 1000
	maxBudgetLen = 50
	maxNoteLen   = 1000
	// A person sees who they are buying for only after adding this many gifts.
	minGifts = 3
)

type server struct {
	store *store.Store
	pages map[string]*template.Template
	log   *slog.Logger
}

// New returns the app's HTTP handler, serving pages backed by st.
func New(st *store.Store, log *slog.Logger) (http.Handler, error) {
	pages, err := parseTemplates()
	if err != nil {
		return nil, err
	}
	s := &server{store: st, pages: pages, log: log}
	return s.routes(), nil
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(web.FS))
	mux.HandleFunc("GET /favicon.ico", s.staticFile("static/favicon.svg"))
	mux.HandleFunc("GET /robots.txt", s.staticFile("static/robots.txt"))

	mux.HandleFunc("GET /{$}", s.home)
	mux.HandleFunc("POST /p", s.createParty)

	// Organizer page: /p/{party}.
	mux.HandleFunc("GET /p/{party}", s.showParty)
	mux.HandleFunc("POST /p/{party}/name", s.renameParty)
	mux.HandleFunc("POST /p/{party}/details", s.updateDetails)
	mux.HandleFunc("POST /p/{party}/people", s.addPerson)
	mux.HandleFunc("POST /p/{party}/people/{person}/name", s.renamePerson)
	mux.HandleFunc("POST /p/{party}/shuffle", s.shuffle)

	// A person's own page, shared with them by the organizer: /g/{person}.
	mux.HandleFunc("GET /g/{person}", s.showGifts)
	mux.HandleFunc("POST /g/{person}", s.addGift)
	mux.HandleFunc("POST /g/{person}/{gift}", s.updateGift)
	mux.HandleFunc("POST /g/{person}/{gift}/delete", s.deleteGift)

	mux.HandleFunc("/", s.notFound)

	return s.logRequests(http.NewCrossOriginProtection().Handler(limitBody(mux)))
}

// Page data

// page holds what the layout needs; every page's data embeds it.
type page struct {
	Title string
	Lang  string
	Path  string // canonical URL of the page, used by the language switch
	Flash *flash
}

type field struct {
	Name        string
	Label       string
	Value       string
	Placeholder string
	Help        string
	Error       string
	MaxLen      int
}

type homePage struct {
	page
	Name field
}

// Draw states on the organizer page.
const (
	drawTooFew  = "too-few" // fewer than two people
	drawReady   = "ready"   // nobody has a pair yet
	drawPartial = "partial" // someone added after the last draw
	drawDone    = "done"    // everyone has a pair
)

type partyPage struct {
	page
	partyForms
	Party    store.Party
	People   []store.Person
	Draw     string
	MinGifts int
}

// partyForms holds the organizer page's forms, filled with what was posted
// when one of them comes back with an error.
type partyForms struct {
	PartyName  field
	PersonName field
	Details    detailsForm
	// EditPerson is the ID of the person whose name is being edited, or 0.
	EditPerson int64
	EditName   field
}

type detailsForm struct {
	Budget field
	Date   field
	Note   field
	Open   bool // keep the form open: it came back with errors
}

type giftsPage struct {
	page
	Party         store.Party
	Person        store.Person
	Number        int // the person's door number, as on the organizer page
	Gifts         []store.Gift
	MinGifts      int
	GiftsLeft     int
	Receiver      *store.Person
	ReceiverGifts []store.Gift
	JustOpened    bool // the door opens with an animation right after the last required gift
	Autofocus     bool
	NewGift       field
	// EditGift is the ID of the gift being edited, or 0.
	EditGift     int64
	EditGiftText field
}

// Home and party creation

func (s *server) home(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	s.render(w, http.StatusOK, "home", homePage{
		page: s.page(w, r, lang, translate(lang, "New gift draw"), "/"),
		Name: partyNameField(lang, "", ""),
	})
}

func (s *server) createParty(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxNameLen)
	if problem != "" {
		s.render(w, http.StatusUnprocessableEntity, "home", homePage{
			page: s.page(w, r, lang, translate(lang, "New gift draw"), "/"),
			Name: partyNameField(lang, name, problem),
		})
		return
	}

	party, err := s.store.CreateParty(r.Context(), name)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, partyPath(party), http.StatusSeeOther)
}

// Organizer page

func (s *server) showParty(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}
	forms := defaultPartyForms(lang, party)
	forms.EditPerson = queryID(r, "edit") // renderParty fills in the name
	s.renderParty(w, r, http.StatusOK, lang, party, forms)
}

func defaultPartyForms(lang string, party store.Party) partyForms {
	return partyForms{
		PartyName:  partyNameField(lang, party.Name, ""),
		PersonName: personNameField(lang, "", ""),
		Details:    newDetailsForm(lang, party.PartyDetails, map[string]string{}),
	}
}

func (s *server) renameParty(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}

	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxNameLen)
	if problem != "" {
		forms := defaultPartyForms(lang, party)
		forms.PartyName = partyNameField(lang, name, problem)
		s.renderParty(w, r, http.StatusUnprocessableEntity, lang, party, forms)
		return
	}
	if err := s.store.RenameParty(r.Context(), party.ID, name); err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, partyPath(party), "info", translate(lang, "Party renamed."))
}

func (s *server) addPerson(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}

	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxNameLen)
	if problem != "" {
		forms := defaultPartyForms(lang, party)
		forms.PersonName = personNameField(lang, name, problem)
		s.renderParty(w, r, http.StatusUnprocessableEntity, lang, party, forms)
		return
	}
	if _, err := s.store.CreatePerson(r.Context(), party.ID, name); err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, partyPath(party), "info", translate(lang, "Added %s.", name))
}

func (s *server) renamePerson(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}
	personID, err := strconv.ParseInt(r.PathValue("person"), 10, 64)
	if err != nil {
		s.notFound(w, r)
		return
	}

	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxNameLen)
	if problem != "" {
		forms := defaultPartyForms(lang, party)
		forms.EditPerson = personID
		forms.EditName = personNameField(lang, name, problem)
		s.renderParty(w, r, http.StatusUnprocessableEntity, lang, party, forms)
		return
	}
	// The store only renames people in this party.
	err = s.store.RenamePerson(r.Context(), party.ID, personID, name)
	if errors.Is(err, store.ErrNotFound) {
		s.notFound(w, r)
		return
	}
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, partyPath(party), "info", translate(lang, "Saved: %s.", name))
}

func (s *server) updateDetails(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}
	if _, ok := formValue(w, r, "budget"); !ok {
		return
	}

	details, problems := validateDetails(r.PostForm.Get("budget"), r.PostForm.Get("date"), r.PostForm.Get("note"))
	if len(problems) > 0 {
		forms := defaultPartyForms(lang, party)
		forms.Details = newDetailsForm(lang, details, problems)
		s.renderParty(w, r, http.StatusUnprocessableEntity, lang, party, forms)
		return
	}
	if err := s.store.UpdatePartyDetails(r.Context(), party.ID, details); err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, partyPath(party), "info", translate(lang, "Party details saved."))
}

// validateDetails trims the party details and checks them. Every field may be
// left empty. Problems are keyed by field name.
func validateDetails(budget, date, note string) (store.PartyDetails, map[string]string) {
	problems := map[string]string{}
	d := store.PartyDetails{
		Budget:       strings.TrimSpace(budget),
		ExchangeDate: strings.TrimSpace(date),
		Note:         strings.TrimSpace(note),
	}
	if utf8.RuneCountInString(d.Budget) > maxBudgetLen {
		d.Budget = string([]rune(d.Budget)[:maxBudgetLen])
		problems["budget"] = "This is too long."
	}
	if utf8.RuneCountInString(d.Note) > maxNoteLen {
		d.Note = string([]rune(d.Note)[:maxNoteLen])
		problems["note"] = "This is too long."
	}
	if d.ExchangeDate != "" {
		if _, err := time.Parse(time.DateOnly, d.ExchangeDate); err != nil {
			d.ExchangeDate = ""
			problems["date"] = "Pick a date from the calendar."
		}
	}
	return d, problems
}

func (s *server) shuffle(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, ok := s.party(w, r, lang)
	if !ok {
		return
	}

	err := s.store.ShufflePeople(r.Context(), party.ID)
	if errors.Is(err, store.ErrAlreadyDrawn) {
		s.redirectWithFlash(w, r, partyPath(party), "error", translate(lang, "Everyone already has someone to buy for."))
		return
	}
	if errors.Is(err, store.ErrNotEnoughPeople) {
		s.redirectWithFlash(w, r, partyPath(party), "error", translate(lang, "Add at least two people before drawing."))
		return
	}
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, partyPath(party), "info", translate(lang, "Pairs drawn!"))
}

func (s *server) renderParty(w http.ResponseWriter, r *http.Request, status int, lang string, party store.Party, forms partyForms) {
	people, err := s.store.People(r.Context(), party.ID)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	if forms.EditPerson != 0 && forms.EditName.Name == "" {
		for _, p := range people {
			if p.ID == forms.EditPerson {
				forms.EditName = personNameField(lang, p.Name, "")
			}
		}
	}
	s.render(w, status, "party", partyPage{
		page:       s.page(w, r, lang, translate(lang, "Organizer"), partyPath(party)),
		partyForms: forms,
		Party:      party,
		People:     people,
		Draw:       drawState(people),
		MinGifts:   minGifts,
	})
}

func drawState(people []store.Person) string {
	if len(people) < 2 {
		return drawTooFew
	}
	paired := 0
	for _, p := range people {
		if p.ReceiverID.Valid {
			paired++
		}
	}
	switch paired {
	case 0:
		return drawReady
	case len(people):
		return drawDone
	default:
		return drawPartial
	}
}

// Person page

func (s *server) showGifts(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, person, ok := s.partyAndPerson(w, r, lang)
	if !ok {
		return
	}
	s.renderGifts(w, r, http.StatusOK, lang, party, person, giftForms{
		New:    giftField(lang, "", ""),
		EditID: queryID(r, "edit"),
	})
}

// giftForms holds the person page's forms: a new gift, and the gift being
// edited, if any.
type giftForms struct {
	New    field
	EditID int64
	Edit   field // empty until filled in from the gift, unless it came back with an error
}

func (s *server) addGift(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, person, ok := s.partyAndPerson(w, r, lang)
	if !ok {
		return
	}

	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxGiftLen)
	if problem != "" {
		s.renderGifts(w, r, http.StatusUnprocessableEntity, lang, party, person, giftForms{New: giftField(lang, name, problem)})
		return
	}
	if err := s.store.CreateGift(r.Context(), person.ID, name); err != nil {
		s.serverError(w, r, err)
		return
	}

	gifts, err := s.store.Gifts(r.Context(), person.ID)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	if len(gifts) == minGifts && person.ReceiverID.Valid {
		setShortCookie(w, revealCookie, "1")
	}
	s.redirectWithFlash(w, r, giftsPath(person), "info", translate(lang, "Gift added."))
}

func (s *server) updateGift(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	party, person, ok := s.partyAndPerson(w, r, lang)
	if !ok {
		return
	}
	giftID, err := strconv.ParseInt(r.PathValue("gift"), 10, 64)
	if err != nil {
		s.notFound(w, r)
		return
	}

	value, ok := formValue(w, r, "name")
	if !ok {
		return
	}
	name, problem := validateText(value, maxGiftLen)
	if problem != "" {
		s.renderGifts(w, r, http.StatusUnprocessableEntity, lang, party, person, giftForms{
			New:    giftField(lang, "", ""),
			EditID: giftID,
			Edit:   editGiftField(lang, name, problem),
		})
		return
	}
	// The store only changes the gift if it belongs to this person.
	err = s.store.UpdateGift(r.Context(), person.ID, giftID, name)
	if errors.Is(err, store.ErrNotFound) {
		s.redirectWithFlash(w, r, giftsPath(person), "error", translate(lang, "This gift was already deleted."))
		return
	}
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, giftsPath(person), "info", translate(lang, "Gift saved."))
}

func (s *server) deleteGift(w http.ResponseWriter, r *http.Request) {
	lang := locale(w, r)
	_, person, ok := s.partyAndPerson(w, r, lang)
	if !ok {
		return
	}

	giftID, err := strconv.ParseInt(r.PathValue("gift"), 10, 64)
	if err != nil {
		s.notFound(w, r)
		return
	}
	// The store only deletes the gift if it belongs to this person.
	err = s.store.DeleteGift(r.Context(), person.ID, giftID)
	if errors.Is(err, store.ErrNotFound) {
		s.redirectWithFlash(w, r, giftsPath(person), "error", translate(lang, "This gift was already deleted."))
		return
	}
	if err != nil {
		s.serverError(w, r, err)
		return
	}
	s.redirectWithFlash(w, r, giftsPath(person), "info", translate(lang, "Gift deleted."))
}

func (s *server) renderGifts(w http.ResponseWriter, r *http.Request, status int, lang string, party store.Party, person store.Person, forms giftForms) {
	ctx := r.Context()
	gifts, err := s.store.Gifts(ctx, person.ID)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	number, err := s.store.PersonNumber(ctx, person)
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	pg := s.page(w, r, lang, translate(lang, "Your wishlist"), giftsPath(person))
	data := giftsPage{
		page:      pg,
		Party:     party,
		Person:    person,
		Number:    number,
		Gifts:     gifts,
		MinGifts:  minGifts,
		GiftsLeft: max(0, minGifts-len(gifts)),
		NewGift:   forms.New,
		// Keep the keyboard up while the person is still filling the required gifts.
		Autofocus: pg.Flash != nil && len(gifts) < minGifts,
	}
	for _, g := range gifts {
		if g.ID == forms.EditID {
			data.EditGift = g.ID
			data.EditGiftText = forms.Edit
			if data.EditGiftText.Name == "" {
				data.EditGiftText = editGiftField(lang, g.Name, "")
			}
		}
	}

	if data.GiftsLeft == 0 && person.ReceiverID.Valid {
		receiver, err := s.store.PersonByID(ctx, person.ReceiverID.V)
		if err != nil {
			s.serverError(w, r, err)
			return
		}
		data.Receiver = &receiver
		if data.ReceiverGifts, err = s.store.Gifts(ctx, receiver.ID); err != nil {
			s.serverError(w, r, err)
			return
		}
		data.JustOpened = popShortCookie(w, r, revealCookie) != ""
		if !person.DrawSeen {
			if err := s.store.MarkDrawSeen(ctx, person.ID); err != nil {
				s.serverError(w, r, err)
				return
			}
		}
	}
	s.render(w, status, "gifts", data)
}

// Lookups

// party loads the party named in the URL, writing a 404 response if there is none.
func (s *server) party(w http.ResponseWriter, r *http.Request, lang string) (store.Party, bool) {
	party, err := s.store.PartyByPublicID(r.Context(), r.PathValue("party"))
	if errors.Is(err, store.ErrNotFound) {
		s.renderNotFound(w, r, lang)
		return store.Party{}, false
	}
	if err != nil {
		s.serverError(w, r, err)
		return store.Party{}, false
	}
	return party, true
}

func (s *server) partyAndPerson(w http.ResponseWriter, r *http.Request, lang string) (store.Party, store.Person, bool) {
	person, err := s.store.PersonByPublicID(r.Context(), r.PathValue("person"))
	if errors.Is(err, store.ErrNotFound) {
		s.renderNotFound(w, r, lang)
		return store.Party{}, store.Person{}, false
	}
	if err != nil {
		s.serverError(w, r, err)
		return store.Party{}, store.Person{}, false
	}
	party, err := s.store.PartyByID(r.Context(), person.PartyID)
	if err != nil {
		s.serverError(w, r, err)
		return store.Party{}, store.Person{}, false
	}
	return party, person, true
}

// Forms

// maxFormBytes caps request bodies. The longest field, a 1000-character gift,
// is at most about 12 KB even with every byte percent-encoded.
const maxFormBytes = 32 << 10

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
		next.ServeHTTP(w, r)
	})
}

// formValue returns a field of the posted form. When the form can't be read,
// it answers 413 (too large) or 400 itself and returns false.
func formValue(w http.ResponseWriter, r *http.Request, key string) (string, bool) {
	if err := r.ParseForm(); err != nil {
		if _, tooLarge := errors.AsType[*http.MaxBytesError](err); tooLarge {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		return "", false
	}
	return r.PostForm.Get(key), true
}

// validateText trims a form value and returns it along with a problem
// description, which is empty when the value is valid. A value that is too
// long comes back cut to maxLen, so the form never echoes more than that.
func validateText(value string, maxLen int) (string, string) {
	value = strings.TrimSpace(value)
	switch {
	case value == "":
		return value, "Please fill this in."
	case utf8.RuneCountInString(value) > maxLen:
		return string([]rune(value)[:maxLen]), "This is too long."
	}
	return value, ""
}

func partyNameField(lang, value, problem string) field {
	return field{
		Name:        "name",
		Label:       translate(lang, "Party name"),
		Placeholder: translate(lang, "e.g. Christmas Eve at the Smiths'"),
		Value:       value,
		Error:       translate(lang, problem),
		MaxLen:      maxNameLen,
	}
}

func personNameField(lang, value, problem string) field {
	return field{
		Name:   "name",
		Label:  translate(lang, "Name"),
		Value:  value,
		Error:  translate(lang, problem),
		MaxLen: maxNameLen,
	}
}

func giftField(lang, value, problem string) field {
	return field{
		Name:        "name",
		Label:       translate(lang, "What would you like?"),
		Placeholder: translate(lang, "e.g. a warm scarf, a book, or a link to a shop"),
		Help:        translate(lang, "One gift at a time. Shop links become clickable."),
		Value:       value,
		Error:       translate(lang, problem),
		MaxLen:      maxGiftLen,
	}
}

func editGiftField(lang, value, problem string) field {
	return field{
		Name:   "name",
		Label:  translate(lang, "Edit gift"),
		Value:  value,
		Error:  translate(lang, problem),
		MaxLen: maxGiftLen,
	}
}

func newDetailsForm(lang string, d store.PartyDetails, problems map[string]string) detailsForm {
	return detailsForm{
		Budget: field{
			Name:        "budget",
			Label:       translate(lang, "Budget per gift"),
			Placeholder: translate(lang, "e.g. 100 zł"),
			Value:       d.Budget,
			Error:       translate(lang, problems["budget"]),
			MaxLen:      maxBudgetLen,
		},
		Date: field{
			Name:  "date",
			Label: translate(lang, "Gift exchange day"),
			Value: d.ExchangeDate,
			Error: translate(lang, problems["date"]),
		},
		Note: field{
			Name:        "note",
			Label:       translate(lang, "Message for everyone"),
			Placeholder: translate(lang, "e.g. We meet at 6 pm. Wrap your gift and write the name on it."),
			Value:       d.Note,
			Error:       translate(lang, problems["note"]),
			MaxLen:      maxNoteLen,
		},
		Open: len(problems) > 0,
	}
}

// queryID reads a positive integer ID from the query string, or returns 0.
func queryID(r *http.Request, key string) int64 {
	id, err := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	if err != nil || id < 0 {
		return 0
	}
	return id
}

func partyPath(p store.Party) string {
	return "/p/" + p.PublicID
}

func giftsPath(person store.Person) string {
	return "/g/" + person.PublicID
}

// Short-lived cookies carry a flash message, or the "door just opened" signal,
// across one redirect.

const (
	flashCookie  = "flash"
	revealCookie = "reveal"
)

type flash struct {
	Kind    string // "info" or "error"
	Message string
}

func setShortCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		Path:     "/",
		MaxAge:   60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// popShortCookie returns the cookie's value and clears it.
func popShortCookie(w http.ResponseWriter, r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{Name: name, Path: "/", MaxAge: -1})
	value, err := url.QueryUnescape(c.Value)
	if err != nil {
		return ""
	}
	return value
}

func (s *server) redirectWithFlash(w http.ResponseWriter, r *http.Request, to, kind, message string) {
	setShortCookie(w, flashCookie, kind+"|"+message)
	http.Redirect(w, r, to, http.StatusSeeOther)
}

// popFlash returns the pending flash message, if any, and clears it.
func popFlash(w http.ResponseWriter, r *http.Request) *flash {
	kind, message, _ := strings.Cut(popShortCookie(w, r, flashCookie), "|")
	if kind != "info" && kind != "error" {
		return nil
	}
	return &flash{Kind: kind, Message: message}
}

// Responses

func (s *server) page(w http.ResponseWriter, r *http.Request, lang, title, path string) page {
	return page{Title: title, Lang: lang, Path: path, Flash: popFlash(w, r)}
}

func (s *server) render(w http.ResponseWriter, status int, name string, data any) {
	tmpl, ok := s.pages[name]
	if !ok {
		s.log.Error("unknown template", "name", name)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Render into a buffer first so a template error doesn't leave a half-written page.
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		s.log.Error("render template", "name", name, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Pages show private links, wishlists and draws: keep them out of every cache.
	w.Header().Set("Cache-Control", "private, no-store")
	if status == http.StatusUnprocessableEntity {
		// A form re-rendered with errors: keep the form page's URL in the address bar
		// instead of the URL the form posted to.
		w.Header().Set("HX-Push-Url", "false")
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (s *server) notFound(w http.ResponseWriter, r *http.Request) {
	s.renderNotFound(w, r, locale(w, r))
}

func (s *server) renderNotFound(w http.ResponseWriter, r *http.Request, lang string) {
	s.render(w, http.StatusNotFound, "error", page{Title: translate(lang, "Not found"), Lang: lang, Path: "/"})
}

func (s *server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Error("request failed", "route", route(r), "err", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *server) staticFile(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.FS, name)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// route names the matched route for logs, such as "GET /p/{party}". Real paths
// are never logged: the IDs in them are the keys to organizer and person pages.
func route(r *http.Request) string {
	if r.Pattern == "" {
		return r.Method + " (unmatched)"
	}
	return r.Pattern
}

func (s *server) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		s.log.Info("request",
			"route", route(r), // set by the ServeMux on this same request
			"status", rec.status,
			"duration", time.Since(start).Round(time.Microsecond),
		)
	})
}
