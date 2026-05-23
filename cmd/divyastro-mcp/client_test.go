package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestAPIClient_GetHappyPath verifies a 2xx JSON response decodes
// into the caller's struct, with auth + accept headers set correctly.
func TestAPIClient_GetHappyPath(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth + accept + UA headers.
		if got := r.Header.Get("Authorization"); got != "Bearer dv_test_xxx" {
			t.Errorf("Authorization = %q, want 'Bearer dv_test_xxx'", got)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want application/json", got)
		}
		if got := r.Header.Get("User-Agent"); !strings.HasPrefix(got, "divyastro-mcp/") {
			t.Errorf("User-Agent = %q, want divyastro-mcp/* prefix", got)
		}
		// Echo back a fixed payload.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"foo": "bar", "n": 42}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "dv_test_xxx", "test")
	var out struct {
		Foo string `json:"foo"`
		N   int    `json:"n"`
	}
	if err := c.get(context.Background(), "/v1/anything", nil, &out); err != nil {
		t.Fatalf("get: %v", err)
	}
	if out.Foo != "bar" || out.N != 42 {
		t.Errorf("decoded payload = %+v, want {Foo:bar N:42}", out)
	}
}

// TestAPIClient_GetSurfacesAPIErrors verifies that a 4xx with the
// API's structured error JSON produces a useful Go error message,
// not a generic "request failed".
func TestAPIClient_GetSurfacesAPIErrors(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid date","code":"invalid_param","details":"date must be YYYY-MM-DD"}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "dv_test_xxx", "test")
	var out map[string]any
	err := c.get(context.Background(), "/v1/whatever", nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"400", "invalid_param", "invalid date", "YYYY-MM-DD"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message %q missing %q", msg, want)
		}
	}
}

// TestAPIClient_GetEncodesQueryParams verifies the query string is
// built correctly — important because all our tools pass params via
// url.Values and a serialization bug would leak silently into
// production tool responses.
func TestAPIClient_GetEncodesQueryParams(t *testing.T) {
	t.Parallel()
	var seenQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "dv_test_xxx", "test")
	q := url.Values{}
	q.Set("lat", "19.076")
	q.Set("lon", "72.8777")
	q.Set("date", "1990-01-15")
	q.Set("tz", "Asia/Kolkata")
	if err := c.get(context.Background(), "/v1/chart/ascendant", q, nil); err != nil {
		t.Fatalf("get: %v", err)
	}
	for k, want := range map[string]string{
		"lat":  "19.076",
		"lon":  "72.8777",
		"date": "1990-01-15",
		"tz":   "Asia/Kolkata",
	} {
		if got := seenQuery.Get(k); got != want {
			t.Errorf("query[%q] = %q, want %q", k, got, want)
		}
	}
}

// TestAPIClient_GetAddsLeadingSlash ensures path normalization works
// — handlers can pass "v1/foo" or "/v1/foo" and both reach the right
// upstream URL.
func TestAPIClient_GetAddsLeadingSlash(t *testing.T) {
	t.Parallel()
	var seenPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "dv_test_xxx", "test")
	if err := c.get(context.Background(), "v1/no/leading/slash", nil, nil); err != nil {
		t.Fatalf("get: %v", err)
	}
	if seenPath != "/v1/no/leading/slash" {
		t.Errorf("seen path = %q, want /v1/no/leading/slash", seenPath)
	}
}

// TestAPIClient_GetTrimsTrailingSlashFromBaseURL ensures "https://api.x.com/"
// and "https://api.x.com" both work — common copy-paste mistake.
func TestAPIClient_GetTrimsTrailingSlashFromBaseURL(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL+"/", "dv_test_xxx", "test")
	if err := c.get(context.Background(), "/v1/x", nil, nil); err != nil {
		t.Fatalf("get: %v", err)
	}
}

// TestAPIClient_GetUnknownErrorBodyFallback verifies that non-JSON
// 4xx/5xx bodies still produce a useful error message (e.g. when
// nginx returns an HTML error page).
func TestAPIClient_GetUnknownErrorBodyFallback(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>502 Bad Gateway</html>"))
	}))
	defer srv.Close()

	c := newAPIClient(srv.URL, "dv_test_xxx", "test")
	err := c.get(context.Background(), "/v1/x", nil, nil)
	if err == nil {
		t.Fatal("expected error for 502, got nil")
	}
	if !strings.Contains(err.Error(), "502") {
		t.Errorf("error %q should mention 502", err.Error())
	}
}

// fakeAPIServer returns an httptest.Server that replies to every
// request with the supplied status + body. The captured *url.URL of
// each incoming request is appended to the returned *requestLog so
// tests can assert path + query values.
type requestLog struct {
	paths   []string
	queries []url.Values
}

func fakeAPIServer(t *testing.T, status int, body string) (*httptest.Server, *requestLog) {
	t.Helper()
	log := &requestLog{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.paths = append(log.paths, r.URL.Path)
		log.queries = append(log.queries, r.URL.Query())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, log
}

// canned upstream payload the tool-handler tests reuse.
const cannedPanchangResponse = `{
  "tithi": "Krishna Saptami",
  "nakshatra": "Mrigashira",
  "rahu_kaal": "10:30-12:00"
}`

// TestBirthInputToQuery covers the most-shared helper — a regression
// here would corrupt every chart-based tool simultaneously.
func TestBirthInputToQuery(t *testing.T) {
	t.Parallel()
	in := BirthInput{
		Lat: 19.076, Lon: 72.8777,
		Date: "1990-01-15", Time: "10:30",
		Tz: "Asia/Kolkata",
	}
	q := in.toQuery()
	for k, want := range map[string]string{
		"lat":  "19.076",
		"lon":  "72.8777",
		"date": "1990-01-15",
		"time": "10:30",
		"tz":   "Asia/Kolkata",
	} {
		if got := q.Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
}

// jsonHelper checks that a json.RawMessage decodes into the expected
// shape — used by tool tests that want to verify the upstream JSON
// flowed through unchanged.
func mustDecode(t *testing.T, raw string) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("decode test fixture: %v", err)
	}
	return out
}

func TestCannedFixtureValid(t *testing.T) {
	t.Parallel()
	got := mustDecode(t, cannedPanchangResponse)
	if got["tithi"] != "Krishna Saptami" {
		t.Errorf("fixture has wrong tithi: %v", got["tithi"])
	}
}
