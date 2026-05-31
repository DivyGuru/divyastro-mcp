package main

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV05Tools wires 10 new tools added in v0.5:
//
//   Composite Panchang (2):
//     panchang_basic      — all 5 elements + rise/set times in one call
//     panchang_advanced   — full panchang with muhurta, choghadiya, disha shool,
//                           samvatsara, same-day transitions
//
//   Panchang extras (3):
//     panchang_lagna_table      — 12 lagna rise-times for the day
//     panchang_monthly          — month-calendar (tithi/nak/yoga/vara per day)
//     panchang_nakshatra_prediction — daily Tara-Bala based prediction
//
//   Birth chart composite (1):
//     chart_astro_details  — avakhada + planets + birth panchang in one call
//
//   Varshaphal (2):
//     varshaphal_harsha_bala — Tajika planetary strengths
//     varshaphal_mudda_dasha — annual dasha periods
//
//   Yearly horoscope (1):
//     narrative_yearly_bhavishyafal — kundli-based yearly horoscope (6000+ rules)
//
//   Geo (1):
//     geo_reverse — lat/lon → nearest city name + timezone
//
// Total new: 10 tools (v0.5 brings total to 136 tools).

func registerV05Tools(s *mcp.Server, c *apiClient) {
	registerPanchangBasic(s, c)
	registerPanchangAdvanced(s, c)
	registerPanchangLagnaTable(s, c)
	registerPanchangMonthly(s, c)
	registerPanchangNakshatraPrediction(s, c)
	registerChartAstroDetails(s, c)
	registerVarshaphalHarshaBala(s, c)
	registerVarshaphalMuddaDasha(s, c)
	registerNarrativeYearlyBhavishyafal(s, c)
	registerGeoReverse(s, c)
}

// =====================================================================
// 1. /v1/panchang/basic — composite: 5 elements + rise/set (3 credits)
// =====================================================================

func registerPanchangBasic(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_basic",
		Description: "Get the complete basic panchang for any date and location in one call: tithi (with deity, type, end-time), nakshatra (with ruling planet, deity, gana, end-time), yoga (with auspiciousness), karana (2 slots with deity and special note), vara (weekday + lord), sunrise, sunset, moonrise, moonset. Replaces 7 separate atomic calls. Use for 'today's panchang', 'what is the tithi', 'when does nakshatra change', 'panchang for Mumbai on 15 June'.",
		Title:       "Panchang Basic (Composite)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/basic", in.toQuery())
	})
}

// =====================================================================
// 2. /v1/panchang/advanced — full daily panchang (5 credits)
// =====================================================================

func registerPanchangAdvanced(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_advanced",
		Description: "Get the complete advanced panchang — everything in basic plus: all muhurta windows (Rahu Kaal, Yamaganda, Gulika, Abhijit, Brahma Muhurat, Durmuhurta), Varjyam, Amrit Kalam, full 16-slot Choghadiya (day + night), Hindu month with Amanta/Purnimanta, Vikram Samvat + Samvatsara name, Ritu (season), Ayana, Sun/Moon sign, Disha Shool, Nak Shool, Moon Nivas, active Panchang Yogas (Amrit Siddhi etc.), and same-day element transitions (next tithi if it changes today). Best single call for muhurta apps. Use for 'full daily panchang', 'Rahu Kaal today', 'choghadiya for today', 'is Amrit Siddhi yoga today', 'disha shool for Monday'.",
		Title:       "Panchang Advanced (Full)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/advanced", in.toQuery())
	})
}

// =====================================================================
// 3. /v1/panchang/lagna-table — 12 lagna rise-times (2 credits)
// =====================================================================

func registerPanchangLagnaTable(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_lagna_table",
		Description: "Get the lagna (ascendant sign) table for a day — shows when each of the 12 zodiac signs rises on the eastern horizon with precise start/end times. Each lagna lasts roughly 1.5–2.5 hours. Used for muhurta selection (choosing an auspicious rising sign for an event). Use for 'when does Taurus rise today', 'lagna table for Mumbai', 'which lagna is rising at 4pm', 'auspicious lagna for marriage'.",
		Title:       "Panchang Lagna Table",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/lagna-table", in.toQuery())
	})
}

// =====================================================================
// 4. /v1/panchang/monthly — month calendar (5 credits)
// =====================================================================

// PanchangMonthlyInput uses TzOffsetHours (a decimal number like 5.5 for IST)
// rather than an IANA timezone string. This is intentional: the /panchang/monthly
// endpoint accepts only a numeric offset, not an IANA name. Use 5.5 for Asia/Kolkata,
// -5.0 for America/New_York, 0.0 for UTC. Do NOT pass "Asia/Kolkata" here.
type PanchangMonthlyInput struct {
	Lat          float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon          float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	TzOffsetHours float64 `json:"tz_offset_hours" jsonschema:"UTC offset as decimal hours — use a number, NOT an IANA string. Examples: 5.5 for India (IST +05:30), -5.0 for US Eastern (EST), 0.0 for UTC, 5.75 for Nepal (+05:45)"`
	Year         int     `json:"year" jsonschema:"Gregorian year (e.g. 2026)"`
	Month        int     `json:"month" jsonschema:"month number 1–12"`
}

func registerPanchangMonthly(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_monthly",
		Description: "Get a full month panchang calendar — one entry per day with: date, vara (weekday), tithi (name + paksha), nakshatra (name + pada), yoga, sunrise, sunset. Ideal for building calendar views. Use for 'panchang for June 2026', 'show all Ekadashi days this month', 'when is Purnima this month', 'calendar for July'.",
		Title:       "Panchang Monthly Calendar",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangMonthlyInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", strconv.FormatFloat(in.TzOffsetHours, 'f', -1, 64))
		q.Set("year", strconv.Itoa(in.Year))
		q.Set("month", strconv.Itoa(in.Month))
		return callPassthrough(ctx, c, "/v1/panchang/monthly", q)
	})
}

// =====================================================================
// 5. /v1/panchang/nakshatra-prediction — Tara Bala daily (1 credit)
// =====================================================================

type NakshatraPredictionInput struct {
	PanchangAtInput
	BirthNak int `json:"birth_nak" jsonschema:"birth (janma) nakshatra index 0–26 (0=Ashwini, 7=Pushya, 16=Anuradha, 26=Revati). Get from chart_nakshatra if unknown."`
}

func registerPanchangNakshatraPrediction(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_nakshatra_prediction",
		Description: "Get a personalised daily prediction based on Tara Bala — the position of today's Moon nakshatra relative to your birth nakshatra. Returns: today's Moon nakshatra, tara type (Janma/Sampat/Vipat/Kshema/Pratyak/Sadhana/Vadha/Mitra/Param-Mitra), auspiciousness score, and specific guidance for health, work, relationships, finance, and travel. Use for 'is today good for me', 'daily nakshatra prediction', 'tara bala today', 'should I travel today based on my nakshatra'.",
		Title:       "Daily Nakshatra Prediction (Tara Bala)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NakshatraPredictionInput) (*mcp.CallToolResult, any, error) {
		q := in.PanchangAtInput.toQuery()
		q.Set("birth_nak", strconv.Itoa(in.BirthNak))
		return callPassthrough(ctx, c, "/v1/panchang/nakshatra-prediction", q)
	})
}

// =====================================================================
// 6. /v1/chart/astro-details — composite birth chart (1 credit)
// =====================================================================

func registerChartAstroDetails(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_astro_details",
		Description: "Get the complete birth chart snapshot in one call: Avakhada Chakra (Varna, Vashya, Yoni, Gana, Nadi, Paya, Tatva), Ascendant (sign, nakshatra, pada), all 9 planet positions (sign, house, nakshatra, retrograde), Sun sign, Moon sign, and complete birth panchang (tithi with deity + type, nakshatra with lord + deity, yoga, vara). Replaces 4–5 separate calls. Use for 'full birth chart details', 'avakhada of my chart', 'birth panchang', 'kundli summary'.",
		Title:       "Astro Details (Full Birth Chart Snapshot)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/astro-details", in.toQuery())
	})
}

// =====================================================================
// 7. /v1/varshaphal/harsha-bala — Tajika planetary strength (1 credit)
// =====================================================================

func registerVarshaphalHarshaBala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_harsha_bala",
		Description: "Compute Tajika Harsha Bala (planetary strength) for the annual solar-return chart. Returns five components per planet: Stana Bala (positional), Uchcha Bala (exaltation proximity), Ooja-Yugma Bala (gender), Kendradi Bala (angular/succedent/cadent), and Drekkana Bala (decanate). Higher total = stronger planet in the annual chart. Use for 'which planet is strongest in my annual chart', 'harsha bala 2026', 'varshaphal planet strengths'.",
		Title:       "Varshaphal Harsha Bala",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/harsha-bala", q)
	})
}

// =====================================================================
// 8. /v1/varshaphal/mudda-dasha — annual dasha periods (1 credit)
// =====================================================================

func registerVarshaphalMuddaDasha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_mudda_dasha",
		Description: "Compute the Varshaphal Mudda Dasha (annual dasha) periods for the solar-return year. The 9 planets each rule a fixed number of days (Sun=11, Moon=24, Mars=6, Mercury=34, Jupiter=48, Venus=20, Saturn=19, Rahu=12, Ketu=6); the sequence starts from the Varsha Lord. Returns 18 periods covering the full year (2 cycles of 180 days each). Use for 'annual dasha 2027', 'mudda dasha periods', 'which planet period am I in this year'.",
		Title:       "Varshaphal Mudda Dasha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/mudda-dasha", q)
	})
}

// =====================================================================
// 9. /v1/vedic/narrative/yearly-bhavishyafal — yearly horoscope (1 credit)
// =====================================================================

type YearlyBhavishyafalInput struct {
	BirthInput
	Year int `json:"year" jsonschema:"target year for the horoscope (e.g. 2027)"`
}

func registerNarrativeYearlyBhavishyafal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_yearly_bhavishyafal",
		Description: "Generate a personalised yearly horoscope (Bhavishyafal) for any year using DivyAstroAPI's 6000+ authored Vedic rules. Returns up to 14 sections: Varshaphal themes (Muntha + Varsha Lord), active Vimshottari dasha at solar return, Sade Sati phase (if applicable), slow-mover transit narratives (Jupiter, Saturn, Rahu/Ketu, Mars), and life-area outlooks (marriage, career, finance). All anchored to the solar return JD — the true astrological new year. Far richer than generic sun-sign yearly horoscopes. Use for 'yearly horoscope 2027', 'annual prediction', 'what does 2027 look like for me', 'bhavishyafal for 2027'.",
		Title:       "Yearly Bhavishyafal (Personalised Annual Horoscope)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in YearlyBhavishyafalInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth.date", in.Date)
		q.Set("birth.time", in.Time)
		q.Set("birth.tz", in.Tz)
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/vedic/narrative/yearly-bhavishyafal", q)
	})
}

// =====================================================================
// 10. /v1/geo/reverse — lat/lon → nearest city (0 credits)
// =====================================================================

type GeoReverseInput struct {
	Lat   float64 `json:"lat" jsonschema:"latitude in decimal degrees (-90 to 90)"`
	Lon   float64 `json:"lon" jsonschema:"longitude in decimal degrees (-180 to 180)"`
	Limit int     `json:"limit,omitempty" jsonschema:"max results to return (1–10, default 5)"`
}

func registerGeoReverse(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "geo_reverse",
		Description: "Reverse geocoding — find the nearest city/place name for a given latitude and longitude. Returns place name, state, district, country code, and IANA timezone. Free (0 credits). Use this when you have GPS coordinates and need the city name or timezone. Use for 'what city is at 24.6°N 74.8°E', 'nearest city to these coordinates', 'what timezone is lat 19 lon 72'.",
		Title:       "Geo Reverse (Coordinates → City)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in GeoReverseInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		if in.Limit > 0 {
			q.Set("limit", strconv.Itoa(in.Limit))
		}
		return callPassthrough(ctx, c, "/v1/geo/reverse", q)
	})
}
