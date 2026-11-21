package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/url"
	"path"
	"regexp"
	"strings"

	"xmas-poll/web"
)

// parseTemplates builds one template set per file in templates/pages,
// each combined with the shared layout.
func parseTemplates() (map[string]*template.Template, error) {
	icons, err := loadIcons()
	if err != nil {
		return nil, err
	}
	funcs := template.FuncMap{
		"t":       translate,
		"plural":  plural,
		"linkify": linkify,
		"inc":     func(i int) int { return i + 1 },
		"pips":    pips,
		"icon": func(name string) (template.HTML, error) {
			svg, ok := icons[name]
			if !ok {
				return "", fmt.Errorf("unknown icon %q", name)
			}
			return svg, nil
		},
	}

	files, err := fs.Glob(web.FS, "templates/pages/*.html")
	if err != nil {
		return nil, err
	}
	pages := make(map[string]*template.Template, len(files))
	for _, file := range files {
		tmpl, err := template.New("").Funcs(funcs).ParseFS(web.FS, "templates/layout.html", file)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", file, err)
		}
		pages[strings.TrimSuffix(path.Base(file), ".html")] = tmpl
	}
	return pages, nil
}

var (
	svgComment = regexp.MustCompile(`<!--.*?-->\s*`)
	svgClass   = regexp.MustCompile(`class="[^"]*"`)
	svgSize    = regexp.MustCompile(`\s(width|height)="24"`)
	svgSpace   = regexp.MustCompile(`>\s+<`)
)

// loadIcons reads the Lucide icons in templates/icons, ready to inline.
// They are sized by the .icon class (1.25em) and hidden from screen readers;
// buttons that only show an icon carry their own aria-label.
func loadIcons() (map[string]template.HTML, error) {
	files, err := fs.Glob(web.FS, "templates/icons/*.svg")
	if err != nil {
		return nil, err
	}
	icons := make(map[string]template.HTML, len(files))
	for _, file := range files {
		b, err := fs.ReadFile(web.FS, file)
		if err != nil {
			return nil, err
		}
		svg := svgComment.ReplaceAllString(string(b), "")
		svg = svgSize.ReplaceAllString(svg, "")
		svg = svgClass.ReplaceAllString(svg, `class="icon" aria-hidden="true" focusable="false"`)
		svg = svgSpace.ReplaceAllString(strings.Join(strings.Fields(svg), " "), "><")
		icons[strings.TrimSuffix(path.Base(file), ".svg")] = template.HTML(svg)
	}
	return icons, nil
}

var urlRe = regexp.MustCompile(`https?://[^\s<>"']+`)

// linkify escapes text and turns http(s) URLs in it into links that open in a
// new tab, shown as host plus a shortened path.
func linkify(text string) template.HTML {
	var b strings.Builder
	last := 0
	for _, m := range urlRe.FindAllStringIndex(text, -1) {
		b.WriteString(template.HTMLEscapeString(text[last:m[0]]))
		raw := strings.TrimRight(text[m[0]:m[1]], ".,;:!?)")
		last = m[0] + len(raw)

		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			b.WriteString(template.HTMLEscapeString(raw))
			continue
		}
		label := strings.TrimPrefix(u.Host, "www.")
		if p := strings.TrimSuffix(u.Path, "/"); p != "" {
			if r := []rune(p); len(r) > 24 {
				p = string(r[:24]) + "…"
			}
			label += p
		}
		fmt.Fprintf(&b, `<a href="%s" target="_blank" rel="noopener noreferrer" class="link">%s</a>`,
			template.HTMLEscapeString(u.String()), template.HTMLEscapeString(label))
	}
	b.WriteString(template.HTMLEscapeString(text[last:]))
	return template.HTML(b.String())
}

// pips returns one mark per required gift: true for each gift already written.
func pips(count int) []bool {
	marks := make([]bool, minGifts)
	for i := range marks {
		marks[i] = i < count
	}
	return marks
}
