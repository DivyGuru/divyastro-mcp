package main

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV04Tools wires the additional tools added in v0.4, bringing the
// total MCP server to 72 + (count below) tools.
//
// Coverage:
//   - Western additional natal & transits      (6)
//   - Western relationship depth (Davison+)    (6)
//   - Western progressions & returns depth     (6)
//   - Western narrative layer                  (11)
//   - Western Hellenistic depth                (7)
//   - Western extended bodies depth            (7)
//   - Western cosmobiology depth               (2)
//   - Western analytics & synthesis            (6)
//   - Western primary directions               (3)
//   - Numerology advanced                      (3)
//   - Tarot                                    (5)
func registerV04Tools(s *mcp.Server, c *apiClient) {
	// Western additional natal & transits
	registerWesternTransitsExact(s, c)
	registerWesternTransitsCalendar(s, c)
	registerWesternSynastryGrid(s, c)
	registerWesternLunarReturn(s, c)
	registerWesternSolarReturnAspects(s, c)
	registerWesternLunarReturnAspects(s, c)

	// Western relationship depth
	registerWesternCompositeAspects(s, c)
	registerWesternDavisonChart(s, c)
	registerWesternDavisonAspects(s, c)
	registerWesternProgressionsAspectsToNatal(s, c)
	registerWesternProgressionsLunation(s, c)
	registerWesternSolarArcAspectsToNatal(s, c)

	// Western narrative layer
	registerWesternHoroscopeDaily(s, c)
	registerWesternHoroscopeWeekly(s, c)
	registerWesternHoroscopeMonthly(s, c)
	registerWesternNatalSummary(s, c)
	registerWesternTransitSummary(s, c)
	registerWesternNarrativeYearly(s, c)
	registerWesternNarrativeDecan(s, c)
	registerWesternNarrativeSabian(s, c)
	registerWesternNarrativeFixedStar(s, c)
	registerWesternNarrativeTimelordProfection(s, c)
	registerWesternNarrativeTimelordZR(s, c)

	// Western Hellenistic depth
	registerWesternProfectionsMonthly(s, c)
	registerWesternProfectionsDaily(s, c)
	registerWesternZodiacalReleasingSpirit(s, c)
	registerWesternDignitiesAlmuten(s, c)
	registerWesternDignityTriplicity(s, c)
	registerWesternDignityBounds(s, c)
	registerWesternDignityFace(s, c)

	// Western extended bodies depth
	registerWesternNatalCentaurs(s, c)
	registerWesternNatalTNOs(s, c)
	registerWesternNatalLilithSet(s, c)
	registerWesternNatalUranians(s, c)
	registerWesternNatalParans(s, c)
	registerWesternNatalVertexSet(s, c)
	registerWesternNatalHarmonicAspects(s, c)

	// Western analytics & synthesis
	registerWesternCompatibilitySunSign(s, c)
	registerWesternLunarPhase(s, c)
	registerWesternNatalAspectPatterns(s, c)
	registerWesternNatalDominant(s, c)
	registerWesternNatalChartShape(s, c)
	registerWesternEclipsesVisibility(s, c)

	// Western primary directions
	registerWesternPrimaryDirectionsPlacidusZodiacal(s, c)
	registerWesternPrimaryDirectionsPlacidusMundane(s, c)
	registerWesternPrimaryDirectionsRegiomontanus(s, c)

	// Numerology advanced
	registerNumerologyPersonalPeriods(s, c)
	registerNumerologyChallenges(s, c)
	registerNumerologyAdvanced(s, c)

	// Tarot
	registerTarotCards(s, c)
	registerTarotCard(s, c)
	registerTarotDaily(s, c)
	registerTarotSpread(s, c)
	registerTarotYesNo(s, c)
}

// =====================================================================
// Additional input types
// =====================================================================

// WesternBirthReturnLunarInput holds birth + from-date for lunar return.
type WesternBirthReturnLunarInput struct {
	WesternBirthInput
	From string `json:"from" jsonschema:"start date-time for lunar return search in RFC3339 UTC (e.g. 2026-06-01T00:00:00Z)"`
}

func (i WesternBirthReturnLunarInput) toQuery() url.Values {
	q := i.WesternBirthInput.toQuery()
	q.Set("from", i.From)
	return q
}

// WesternBirthTransitCalendarInput holds birth + date-range for transit calendar.
type WesternBirthTransitCalendarInput struct {
	BirthLat  float64 `json:"birth_lat" jsonschema:"birth latitude in decimal degrees"`
	BirthLon  float64 `json:"birth_lon" jsonschema:"birth longitude in decimal degrees"`
	BirthDate string  `json:"birth_date" jsonschema:"birth date YYYY-MM-DD"`
	BirthTime string  `json:"birth_time" jsonschema:"birth time HH:MM"`
	BirthTz   string  `json:"birth_tz" jsonschema:"birth IANA timezone"`
	From      string  `json:"from" jsonschema:"start date YYYY-MM-DD (first day of range)"`
	To        string  `json:"to" jsonschema:"end date YYYY-MM-DD (last day inclusive, max 30 days from start)"`
}

func (i WesternBirthTransitCalendarInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("birth.lat", strconv.FormatFloat(i.BirthLat, 'f', -1, 64))
	q.Set("birth.lon", strconv.FormatFloat(i.BirthLon, 'f', -1, 64))
	q.Set("birth.date", i.BirthDate)
	q.Set("birth.time", i.BirthTime)
	q.Set("birth.tz", i.BirthTz)
	q.Set("from", i.From)
	q.Set("to", i.To)
	return q
}

// WesternBirthAgeInput holds birth + age_in_years for progressions.
type WesternBirthAgeInput struct {
	WesternBirthInput
	AgeInYears float64 `json:"age_in_years" jsonschema:"age in decimal years (0–120) at which to compute the progressed chart"`
}

func (i WesternBirthAgeInput) toQuery() url.Values {
	q := i.WesternBirthInput.toQuery()
	q.Set("age_in_years", strconv.FormatFloat(i.AgeInYears, 'f', -1, 64))
	return q
}

// DignityLonInput holds a tropical longitude for granular dignity endpoints.
type DignityLonInput struct {
	Lon      float64 `json:"lon" jsonschema:"tropical ecliptic longitude 0–360 (e.g. 15.5 for Aries 15°30')"`
	Scheme   string  `json:"scheme,omitempty" jsonschema:"optional dignity scheme: dorothean (default), ptolemaic, lilly"`
	Diurnal  string  `json:"diurnal,omitempty" jsonschema:"optional chart sect: true for day chart (default), false for night chart"`
}

func (d DignityLonInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("lon", strconv.FormatFloat(d.Lon, 'f', -1, 64))
	if d.Scheme != "" {
		q.Set("scheme", d.Scheme)
	}
	if d.Diurnal != "" {
		q.Set("diurnal", d.Diurnal)
	}
	return q
}

// SignInput holds a single zodiac sign name for horoscope/narrative endpoints.
type SignInput struct {
	Sign string `json:"sign" jsonschema:"tropical zodiac sign name (Aries, Taurus, Gemini, Cancer, Leo, Virgo, Libra, Scorpio, Sagittarius, Capricorn, Aquarius, Pisces)"`
}

// TwoSignInput holds two sign names for compatibility.
type TwoSignInput struct {
	SignA string `json:"sign_a" jsonschema:"first tropical zodiac sign (e.g. Aries, Leo, Scorpio)"`
	SignB string `json:"sign_b" jsonschema:"second tropical zodiac sign"`
}

// SabianInput holds an absolute zodiac degree (1–360) for Sabian symbols.
type SabianInput struct {
	Degree int    `json:"degree" jsonschema:"absolute zodiac degree 1–360 (1 = Aries 1°, 30 = Aries 30°, 31 = Taurus 1°, … 360 = Pisces 30°)"`
	Lang   string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
}

// WesternPrimaryDirectionsInput holds birth + key + max_age for primary directions.
type WesternPrimaryDirectionsInput struct {
	WesternBirthInput
	Key     string `json:"key,omitempty" jsonschema:"time key: naibod (default, 0.985647°/year) or ptolemy (1°/year)"`
	MaxAge  int    `json:"max_age,omitempty" jsonschema:"maximum age in years to compute directions for (default 80, max 120)"`
	Converse bool  `json:"converse,omitempty" jsonschema:"include converse directions (default false)"`
}

func (p WesternPrimaryDirectionsInput) toQuery() url.Values {
	q := p.WesternBirthInput.toQuery()
	if p.Key != "" {
		q.Set("key", p.Key)
	}
	if p.MaxAge > 0 {
		q.Set("max_age", strconv.Itoa(p.MaxAge))
	}
	if p.Converse {
		q.Set("converse", "true")
	}
	return q
}

// NumerologyDOBLangInput holds DOB + optional lang for numerology endpoints.
type NumerologyDOBLangInput struct {
	DOB     string `json:"dob" jsonschema:"date of birth in YYYY-MM-DD format"`
	RefDate string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today) for personal-periods"`
	Lang    string `json:"lang,omitempty" jsonschema:"optional language for lucky factors: en (default), hi, mr"`
}

func (n NumerologyDOBLangInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("dob", n.DOB)
	if n.RefDate != "" {
		q.Set("ref_date", n.RefDate)
	}
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

// NumerologyAdvancedInput holds name + optional DOB + lang for advanced numerology.
type NumerologyAdvancedInput struct {
	Name string `json:"name" jsonschema:"full name (used for Karmic Lessons, Hidden Passion, Balance, Sub-conscious Self)"`
	DOB  string `json:"dob,omitempty" jsonschema:"optional date of birth YYYY-MM-DD (required for Maturity Number)"`
	Lang string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
}

func (n NumerologyAdvancedInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("name", n.Name)
	if n.DOB != "" {
		q.Set("dob", n.DOB)
	}
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

// TarotSpreadInput holds parameters for a tarot spread reading.
type TarotSpreadInput struct {
	Type  string `json:"type" jsonschema:"spread type: single (1 card), three (past/present/future), celtic (10-card Celtic Cross)"`
	Topic string `json:"topic,omitempty" jsonschema:"optional reading focus: general (default), love, career, finance, spiritual, health"`
	Seed  int64  `json:"seed,omitempty" jsonschema:"optional reproducibility seed (positive integer; omit for a fresh time-based draw)"`
}

func (t TarotSpreadInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("type", t.Type)
	if t.Topic != "" {
		q.Set("topic", t.Topic)
	}
	if t.Seed > 0 {
		q.Set("seed", strconv.FormatInt(t.Seed, 10))
	}
	return q
}

// TarotSeedInput holds an optional seed for yes-no readings.
type TarotSeedInput struct {
	Seed int64 `json:"seed,omitempty" jsonschema:"optional reproducibility seed (positive integer; omit for a fresh time-based draw)"`
}

func (t TarotSeedInput) toQuery() url.Values {
	q := url.Values{}
	if t.Seed > 0 {
		q.Set("seed", strconv.FormatInt(t.Seed, 10))
	}
	return q
}

// =====================================================================
// Western additional transits
// =====================================================================

func registerWesternTransitsExact(s *mcp.Server, c *apiClient) {
	type Input struct {
		WesternBirthInput
		TransitPlanet string `json:"transit_planet" jsonschema:"planet doing the transiting (Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto)"`
		NatalPlanet   string `json:"natal_planet" jsonschema:"natal planet being aspected"`
		Aspect        string `json:"aspect" jsonschema:"aspect to find: conjunction, opposition, trine, square, sextile"`
		From          string `json:"from" jsonschema:"search window start in RFC3339 UTC (e.g. 2026-01-01T00:00:00Z)"`
		To            string `json:"to" jsonschema:"search window end in RFC3339 UTC (e.g. 2027-01-01T00:00:00Z)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_transits_exact",
		Description: "Find the exact moment a transiting planet will aspect a natal planet — useful for precise transit timing. Returns the Julian Day and UTC time of the exact aspect. Use for 'when does Saturn exactly conjunct my natal Sun', 'exact date of Jupiter trine natal Moon', 'precise transit timing'.",
		Title:       "Western Exact Transit Timing",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("transit_planet", in.TransitPlanet)
		q.Set("natal_planet", in.NatalPlanet)
		q.Set("aspect", in.Aspect)
		q.Set("from", in.From)
		q.Set("to", in.To)
		return callPassthrough(ctx, c, "/v1/western/transits/exact", q)
	})
}

func registerWesternTransitsCalendar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_transits_calendar",
		Description: "Get a day-by-day calendar of transiting planet aspects to a natal chart over a date range (up to 30 days). Each day shows all active transit-to-natal aspects computed at noon UTC. Use for 'show me my transits for June 2026', 'what are my weekly transit highlights', 'transit calendar for the next month'.",
		Title:       "Western Transit Calendar",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthTransitCalendarInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/transits/calendar", in.toQuery())
	})
}

func registerWesternSynastryGrid(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_synastry_grid",
		Description: "Get the full N×M synastry aspect grid between two natal charts — every planet in chart A aspecting every planet in chart B. Use for 'full synastry grid', 'all aspects between our charts', 'complete synastry table'.",
		Title:       "Western Synastry Grid",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/synastry/grid", in.toQuery())
	})
}

func registerWesternLunarReturn(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_lunar_return",
		Description: "Compute the Lunar Return chart — the chart for the exact moment the transiting Moon returns to its natal degree (~every 27.3 days). Returns planets + houses for that return moment. Use for 'lunar return chart', 'monthly return chart', 'moon return'.",
		Title:       "Western Lunar Return",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthReturnLunarInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/returns/lunar", in.toQuery())
	})
}

func registerWesternSolarReturnAspects(s *mcp.Server, c *apiClient) {
	type Input struct {
		WesternBirthInput
		ReturnYear int `json:"return_year" jsonschema:"year for which to compute the solar return (e.g. 2026)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_solar_return_aspects",
		Description: "Get the aspects formed between planets in the solar return chart — shows which planets are strongly configured in the return year. Use for 'aspects in my solar return', 'solar return chart aspects', 'what aspects are in my birthday chart'.",
		Title:       "Western Solar Return Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("return_year", strconv.Itoa(in.ReturnYear))
		return callPassthrough(ctx, c, "/v1/western/returns/solar/aspects", q)
	})
}

func registerWesternLunarReturnAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_lunar_return_aspects",
		Description: "Get the aspects formed between planets in the lunar return chart. Use for 'aspects in my lunar return', 'moon return chart aspects'.",
		Title:       "Western Lunar Return Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthReturnLunarInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/returns/lunar/aspects", in.toQuery())
	})
}

// =====================================================================
// Western relationship depth
// =====================================================================

func registerWesternCompositeAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_composite_aspects",
		Description: "Get aspects in the composite (midpoint) chart for two people — shows how the relationship's combined energy is configured. Use for 'composite chart aspects', 'aspects in our composite chart', 'relationship chart aspects'.",
		Title:       "Western Composite Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/composite/aspects", in.toQuery())
	})
}

func registerWesternDavisonChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_davison_chart",
		Description: "Compute the Davison relationship chart — an alternative to the composite chart that uses the midpoint of the two birth times rather than midpoint longitudes. Returns planets + houses for the Davison chart. Use for 'Davison chart', 'Davison relationship chart', 'Davison horoscope'.",
		Title:       "Western Davison Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/davison/chart", in.toQuery())
	})
}

func registerWesternDavisonAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_davison_aspects",
		Description: "Get aspects in the Davison relationship chart. Use for 'Davison chart aspects', 'aspects in our Davison chart'.",
		Title:       "Western Davison Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/davison/aspects", in.toQuery())
	})
}

func registerWesternProgressionsAspectsToNatal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_progressions_aspects_to_natal",
		Description: "Get secondary progressed planet aspects to natal planets — the slow symbolic evolution of the birth chart over time. Use for 'progressed aspects to natal', 'which natal planets are my progressed planets aspecting', 'secondary progression aspects'.",
		Title:       "Western Progressions Aspects to Natal",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthAgeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/progressions/aspects-to-natal", in.toQuery())
	})
}

func registerWesternProgressionsLunation(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_progressions_lunation",
		Description: "Get the progressed Moon phase — shows what stage of the 29-year progressed lunation cycle you are in (new moon, crescent, first quarter, gibbous, full moon, disseminating, last quarter, balsamic). Use for 'progressed moon phase', 'what progressed lunar cycle am I in', 'progressed new moon'.",
		Title:       "Western Progressed Lunation",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthAgeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/progressions/lunation", in.toQuery())
	})
}

func registerWesternSolarArcAspectsToNatal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_solar_arc_aspects_to_natal",
		Description: "Get solar arc directed planet aspects to natal planets — a predictive technique where all planets advance at the rate of the Sun (~1°/year). Use for 'solar arc aspects to natal', 'which natal planets have solar arc hits', 'solar arc direction aspects'.",
		Title:       "Western Solar Arc Aspects to Natal",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthAgeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/solar-arc/aspects-to-natal", in.toQuery())
	})
}

// =====================================================================
// Western narrative layer
// =====================================================================

func registerWesternHoroscopeDaily(s *mcp.Server, c *apiClient) {
	type Input struct {
		Sign string `json:"sign" jsonschema:"tropical zodiac sign (Aries, Taurus, Gemini, Cancer, Leo, Virgo, Libra, Scorpio, Sagittarius, Capricorn, Aquarius, Pisces)"`
		Date string `json:"date" jsonschema:"date in YYYY-MM-DD format"`
		Tone string `json:"tone,omitempty" jsonschema:"optional tone: spiritual (default), motivational, romantic, practical"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_horoscope_daily",
		Description: "Get the daily Western horoscope for a sun sign. Uses Moon transit position for perpetual coverage — every date has a unique, astrologically grounded reading. Use for 'Aries horoscope today', 'daily horoscope for Scorpio', 'what does today look like for Taurus'.",
		Title:       "Western Daily Horoscope",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign", in.Sign)
		q.Set("date", in.Date)
		if in.Tone != "" {
			q.Set("tone", in.Tone)
		}
		return callPassthrough(ctx, c, "/v1/western/horoscope/daily", q)
	})
}

func registerWesternHoroscopeWeekly(s *mcp.Server, c *apiClient) {
	type Input struct {
		Sign      string `json:"sign" jsonschema:"tropical zodiac sign name"`
		WeekStart string `json:"week_start" jsonschema:"any date in the desired week in YYYY-MM-DD format"`
		Tone      string `json:"tone,omitempty" jsonschema:"optional tone: spiritual (default), motivational, romantic, practical"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_horoscope_weekly",
		Description: "Get the weekly Western horoscope for a sun sign. Uses Moon's weekly position for distinct, astrologically meaningful weekly content. Use for 'Gemini horoscope this week', 'weekly outlook for Capricorn', 'week ahead for Libra'.",
		Title:       "Western Weekly Horoscope",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign", in.Sign)
		q.Set("week_start", in.WeekStart)
		if in.Tone != "" {
			q.Set("tone", in.Tone)
		}
		return callPassthrough(ctx, c, "/v1/western/horoscope/weekly", q)
	})
}

func registerWesternHoroscopeMonthly(s *mcp.Server, c *apiClient) {
	type Input struct {
		Sign      string `json:"sign" jsonschema:"tropical zodiac sign name"`
		YearMonth string `json:"year_month" jsonschema:"month in YYYY-MM format (e.g. 2026-07)"`
		Tone      string `json:"tone,omitempty" jsonschema:"optional tone: spiritual (default), motivational, romantic, practical"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_horoscope_monthly",
		Description: "Get the monthly Western horoscope for a sun sign. June–December 2026 has authored content; other months use the annual solar theme. Use for 'Scorpio horoscope for July 2026', 'monthly outlook for Pisces', 'what does this month look like for Leo'.",
		Title:       "Western Monthly Horoscope",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign", in.Sign)
		q.Set("year_month", in.YearMonth)
		if in.Tone != "" {
			q.Set("tone", in.Tone)
		}
		return callPassthrough(ctx, c, "/v1/western/horoscope/monthly", q)
	})
}

func registerWesternNatalSummary(s *mcp.Server, c *apiClient) {
	type Input struct {
		WesternBirthInput
		Verbosity string `json:"verbosity,omitempty" jsonschema:"optional depth: brief (3 sections), medium (6, default), full (12)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_summary",
		Description: "Get a composed narrative summary of a Western natal chart — authored prose covering Ascendant, Sun, Moon placements and key aspect highlights. Use for 'interpret my Western natal chart', 'natal chart reading', 'describe my birth chart in words'.",
		Title:       "Western Natal Summary Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		if in.Verbosity != "" {
			q.Set("verbosity", in.Verbosity)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/natal-summary", q)
	})
}

func registerWesternTransitSummary(s *mcp.Server, c *apiClient) {
	type Input struct {
		WesternTransitToNatalInput
		Verbosity string `json:"verbosity,omitempty" jsonschema:"optional depth: brief (3 sections), medium (6, default), full (12)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_transit_summary",
		Description: "Get a composed narrative summary of current transits to a natal chart — authored prose for the most exact active transits at the given moment. Use for 'interpret my current transits', 'what are my transits saying', 'describe the astrology affecting me right now'.",
		Title:       "Western Transit Summary Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := in.WesternTransitToNatalInput.toQuery()
		if in.Verbosity != "" {
			q.Set("verbosity", in.Verbosity)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/transit-summary", q)
	})
}

func registerWesternNarrativeYearly(s *mcp.Server, c *apiClient) {
	type Input struct {
		Sign string `json:"sign" jsonschema:"tropical zodiac sign (Aries through Pisces)"`
		Lang string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_yearly",
		Description: "Get the annual solar-return theme for a sun sign — a deep interpretive essay on the year's archetypal themes for that sign. Available for all 12 signs in English, Hindi, and Marathi. Use for 'annual theme for Aries', 'yearly horoscope interpretation for Virgo', 'solar year themes for Scorpio'.",
		Title:       "Western Yearly Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign", in.Sign)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/yearly", q)
	})
}

func registerWesternNarrativeDecan(s *mcp.Server, c *apiClient) {
	type Input struct {
		Sign  string `json:"sign" jsonschema:"tropical zodiac sign"`
		Decan int    `json:"decan" jsonschema:"decan number: 1 (degrees 0–10), 2 (10–20), or 3 (20–30)"`
		Lang  string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_decan",
		Description: "Get the authored narrative for a zodiac decan (10° division of a sign). Each sign has 3 decans ruled by different planets. Use for 'interpret Aries first decan', 'Scorpio second decan meaning', 'decan interpretation for my Sun'.",
		Title:       "Western Decan Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign", in.Sign)
		q.Set("decan", strconv.Itoa(in.Decan))
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/decan", q)
	})
}

func registerWesternNarrativeSabian(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_sabian",
		Description: "Get the authored Sabian symbol and interpretation for a specific zodiac degree (1–360, where 1 = Aries 1° and 360 = Pisces 30°). Use for 'Sabian symbol for degree 15 of Leo', 'what is the Sabian for my Sun', 'Sabian symbol interpretation'.",
		Title:       "Western Sabian Symbol",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in SabianInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("degree", strconv.Itoa(in.Degree))
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/sabian", q)
	})
}

func registerWesternNarrativeFixedStar(s *mcp.Server, c *apiClient) {
	type Input struct {
		Star string `json:"star" jsonschema:"fixed star name (e.g. Algol, Regulus, Sirius, Spica, Antares, Aldebaran, Vega, Arcturus)"`
		Lang string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_fixed_star",
		Description: "Get the authored narrative interpretation for a named fixed star in a natal chart. Use for 'what does Algol conjunct my natal Mars mean', 'Regulus interpretation', 'meaning of Sirius in my chart'.",
		Title:       "Western Fixed Star Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/fixed-star/"+url.PathEscape(in.Star), q)
	})
}

func registerWesternNarrativeTimelordProfection(s *mcp.Server, c *apiClient) {
	type Input struct {
		House int    `json:"house" jsonschema:"activated profection house 1–12 (from western_profections_annual)"`
		Lang  string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_timelord_profection",
		Description: "Get the authored narrative for the annual profection house — a description of the life themes activated in this profection year. Use after western_profections_annual to get the interpretation. Use for 'what does my 7th house profection mean', 'interpretation of my profection year', 'annual profection reading'.",
		Title:       "Western Profection Year Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("house", strconv.Itoa(in.House))
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/timelord/profection", q)
	})
}

func registerWesternNarrativeTimelordZR(s *mcp.Server, c *apiClient) {
	type Input struct {
		Lord string `json:"lord" jsonschema:"Zodiacal Releasing lord name (Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn)"`
		Lang string `json:"lang,omitempty" jsonschema:"optional language: en (default), hi, mr"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_timelord_zr",
		Description: "Get the authored narrative for a Zodiacal Releasing planetary lord — describes the life-period themes when a given planet is lord in the ZR system. Use for 'what does Saturn as ZR lord mean', 'Zodiacal Releasing Saturn period interpretation'.",
		Title:       "Western Zodiacal Releasing Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lord", in.Lord)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/timelord/zr", q)
	})
}

// =====================================================================
// Western Hellenistic depth
// =====================================================================

func registerWesternProfectionsMonthly(s *mcp.Server, c *apiClient) {
	type Input struct {
		BirthDate string `json:"birth_date" jsonschema:"birth date YYYY-MM-DD"`
		RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
		AscSign   string `json:"asc_sign" jsonschema:"natal Ascendant sign (Aries, Taurus, … Pisces) or 0-based index 0–11"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_profections_monthly",
		Description: "Compute the monthly profection sub-period — which house and lord are activated this month within the annual profection cycle. Use for 'monthly profection', 'sub-period lord this month', 'monthly time lord'.",
		Title:       "Western Monthly Profection",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		q.Set("asc_sign", in.AscSign)
		return callPassthrough(ctx, c, "/v1/western/profections/monthly", q)
	})
}

func registerWesternProfectionsDaily(s *mcp.Server, c *apiClient) {
	type Input struct {
		BirthDate string `json:"birth_date" jsonschema:"birth date YYYY-MM-DD"`
		RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
		AscSign   string `json:"asc_sign" jsonschema:"natal Ascendant sign or 0-based index 0–11"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_profections_daily",
		Description: "Compute the daily profection drill-down — the most granular level of the profection system. Use for 'daily profection', 'today's profection lord', 'daily time lord'.",
		Title:       "Western Daily Profection",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		q.Set("asc_sign", in.AscSign)
		return callPassthrough(ctx, c, "/v1/western/profections/daily", q)
	})
}

func registerWesternZodiacalReleasingSpirit(s *mcp.Server, c *apiClient) {
	type Input struct {
		BirthDate string `json:"birth_date" jsonschema:"birth date YYYY-MM-DD"`
		RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
		LotSign   string `json:"lot_sign" jsonschema:"sign of Lot of Spirit (from western_natal_lots) — sign name or 0-based index 0–11"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_zodiacal_releasing_spirit",
		Description: "Compute Zodiacal Releasing from the Lot of Spirit — indicates career and achievement periods. Complement to western_zodiacal_releasing (Fortune). Use for 'ZR from Spirit', 'career time lords', 'professional periods in Hellenistic astrology'.",
		Title:       "Western Zodiacal Releasing from Spirit",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		q.Set("lot_sign", in.LotSign)
		return callPassthrough(ctx, c, "/v1/western/zodiacal-releasing/spirit", q)
	})
}

func registerWesternDignitiesAlmuten(s *mcp.Server, c *apiClient) {
	type Input struct {
		AscLon  float64 `json:"asc_lon" jsonschema:"Ascendant tropical longitude 0–360"`
		Diurnal string  `json:"diurnal,omitempty" jsonschema:"chart sect: true for day chart (default), false for night chart"`
		Scheme  string  `json:"scheme,omitempty" jsonschema:"dignity scheme: dorothean (default), ptolemaic, lilly"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_almuten",
		Description: "Compute the Almuten Figuris — the planet with the highest combined essential dignity score at the Ascendant degree. The 'ruler of the chart' in classical astrology. Use for 'almuten figuris', 'chart almuten', 'which planet rules my chart classically'.",
		Title:       "Western Almuten Figuris",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("asc_lon", strconv.FormatFloat(in.AscLon, 'f', -1, 64))
		if in.Diurnal != "" {
			q.Set("diurnal", in.Diurnal)
		}
		if in.Scheme != "" {
			q.Set("scheme", in.Scheme)
		}
		return callPassthrough(ctx, c, "/v1/western/dignities/almuten", q)
	})
}

func registerWesternDignityTriplicity(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_triplicity",
		Description: "Get the triplicity (element) ruler for a zodiac degree — shows which planet rules by fire/earth/air/water triplicity, with day and night rulers. Use for 'triplicity ruler of my Sun', 'what planet rules Aries triplicity', 'fire triplicity lord'.",
		Title:       "Western Triplicity Dignity",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DignityLonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/dignities/triplicity", in.toQuery())
	})
}

func registerWesternDignityBounds(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_bounds",
		Description: "Get the bound (term) ruler for a zodiac degree and the full sign breakdown — shows which planet rules each of the five bounds within a sign. Egyptian or Ptolemaic tables. Use for 'bound ruler of my Venus', 'which term does my Sun fall in', 'Egyptian terms'.",
		Title:       "Western Bounds/Terms Dignity",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DignityLonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/dignities/bounds", in.toQuery())
	})
}

func registerWesternDignityFace(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_face",
		Description: "Get the face (decan) ruler for a zodiac degree using the Chaldean order — shows which of the 36 faces the degree falls in and its planetary ruler. Use for 'face ruler of my Moon', 'decan ruler for my Ascendant degree', 'Chaldean face'.",
		Title:       "Western Face/Decan Dignity",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DignityLonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/dignities/face", in.toQuery())
	})
}

// =====================================================================
// Western extended bodies depth
// =====================================================================

func registerWesternNatalCentaurs(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_centaurs",
		Description: "Get natal positions for the centaur bodies: Pholus, Nessus, Asbolus, and Chariklo. Use for 'centaur positions in my natal chart', 'where is Pholus', 'Nessus natal placement'.",
		Title:       "Western Natal Centaurs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/centaurs", in.toQuery())
	})
}

func registerWesternNatalTNOs(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_tnos",
		Description: "Get natal positions for trans-Neptunian objects: Eris, Sedna, Haumea, Makemake, Quaoar, Orcus, and Varuna. Use for 'TNO positions in my chart', 'where is Eris natally', 'trans-Neptunian placements'.",
		Title:       "Western Natal TNOs",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/tnos", in.toQuery())
	})
}

func registerWesternNatalLilithSet(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_lilith_set",
		Description: "Get all Lilith positions: Mean Black Moon Lilith, True Black Moon Lilith, Osculating (Interpolated) Lilith, and asteroid 1181 Lilith. Use for 'Black Moon Lilith in my chart', 'all Lilith placements', 'Mean vs True Lilith'.",
		Title:       "Western Natal Lilith Set",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/lilith-set", in.toQuery())
	})
}

func registerWesternNatalUranians(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_uranian_hypotheticals",
		Description: "Get positions for the eight Uranian/Hamburg school hypothetical planets: Cupido, Hades, Zeus, Kronos, Apollon, Admetos, Vulkanus, and Poseidon. Use for 'Uranian planets in my chart', 'Hamburg school hypotheticals', 'Kronos placement'.",
		Title:       "Western Natal Uranian Hypotheticals",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/uranian-hypotheticals", in.toQuery())
	})
}

func registerWesternNatalParans(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_parans",
		Description: "Compute paran aspects — fixed stars crossing the 4 angles (Ascendant, MC, Descendant, IC) in paral­lel with natal planets. Use for 'parans in my natal chart', 'which fixed stars are in paran with my planets', 'angular fixed star contacts'.",
		Title:       "Western Natal Parans",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/parans", in.toQuery())
	})
}

func registerWesternNatalVertexSet(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_vertex_set",
		Description: "Get the Vertex, Anti-Vertex, and East Point for a natal chart. The Vertex is an electrically charged destiny point; contacts to it from transits or synastry indicate fated encounters. Use for 'Vertex in my natal chart', 'what sign is my Vertex in', 'anti-vertex placement'.",
		Title:       "Western Natal Vertex Set",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/vertex-set", in.toQuery())
	})
}

func registerWesternNatalHarmonicAspects(s *mcp.Server, c *apiClient) {
	type Input struct {
		WesternBirthInput
		Harmonics string `json:"harmonics,omitempty" jsonschema:"optional comma-separated harmonics to check (e.g. '5,7,9,11'); defaults to all major harmonics"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_harmonic_aspects",
		Description: "Find harmonic aspects across multiple harmonics simultaneously — quintiles (H5), septiles (H7), noviles (H9), etc. Use for 'harmonic aspects in my chart', 'quintile aspects', 'septile configuration', 'higher harmonic contacts'.",
		Title:       "Western Natal Harmonic Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		if in.Harmonics != "" {
			q.Set("harmonics", in.Harmonics)
		}
		return callPassthrough(ctx, c, "/v1/western/natal/harmonic-aspects", q)
	})
}

// =====================================================================
// Western analytics & synthesis
// =====================================================================

func registerWesternCompatibilitySunSign(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_compatibility_sun_sign",
		Description: "Get the Western sun-sign compatibility between two zodiac signs — score 0–100, band (excellent/good/fair/challenging), aspect basis (trine/sextile/square etc.), element interaction, strengths, and challenges. No birth data needed. Use for 'compatibility between Aries and Leo', 'Scorpio and Aquarius compatibility', 'are Taurus and Cancer compatible'.",
		Title:       "Western Sun-Sign Compatibility",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoSignInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("sign_a", in.SignA)
		q.Set("sign_b", in.SignB)
		return callPassthrough(ctx, c, "/v1/western/compatibility/sun-sign", q)
	})
}

func registerWesternLunarPhase(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_lunar_phase",
		Description: "Get the natal Moon phase at birth — phase angle, illumination percentage, 8-named phase (new moon through waning crescent), and age in days since last new moon. Use for 'what phase was the Moon in when I was born', 'natal lunar phase', 'birth Moon illumination'.",
		Title:       "Western Natal Lunar Phase",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/lunar/phase", in.toQuery())
	})
}

func registerWesternNatalAspectPatterns(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_aspect_patterns",
		Description: "Detect classical multi-planet aspect patterns in a natal chart — Grand Trine, T-Square, Grand Cross, Yod (Finger of God), Kite, Mystic Rectangle, Stellium, Grand Sextile. Use for 'what aspect patterns do I have', 'Grand Trine in my chart', 'is there a Yod in my natal chart', 'aspect configurations'.",
		Title:       "Western Natal Aspect Patterns",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/aspect-patterns", in.toQuery())
	})
}

func registerWesternNatalDominant(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_dominant",
		Description: "Get the chart dominance profile — dominant element (fire/earth/air/water), dominant modality (cardinal/fixed/mutable), most aspected planet, Lilly dignity ranking of traditional planets, and hemisphere balance. Use for 'what element dominates my chart', 'most aspected planet in my natal chart', 'chart dominance analysis', 'what is my dominant planet'.",
		Title:       "Western Natal Dominance Profile",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/dominant", in.toQuery())
	})
}

func registerWesternNatalChartShape(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_chart_shape",
		Description: "Classify the natal chart into one of the seven Jones planetary distribution patterns: Splash, Bundle, Locomotive, Bowl, Bucket, See-Saw, or Splay. Use for 'what shape is my chart', 'Jones pattern in my natal chart', 'is my chart a Bowl or Bundle', 'planetary distribution pattern'.",
		Title:       "Western Natal Chart Shape",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/chart-shape", in.toQuery())
	})
}

func registerWesternEclipsesVisibility(s *mcp.Server, c *apiClient) {
	type Input struct {
		Lat  float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
		Lon  float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
		From string  `json:"from,omitempty" jsonschema:"optional start date YYYY-MM-DD (defaults to today)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_eclipses_visibility",
		Description: "Check whether upcoming eclipses are visible from a specific location — returns visibility type (total, partial, annular, not visible) and times of contact. Use for 'is the next solar eclipse visible from Mumbai', 'eclipse visibility from my location', 'can I see the lunar eclipse'.",
		Title:       "Western Eclipse Visibility",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		if in.From != "" {
			q.Set("from", in.From)
		}
		return callPassthrough(ctx, c, "/v1/western/eclipses/visibility", q)
	})
}

// =====================================================================
// Western Primary Directions
// =====================================================================

func registerWesternPrimaryDirectionsPlacidusZodiacal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_primary_directions_placidus_zodiacal",
		Description: "Compute Placidus zodiacal primary directions — the oldest predictive technique in Western astrology (1° ≈ 1 year). Returns all direction events sorted chronologically: which natal planets are being directed and what they aspect. Use for 'primary directions for my chart', 'Placidus directions', 'what primary directions am I under', 'direction to natal Sun'.",
		Title:       "Western Primary Directions (Placidus Zodiacal)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternPrimaryDirectionsInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/primary-directions/placidus-zodiacal", in.toQuery())
	})
}

func registerWesternPrimaryDirectionsPlacidusMundane(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_primary_directions_placidus_mundane",
		Description: "Compute Placidus mundane primary directions — the promissor moves through the Placidus mundane house framework. An alternative calculation method to zodiacal directions. Use for 'Placidus mundane directions', 'mundane primary directions'.",
		Title:       "Western Primary Directions (Placidus Mundane)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternPrimaryDirectionsInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/primary-directions/placidus-mundane", in.toQuery())
	})
}

func registerWesternPrimaryDirectionsRegiomontanus(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_primary_directions_regiomontanus",
		Description: "Compute Regiomontanus primary directions — uses equatorial Right Ascension arc, the simplest and fastest calculation method. Use for 'Regiomontanus directions', 'RA-based primary directions'.",
		Title:       "Western Primary Directions (Regiomontanus)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternPrimaryDirectionsInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/primary-directions/regiomontanus", in.toQuery())
	})
}

// =====================================================================
// Numerology Advanced
// =====================================================================

func registerNumerologyPersonalPeriods(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_personal_periods",
		Description: "Compute the Personal Year, Personal Month, and Personal Day numbers for any date — the time-transiting layer of Vedic numerology that shows the numerological climate for a specific period. Each number comes with lucky factors (planet, day, color, gemstone). Use for 'what is my personal year number', 'personal month number for June', 'numerology forecast for today'.",
		Title:       "Numerology Personal Periods",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyDOBLangInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/personal-periods", in.toQuery())
	})
}

func registerNumerologyChallenges(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_challenges",
		Description: "Compute the four Vedic Challenge Numbers from a birth date — the life obstacles and developmental lessons to be mastered in different life stages (CH1 early life, CH2 middle life, CH3 major challenge, CH4 later life). Use for 'numerology challenge numbers', 'what are my life challenges', 'Vedic challenge number from birth date'.",
		Title:       "Numerology Challenge Numbers",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyDOBLangInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/challenges", in.toQuery())
	})
}

func registerNumerologyAdvanced(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_advanced",
		Description: "Compute the advanced Vedic numerology profile from a name (and optionally DOB): Maturity Number (life purpose emerging after 40), Karmic Lessons (missing digits = undeveloped qualities), Sub-conscious Self (9 minus karmic lessons), Hidden Passion (dominant talent), and Balance Number (crisis response quality). Use for 'karmic lessons from my name', 'hidden passion number', 'maturity number', 'advanced numerology profile'.",
		Title:       "Numerology Advanced Profile",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyAdvancedInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/advanced", in.toQuery())
	})
}

// =====================================================================
// Tarot
// =====================================================================

func registerTarotCards(s *mcp.Server, c *apiClient) {
	type Input struct {
		Suit   string `json:"suit,omitempty" jsonschema:"optional filter by suit: Major Arcana, Wands, Cups, Pentacles, Swords"`
		Search string `json:"search,omitempty" jsonschema:"optional keyword search across card names and keywords (cannot combine with suit)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "tarot_cards",
		Description: "List tarot cards from the 78-card Rider-Waite deck with upright/reversed meanings, keywords, and yes/no tendency. Filter by suit or search by keyword. Use for 'show me all Cups cards', 'tarot cards about love', 'what cards are in the Major Arcana', 'list Swords suit'.",
		Title:       "Tarot Card List",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Suit != "" {
			q.Set("suit", in.Suit)
		}
		if in.Search != "" {
			q.Set("search", in.Search)
		}
		return callPassthrough(ctx, c, "/v1/tarot/cards", q)
	})
}

func registerTarotCard(s *mcp.Server, c *apiClient) {
	type Input struct {
		ID int `json:"id" jsonschema:"card ID 0–77 (0 = The Fool, 21 = The World, 22 = Ace of Wands, 77 = King of Swords)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "tarot_card",
		Description: "Get the complete details for a single tarot card by ID — name, suit, element, archetype, upright and reversed meanings, keywords, and yes/no tendency. Use for 'tell me about The Tower card', 'what does the Ten of Cups mean', 'Death card meaning', 'High Priestess interpretation'.",
		Title:       "Tarot Single Card",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/tarot/card/"+strconv.Itoa(in.ID), url.Values{})
	})
}

func registerTarotDaily(s *mcp.Server, c *apiClient) {
	type Input struct {
		Date string `json:"date,omitempty" jsonschema:"optional date YYYY-MM-DD (defaults to today UTC) — the same card is returned for all callers on the same date"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "tarot_daily",
		Description: "Get the card of the day — a deterministic daily tarot card that is the same for everyone on a given date. Includes the card's upright or reversed meaning as a daily message. Use for 'tarot card of the day', 'daily tarot', 'what card represents today', 'daily tarot message'.",
		Title:       "Tarot Daily Card",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in Input) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		return callPassthrough(ctx, c, "/v1/tarot/daily", q)
	})
}

func registerTarotSpread(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "tarot_spread",
		Description: "Draw a tarot spread with a chosen layout and optional topic focus. Single card for a focused message; three-card for past/present/future; Celtic Cross for a full 10-card reading. Results are reproducible with a seed. Use for 'do a 3-card tarot reading', 'Celtic Cross spread for my relationship', 'pull a tarot card about my career', 'single card reading about love'.",
		Title:       "Tarot Spread Reading",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TarotSpreadInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/tarot/spread", in.toQuery())
	})
}

func registerTarotYesNo(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "tarot_yes_no",
		Description: "Get a yes/no tarot answer — draws one card, returns yes/no/maybe based on the card's natural tendency modified by orientation, with reasoning. Use for 'tarot yes or no', 'should I take this job — tarot answer', 'will my project succeed — tarot reading', 'quick tarot answer'.",
		Title:       "Tarot Yes/No Reading",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TarotSeedInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/tarot/yes-no", in.toQuery())
	})
}

// url.PathEscape is used by fixed-star tool — ensure it is imported.
var _ = strings.TrimSpace
