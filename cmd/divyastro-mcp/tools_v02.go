package main

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV02Tools wires the 30 additional tools shipped in v0.2.
// Grouped by domain in the order they appear in the API catalog so
// reviewers can spot-check coverage at a glance:
//
//   - Panchang depth (5)
//   - Chart depth (7)
//   - Dasha depth (2)
//   - Milan depth (3)
//   - Muhurta (3)
//   - Varshaphal (1)
//   - Numerology (1)
//   - Eclipses (2)
//   - Festivals (1)
//   - Planet-moments (2)
//   - Prashna (1)
//   - Western natal (1)
//   - Horoscope (1)
//
// Total: 30, bringing the MCP server to 42 tools combined with v0.1.
func registerV02Tools(s *mcp.Server, c *apiClient) {
	// Panchang depth
	registerPanchangTithi(s, c)
	registerPanchangChoghadiya(s, c)
	registerPanchangHora(s, c)
	registerPanchangRahuKaal(s, c)
	registerPanchangBrahmaMuhurat(s, c)

	// Chart depth
	registerChartHouses(s, c)
	registerChartAspects(s, c)
	registerChartDignity(s, c)
	registerChartDivisional(s, c)
	registerChartAvakhada(s, c)
	registerChartShadbala(s, c)
	registerChartAshtakavarga(s, c)

	// Dasha depth
	registerDashaVimshottariFull(s, c)
	registerDashaYoginiCurrent(s, c)

	// Milan depth
	registerNadiDosha(s, c)
	registerVivahPhal(s, c)
	registerAshtakootaBreakdown(s, c)

	// Muhurta
	registerMuhurtaVivah(s, c)
	registerMuhurtaNaamkaran(s, c)
	registerMuhurtaBestTime(s, c)

	// Varshaphal
	registerVarshaphalChart(s, c)

	// Numerology
	registerNumerologyFull(s, c)

	// Eclipses
	registerEclipsesSolar(s, c)
	registerEclipsesLunar(s, c)

	// Festivals
	registerFestivalsMonth(s, c)

	// Planet-moments
	registerPlanetRetrogradeWindow(s, c)
	registerPlanetIngress(s, c)

	// Prashna
	registerPrashnaAnswer(s, c)

	// Western
	registerWesternNatalChart(s, c)

	// Horoscope
	registerHoroscopeMonthly(s, c)
}

// =====================================================================
// Shared input types for v0.2 tools
// =====================================================================

// PanchangAtInput is the standard panchang query: location + IANA
// timezone + optional date. Mirrors PanchangTodayInput but lets the
// caller pick a specific date.
type PanchangAtInput struct {
	Lat  float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon  float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	Tz   string  `json:"tz" jsonschema:"IANA timezone of the observer (e.g. Asia/Kolkata)"`
	Date string  `json:"date,omitempty" jsonschema:"optional date in YYYY-MM-DD; defaults to today in the observer's timezone"`
}

func (p PanchangAtInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(p.Lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(p.Lon, 'f', -1, 64))
	q.Set("tz", p.Tz)
	if p.Date != "" {
		q.Set("date", p.Date)
	}
	return q
}

// =====================================================================
// Panchang depth (5 tools)
// =====================================================================

func registerPanchangTithi(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_tithi",
		Description: "Get the lunar day (tithi) for a date and location — Krishna/Shukla paksha, tithi name (Pratipada/Dwitiya/.../Amavasya/Purnima), and end-of-tithi timestamp. Use for 'what tithi is today', 'is today Ekadashi'.",
		Title:       "Tithi (Lunar Day)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/tithi", in.toQuery())
	})
}

func registerPanchangChoghadiya(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_choghadiya",
		Description: "Get the day's Choghadiya muhurta windows — Amrit, Shubh, Labh (auspicious) and Rog, Kaal, Udveg (inauspicious). Returns 8 day-time + 8 night-time windows with start/end. Use for 'is now a good time', 'choghadiya for today'.",
		Title:       "Choghadiya Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/choghadiya", in.toQuery())
	})
}

func registerPanchangHora(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_hora",
		Description: "Get the planetary hora (planet-hour) schedule for a date — each hour ruled by Sun/Moon/Mars/Mercury/Jupiter/Venus/Saturn in the classical sequence. Use for 'best hora for business', 'planet hour right now'.",
		Title:       "Hora (Planetary Hours)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/hora", in.toQuery())
	})
}

func registerPanchangRahuKaal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_rahu_kaal",
		Description: "Get today's Rahu Kaal window — the inauspicious 90-minute span avoided for new ventures. Returns start/end timestamps in the observer's timezone. Use for 'when is Rahu Kaal today', 'avoid Rahu Kaal'.",
		Title:       "Rahu Kaal",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/rahu-kaal", in.toQuery())
	})
}

func registerPanchangBrahmaMuhurat(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_brahma_muhurat",
		Description: "Get the Brahma Muhurat window — the 96-minute span before sunrise considered most auspicious for sadhana, study, and meditation. Returns start/end timestamps. Use for 'when is Brahma Muhurat', 'best time to wake for meditation'.",
		Title:       "Brahma Muhurat",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/brahma-muhurat", in.toQuery())
	})
}

// =====================================================================
// Chart depth (7 tools)
// =====================================================================

func registerChartHouses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_houses",
		Description: "Get the 12 house cusps of a Vedic birth chart — sign on each house, lord, and degrees. Whole-sign houses by default. Use for 'what's in my 7th house', 'house cusps for my chart'.",
		Title:       "Birth Chart Houses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/houses", in.toQuery())
	})
}

func registerChartAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_aspects",
		Description: "Get all Vedic graha drishti (planetary aspects) on a birth chart — full + special aspects of Mars/Jupiter/Saturn/Rahu/Ketu, plus the 7th-house aspect of every planet. Use for 'who aspects my Sun', 'planet aspects in my chart'.",
		Title:       "Birth Chart Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/aspects", in.toQuery())
	})
}

func registerChartDignity(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_dignity",
		Description: "Get classical dignity status of all 9 grahas in a birth chart — exalted, debilitated, mool trikona, own sign, friend's house, neutral, enemy's, or great enemy's. Use for 'is my Saturn debilitated', 'planet dignity'.",
		Title:       "Planetary Dignity",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/dignity", in.toQuery())
	})
}

// DivisionalChartInput extends BirthInput with the varga selector so
// callers can request any divisional chart from D1 to D60 by
// shorthand ('D9', '9', or 'navamsa') — the API normalizes all three.
type DivisionalChartInput struct {
	BirthInput
	Varga string `json:"varga" jsonschema:"divisional chart name or divisor — 'D9' / 'navamsa' (marriage), 'D10' / 'dashamsa' (career), 'D7' / 'saptamsa' (children), 'D2' / 'hora' (wealth), 'D3' (siblings), 'D12' / 'dwadashamsa' (parents), 'D16', 'D20', 'D24', 'D27', 'D30' / 'trimshamsa' (misfortunes), 'D40', 'D45', 'D60'"`
}

func registerChartDivisional(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_divisional",
		Description: "Compute any divisional chart (varga) from D1 to D60 — for marriage (D9), career (D10), children (D7), wealth (D2), parents (D12), misfortunes (D30), and the rest of the classical 16. Returns divisional lagna and each planet's divisional sign + house. Use 'chart_navamsa' for the dedicated D9 entry point if a separate marriage tool is needed; this is the general-purpose varga endpoint.",
		Title:       "Divisional Chart (Any Varga)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DivisionalChartInput) (*mcp.CallToolResult, any, error) {
		path := "/v1/chart/divisional/" + in.Varga
		return callPassthrough(ctx, c, path, in.BirthInput.toQuery())
	})
}

func registerChartAvakhada(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_avakhada",
		Description: "Get the Avakhada Chakra — the traditional 8-element birth identity card: varna, vashya, yoni, gana, nadi, paya, and tatva. Used in Vedic naming, marriage matching, and personality analysis. Use for 'what's my yoni', 'avakhada for my birth'.",
		Title:       "Avakhada Chakra",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/avakhada", in.toQuery())
	})
}

func registerChartShadbala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_shadbala",
		Description: "Compute Shadbala — the six-fold strength of each graha (sthana / dig / kala / chesta / naisargika / drik bala) plus the total. Identifies the strongest and weakest planets in the chart. Use for 'which is my strongest planet', 'shadbala scores'.",
		Title:       "Shadbala (Six-fold Strength)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/shadbala", in.toQuery())
	})
}

func registerChartAshtakavarga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_ashtakavarga",
		Description: "Compute the Sarvashtakavarga — the 12-sign benefic-point grid built from each graha's individual Ashtakavarga (BAV). Used for transit timing and house strength analysis. Use for 'ashtakavarga of my chart', 'SAV grid'.",
		Title:       "Sarvashtakavarga",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/ashtakvarga", in.toQuery())
	})
}

// =====================================================================
// Dasha depth (2 tools)
// =====================================================================

func registerDashaVimshottariFull(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_vimshottari_full",
		Description: "Get the full 120-year Vimshottari Mahadasha sequence anchored at birth — all 9 lords, with start/end dates and durations. Use for 'show all my mahadashas', 'full vimshottari timeline'.",
		Title:       "Vimshottari Mahadashas (Full)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari", in.toQuery())
	})
}

func registerDashaYoginiCurrent(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_yogini_current",
		Description: "Get the currently-running Yogini Mahadasha — the 8-period (36-year cycle) secondary system used alongside Vimshottari for cross-validation. Returns the active Yogini lord with start/end. Use for 'yogini dasha now', 'which yogini is running'.",
		Title:       "Current Yogini Dasha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaCurrentInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/dasha/yogini/current", q)
	})
}

// =====================================================================
// Milan depth (3 tools)
// =====================================================================

func registerNadiDosha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "nadi_dosha",
		Description: "Check Nadi Dosha — the 8-point koota in Ashtakoota matching that fails when both partners share the same Nadi (Adi/Madhya/Antya). Returns presence + cancellation rules + remedy notes. Use for 'is there nadi dosha', 'nadi compatibility'.",
		Title:       "Nadi Dosha Check",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/nadi-dosha", boyGirlQuery(in))
	})
}

func registerVivahPhal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "vivah_phal",
		Description: "Get Vivah Phal — marriage-outcome assessment for a single person's chart: longevity of the union, bhakoot/gana indicators, and qualitative outlook derived from lagna, Mars, Jupiter, and 7th lord positions. Use for 'vivah phal for my chart', 'marriage outcome', 'will my marriage last'.",
		Title:       "Vivah Phal (Marriage Outcome)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/vivah-phal", birthPrefixQuery(in))
	})
}

func registerAshtakootaBreakdown(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ashtakoota_breakdown",
		Description: "Get the full Ashtakoota match-making table — all 8 kootas (varna, vashya, tara, yoni, graha-maitri, gana, bhakoot, nadi) individually with each koota's score, max, reasoning, and any dosha flags. Use for 'detailed kundli matching', 'koota-by-koota analysis'. Pair with match_making_score for the headline total.",
		Title:       "Ashtakoota Full Breakdown",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/ashtakoota", boyGirlQuery(in))
	})
}

// boyGirlQuery builds the boy.*/girl.* URL params that all two-person
// milan endpoints read via parseBoyGirlMoonInfo / parsePrefixedMoment.
// API convention: boy = groom (male), girl = bride (female).
// AshtakootaInput uses Bride*/Groom* names from the user-facing perspective,
// so Bride → girl.* and Groom → boy.* here.
func boyGirlQuery(in AshtakootaInput) url.Values {
	q := url.Values{}
	q.Set("girl.lat", strconv.FormatFloat(in.BrideLat, 'f', -1, 64))
	q.Set("girl.lon", strconv.FormatFloat(in.BrideLon, 'f', -1, 64))
	q.Set("girl.date", in.BrideDate)
	q.Set("girl.time", in.BrideTime)
	q.Set("girl.tz", in.BrideTz)
	q.Set("boy.lat", strconv.FormatFloat(in.GroomLat, 'f', -1, 64))
	q.Set("boy.lon", strconv.FormatFloat(in.GroomLon, 'f', -1, 64))
	q.Set("boy.date", in.GroomDate)
	q.Set("boy.time", in.GroomTime)
	q.Set("boy.tz", in.GroomTz)
	return q
}

// birthPrefixQuery builds the birth.* URL params that single-person
// milan endpoints read via parseSingleBirthMoonInfo / parsePrefixedMoment.
func birthPrefixQuery(in BirthInput) url.Values {
	q := url.Values{}
	q.Set("birth.lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
	q.Set("birth.lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
	q.Set("birth.date", in.Date)
	q.Set("birth.time", in.Time)
	q.Set("birth.tz", in.Tz)
	return q
}

// =====================================================================
// Muhurta (3 tools)
// =====================================================================

// MuhurtaWindowInput defines a date+location window plus optional
// duration filter, shared across the muhurta tools.
type MuhurtaWindowInput struct {
	Lat       float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon       float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	Tz        string  `json:"tz" jsonschema:"IANA timezone of the observer (e.g. Asia/Kolkata)"`
	StartDate string  `json:"start_date" jsonschema:"start of the search window in YYYY-MM-DD"`
	EndDate   string  `json:"end_date" jsonschema:"end of the search window in YYYY-MM-DD (inclusive)"`
}

func (m MuhurtaWindowInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(m.Lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(m.Lon, 'f', -1, 64))
	q.Set("tz", m.Tz)
	q.Set("start_date", m.StartDate)
	q.Set("end_date", m.EndDate)
	return q
}

func registerMuhurtaVivah(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_vivah",
		Description: "Find auspicious wedding (vivah) muhurta windows in a date range — ranks days by tithi, nakshatra, vara, and Guru-Shukra strength. Returns top muhurta windows with reasoning. Use for 'wedding dates in January', 'vivah muhurta for next month'.",
		Title:       "Vivah (Wedding) Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/vivah", in.toQuery())
	})
}

func registerMuhurtaNaamkaran(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_naamkaran",
		Description: "Find auspicious naamkaran (baby naming ceremony) muhurta windows — typically the 11th or 12th day after birth, scored by tithi/nakshatra. Use for 'naming ceremony date', 'naamkaran muhurta for newborn'.",
		Title:       "Naamkaran Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/naamkaran", in.toQuery())
	})
}

func registerMuhurtaBestTime(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_best_time",
		Description: "Find the highest-scoring muhurta window across a date range using a balanced scoring of tithi, nakshatra, vara, choghadiya, and Rahu Kaal exclusion. General-purpose 'when is the best time' query when no specific ceremony type applies. Use for 'best time this week', 'auspicious moment'.",
		Title:       "Best-Time Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/best-time", in.toQuery())
	})
}

// =====================================================================
// Varshaphal (1 tool)
// =====================================================================

// VarshaphalInput is the standard birth quintet plus the year for
// which the annual chart is computed.
type VarshaphalInput struct {
	BirthInput
	Year int `json:"year" jsonschema:"the year (Gregorian) for which to compute the annual / solar-return chart, e.g. 2026"`
}

func registerVarshaphalChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_chart",
		Description: "Compute the Varshaphal (annual solar-return) chart for a given year — Varsha Lagna, muntha placement, and the year-ruler (Varsha Lord). Foundation for any annual prediction. Use for 'my annual chart for 2026', 'solar return 2026'.",
		Title:       "Varshaphal Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/chart", q)
	})
}

// =====================================================================
// Numerology (1 tool — full report covers all 5 numbers)
// =====================================================================

// NumerologyFullInput accepts both name and DOB; the API computes
// whichever numbers are derivable from the inputs supplied (name
// alone → Soul/Personality/Destiny; DOB alone → Driver/Conductor;
// both → all 5).
type NumerologyFullInput struct {
	Name string `json:"name,omitempty" jsonschema:"full name as written on official documents (skipped if computing DOB-only numbers)"`
	Dob  string `json:"dob,omitempty" jsonschema:"date of birth in YYYY-MM-DD (skipped if computing name-only numbers)"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerNumerologyFull(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_full",
		Description: "Compute the 5 core Vedic numerology numbers — Driver (Mulank, from DOB day), Conductor (Bhagyank, from full DOB), Soul (Soul-Urge, from name vowels), Personality (from name consonants), and Destiny (from full name). Master numbers 11/22/33 are preserved (not reduced). Returns lucky factors per number. Use for 'numerology reading', 'driver number', 'soul urge'.",
		Title:       "Numerology — Full Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyFullInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Name != "" {
			q.Set("name", in.Name)
		}
		if in.Dob != "" {
			q.Set("dob", in.Dob)
		}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/numerology/full", q)
	})
}

// =====================================================================
// Eclipses (2 tools)
// =====================================================================

// EclipseRangeInput is a date range + optional location filter for
// eclipse queries. Lat/lon are optional — when set, the API filters
// to eclipses visible from that location.
type EclipseRangeInput struct {
	StartDate string  `json:"start_date" jsonschema:"start of the search window in YYYY-MM-DD"`
	EndDate   string  `json:"end_date" jsonschema:"end of the search window in YYYY-MM-DD"`
	Lat       float64 `json:"lat,omitempty" jsonschema:"optional observer latitude — when set, returns only eclipses visible from this location"`
	Lon       float64 `json:"lon,omitempty" jsonschema:"optional observer longitude (paired with lat)"`
}

func (e EclipseRangeInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("start_date", e.StartDate)
	q.Set("end_date", e.EndDate)
	if e.Lat != 0 || e.Lon != 0 {
		q.Set("lat", strconv.FormatFloat(e.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(e.Lon, 'f', -1, 64))
	}
	return q
}

func registerEclipsesSolar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "eclipses_solar",
		Description: "List solar eclipses in a date range — returns peak time (UT), magnitude, type (total / annular / partial / hybrid), and visibility path. Optional lat/lon filters to eclipses visible from that location. Use for 'when is the next solar eclipse', 'solar eclipses in 2026'.",
		Title:       "Solar Eclipses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in EclipseRangeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/eclipses/solar", in.toQuery())
	})
}

func registerEclipsesLunar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "eclipses_lunar",
		Description: "List lunar eclipses in a date range — returns peak time (UT), magnitude, type (total / partial / penumbral), and visibility. Optional lat/lon filters to eclipses visible from that location. Use for 'next lunar eclipse', 'chandra grahan dates'.",
		Title:       "Lunar Eclipses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in EclipseRangeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/eclipses/lunar", in.toQuery())
	})
}

// =====================================================================
// Festivals (1 tool)
// =====================================================================

type FestivalsMonthInput struct {
	Year   int    `json:"year" jsonschema:"the Gregorian year, e.g. 2026"`
	Month  int    `json:"month" jsonschema:"the Gregorian month 1..12"`
	Tz     string `json:"tz" jsonschema:"IANA timezone for date interpretation (e.g. Asia/Kolkata)"`
	Region string `json:"region,omitempty" jsonschema:"optional region filter: 'north_india', 'south_india', 'maharashtra', 'gujarat', 'bengal' — defaults to all-India"`
}

func registerFestivalsMonth(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "festivals_month",
		Description: "List Hindu festivals and observances in a given calendar month — Diwali, Holi, Navratri, Karva Chauth, Ekadashis, Pradosham, etc., with their tithi-anchored date and significance. Use for 'festivals in November 2026', 'Hindu holidays this month'.",
		Title:       "Hindu Festivals (Month)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in FestivalsMonthInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		q.Set("month", strconv.Itoa(in.Month))
		q.Set("tz", in.Tz)
		if in.Region != "" {
			q.Set("region", in.Region)
		}
		return callPassthrough(ctx, c, "/v1/festivals/month", q)
	})
}

// =====================================================================
// Planet-moments (2 tools)
// =====================================================================

// PlanetWindowInput is the date range + planet identifier shared
// across retrograde-window and ingress queries.
type PlanetWindowInput struct {
	Planet    string `json:"planet" jsonschema:"planet name: 'Sun', 'Moon', 'Mars', 'Mercury', 'Jupiter', 'Venus', 'Saturn', 'Rahu', 'Ketu'"`
	StartDate string `json:"start_date" jsonschema:"start of the search window in YYYY-MM-DD"`
	EndDate   string `json:"end_date" jsonschema:"end of the search window in YYYY-MM-DD"`
}

func (p PlanetWindowInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("planet", p.Planet)
	q.Set("start_date", p.StartDate)
	q.Set("end_date", p.EndDate)
	return q
}

func registerPlanetRetrogradeWindow(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "planet_retrograde_window",
		Description: "Find a planet's retrograde windows in a date range — start station (retrograde turn), end station (direct turn), and the sign each station occurs in. Mercury / Venus / Mars / Jupiter / Saturn retrogrades are most asked about. Use for 'when is Mercury retrograde', 'Saturn retrograde dates'.",
		Title:       "Planet Retrograde Windows",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PlanetWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/planet-moments/retrograde-window", in.toQuery())
	})
}

func registerPlanetIngress(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "planet_ingress",
		Description: "Find a planet's sign ingresses (sankranti / sign changes) in a date range — exact moment the planet crosses from one sidereal sign to the next. Use for 'when does Saturn enter Pisces', 'Jupiter sankranti'.",
		Title:       "Planet Ingress (Sign Change)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PlanetWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/planet-moments/ingress", in.toQuery())
	})
}

// =====================================================================
// Prashna (1 tool)
// =====================================================================

type PrashnaInput struct {
	Lat      float64 `json:"lat" jsonschema:"querent's latitude in decimal degrees"`
	Lon      float64 `json:"lon" jsonschema:"querent's longitude in decimal degrees"`
	Tz       string  `json:"tz" jsonschema:"IANA timezone of the question (e.g. Asia/Kolkata)"`
	At       string  `json:"at,omitempty" jsonschema:"optional moment the question was asked, RFC3339; defaults to now"`
	Question string  `json:"question" jsonschema:"the yes/no or yes/no/uncertain question being asked (free text)"`
}

func registerPrashnaAnswer(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "prashna_answer",
		Description: "Cast a Prashna (horary) chart for the moment a question is asked, and compute the classical yes/no/uncertain answer based on lagna lord, significators, and arudha. Use for 'will I get the job', 'should I make this decision', 'horary question'.",
		Title:       "Prashna Answer",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PrashnaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		q.Set("question", in.Question)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/prashna/answer", q)
	})
}

// =====================================================================
// Western (1 tool — covers natal chart end-to-end)
// =====================================================================

func registerWesternNatalChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_chart",
		Description: "Compute a Western (tropical) natal chart — all planets with sign + degree, ascendant, MC, 12 houses (Placidus / Koch / Whole Sign), and aspect grid. Distinct from the Vedic 'chart_planets' — uses tropical zodiac, not sidereal. Use for 'tropical chart', 'western natal chart', 'placidus houses'.",
		Title:       "Western Natal Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/chart", in.toQuery())
	})
}

// =====================================================================
// Horoscope (1 tool — monthly)
// =====================================================================

func registerHoroscopeMonthly(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "horoscope_monthly",
		Description: "Get the monthly horoscope for a Vedic Moon sign (rashi). Returns the full-month outlook with sankranti shifts, mahadasha-transit interplay, weekly themes, and a recommended monthly focus. Use for 'this month's horoscope', 'monthly rashifal'.",
		Title:       "Monthly Horoscope by Rashi",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in HoroscopeInput) (*mcp.CallToolResult, any, error) {
		// Same lowercase-rashi normalization as daily/weekly horoscope
		// — keeps query consistency across the three.
		in.Rashi = lowerTrim(in.Rashi)
		return callPassthrough(ctx, c, "/v1/horoscope/monthly", in.toQuery())
	})
}

// lowerTrim is a small helper used by horoscope_monthly to mirror the
// daily/weekly normalization without re-importing strings into this
// file's already-busy dependency list.
func lowerTrim(s string) string {
	out := make([]byte, 0, len(s))
	// Trim leading whitespace.
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n') {
		i++
	}
	// Trim trailing whitespace.
	end := len(s)
	for end > i && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	for ; i < end; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		out = append(out, c)
	}
	return string(out)
}
