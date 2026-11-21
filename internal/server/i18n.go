package server

import (
	"fmt"
	"net/http"
)

const (
	defaultLang  = "pl"
	localeCookie = "locale"
)

// translations maps a language to English message -> translated message.
// Messages are fmt format strings; translations must keep the same verbs in the same order.
var translations = map[string]map[string]string{
	"pl": {
		// Layout
		"Close":     "Zamknij",
		"Language":  "Język",
		"Copy link": "Kopiuj link",
		"Copied!":   "Skopiowano!",
		"Couldn't copy, long-press the link instead": "Nie udało się skopiować, przytrzymaj link",
		"Share": "Udostępnij",

		// Home
		"New gift draw": "Nowe losowanie",
		"Family gift draw, without the paper slips":                                                     "Rodzinne losowanie prezentów, bez karteczek",
		"Add your people, draw who buys for whom and send everyone their own link. No accounts needed.": "Dodaj bliskich, wylosuj, kto komu kupuje prezent, i wyślij każdemu jego link. Bez zakładania kont.",
		"Start a new draw":                  "Zacznij nowe losowanie",
		"Party name":                        "Nazwa przyjęcia",
		"e.g. Christmas Eve at the Smiths'": "np. Wigilia u Kowalskich",
		"Create the draw":                   "Utwórz losowanie",
		"How it works":                      "Jak to działa",
		"Add people":                        "Dodaj osoby",
		"Everyone taking part, yourself included.": "Wszystkich uczestników, także siebie.",
		"Draw pairs": "Wylosuj pary",
		"Everyone gets one person to buy for. Nobody draws themselves.": "Każdy dostaje jedną osobę, której kupuje prezent. Nikt nie wylosuje siebie.",
		"Send the links": "Wyślij linki",
		"Everyone writes a wishlist and sees what their person wants.": "Każdy pisze listę prezentów i widzi, czego chce jego osoba.",

		// Organizer
		"Organizer":              "Panel organizatora",
		"Not drawn yet":          "Jeszcze nie wylosowano",
		"Pairs drawn":            "Pary wylosowane",
		"Waiting for a new draw": "Czeka na ponowne losowanie",
		"This is your organizer page. Save its link to come back later.": "To Twój panel organizatora. Zapisz jego link, żeby tu wrócić.",
		"People": "Uczestnicy",
		"No one here yet. Add the first person below, yourself included.": "Nikogo tu jeszcze nie ma. Dodaj pierwszą osobę poniżej, także siebie.",
		"Each person's link appears here after the draw.":                 "Link dla każdej osoby pojawi się tu po losowaniu.",
		"Needs a new draw":                 "Potrzebne ponowne losowanie",
		"%d of %d gifts":                   "%d z %d prezentów",
		"Name":                             "Imię",
		"Add a person":                     "Dodaj osobę",
		"Add":                              "Dodaj",
		"Add at least two people to draw.": "Dodaj co najmniej dwie osoby, żeby losować.",
		"Ready to draw?":                   "Gotowi do losowania?",
		"Draw now":                         "Losuj pary",
		"Send everyone their link from the list above. Only they can see who they're buying for.": "Wyślij każdemu jego link z listy powyżej. Tylko ta osoba zobaczy, komu kupuje prezent.",
		"New people need a pair": "Nowe osoby nie mają pary",
		"Drawing again gives everyone a new person, so let anyone who already looked know.": "Ponowne losowanie zmieni pary wszystkim, więc daj znać osobom, które już sprawdziły.",
		"Draw again": "Losuj ponownie",
		"Draw again? Everyone will get a new person.": "Losować ponownie? Wszyscy dostaną nowe osoby.",
		"Rename the party": "Zmień nazwę przyjęcia",
		"Save name":        "Zapisz nazwę",
		"Party renamed.":   "Zmieniono nazwę.",
		"Added %s.":        "Dodano: %s.",
		"Pairs drawn!":     "Pary wylosowane!",
		"Everyone already has someone to buy for.": "Każdy ma już osobę, której kupuje prezent.",
		"Add at least two people before drawing.":  "Dodaj co najmniej dwie osoby przed losowaniem.",

		// Person page
		"Your wishlist":             "Twoja lista prezentów",
		"Hi, %s!":                   "Cześć, %s!",
		"%s · gift draw":            "%s · losowanie prezentów",
		"Your door is still closed": "Twoje drzwiczki są jeszcze zamknięte",
		"Add %d more %s and it opens to show who you're buying for.": "Dodaj jeszcze %d %s, a otworzą się i zobaczysz, komu kupujesz prezent.",
		"Your list is ready": "Twoja lista jest gotowa",
		"The organizer hasn't drawn the pairs yet. Check back here later.": "Organizator jeszcze nie wylosował par. Zajrzyj tu później.",
		"You're buying a gift for this person. Here's what they'd like:":   "Tej osobie kupujesz prezent. Oto, czego sobie życzy:",
		"Their list is still empty. Check back later.":                     "Ta lista jest jeszcze pusta. Zajrzyj później.",
		"Your list": "Twoja lista",
		"Nothing here yet. What would make you happy this year?": "Jeszcze nic tu nie ma. Co sprawiłoby Ci radość w tym roku?",
		"What would you like?":                             "Czego sobie życzysz?",
		"e.g. a warm scarf, a book, or a link to a shop":   "np. ciepły szalik, książka albo link do sklepu",
		"One gift at a time. Shop links become clickable.": "Jeden prezent naraz. Linki do sklepów będą klikalne.",
		"Add gift":                       "Dodaj prezent",
		"Delete gift":                    "Usuń prezent",
		"Delete this gift?":              "Usunąć ten prezent?",
		"Gift added.":                    "Dodano prezent.",
		"Gift deleted.":                  "Usunięto prezent.",
		"This gift was already deleted.": "Ten prezent został już usunięty.",

		// Errors
		"Please fill this in.":   "Uzupełnij to pole.",
		"This is too long.":      "Ten tekst jest za długi.",
		"Not found":              "Nie znaleziono",
		"This door doesn't open": "Te drzwiczki się nie otwierają",
		"We couldn't find this party or person. Check that you copied the whole link.": "Nie znaleźliśmy tego losowania ani osoby. Sprawdź, czy link został skopiowany w całości.",
		"Go to the start page": "Przejdź na stronę główną",
	},
}

// translate returns msg in lang (falling back to English) formatted with args.
func translate(lang, msg string, args ...any) string {
	if t, ok := translations[lang][msg]; ok {
		msg = t
	}
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// pluralForms holds, per word, the English singular and plural, and the three
// Polish forms: one, few (2-4, 22-24, ...) and many.
var pluralForms = map[string]struct{ en, pl []string }{
	"gift":   {en: []string{"gift", "gifts"}, pl: []string{"prezent", "prezenty", "prezentów"}},
	"person": {en: []string{"person", "people"}, pl: []string{"osoba", "osoby", "osób"}},
}

// plural returns the form of word that goes with the count n.
func plural(lang string, n int, word string) string {
	forms, ok := pluralForms[word]
	if !ok {
		return word
	}
	if lang != "pl" {
		if n == 1 {
			return forms.en[0]
		}
		return forms.en[1]
	}
	switch {
	case n == 1:
		return forms.pl[0]
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 12 || n%100 > 14):
		return forms.pl[1]
	default:
		return forms.pl[2]
	}
}

func supportedLang(lang string) bool {
	return lang == "en" || lang == "pl"
}

// locale picks the language from ?locale=, remembering it in a cookie, then
// from that cookie, then falls back to the default.
func locale(w http.ResponseWriter, r *http.Request) string {
	if lang := r.URL.Query().Get("locale"); supportedLang(lang) {
		http.SetCookie(w, &http.Cookie{
			Name:     localeCookie,
			Value:    lang,
			Path:     "/",
			MaxAge:   365 * 24 * 60 * 60,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		return lang
	}
	if c, err := r.Cookie(localeCookie); err == nil && supportedLang(c.Value) {
		return c.Value
	}
	return defaultLang
}
