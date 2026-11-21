// Package web holds the HTML templates, icons and static files, embedded into the binary.
package web

import "embed"

// FS contains templates/ (layout, pages, icons) and static/ (CSS, JS, fonts, favicon).
//
//go:embed templates static
var FS embed.FS
