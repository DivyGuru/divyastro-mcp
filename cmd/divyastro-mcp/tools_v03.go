package main

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV03Tools wires the 30 additional Western astrology tools in v0.3,
// bringing the MCP server to 72 tools combined across v0.1 + v0.2 + v0.3.
//
// Coverage by group:
//
//   - Western natal depth (4)         — houses, aspects, rulerships, transits
//   - Western relationships (3)       — synastry aspects/score, composite chart
//   - Western projection & returns (3) — progressions, solar arc, solar return
//   - Western transits & narrative (3) — transits-to-natal, aspect/transit interp
//   - Western lots & midpoints (2)    — Arabic Parts, midpoint table
//   - Hellenistic time lords (4)      — profections, ZR, firdaria, dignities
//   - Extended bodies (2)             — fixed stars, big-four asteroids
//   - Cosmobiology (4)                — declinations, harmonic, antiscia, midpoint-tree
//   - Mundane & astrocartography (5)  — heliocentric, eclipses, ingresses, Rx, ACG
func registerV03Tools(s *mcp.Server, c *apiClient) {
	// Western natal depth
	registerWesternNatalHouses(s, c)
	registerWesternNatalAspects(s, c)
	registerWesternNatalRulerships(s, c)
	registerWesternTransitPositions(s, c)

	// Relationships
	registerWesternSynastryAspects(s, c)
	registerWesternSynastryScore(s, c)
	registerWesternCompositeChart(s, c)

	// Projection & returns
	registerWesternProgressionsPlanets(s, c)
	registerWesternSolarArcPlanets(s, c)
	registerWesternSolarReturn(s, c)

	// Transits & narrative
	registerWesternTransitsToNatal(s, c)
	registerWesternInterpretAspect(s, c)
	registerWesternInterpretTransit(s, c)

	// Lots & midpoints
	registerWesternNatalLots(s, c)
	registerWesternNatalMidpoints(s, c)

	// Hellenistic time lords
	registerWesternProfectionsAnnual(s, c)
	registerWesternZodiacalReleasing(s, c)
	registerWesternFirdaria(s, c)
	registerWesternDignities(s, c)

	// Extended bodies
	registerWesternNatalFixedStars(s, c)
	registerWesternNatalBigFourAsteroids(s, c)

	// Cosmobiology
	registerWesternNatalDeclinations(s, c)
	registerWesternNatalHarmonic(s, c)
	registerWesternNatalAntiscia(s, c)
	registerWesternNatalMidpointTree(s, c)

	// Mundane & astrocartography
	registerWesternHeliocentricPlanets(s, c)
	registerWesternEclipsesUpcoming(s, c)
	registerWesternIngresses(s, c)
	registerWesternRetrogradeWindow(s, c)
	registerWesternACGPlanetLines(s, c)
}

// =====================================================================
// Shared input types
// =====================================================================

// WesternBirthInput extends BirthInput with an optional house system.
type WesternBirthInput struct {
	Lat         float64 `json:"lat" jsonschema:"birth latitude in decimal degrees (e.g. 48.8566 for Paris)"`
	Lon         float64 `json:"lon" jsonschema:"birth longitude in decimal degrees (e.g. 2.3522 for Paris)"`
	Date        string  `json:"date" jsonschema:"birth date in YYYY-MM-DD format"`
	Time        string  `json:"time" jsonschema:"birth time in HH:MM 24-hour format"`
	Tz          string  `json:"tz" jsonschema:"IANA timezone (e.g. Europe/Paris, America/New_York)"`
	HouseSystem string  `json:"house_system,omitempty" jsonschema:"optional house system: placidus (default), whole-sign, equal, koch, campanus, regiomontanus, porphyry"`
}

func (b WesternBirthInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(b.Lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(b.Lon, 'f', -1, 64))
	q.Set("date", b.Date)
	q.Set("time", b.Time)
	q.Set("tz", b.Tz)
	if b.HouseSystem != "" {
		q.Set("house_system", b.HouseSystem)
	}
	return q
}

// TwoPersonInput holds birth data for two people (synastry / composite).
type TwoPersonInput struct {
	ALat  float64 `json:"a_lat" jsonschema:"person A latitude"`
	ALon  float64 `json:"a_lon" jsonschema:"person A longitude"`
	ADate string  `json:"a_date" jsonschema:"person A birth date YYYY-MM-DD"`
	ATime string  `json:"a_time" jsonschema:"person A birth time HH:MM"`
	ATz   string  `json:"a_tz" jsonschema:"person A IANA timezone"`
	BLat  float64 `json:"b_lat" jsonschema:"person B latitude"`
	BLon  float64 `json:"b_lon" jsonschema:"person B longitude"`
	BDate string  `json:"b_date" jsonschema:"person B birth date YYYY-MM-DD"`
	BTime string  `json:"b_time" jsonschema:"person B birth time HH:MM"`
	BTz   string  `json:"b_tz" jsonschema:"person B IANA timezone"`
}

func (t TwoPersonInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("personA.lat", strconv.FormatFloat(t.ALat, 'f', -1, 64))
	q.Set("personA.lon", strconv.FormatFloat(t.ALon, 'f', -1, 64))
	q.Set("personA.date", t.ADate)
	q.Set("personA.time", t.ATime)
	q.Set("personA.tz", t.ATz)
	q.Set("personB.lat", strconv.FormatFloat(t.BLat, 'f', -1, 64))
	q.Set("personB.lon", strconv.FormatFloat(t.BLon, 'f', -1, 64))
	q.Set("personB.date", t.BDate)
	q.Set("personB.time", t.BTime)
	q.Set("personB.tz", t.BTz)
	return q
}

// DateRangeInput is used by mundane endpoints (ingresses, retrograde window).
type DateRangeInput struct {
	StartDate string `json:"start_date" jsonschema:"start date in YYYY-MM-DD format"`
	EndDate   string `json:"end_date" jsonschema:"end date in YYYY-MM-DD format (max 5 years ahead)"`
}

func (d DateRangeInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("start_date", d.StartDate)
	q.Set("end_date", d.EndDate)
	return q
}

// WesternMomentInput is a birth-context-less moment (date + time + tz).
// Used for heliocentric and ACG endpoints that don't need lat/lon.
type WesternMomentInput struct {
	Date string `json:"date" jsonschema:"date in YYYY-MM-DD format"`
	Time string `json:"time" jsonschema:"time in HH:MM 24-hour format"`
	Tz   string `json:"tz" jsonschema:"IANA timezone (e.g. UTC, Asia/Kolkata)"`
}

func (m WesternMomentInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("date", m.Date)
	q.Set("time", m.Time)
	q.Set("tz", m.Tz)
	return q
}

// =====================================================================
// 1–4: Western Natal Depth
// =====================================================================

func registerWesternNatalHouses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_houses",
		Description: "Compute Western tropical house cusps for a birth chart. Returns all 12 house cusp degrees plus Ascendant and MC in the chosen house system (default Placidus). Use for 'what houses are in my Western chart', 'show me the house cusps', 'Placidus house degrees'.",
		Title:       "Western Natal Houses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/houses", in.toQuery())
	})
}

func registerWesternNatalAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_aspects",
		Description: "Compute all natal aspects in a Western tropical birth chart — conjunctions, oppositions, trines, squares, sextiles, and minor aspects. Returns orbs and partile flags. Use for 'what aspects does my natal chart have', 'Sun square Moon aspect', 'tight natal aspects'.",
		Title:       "Western Natal Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/aspects", in.toQuery())
	})
}

func registerWesternNatalRulerships(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_rulerships",
		Description: "Return the house ruler for each of the 12 houses in a Western tropical chart — shows which planet rules each house and where that planet is placed. Supports modern and traditional rulership schemes. Use for 'who rules my 7th house', 'chart ruler', 'planetary rulerships'.",
		Title:       "Western Natal Rulerships",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/rulerships", in.toQuery())
	})
}

func registerWesternTransitPositions(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_transit_positions",
		Description: "Get current (or any-date) tropical planetary positions for transit analysis — all planets with sign, degree, speed, and retrograde status. Use for 'where are the planets right now in Western astrology', 'current planetary transits', 'what sign is Saturn transiting'.",
		Title:       "Western Transit Positions",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternMomentInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/transits/positions", in.toQuery())
	})
}

// =====================================================================
// 5–7: Western Relationships
// =====================================================================

func registerWesternSynastryAspects(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_synastry_aspects",
		Description: "Compute the cross-chart aspects between two natal charts (synastry). Returns all significant inter-chart aspects with orbs — essential for relationship compatibility analysis. Use for 'synastry between two charts', 'what aspects connect our charts', 'compatibility aspects'.",
		Title:       "Western Synastry Aspects",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/synastry/aspects", in.toQuery())
	})
}

func registerWesternSynastryScore(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_synastry_score",
		Description: "Score the compatibility between two natal charts (Western synastry). Returns a numeric score and breakdown by aspect quality and house overlays. Use for 'compatibility score between two people', 'how compatible are we astrologically', 'synastry score'.",
		Title:       "Western Synastry Compatibility Score",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/synastry/score", in.toQuery())
	})
}

func registerWesternCompositeChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_composite_chart",
		Description: "Compute the composite chart (midpoint method) for two people — the relationship chart showing the energies of the union itself, not either individual. Returns composite planetary positions and house cusps. Use for 'composite chart for a couple', 'relationship chart', 'midpoint composite'.",
		Title:       "Western Composite Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/composite/chart", in.toQuery())
	})
}

// =====================================================================
// 8–10: Western Projection & Returns
// =====================================================================

// WesternProgressionInput adds age_in_years to birth data.
type WesternProgressionInput struct {
	WesternBirthInput
	AgeInYears float64 `json:"age_in_years" jsonschema:"age at the progression date in decimal years (e.g. 35.5 for 35 years 6 months)"`
}

func registerWesternProgressionsPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_progressions_planets",
		Description: "Compute secondary progressed planetary positions for a birth chart at a given age. Each year of life progresses one day in ephemeris time. Use for 'progressed chart', 'secondary progressions', 'where is my progressed Sun', 'progressed planets at age 40'.",
		Title:       "Western Secondary Progressions",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternProgressionInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("age_in_years", strconv.FormatFloat(in.AgeInYears, 'f', -1, 64))
		return callPassthrough(ctx, c, "/v1/western/progressions/planets", q)
	})
}

func registerWesternSolarArcPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_solar_arc_planets",
		Description: "Compute Solar Arc directed planetary positions — each planet advances by the Sun's progressed arc (~1°/year). Returns directed positions for all planets. Use for 'solar arc directions', 'directed chart', 'solar arc Sun at 40', 'Ebertin solar arc'.",
		Title:       "Western Solar Arc Directions",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternProgressionInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("age_in_years", strconv.FormatFloat(in.AgeInYears, 'f', -1, 64))
		return callPassthrough(ctx, c, "/v1/western/solar-arc/planets", q)
	})
}

// WesternSolarReturnInput adds return_year to birth data.
type WesternSolarReturnInput struct {
	WesternBirthInput
	ReturnYear int `json:"return_year" jsonschema:"year for which to compute the solar return (e.g. 2026)"`
}

func registerWesternSolarReturn(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_solar_return",
		Description: "Compute the Solar Return chart for a given year — the chart for the exact moment the Sun returns to its natal degree. Returns all planets + ascendant for the solar return year. Use for 'solar return for 2026', 'birthday chart', 'annual chart'.",
		Title:       "Western Solar Return Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternSolarReturnInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("return_year", strconv.Itoa(in.ReturnYear))
		return callPassthrough(ctx, c, "/v1/western/returns/solar", q)
	})
}

// =====================================================================
// 11–13: Western Transits & Narrative
// =====================================================================

// WesternTransitToNatalInput holds birth + transit date.
type WesternTransitToNatalInput struct {
	BirthLat  float64 `json:"birth_lat" jsonschema:"birth latitude in decimal degrees"`
	BirthLon  float64 `json:"birth_lon" jsonschema:"birth longitude in decimal degrees"`
	BirthDate string  `json:"birth_date" jsonschema:"birth date YYYY-MM-DD"`
	BirthTime string  `json:"birth_time" jsonschema:"birth time HH:MM"`
	BirthTz   string  `json:"birth_tz" jsonschema:"birth IANA timezone"`
	Date      string  `json:"date" jsonschema:"transit date YYYY-MM-DD"`
	Time      string  `json:"time" jsonschema:"transit time HH:MM"`
	Tz        string  `json:"tz" jsonschema:"transit IANA timezone"`
}

func (t WesternTransitToNatalInput) toQuery() url.Values {
	q := url.Values{}
	q.Set("birth.lat", strconv.FormatFloat(t.BirthLat, 'f', -1, 64))
	q.Set("birth.lon", strconv.FormatFloat(t.BirthLon, 'f', -1, 64))
	q.Set("birth.date", t.BirthDate)
	q.Set("birth.time", t.BirthTime)
	q.Set("birth.tz", t.BirthTz)
	q.Set("date", t.Date)
	q.Set("time", t.Time)
	q.Set("tz", t.Tz)
	return q
}

func registerWesternTransitsToNatal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_transits_to_natal",
		Description: "Compute all current transiting planet aspects to natal planets — the real-time influences on a birth chart. Returns each transit planet's aspect to each natal planet with orb and applying/separating status. Use for 'current transits to my natal chart', 'what transits am I having', 'Saturn transiting my Sun'.",
		Title:       "Western Transits to Natal",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternTransitToNatalInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/transits/to-natal", in.toQuery())
	})
}

// AspectInterpInput is for the static aspect interpretation endpoint.
type AspectInterpInput struct {
	PlanetA string `json:"planet_a" jsonschema:"first planet name (e.g. Sun, Moon, Mars, Venus, Jupiter, Saturn, Uranus, Neptune, Pluto)"`
	PlanetB string `json:"planet_b" jsonschema:"second planet name"`
	Aspect  string `json:"aspect" jsonschema:"aspect type: conjunction, opposition, trine, square, sextile"`
	Context string `json:"context,omitempty" jsonschema:"optional chart context: natal (default) or synastry"`
}

func registerWesternInterpretAspect(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_interpretation_aspect",
		Description: "Get a narrative interpretation for a specific natal or synastry aspect (e.g. Sun square Mars). Returns headline, body text, and caveats. Use for 'what does Sun square Mars mean', 'interpret Venus trine Jupiter', 'meaning of Moon opposite Saturn'.",
		Title:       "Western Aspect Interpretation",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AspectInterpInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("planet_a", in.PlanetA)
		q.Set("planet_b", in.PlanetB)
		q.Set("aspect", in.Aspect)
		if in.Context != "" {
			q.Set("context", in.Context)
		}
		return callPassthrough(ctx, c, "/v1/western/interpretation/aspect", q)
	})
}

// TransitInterpInput is for the transit interpretation endpoint.
type TransitInterpInput struct {
	TransitPlanet string `json:"transit_planet" jsonschema:"the transiting planet (e.g. Saturn, Jupiter, Mars)"`
	NatalPlanet   string `json:"natal_planet" jsonschema:"the natal planet being transited (e.g. Sun, Moon, Venus)"`
	Aspect        string `json:"aspect" jsonschema:"aspect type: conjunction, opposition, trine, square, sextile"`
}

func registerWesternInterpretTransit(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_interpretation_transit",
		Description: "Get a narrative interpretation for a transiting planet aspecting a natal planet (e.g. Saturn transiting square natal Sun). Returns themes, timing advice, and caveats. Use for 'what does Saturn square natal Sun mean', 'Jupiter conjunct natal Moon meaning', 'transit interpretation'.",
		Title:       "Western Transit Interpretation",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TransitInterpInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("transit_planet", in.TransitPlanet)
		q.Set("natal_planet", in.NatalPlanet)
		q.Set("aspect", in.Aspect)
		return callPassthrough(ctx, c, "/v1/western/interpretation/transit", q)
	})
}

// =====================================================================
// 14–15: Western Lots & Midpoints
// =====================================================================

func registerWesternNatalLots(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_lots",
		Description: "Compute all 12 Hellenistic Arabic Parts (Lots) for a Western tropical natal chart — Fortune, Spirit, Eros, Necessity, Courage, Victory, and more. Returns each lot's degree, sign, and house. Use for 'Part of Fortune', 'Arabic Parts', 'Lot of Spirit in my chart'.",
		Title:       "Western Natal Arabic Parts / Lots",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/lots", in.toQuery())
	})
}

func registerWesternNatalMidpoints(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_midpoints",
		Description: "Compute all planet-pair midpoints in a Western natal chart and identify which planets activate each midpoint within orb. Essential for cosmobiology / Ebertin analysis. Use for 'midpoints in my chart', 'Sun/Moon midpoint', 'activated midpoints'.",
		Title:       "Western Natal Midpoints",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/midpoints", in.toQuery())
	})
}

// =====================================================================
// 16–19: Hellenistic Time Lords
// =====================================================================

// ProfectionInput holds birth date + optional reference date + asc sign.
type ProfectionInput struct {
	BirthDate string `json:"birth_date" jsonschema:"birth date in YYYY-MM-DD format"`
	AscSign   string `json:"asc_sign" jsonschema:"ascending sign at birth: aries/taurus/gemini/cancer/leo/virgo/libra/scorpio/sagittarius/capricorn/aquarius/pisces (or 0-11)"`
	RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
}

func registerWesternProfectionsAnnual(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_profections_annual",
		Description: "Compute Annual Profections — the Hellenistic annual time-lord technique where the Ascendant advances one sign per year, activating that sign's ruler as the year's lord. Returns the profected sign, lord, and activated house. Use for 'annual profections', 'profection for age 35', 'year lord', 'time lord techniques'.",
		Title:       "Western Annual Profections",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ProfectionInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		q.Set("asc_sign", in.AscSign)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		return callPassthrough(ctx, c, "/v1/western/profections/annual", q)
	})
}

// ZRInput holds birth data for Zodiacal Releasing.
type ZRInput struct {
	BirthDate string `json:"birth_date" jsonschema:"birth date in YYYY-MM-DD format"`
	LotSign   string `json:"lot_sign" jsonschema:"sign of the lot to release from — typically Lot of Fortune or Spirit sign (e.g. cancer, scorpio)"`
	RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
}

func registerWesternZodiacalReleasing(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_zodiacal_releasing",
		Description: "Compute Zodiacal Releasing from the Lot of Fortune — the Hellenistic predictive technique assigning planetary periods by releasing from a lot sign through the zodiac. Returns current Level 1 and Level 2 time-lord periods. Use for 'zodiacal releasing', 'Lot of Fortune periods', 'peak periods from Fortune'.",
		Title:       "Western Zodiacal Releasing (Fortune)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ZRInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		q.Set("lot_sign", in.LotSign)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		return callPassthrough(ctx, c, "/v1/western/zodiacal-releasing/fortune", q)
	})
}

// FirdariaInput holds only birth date and optional reference date —
// Firdaria does not require ascendant sign or geographic coordinates.
type FirdariaInput struct {
	BirthDate string `json:"birth_date" jsonschema:"birth date in YYYY-MM-DD format"`
	RefDate   string `json:"ref_date,omitempty" jsonschema:"optional reference date YYYY-MM-DD (defaults to today)"`
}

func registerWesternFirdaria(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_firdaria",
		Description: "Compute Firdaria planetary periods — the Medieval/Persian 75-year planetary time-lord cycle. Each planet rules a period proportional to its Chaldean years. Returns the current major and minor Firdaria lords. Use for 'firdaria periods', 'planetary periods', 'current firdaria lord'.",
		Title:       "Western Firdaria Periods",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in FirdariaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("birth_date", in.BirthDate)
		if in.RefDate != "" {
			q.Set("ref_date", in.RefDate)
		}
		return callPassthrough(ctx, c, "/v1/western/firdaria", q)
	})
}

// DignityInput holds planet + birth context for dignities.
type DignityInput struct {
	WesternBirthInput
	Planet string `json:"planet" jsonschema:"planet to assess: Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn"`
	Scheme string `json:"scheme,omitempty" jsonschema:"dignity scheme: dorothean (default), ptolemaic, or lilly"`
}

func registerWesternDignities(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_planet",
		Description: "Compute all essential dignities for a planet in a natal chart — domicile, exaltation, detriment, fall, triplicity lord, bound lord, and face lord. Returns an Almuten score. Use for 'is my Venus dignified', 'essential dignities', 'planetary strength', 'Almuten of planet'.",
		Title:       "Western Essential Dignities",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DignityInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("planet", in.Planet)
		if in.Scheme != "" {
			q.Set("scheme", in.Scheme)
		}
		return callPassthrough(ctx, c, "/v1/western/dignities/planet", q)
	})
}

// =====================================================================
// 20–21: Extended Bodies
// =====================================================================

func registerWesternNatalFixedStars(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_fixed_stars",
		Description: "Compute conjunctions between natal planets and major fixed stars (Brady canon of 47 stars — Algol, Regulus, Spica, Antares, etc.) within a 1° orb. Use for 'fixed stars in my chart', 'Regulus conjunct my Sun', 'Algol conjunction', 'natal fixed star analysis'.",
		Title:       "Western Natal Fixed Stars",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/fixed-stars", in.toQuery())
	})
}

func registerWesternNatalBigFourAsteroids(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_big_four_asteroids",
		Description: "Compute natal positions of the four major asteroids — Ceres, Pallas, Juno, Vesta — plus Chiron (Wounded Healer). Returns sign, degree, and house for each. Use for 'asteroid positions in my chart', 'Chiron placement', 'Ceres sign', 'big four asteroids'.",
		Title:       "Western Big Four Asteroids + Chiron",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/asteroids/big-four-plus", in.toQuery())
	})
}

// =====================================================================
// 22–25: Cosmobiology
// =====================================================================

func registerWesternNatalDeclinations(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_declinations",
		Description: "Compute natal planetary declinations — the celestial latitude north or south of the equator. Identifies parallel aspects (same declination = same-sign energy) and contra-parallel aspects (opposite declination = opposition energy), plus out-of-bounds planets beyond ±23.4°. Use for 'declinations in my chart', 'parallel aspects', 'out-of-bounds planets'.",
		Title:       "Western Natal Declinations",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/declinations", in.toQuery())
	})
}

// HarmonicInput adds harmonic number to birth data.
type HarmonicInput struct {
	WesternBirthInput
	N int `json:"n" jsonschema:"harmonic number (2-32); 4th harmonic reveals squares, 5th quintiles, 7th septiles, etc."`
}

func registerWesternNatalHarmonic(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_harmonic",
		Description: "Compute the Nth harmonic chart (all planetary positions multiplied by N modulo 360°). Harmonics reveal hidden aspect patterns: H4 = squares, H5 = quintiles, H7 = septiles. Use for 'harmonic chart', '5th harmonic', 'quintile chart', 'John Addey harmonics'.",
		Title:       "Western Natal Harmonic Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in HarmonicInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		path := "/v1/western/natal/harmonic/" + strconv.Itoa(in.N)
		return callPassthrough(ctx, c, path, q)
	})
}

func registerWesternNatalAntiscia(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_antiscia",
		Description: "Compute natal antiscia (solstice points) and contra-antiscia — mirror positions across the Cancer/Capricorn and Aries/Libra axes. Antiscia contacts between planets indicate hidden connections. Use for 'antiscia in my chart', 'solstice points', 'contra-antiscia aspects'.",
		Title:       "Western Natal Antiscia",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/antiscia", in.toQuery())
	})
}

func registerWesternNatalMidpointTree(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_midpoint_tree",
		Description: "Compute the Ebertin-style 90° dial midpoint tree — identifies all activated midpoints (planet on the midpoint of two others within orb). Returns a tree structure with Ebertin keyword weights. Use for 'midpoint tree', '90 degree dial', 'Ebertin midpoints', 'cosmobiology midpoints'.",
		Title:       "Western Natal Midpoint Tree (Ebertin 90° Dial)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/midpoint-tree", in.toQuery())
	})
}

// =====================================================================
// 26–30: Mundane Astrology & Astrocartography
// =====================================================================

func registerWesternHeliocentricPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_heliocentric_planets",
		Description: "Compute heliocentric (Sun-centred) planetary positions — where the planets actually are in the solar system as seen from the Sun, not Earth. Returns Earth's position too. Use for 'heliocentric chart', 'Sun-centred positions', 'heliocentric astrology'.",
		Title:       "Western Heliocentric Planets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternMomentInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/heliocentric/planets", in.toQuery())
	})
}

// EclipsesUpcomingInput has optional months param.
type EclipsesUpcomingInput struct {
	Months int `json:"months,omitempty" jsonschema:"optional number of months ahead to search (1-60, default 12)"`
}

func registerWesternEclipsesUpcoming(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_eclipses_upcoming",
		Description: "Find upcoming solar and lunar eclipses from today forward (default 12 months ahead). Returns eclipse type, UTC peak time, and whether it is total, partial, annular, or penumbral. Use for 'upcoming eclipses', 'next solar eclipse', 'eclipse dates', 'when is the next lunar eclipse'.",
		Title:       "Western Upcoming Eclipses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in EclipsesUpcomingInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Months > 0 {
			q.Set("months", strconv.Itoa(in.Months))
		}
		return callPassthrough(ctx, c, "/v1/western/eclipses/upcoming", q)
	})
}

func registerWesternIngresses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_ingresses",
		Description: "Find all planetary sign-change ingresses in a date range — when each planet enters a new tropical zodiac sign, including retrograde re-entries. Use for 'when does Saturn change signs', 'ingresses this year', 'Jupiter entering Gemini', 'sign changes'.",
		Title:       "Western Planetary Ingresses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DateRangeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/ingresses", in.toQuery())
	})
}

// RetrogradeWindowInput adds planet name to date range.
type RetrogradeWindowInput struct {
	DateRangeInput
	Planet string `json:"planet" jsonschema:"planet name: Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, or Pluto"`
}

func registerWesternRetrogradeWindow(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_retrograde_window",
		Description: "Find all retrograde station dates (station retrograde → station direct) for a planet within a date range. Returns exact station moments and the sign during the retrograde peak. Use for 'Mercury retrograde dates', 'when is Mercury retrograde', 'Saturn retrograde this year', 'retrograde periods'.",
		Title:       "Western Retrograde Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RetrogradeWindowInput) (*mcp.CallToolResult, any, error) {
		q := in.DateRangeInput.toQuery()
		path := "/v1/western/retrograde-window/" + in.Planet
		return callPassthrough(ctx, c, path, q)
	})
}

// ACGInput holds natal moment and optional lat_step.
type ACGInput struct {
	WesternMomentInput
	LatStep float64 `json:"lat_step,omitempty" jsonschema:"optional latitude sampling step in degrees (0.5-5.0, default 2.0 — higher is faster, lower is more precise)"`
}

func registerWesternACGPlanetLines(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_astrocartography_planet_lines",
		Description: "Compute Astrocartography (ACG) planet lines for a natal moment — traces the geographic lines where each planet rises (ASC), sets (DSC), culminates (MC), or is at nadir (IC) at birth. Returns polylines in lat/lon for mapping. Use for 'astrocartography chart', 'relocation astrology', 'where should I live', 'ACG lines'.",
		Title:       "Western Astrocartography Planet Lines",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ACGInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternMomentInput.toQuery()
		if in.LatStep > 0 {
			q.Set("lat_step", strconv.FormatFloat(in.LatStep, 'f', -1, 64))
		}
		return callPassthrough(ctx, c, "/v1/western/astrocartography/planet-lines", q)
	})
}
