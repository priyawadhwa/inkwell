package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func render(t *testing.T, flags Flags) string {
	t.Helper()
	s := &server{flags: flags, revision: "test", now: func() time.Time { return time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC) }}
	srv := httptest.NewServer(s.routes())
	t.Cleanup(srv.Close)

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	return string(body)
}

// TestHomePage renders the dashboard with the flags we ship.
func TestHomePage(t *testing.T) {
	flags, err := loadFlags()
	if err != nil {
		t.Fatal(err)
	}
	body := render(t, flags)

	for _, want := range []string{
		"Good morning, Linky",
		`<span class="amount">$12,480.22</span>`,  // Everyday Checking
		`<span class="amount">$142,299.41</span>`, // total balance
		"Coral Café",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home page missing %q", want)
		}
	}
}

func TestCamouflageMode(t *testing.T) {
	body := render(t, Flags{CamouflageMode: true})

	if strings.Contains(body, `<span class="amount">`) {
		t.Error("camouflage mode rendered an unmasked amount")
	}
	if !strings.Contains(body, `data-amount="$12,480.22"`) {
		t.Error("camouflage mode lost the value needed to reveal the balance")
	}
	if !strings.Contains(body, "Camouflage on") {
		t.Error("camouflage mode badge missing")
	}
}

func TestInkCloudCardFreeze(t *testing.T) {
	if body := render(t, Flags{}); strings.Contains(body, `class="freeze"`) {
		t.Error("freeze button shown with the flag off")
	}
	if body := render(t, Flags{InkCloudCardFreeze: true}); !strings.Contains(body, `class="freeze"`) {
		t.Error("freeze button missing with the flag on")
	}
}

func TestCents(t *testing.T) {
	for c, want := range map[Cents]string{
		0:         "$0.00",
		675:       "$6.75",
		-8412:     "-$84.12",
		1248022:   "$12,480.22",
		914241900: "$9,142,419.00",
	} {
		if got := c.String(); got != want {
			t.Errorf("Cents(%d) = %q, want %q", c, got, want)
		}
	}
}

func TestVersion(t *testing.T) {
	s := &server{revision: "inkwell-00042-abc", now: time.Now}
	rec := httptest.NewRecorder()
	s.routes().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	if got := rec.Body.String(); !strings.Contains(got, "inkwell-00042-abc") {
		t.Errorf("/version = %q", got)
	}
}
