package main

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV10Tools closes the coverage gap found in the 2026-07-03 API↔MCP
// sync audit — endpoints added to the API after the v0.9 scan:
//
//   Char Dasha / Jaimini (5): /v1/dasha/char, /char/major, /char/current,
//     /char/{md}, /char/{md}/{ad}
//   Lal Kitab core (3): /v1/lalkitab/chart, /dasha, /dasha/current
//   Remedies (5): /v1/vedic/remedies/daan/{id}, /rudraksha/{id},
//     /yantra/{id}, /lal-kitab (list), /lal-kitab/{planet}/{house}
//
// Total new: 13 tools.
//
// NOTE: new tools send ?locale= (the API-wide documented param). The API
// also aliases ?lang= for backwards compatibility, so older v0.4-v0.8
// tools keep working.

func registerV10Tools(s *mcp.Server, c *apiClient) {
	// Char Dasha (Jaimini)
	registerDashaCharFull(s, c)
	registerDashaCharMajor(s, c)
	registerDashaCharCurrent(s, c)
	registerDashaCharMD(s, c)
	registerDashaCharMDAD(s, c)

	// Lal Kitab core
	registerLalKitabChart(s, c)
	registerLalKitabDasha(s, c)
	registerLalKitabDashaCurrent(s, c)

	// Remedy catalog additions
	registerRemedyDaanByID(s, c)
	registerRemedyRudrakshaByID(s, c)
	registerRemedyYantraByID(s, c)
	registerRemedyLalKitabList(s, c)
	registerRemedyLalKitabLookup(s, c)
}

// =====================================================================
// Char Dasha (Jaimini sign-based dasha) — 5 tools
// =====================================================================

func registerDashaCharFull(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_char_full",
		Description: "Get the full Char Dasha (Jaimini sign-based dasha) schedule for a birth chart — every mahadasha sign period with start/end dates across the full cycle. Use for 'Char Dasha periods', 'Jaimini dasha timeline', 'sign dasha schedule'.",
		Title:       "Char Dasha Full Schedule",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/char", in.toQuery())
	})
}

func registerDashaCharMajor(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_char_major",
		Description: "Get the Char Dasha (Jaimini) mahadasha-level periods only — one row per sign period with start/end dates. Lighter than dasha_char_full. Use for 'major Char Dasha periods', 'which sign mahadashas do I get and when'.",
		Title:       "Char Dasha Major Periods",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/char/major", in.toQuery())
	})
}

func registerDashaCharCurrent(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_char_current",
		Description: "Get the currently running Char Dasha (Jaimini) period stack for a birth chart — the active mahadasha sign and antardasha sign with their date windows. Use for 'which Char Dasha is running now', 'current Jaimini dasha'.",
		Title:       "Current Char Dasha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/char/current", in.toQuery())
	})
}

// CharDashaMDInput drills one level into a specific Char Dasha mahadasha.
type CharDashaMDInput struct {
	Lat  float64 `json:"lat" jsonschema:"latitude in decimal degrees, north positive (e.g. 19.0760 for Mumbai)"`
	Lon  float64 `json:"lon" jsonschema:"longitude in decimal degrees, east positive (e.g. 72.8777 for Mumbai)"`
	Date string  `json:"date" jsonschema:"birth date in YYYY-MM-DD format (e.g. 1990-01-15)"`
	Time string  `json:"time" jsonschema:"birth time in HH:MM 24-hour format (e.g. 10:30)"`
	Tz   string  `json:"tz" jsonschema:"IANA timezone of the birth (e.g. Asia/Kolkata)"`
	Md   string  `json:"md" jsonschema:"mahadasha sign name in English (e.g. aries, taurus, gemini, ... pisces)"`
}

func (in CharDashaMDInput) toQuery() url.Values {
	return BirthInput{Lat: in.Lat, Lon: in.Lon, Date: in.Date, Time: in.Time, Tz: in.Tz}.toQuery()
}

func registerDashaCharMD(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_char_antardashas",
		Description: "Drill into one Char Dasha (Jaimini) mahadasha sign and list its antardasha sub-periods with dates. Use for 'antardashas inside my Leo Char Dasha', 'sub-periods of a Jaimini sign dasha'.",
		Title:       "Char Dasha Antardashas",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CharDashaMDInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/char/"+url.PathEscape(in.Md), in.toQuery())
	})
}

// CharDashaMDADInput drills two levels into a Char Dasha mahadasha/antardasha.
type CharDashaMDADInput struct {
	Lat  float64 `json:"lat" jsonschema:"latitude in decimal degrees, north positive"`
	Lon  float64 `json:"lon" jsonschema:"longitude in decimal degrees, east positive"`
	Date string  `json:"date" jsonschema:"birth date in YYYY-MM-DD format"`
	Time string  `json:"time" jsonschema:"birth time in HH:MM 24-hour format"`
	Tz   string  `json:"tz" jsonschema:"IANA timezone of the birth (e.g. Asia/Kolkata)"`
	Md   string  `json:"md" jsonschema:"mahadasha sign name in English (e.g. leo)"`
	Ad   string  `json:"ad" jsonschema:"antardasha sign name in English (e.g. scorpio)"`
}

func (in CharDashaMDADInput) toQuery() url.Values {
	return BirthInput{Lat: in.Lat, Lon: in.Lon, Date: in.Date, Time: in.Time, Tz: in.Tz}.toQuery()
}

func registerDashaCharMDAD(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_char_pratyantardashas",
		Description: "Drill into a specific Char Dasha (Jaimini) mahadasha + antardasha pair and list the next-level sub-periods with dates. Use for 'sub-sub-periods of Leo-Scorpio Char Dasha'.",
		Title:       "Char Dasha Pratyantardashas",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CharDashaMDADInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/char/"+url.PathEscape(in.Md)+"/"+url.PathEscape(in.Ad), in.toQuery())
	})
}

// =====================================================================
// Lal Kitab core — 3 tools
// =====================================================================

func registerLalKitabChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "lalkitab_chart",
		Description: "Get the Lal Kitab chart for a birth — planets mapped into Lal Kitab's fixed Aries-based houses, with per-planet house placements used by Lal Kitab remedies. Use for 'Lal Kitab kundli', 'Lal Kitab planet positions'.",
		Title:       "Lal Kitab Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/lalkitab/chart", in.toQuery())
	})
}

func registerLalKitabDasha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "lalkitab_dasha",
		Description: "Get the full 35-year Lal Kitab dasha schedule for a birth — the planet-ruled periods unique to the Lal Kitab system with start/end dates. Use for 'Lal Kitab dasha periods', '35-year dasha cycle'.",
		Title:       "Lal Kitab Dasha Schedule",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/lalkitab/dasha", in.toQuery())
	})
}

func registerLalKitabDashaCurrent(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "lalkitab_dasha_current",
		Description: "Get the currently active Lal Kitab dasha planet and period window for a birth chart. Use for 'which Lal Kitab dasha is running', 'current Lal Kitab period'.",
		Title:       "Current Lal Kitab Dasha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/lalkitab/dasha/current", in.toQuery())
	})
}

// =====================================================================
// Remedy catalog additions — 5 tools
// =====================================================================

// RemedyByIDLocaleInput mirrors RemedyByIDInput but sends ?locale= (the
// API-wide documented param).
type RemedyByIDLocaleInput struct {
	ID     string `json:"id" jsonschema:"remedy item ID — discover IDs via the remedy list tools"`
	Locale string `json:"locale,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func (in RemedyByIDLocaleInput) toQuery() url.Values {
	q := url.Values{}
	if in.Locale != "" {
		q.Set("locale", in.Locale)
	}
	return q
}

func registerRemedyDaanByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_daan_by_id",
		Description: "Get the full details of a specific daan (charity/donation) remedy by ID — what to donate, to whom, on which day/nakshatra, and the planetary affliction it addresses. Use for 'daan remedy details', 'what should I donate for Saturn'.",
		Title:       "Daan Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDLocaleInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/remedies/daan/"+url.PathEscape(in.ID), in.toQuery())
	})
}

func registerRemedyRudrakshaByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_rudraksha_by_id",
		Description: "Get the full details of a specific rudraksha remedy by ID — mukhi count, ruling planet, wearing procedure, day and mantra for energization. Use for 'which rudraksha for Mars', 'rudraksha remedy details'.",
		Title:       "Rudraksha Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDLocaleInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/remedies/rudraksha/"+url.PathEscape(in.ID), in.toQuery())
	})
}

func registerRemedyYantraByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_yantra_by_id",
		Description: "Get the full details of a specific yantra remedy by ID — the yantra's purpose, installation procedure, direction, day, and associated mantra. Use for 'yantra remedy details', 'which yantra for Rahu'.",
		Title:       "Yantra Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDLocaleInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/remedies/yantra/"+url.PathEscape(in.ID), in.toQuery())
	})
}

// LalKitabRemedyListInput filters the Lal Kitab remedy catalog.
type LalKitabRemedyListInput struct {
	Planet string `json:"planet,omitempty" jsonschema:"optional planet filter in English (sun, moon, mars, mercury, jupiter, venus, saturn, rahu, ketu)"`
	Locale string `json:"locale,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerRemedyLalKitabList(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_lalkitab_list",
		Description: "List Lal Kitab remedies — the simple, household-item upay (remedies) Lal Kitab is famous for, optionally filtered by planet. Use for 'Lal Kitab remedies for Saturn', 'list Lal Kitab upay'.",
		Title:       "Lal Kitab Remedy List",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in LalKitabRemedyListInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Planet != "" {
			q.Set("planet", in.Planet)
		}
		if in.Locale != "" {
			q.Set("locale", in.Locale)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/lal-kitab", q)
	})
}

// LalKitabRemedyLookupInput fetches remedies for a planet-in-house placement.
type LalKitabRemedyLookupInput struct {
	Planet string `json:"planet" jsonschema:"planet in English (sun, moon, mars, mercury, jupiter, venus, saturn, rahu, ketu)"`
	House  string `json:"house" jsonschema:"Lal Kitab house number 1-12"`
	Locale string `json:"locale,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerRemedyLalKitabLookup(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_lalkitab_for_placement",
		Description: "Get the Lal Kitab remedies for a specific planet-in-house placement (e.g. Saturn in house 8) — the exact upay Lal Kitab prescribes for that placement. Combine with lalkitab_chart to find placements first. Use for 'Lal Kitab remedy for Saturn in 8th house'.",
		Title:       "Lal Kitab Remedy for Placement",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in LalKitabRemedyLookupInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Locale != "" {
			q.Set("locale", in.Locale)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/lal-kitab/"+url.PathEscape(in.Planet)+"/"+url.PathEscape(in.House), q)
	})
}
