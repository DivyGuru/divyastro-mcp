package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// callPassthrough invokes a GET on the upstream API and returns its
// raw JSON body verbatim as a tool result. We use json.RawMessage
// (not map[string]any) so the helper works regardless of whether
// the response is an object, array, or scalar — some endpoints
// (transit/positions, chart/planets) return arrays.
func callPassthrough(ctx context.Context, c *apiClient, path string, q url.Values) (*mcp.CallToolResult, any, error) {
	var raw json.RawMessage
	if err := c.get(ctx, path, q, &raw); err != nil {
		return errResult(err)
	}
	// Re-marshal with indent so the AI client's text view is readable.
	var pretty any
	if err := json.Unmarshal(raw, &pretty); err != nil {
		return errResult(fmt.Errorf("decode upstream JSON: %w", err))
	}
	return rawJSONResult(pretty)
}

// registerAllTools wires every tool exposed by this MCP server.
//
// v0.1 (12 tools) shipped a curated subset for breadth across the
// catalog. v0.2 (30 additions, total 42) covers depth — the longer
// tail of compute the AI client wants when the user asks specific
// questions like "what's my D10", "find a vivah muhurta in June", or
// "is the next solar eclipse visible from Mumbai".
//
// AI assistants pick the right tool based on the description alone,
// so each Description is written like an English sentence describing
// what the user wants, not the API path.
//
// Adding a new tool: define an Input struct (json tags + jsonschema
// tags for description), write a handler with signature
// (ctx, *mcp.CallToolRequest, Input) → (*mcp.CallToolResult, Output, error),
// and call mcp.AddTool. The schema is auto-inferred from struct
// tags; no manual JSON Schema authoring.
func registerAllTools(server *mcp.Server, client *apiClient) {
	// v0.1 — the original 12 tools (panchang/today, chart basics,
	// dasha current, mangal/match, sade-sati, horoscope d/w, transits).
	registerPanchangToday(server, client)
	registerChartAscendant(server, client)
	registerChartPlanets(server, client)
	registerChartNavamsa(server, client)
	registerChartNakshatra(server, client)
	registerDashaCurrent(server, client)
	registerMangalDosha(server, client)
	registerAshtakootaTotal(server, client)
	registerSadeSati(server, client)
	registerHoroscopeDaily(server, client)
	registerHoroscopeWeekly(server, client)
	registerTransitPositions(server, client)

	// v0.2 — 30 additional tools across panchang depth, chart deep
	// (houses / aspects / dignity / divisionals / shadbala /
	// ashtakavarga / avakhada), dasha (full vimshottari + yogini),
	// milan deep (nadi / vivah phal / ashtakoota breakdown), muhurta
	// (vivah / naamkaran / best-time), varshaphal, numerology, eclipses,
	// festivals, planet-moments (retrograde / ingress), prashna,
	// western natal, monthly horoscope.
	registerV02Tools(server, client)

	// v0.3 — 30 Western astrology tools (natal depth, synastry, composite,
	// progressions, solar arc, returns, transits-to-natal, interpretation,
	// lots, midpoints, Hellenistic time lords, extended bodies, cosmobiology,
	// heliocentric, eclipses, ingresses, retrograde window, astrocartography).
	// Total: 72 tools.
	registerV03Tools(server, client)
}

// rawJSONResult turns an arbitrary JSON-decoded payload into an
// MCP tool result. We return the upstream JSON verbatim as a single
// text content block — AI assistants can parse JSON natively, and
// preserving the raw shape means schema changes on the server side
// don't need an MCP client patch.
func rawJSONResult(payload any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshal upstream payload: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}, nil, nil
}

// errResult formats an upstream API error as an MCP tool error so the
// AI client can surface it to the user. We mark IsError=true so the
// AI knows the call failed (vs returning empty data).
func errResult(err error) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
	}, nil, nil
}

// =====================================================================
// Tool inputs
// =====================================================================

// BirthInput is the standard birth-data quintet used by most chart
// and compatibility tools. Date is YYYY-MM-DD, time is HH:MM (24h),
// tz is an IANA timezone like "Asia/Kolkata". Lat/lon are in decimal
// degrees (north + / east +).
type BirthInput struct {
	Lat  float64 `json:"lat" jsonschema:"latitude in decimal degrees, north positive (e.g. 19.0760 for Mumbai)"`
	Lon  float64 `json:"lon" jsonschema:"longitude in decimal degrees, east positive (e.g. 72.8777 for Mumbai)"`
	Date string  `json:"date" jsonschema:"birth date in YYYY-MM-DD format (e.g. 1990-01-15)"`
	Time string  `json:"time" jsonschema:"birth time in HH:MM 24-hour format (e.g. 10:30)"`
	Tz   string  `json:"tz" jsonschema:"IANA timezone of the birth (e.g. Asia/Kolkata, America/New_York)"`
}

func (b BirthInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(b.Lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(b.Lon, 'f', -1, 64))
	q.Set("date", b.Date)
	q.Set("time", b.Time)
	q.Set("tz", b.Tz)
	return q
}

// =====================================================================
// 1. /v1/panchang/today — daily Hindu calendar
// =====================================================================

type PanchangTodayInput struct {
	Lat  float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon  float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	Tz   string  `json:"tz" jsonschema:"IANA timezone of the observer (e.g. Asia/Kolkata)"`
	Date string  `json:"date,omitempty" jsonschema:"optional date in YYYY-MM-DD; defaults to today in the observer's timezone"`
}

func registerPanchangToday(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_today",
		Description: "Get the Hindu daily panchang for a location: tithi (lunar day), vara (weekday), nakshatra (lunar mansion), yoga, karana, sunrise/sunset, and Rahu Kaal / Yamaganda / Gulika muhurat windows. Use for 'what's the panchang today', 'is today auspicious', 'when is Rahu Kaal'.",
		Title:       "Daily Panchang",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangTodayInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		return callPassthrough(ctx, c, "/v1/panchang/summary", q)
	})
}

// =====================================================================
// 2. /v1/chart/ascendant — birth lagna
// =====================================================================

func registerChartAscendant(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_ascendant",
		Description: "Compute the rising sign (lagna / ascendant) of a Vedic birth chart. Returns sign, degree, nakshatra and pada. Use for 'what is my lagna', 'rising sign for birth on X date'.",
		Title:       "Birth Chart Ascendant (Lagna)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/ascendant", in.toQuery())
	})
}

// =====================================================================
// 3. /v1/chart/planets — full planet placements
// =====================================================================

func registerChartPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_planets",
		Description: "Compute the placement of all 9 grahas (Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu) in a Vedic birth chart: sign, degree, nakshatra, house, retrograde flag. Use whenever someone asks 'where is my Saturn', 'what's my Moon sign', or wants the full birth chart.",
		Title:       "Birth Chart Planet Positions",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/planets", in.toQuery())
	})
}

// =====================================================================
// 4. /v1/chart/divisional/D9 — Navamsa
// =====================================================================

func registerChartNavamsa(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_navamsa",
		Description: "Compute the Navamsa (D9) divisional chart, used in Vedic astrology for marriage, dharma, and the inner spouse. Returns Navamsa lagna and each planet's Navamsa sign + house. Use for 'show my Navamsa', 'D9 chart', 'marriage chart'.",
		Title:       "Navamsa (D9) Divisional Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/divisional/D9", in.toQuery())
	})
}

// =====================================================================
// 5. /v1/chart/nakshatra — birth nakshatra
// =====================================================================

func registerChartNakshatra(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_nakshatra",
		Description: "Get the birth (Janma) nakshatra and pada — the lunar mansion the Moon was in at birth, used for naming, muhurta, and Vimshottari dasha. Returns nakshatra name (English + Devanagari), pada (1-4), and ruling planet. Use for 'what's my birth nakshatra', 'janma nakshatra'.",
		Title:       "Birth Nakshatra",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/nakshatra", in.toQuery())
	})
}

// =====================================================================
// 6. /v1/dasha/vimshottari/current — currently active dasha stack
// =====================================================================

type DashaCurrentInput struct {
	BirthInput
	At string `json:"at,omitempty" jsonschema:"optional date-time to evaluate the dasha at, in RFC3339 (e.g. 2026-05-10T12:00:00Z); defaults to now"`
}

func registerDashaCurrent(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_current",
		Description: "Get the currently running 5-level Vimshottari Mahadasha / Antardasha / Pratyantardasha / Sookshma / Prana stack for a birth chart. Use for 'what's my current mahadasha', 'which dasha am I in', 'what dasha is running on date X'.",
		Title:       "Current Vimshottari Dasha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaCurrentInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari/current", q)
	})
}

// =====================================================================
// 7. /v1/milan/mangal-dosha — Manglik check
// =====================================================================

func registerMangalDosha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "mangal_dosha",
		Description: "Check whether a birth chart has Mangal Dosha (Manglik condition) — Mars in a position that classically affects marriage. Returns presence flag, sub-rule (BPHS strict / modern North / from Moon), cancellation status, and severity. Use for 'am I Manglik', 'check Mars dosha', 'Manglik or not'.",
		Title:       "Mangal Dosha Check",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/mangal-dosha", in.toQuery())
	})
}

// =====================================================================
// 8. /v1/milan/ashtakoota/total — match-making compatibility score
// =====================================================================

type AshtakootaInput struct {
	// Bride
	BrideLat  float64 `json:"bride_lat" jsonschema:"bride latitude in decimal degrees"`
	BrideLon  float64 `json:"bride_lon" jsonschema:"bride longitude in decimal degrees"`
	BrideDate string  `json:"bride_date" jsonschema:"bride birth date YYYY-MM-DD"`
	BrideTime string  `json:"bride_time" jsonschema:"bride birth time HH:MM 24h"`
	BrideTz   string  `json:"bride_tz" jsonschema:"bride birth IANA timezone (e.g. Asia/Kolkata)"`
	// Groom
	GroomLat  float64 `json:"groom_lat" jsonschema:"groom latitude in decimal degrees"`
	GroomLon  float64 `json:"groom_lon" jsonschema:"groom longitude in decimal degrees"`
	GroomDate string  `json:"groom_date" jsonschema:"groom birth date YYYY-MM-DD"`
	GroomTime string  `json:"groom_time" jsonschema:"groom birth time HH:MM 24h"`
	GroomTz   string  `json:"groom_tz" jsonschema:"groom birth IANA timezone"`
}

func registerAshtakootaTotal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "match_making_score",
		Description: "Compute the Ashtakoota (Guna Milan) compatibility score between a bride and groom — the traditional 36-point Vedic match-making system covering varna, vashya, tara, yoni, graha-maitri, gana, bhakoot, and nadi koota. Returns total score and per-koota breakdown. Use for 'kundli matching', 'match score', 'Guna Milan'.",
		Title:       "Match Making (Ashtakoota / Guna Milan)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("bride_lat", strconv.FormatFloat(in.BrideLat, 'f', -1, 64))
		q.Set("bride_lon", strconv.FormatFloat(in.BrideLon, 'f', -1, 64))
		q.Set("bride_date", in.BrideDate)
		q.Set("bride_time", in.BrideTime)
		q.Set("bride_tz", in.BrideTz)
		q.Set("groom_lat", strconv.FormatFloat(in.GroomLat, 'f', -1, 64))
		q.Set("groom_lon", strconv.FormatFloat(in.GroomLon, 'f', -1, 64))
		q.Set("groom_date", in.GroomDate)
		q.Set("groom_time", in.GroomTime)
		q.Set("groom_tz", in.GroomTz)
		return callPassthrough(ctx, c, "/v1/milan/ashtakoota/total", q)
	})
}

// =====================================================================
// 9. /v1/transit/sade-sati — Saturn 7.5-year transit
// =====================================================================

func registerSadeSati(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "sade_sati",
		Description: "Determine whether a person is currently in Saturn's Sade Sati (the 7.5-year transit through the 12th, 1st, and 2nd house from natal Moon) and list past + upcoming Sade Sati windows. Use for 'am I in Sade Sati', 'when does Saturn finish', 'Saturn transit phase'.",
		Title:       "Sade Sati Status",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/sade-sati", in.toQuery())
	})
}

// =====================================================================
// 10. /v1/horoscope/daily — daily Moon-sign horoscope
// =====================================================================

type HoroscopeInput struct {
	Rashi string `json:"rashi" jsonschema:"Moon sign / rashi name in English (Aries, Taurus, ..., Pisces) or Hindi (Mesha, Vrishabha, ...). Case-insensitive."`
	Date  string `json:"date,omitempty" jsonschema:"optional YYYY-MM-DD; defaults to today UTC"`
	Lang  string `json:"lang,omitempty" jsonschema:"optional response language: en (default), hi, mr"`
}

func (h HoroscopeInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("rashi", h.Rashi)
	if h.Date != "" {
		q.Set("date", h.Date)
	}
	if h.Lang != "" {
		q.Set("lang", h.Lang)
	}
	return q
}

func registerHoroscopeDaily(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "horoscope_daily",
		Description: "Get the daily horoscope for a Vedic Moon sign (rashi) — covers career, relationships, health, finance themes for the day plus daily lucky factors (color, number, direction). Use for 'today's horoscope for Aries', 'rashifal for Vrishabha'.",
		Title:       "Daily Horoscope by Rashi",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in HoroscopeInput) (*mcp.CallToolResult, any, error) {
		// Normalize rashi to lowercase to match API expectations.
		in.Rashi = strings.ToLower(strings.TrimSpace(in.Rashi))
		return callPassthrough(ctx, c, "/v1/horoscope/daily", in.toQuery())
	})
}

// =====================================================================
// 11. /v1/horoscope/weekly — weekly Moon-sign horoscope
// =====================================================================

func registerHoroscopeWeekly(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "horoscope_weekly",
		Description: "Get the weekly horoscope for a Vedic Moon sign (rashi). Returns the 7-day outlook with major planetary themes, week-spanning transits, and a focus recommendation. Use for 'this week's horoscope', 'weekly rashifal'.",
		Title:       "Weekly Horoscope by Rashi",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in HoroscopeInput) (*mcp.CallToolResult, any, error) {
		in.Rashi = strings.ToLower(strings.TrimSpace(in.Rashi))
		return callPassthrough(ctx, c, "/v1/horoscope/weekly", in.toQuery())
	})
}

// =====================================================================
// 12. /v1/transit/positions — current planet transits
// =====================================================================

type TransitPositionsInput struct {
	At string `json:"at,omitempty" jsonschema:"optional date-time to evaluate transits at, RFC3339 (e.g. 2026-05-10T12:00:00Z); defaults to now"`
	Tz string `json:"tz,omitempty" jsonschema:"optional IANA timezone for date interpretation; defaults to UTC"`
}

func registerTransitPositions(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transits_now",
		Description: "Get the current sidereal positions of all 9 grahas — the live sky positions used to evaluate current transits, gochar effects, and dasha-transit interplay. Returns each planet's sign, degree, nakshatra, and retrograde flag at the requested moment. Use for 'where are the planets now', 'current transits', 'gochar today'.",
		Title:       "Current Planetary Transits",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TransitPositionsInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.At != "" {
			q.Set("at", in.At)
		}
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/transit/positions", q)
	})
}
