// Command inkwell serves the Inkwell Bank customer dashboard.
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

var page = template.Must(template.ParseFS(templateFS, "templates/index.html"))

type server struct {
	flags    Flags
	revision string
	now      func() time.Time
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	static, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(static)))
	mux.HandleFunc("GET /{$}", s.handleIndex)
	mux.HandleFunc("GET /version", s.handleVersion)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	return mux
}

type pageData struct {
	Customer Customer
	Flags    Flags
	Revision string
}

// Amount renders money, hidden behind the camouflage shimmer when
// CamouflageMode is on. The real value rides along so a tap can reveal it.
func (d pageData) Amount(c Cents) template.HTML {
	v := template.HTMLEscapeString(c.String())
	if d.Flags.CamouflageMode {
		return template.HTML(fmt.Sprintf(`<span class="amount camo" data-amount="%s" title="Tap to reveal">•••••</span>`, v))
	}
	return template.HTML(fmt.Sprintf(`<span class="amount">%s</span>`, v))
}

func (s *server) handleIndex(w http.ResponseWriter, _ *http.Request) {
	data := pageData{Customer: demoCustomer(s.now()), Flags: s.flags, Revision: s.revision}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, data); err != nil {
		slog.Error("rendering index", "err", err)
	}
}

// handleVersion lets open browsers notice a new deploy and reload themselves.
func (s *server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]string{"revision": s.revision})
}

// revision prefers Cloud Run's revision name, which changes on every deploy,
// and falls back to the VCS commit Go stamped into the binary.
func revision() string {
	if r := os.Getenv("K_REVISION"); r != "" {
		return r
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return s.Value[:7]
			}
		}
	}
	return "dev"
}

func main() {
	flags, err := loadFlags()
	if err != nil {
		slog.Error("loading flags", "err", err)
		os.Exit(1)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	s := &server{flags: flags, revision: revision(), now: time.Now}
	slog.Info("serving", "port", port, "revision", s.revision, "flags", flags)
	if err := http.ListenAndServe(":"+port, s.routes()); err != nil {
		slog.Error("serving", "err", err)
		os.Exit(1)
	}
}
