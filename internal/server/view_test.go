package server

import (
	"regexp"
	"slices"
	"testing"
)

var verbRe = regexp.MustCompile(`%[a-z]`)

func TestTranslationsKeepFormatVerbs(t *testing.T) {
	for lang, msgs := range translations {
		for en, tr := range msgs {
			want := verbRe.FindAllString(en, -1)
			got := verbRe.FindAllString(tr, -1)
			if !slices.Equal(want, got) {
				t.Errorf("%s: %q has verbs %v, translation %q has %v", lang, en, want, tr, got)
			}
		}
	}
}
