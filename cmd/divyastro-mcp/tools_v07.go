package main

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV07Tools adds the Vedic Narrative family, Reports family,
// Vedic Remedies, and remaining Western astrology endpoints.
//
//   Vedic Narrative (32): career-outlook, dasha-phal, dasha-tree-phal,
//     doshas, finance-outlook, horoscope/daily-by-lagna,
//     horoscope/daily-by-moon, horoscope/daily-tamil,
//     horoscope/weekly-by-lagna, horoscope/weekly-by-moon,
//     house-lord, house-lords-overview, karakas, lagna,
//     marriage-outlook, milan, moon-sign, nakshatra, planet-in-house,
//     planet-in-sign, planets-in-houses, planets-in-signs, profile,
//     sade-sati-history, sade-sati-phase, sade-sati-status, sun-sign,
//     transit-ashtakavarga, transit-double, transit-phal,
//     varshaphal-themes, yogas
//   Reports (12): dasha-analysis, horoscope/daily, horoscope/monthly,
//     horoscope/weekly, kundli/brihad, kundli/detailed, kundli/lite,
//     mangal-dosha, match-making, numerology, sade-sati, varshaphal
//   Vedic Remedies (1): full family endpoint
//   Western (15): astrocartography/local-space, composite/houses,
//     composite/planets, davison/houses, davison/planets,
//     dignities/receptions, narrative/timelord/firdaria,
//     natal/ascendant, natal/mc, natal/out-of-bounds, natal/planets,
//     returns/solar/exact-jd, narrative/transit-summary,
//     horoscope/monthly, horoscope/weekly
//
// Total new: 60 tools.

func registerV07Tools(s *mcp.Server, c *apiClient) {
	// Vedic Narrative
	registerNarrativeCareerOutlook(s, c)
	registerNarrativeDashaPhal(s, c)
	registerNarrativeDashaTreePhal(s, c)
	registerNarrativeDoshas(s, c)
	registerNarrativeFinanceOutlook(s, c)
	registerNarrativeHoroscopeDailyByLagna(s, c)
	registerNarrativeHoroscopeDailyByMoon(s, c)
	registerNarrativeHoroscopeDailyTamil(s, c)
	registerNarrativeHoroscopeWeeklyByLagna(s, c)
	registerNarrativeHoroscopeWeeklyByMoon(s, c)
	registerNarrativeHouseLord(s, c)
	registerNarrativeHouseLordsOverview(s, c)
	registerNarrativeKarakas(s, c)
	registerNarrativeLagna(s, c)
	registerNarrativeMarriageOutlook(s, c)
	registerNarrativeMilan(s, c)
	registerNarrativeMoonSign(s, c)
	registerNarrativeNakshatra(s, c)
	registerNarrativePlanetInHouse(s, c)
	registerNarrativePlanetInSign(s, c)
	registerNarrativePlanetsInHouses(s, c)
	registerNarrativePlanetsInSigns(s, c)
	registerNarrativeProfile(s, c)
	registerNarrativeSadeSatiHistory(s, c)
	registerNarrativeSadeSatiPhase(s, c)
	registerNarrativeSadeSatiStatus(s, c)
	registerNarrativeSunSign(s, c)
	registerNarrativeTransitAshtakavarga(s, c)
	registerNarrativeTransitDouble(s, c)
	registerNarrativeTransitPhal(s, c)
	registerNarrativeVarshaphalThemes(s, c)
	registerNarrativeYogas(s, c)

	// Reports
	registerReportDashaAnalysis(s, c)
	registerReportHoroscopeDaily(s, c)
	registerReportHoroscopeMonthly(s, c)
	registerReportHoroscopeWeekly(s, c)
	registerReportKundliBrihad(s, c)
	registerReportKundliDetailed(s, c)
	registerReportKundliLite(s, c)
	registerReportMangalDosha(s, c)
	registerReportMatchMaking(s, c)
	registerReportNumerology(s, c)
	registerReportSadeSati(s, c)
	registerReportVarshaphal(s, c)

	// Remedies
	registerVedicRemedies(s, c)

	// Western missing
	registerWesternAstrocartographyLocalSpace(s, c)
	registerWesternCompositeHouses(s, c)
	registerWesternCompositePlanets(s, c)
	registerWesternDavisonHouses(s, c)
	registerWesternDavisonPlanets(s, c)
	registerWesternDignitiesReceptions(s, c)
	registerWesternNarrativeFirdaria(s, c)
	registerWesternNatalAscendant(s, c)
	registerWesternNatalMC(s, c)
	registerWesternNatalOutOfBounds(s, c)
	registerWesternNatalPlanets(s, c)
	registerWesternSolarReturnExactJD(s, c)
	// western_transit_summary already registered correctly in registerV04Tools.
	// western_horoscope_monthly + weekly already registered in registerV04Tools.
}

// =====================================================================
// Shared input helpers for narrative + reports
// =====================================================================

type NarrativeBirthInput struct {
	BirthInput
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func (n NarrativeBirthInput) toQuery() url.Values {
	q := n.BirthInput.toQuery()
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

// toDualQuery emits birth.* prefixed params for endpoints that use
// serveVedicDualMoment (which calls parsePrefixedMoment(r, "birth.")).
func (n NarrativeBirthInput) toDualQuery() url.Values {
	q := birthPrefixQuery(n.BirthInput)
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

type NarrativePlanetInput struct {
	BirthInput
	Planet string `json:"planet" jsonschema:"planet name: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu"`
	Lang   string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func (n NarrativePlanetInput) toQuery() url.Values {
	q := n.BirthInput.toQuery()
	q.Set("planet", n.Planet)
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

type NarrativeHouseInput struct {
	BirthInput
	House int    `json:"house" jsonschema:"house number 1–12"`
	Lang  string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func (n NarrativeHouseInput) toQuery() url.Values {
	q := n.BirthInput.toQuery()
	q.Set("house", strconv.Itoa(n.House))
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

// =====================================================================
// Vedic Narrative endpoints (32 tools)
// =====================================================================

func registerNarrativeCareerOutlook(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_career_outlook",
		Description: "Get a multi-paragraph career and profession narrative for a Vedic birth chart — analyses the 10th house lord, planets in the 10th, D10 chart, and dashas to describe career strengths, suitable professions, and current period's impact on work. Use for 'career reading', 'what profession suits me', 'career analysis from my chart'.",
		Title:       "Vedic Career Outlook Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/career-outlook", in.toQuery())
	})
}

func registerNarrativeDashaPhal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_dasha_phal",
		Description: "Get the Dasha Phal (period interpretation) narrative for the currently running Vimshottari Mahadasha + Antardasha combination — a human-readable explanation of what this dasha pair typically brings and how it interacts with the natal chart. Use for 'what does my current dasha mean', 'dasha phal', 'interpret my mahadasha and antardasha'.",
		Title:       "Dasha Phal Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/dasha-phal", in.toDualQuery())
	})
}

type NarrativeDashaTreeInput struct {
	BirthInput
	Maha  string `json:"maha" jsonschema:"Mahadasha lord (e.g. Saturn, Jupiter, Rahu)"`
	Antar string `json:"antar,omitempty" jsonschema:"optional Antardasha lord; if omitted returns all antardashas for the given Mahadasha"`
	Lang  string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func registerNarrativeDashaTreePhal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_dasha_tree_phal",
		Description: "Get narrative interpretations for a specific dasha combination — pass mahadasha lord (and optionally antardasha lord) to get the meaning of that specific period for the natal chart. Useful for asking about a future dasha period. Use for 'what will Saturn mahadasha bring', 'Jupiter-Rahu period meaning for my chart'.",
		Title:       "Dasha Tree Phal Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeDashaTreeInput) (*mcp.CallToolResult, any, error) {
		q := birthPrefixQuery(in.BirthInput)
		q.Set("maha", in.Maha)
		if in.Antar != "" {
			q.Set("antar", in.Antar)
		}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/dasha-tree-phal", q)
	})
}

func registerNarrativeDoshas(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_doshas",
		Description: "Get the Doshas narrative for a birth chart — comprehensive analysis of all major afflictions (Mangal Dosha, Kaal Sarp Dosha, Pitru Dosha, Shrapit Dosha, Guru Chandal Yoga, Grahan Yoga) with presence, severity, cancellations, and remedial guidance. Use for 'what doshas do I have', 'check all doshas in my chart', 'dosha analysis'.",
		Title:       "Doshas Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/doshas", in.toQuery())
	})
}

func registerNarrativeFinanceOutlook(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_finance_outlook",
		Description: "Get a finance and wealth narrative for a Vedic birth chart — analyses the 2nd and 11th house lords, Jupiter, Venus, and wealth yogas to describe financial strengths, challenges, and periods of prosperity. Use for 'financial analysis of my chart', 'wealth reading', 'when will I earn more'.",
		Title:       "Vedic Finance Outlook Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/finance-outlook", in.toQuery())
	})
}

type NarrativeHoroscopeInput struct {
	BirthInput
	Date string `json:"date,omitempty" jsonschema:"optional date YYYY-MM-DD; defaults to today"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func (n NarrativeHoroscopeInput) toQuery() url.Values {
	q := n.BirthInput.toQuery()
	if n.Date != "" {
		q.Set("date", n.Date)
	}
	if n.Lang != "" {
		q.Set("lang", n.Lang)
	}
	return q
}

func registerNarrativeHoroscopeDailyByLagna(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_horoscope_daily_by_lagna",
		Description: "Get a personalised daily horoscope narrative computed from the rising sign (lagna) of a birth chart — more precise than generic Moon-sign horoscopes since it uses the exact birth data. Use for 'daily horoscope by ascendant', 'lagna-based rashifal', 'personalised daily reading'.",
		Title:       "Daily Horoscope by Lagna (Personalised)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeHoroscopeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/horoscope/daily-by-lagna", in.toQuery())
	})
}

func registerNarrativeHoroscopeDailyByMoon(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_horoscope_daily_by_moon",
		Description: "Get a personalised daily horoscope narrative computed from the natal Moon sign with full transit context — richer than a generic rashifal because it uses the birth chart's Moon degree and current transit aspects. Use for 'personalised Moon sign horoscope', 'chandra rashi horoscope with transits'.",
		Title:       "Daily Horoscope by Natal Moon (Personalised)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeHoroscopeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/horoscope/daily-by-moon", in.toQuery())
	})
}

type NarrativeTamilHoroscopeInput struct {
	BirthInput
	Date string `json:"date,omitempty" jsonschema:"optional date YYYY-MM-DD; defaults to today"`
}

func registerNarrativeHoroscopeDailyTamil(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_horoscope_daily_tamil",
		Description: "Get a daily horoscope in Tamil astrology style — uses the Tamil rashi system and traditional Tamil Panchangam for the day's reading. Use for 'Tamil horoscope today', 'தினசரி ராசிபலன்', 'Tamil rasi palan'.",
		Title:       "Daily Tamil Horoscope",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeTamilHoroscopeInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/horoscope/daily-tamil", q)
	})
}

func registerNarrativeHoroscopeWeeklyByLagna(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_horoscope_weekly_by_lagna",
		Description: "Get a personalised weekly horoscope narrative computed from the birth chart's rising sign (lagna) — personalised 7-day outlook with major transit themes. Use for 'personalised weekly horoscope', 'weekly lagna rashifal'.",
		Title:       "Weekly Horoscope by Lagna (Personalised)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeHoroscopeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/horoscope/weekly-by-lagna", in.toQuery())
	})
}

func registerNarrativeHoroscopeWeeklyByMoon(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_horoscope_weekly_by_moon",
		Description: "Get a personalised weekly horoscope narrative computed from the natal Moon sign — personalised weekly reading using the birth chart's Moon degree. Use for 'personalised weekly Moon sign horoscope', 'weekly chandra rashi reading'.",
		Title:       "Weekly Horoscope by Natal Moon (Personalised)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeHoroscopeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/horoscope/weekly-by-moon", in.toQuery())
	})
}

func registerNarrativeHouseLord(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_house_lord",
		Description: "Get a narrative interpretation of a specific house lord's placement — where the lord of house N is placed in the chart and what it means for the areas of life that house governs. Use for 'interpret my 7th lord', '10th house lord analysis', 'what does my lagna lord placement mean'.",
		Title:       "House Lord Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeHouseInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/house-lord", in.toQuery())
	})
}

func registerNarrativeHouseLordsOverview(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_house_lords_overview",
		Description: "Get a comprehensive narrative overview of all 12 house lords' placements in a birth chart — a full bhavesh analysis summarising how each house lord's position shapes the chart. Use for 'all house lords analysis', 'complete bhavesh reading', 'house lords overview'.",
		Title:       "House Lords Overview Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/house-lords-overview", in.toQuery())
	})
}

func registerNarrativeKarakas(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_karakas",
		Description: "Get the Jaimini Chara Karakas for a birth chart — the 8 karaka planets (Atmakaraka, Amatyakaraka, Bhratrukaraka, Matrukaraka, Putrakaraka, Gnatikaraka, Darakaraka, Dara/Swamsha) determined by planet degrees, with interpretation. Use for 'what is my atmakaraka', 'chara karakas for my chart', 'soul planet Jaimini'.",
		Title:       "Jaimini Chara Karakas",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/karakas", in.toQuery())
	})
}

func registerNarrativeLagna(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_lagna",
		Description: "Get a narrative interpretation of the birth chart's lagna (ascendant) — describes the physical appearance, personality, health tendencies, and life approach associated with the rising sign, including first-house planets' influence. Use for 'lagna interpretation', 'rising sign meaning', 'ascendant reading', 'what does my lagna say about me'.",
		Title:       "Lagna (Ascendant) Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/lagna", in.toQuery())
	})
}

func registerNarrativeMarriageOutlook(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_marriage_outlook",
		Description: "Get a marriage and relationship narrative for a Vedic birth chart — analyses the 7th house lord, Venus, navamsa, and marriage-timing dashas to describe relationship characteristics, likely spouse traits, and timing of marriage. Use for 'marriage reading', 'when will I get married', 'spouse characteristics from my chart', 'vivah analysis'.",
		Title:       "Vedic Marriage Outlook Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/marriage-outlook", in.toQuery())
	})
}

type NarrativeMilanInput struct {
	BrideBirth BirthInput `json:"bride"`
	GroomBirth BirthInput `json:"groom"`
	Lang       string     `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerNarrativeMilan(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_milan",
		Description: "Get a narrative compatibility analysis for a couple — goes beyond the numeric Ashtakoota score with a written interpretation of compatibility across all factors (guna milan, mangal dosha, dashas, navamsa). Use for 'detailed match making narrative', 'compatibility reading for couple', 'written kundli matching report'.",
		Title:       "Milan (Compatibility) Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeMilanInput) (*mcp.CallToolResult, any, error) {
		// Handler reads boy.date/boy.time/boy.tz and girl.date/girl.time/girl.tz
		// via parsePrefixedMoment(r, "boy.") / parsePrefixedMoment(r, "girl.").
		// The bride_*/groom_* params from matchMakingQuery are never read here.
		q := url.Values{}
		q.Set("boy.date", in.BrideBirth.Date)
		q.Set("boy.time", in.BrideBirth.Time)
		q.Set("boy.tz", in.BrideBirth.Tz)
		q.Set("girl.date", in.GroomBirth.Date)
		q.Set("girl.time", in.GroomBirth.Time)
		q.Set("girl.tz", in.GroomBirth.Tz)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/milan", q)
	})
}

func registerNarrativeMoonSign(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_moon_sign",
		Description: "Get a narrative interpretation of the natal Moon sign (Chandra Rashi) — covers emotional nature, mental tendencies, mother relationship, and instinctive responses associated with the Moon's sign placement. Use for 'Moon sign reading', 'chandra rashi meaning', 'Moon in Scorpio interpretation'.",
		Title:       "Moon Sign Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/moon-sign", in.toQuery())
	})
}

func registerNarrativeNakshatra(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_nakshatra",
		Description: "Get a narrative interpretation of the birth nakshatra (Janma Nakshatra) — covers personality, strengths, weaknesses, compatible nakshatras, ruling deity, symbol, and life path themes of the nakshatra the Moon occupied at birth. Use for 'nakshatra interpretation', 'what does my janma nakshatra mean', 'Rohini nakshatra reading'.",
		Title:       "Birth Nakshatra Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/nakshatra", in.toQuery())
	})
}

func registerNarrativePlanetInHouse(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_planet_in_house",
		Description: "Get a narrative interpretation of a specific planet's house placement in a Vedic birth chart — describes what that planet's position in that house means for the native's life areas, including classical significations and any classical yoga formation. Use for 'Jupiter in 5th house meaning', 'what does Saturn in 7th house mean', 'planet house interpretation'.",
		Title:       "Planet in House Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativePlanetInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/planet-in-house", in.toQuery())
	})
}

func registerNarrativePlanetInSign(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_planet_in_sign",
		Description: "Get a narrative interpretation of a specific planet's sign placement — describes how the planet expresses itself in that rashi, classical dignity/debility implications, and characteristic traits. Use for 'Saturn in Libra meaning', 'Mars in Aries interpretation', 'planet sign narrative'.",
		Title:       "Planet in Sign Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativePlanetInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/planet-in-sign", in.toQuery())
	})
}

func registerNarrativePlanetsInHouses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_planets_in_houses",
		Description: "Get narrative interpretations for ALL planets' house placements in a birth chart in a single call — the full house-placement reading for every graha. Use when you want a comprehensive house-placement analysis rather than one planet at a time. Use for 'all planets in houses reading', 'complete bhava placement analysis'.",
		Title:       "All Planets in Houses Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/planets-in-houses", in.toQuery())
	})
}

func registerNarrativePlanetsInSigns(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_planets_in_signs",
		Description: "Get narrative interpretations for ALL planets' sign placements in a birth chart in a single call — the full sign-placement reading for every graha. Use for 'all planet signs reading', 'complete rashi placement analysis', 'what signs are all my planets in'.",
		Title:       "All Planets in Signs Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/planets-in-signs", in.toQuery())
	})
}

func registerNarrativeProfile(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_profile",
		Description: "Get a holistic birth chart profile narrative — an executive summary of the entire chart covering lagna, Moon sign, key planetary strengths/weaknesses, dominant yogas, current dasha, and life-theme highlights in a single comprehensive reading. Use for 'full chart reading', 'birth chart interpretation', 'complete kundli analysis', 'summarise my horoscope'.",
		Title:       "Full Birth Chart Profile Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/profile", in.toQuery())
	})
}

func registerNarrativeSadeSatiHistory(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_sade_sati_history",
		Description: "Get a narrative of all past and future Sade Sati periods for a birth chart — written descriptions of each 7.5-year Saturn transit cycle with dates, life themes expected, and which phase (rising/peak/setting) the person experienced. Use for 'all Sade Sati periods', 'Sade Sati history narrative', 'when were all my Sade Sati cycles'.",
		Title:       "Sade Sati History Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/sade-sati-history", in.toDualQuery())
	})
}

func registerNarrativeSadeSatiPhase(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_sade_sati_phase",
		Description: "Get a narrative interpretation of the current Sade Sati phase (rising, peak, or setting) — what the current phase typically brings and remedies. Only meaningful when the person is actually in Sade Sati. Use for 'what does my current Sade Sati phase mean', 'Sade Sati phase interpretation', 'Saturn phase reading'.",
		Title:       "Sade Sati Phase Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/sade-sati-phase", in.toDualQuery())
	})
}

func registerNarrativeSadeSatiStatus(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_sade_sati_status",
		Description: "Get the current Sade Sati status with narrative context — whether the person is in Sade Sati, which phase, start/end dates, and a written description of the current transit's meaning for their Moon sign. Use for 'am I in Sade Sati with explanation', 'Sade Sati status and meaning'.",
		Title:       "Sade Sati Status Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/sade-sati-status", in.toDualQuery())
	})
}

func registerNarrativeSunSign(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_sun_sign",
		Description: "Get a narrative interpretation of the natal Sun sign (Surya Rashi) — covers ego, identity, father, authority, vitality, and the life-purpose themes associated with the Sun's sidereal sign placement. Use for 'Sun sign reading', 'Surya rashi interpretation', 'Sun in Capricorn meaning'.",
		Title:       "Sun Sign Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/sun-sign", in.toQuery())
	})
}

func registerNarrativeTransitAshtakavarga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_transit_ashtakavarga",
		Description: "Get a narrative interpretation of the Ashtakavarga transit scores for a birth chart — written analysis of how current planetary transits score against the natal chart's benefic-point grid, with guidance on which transits are favourable or challenging. Use for 'ashtakavarga transit reading', 'explain my transit scores', 'written ashtakavarga analysis'.",
		Title:       "Transit Ashtakavarga Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/transit-ashtakavarga", in.toDualQuery())
	})
}

func registerNarrativeTransitDouble(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_transit_double",
		Description: "Get a narrative of the Double Transit analysis — written explanation of which life events the current Jupiter-Saturn joint transit is activating in the chart, and when the confluence is strongest. Use for 'double transit reading', 'written double transit analysis', 'when will Jupiter-Saturn double transit trigger events for me'.",
		Title:       "Double Transit Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/transit-double", in.toDualQuery())
	})
}

func registerNarrativeTransitPhal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_transit_phal",
		Description: "Get the Transit Phal (gochar results) narrative — a written analysis of all current major planetary transits over the natal chart, their house positions from Moon, and the classical effects expected from each transit. Use for 'gochar phal', 'transit results for my chart', 'what do current transits mean for me'.",
		Title:       "Transit Phal Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/transit-phal", in.toDualQuery())
	})
}

func registerNarrativeVarshaphalThemes(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_varshaphal_themes",
		Description: "Get a narrative of the key themes for the current Varshaphal (solar-return) year — written analysis of the annual chart's varsha lord, muntha, annual lagna, and tajika yogas as a prediction for the year's major themes. Use for 'varshaphal themes this year', 'annual chart reading', 'solar return year prediction'.",
		Title:       "Varshaphal Year Themes Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/varshaphal-themes", in.toDualQuery())
	})
}

func registerNarrativeYogas(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_yogas",
		Description: "Get a narrative of all major Vedic yogas present in a birth chart — Raja Yoga, Dhana Yoga, Parivartana Yoga, Kemadruma Yoga, Neecha Bhanga, Panchamahapurusha, and others — with presence, strength, and what each yoga means for the native. Use for 'what yogas do I have', 'raj yoga in my chart', 'yoga analysis', 'all yogas reading'.",
		Title:       "Yogas Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/vedic/narrative/yogas", in.toQuery())
	})
}

// =====================================================================
// Reports (12 tools) — comprehensive PDF-ready structured reports
// =====================================================================

type ReportBirthInput struct {
	BirthInput
	Lang   string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
	Format string `json:"format,omitempty" jsonschema:"output format: json (default) or pdf"`
}

func (r ReportBirthInput) toQuery() url.Values {
	q := r.BirthInput.toQuery()
	if r.Lang != "" {
		q.Set("lang", r.Lang)
	}
	if r.Format != "" {
		q.Set("format", r.Format)
	}
	return q
}

func registerReportDashaAnalysis(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_dasha_analysis",
		Description: "Generate a comprehensive Dasha Analysis report — the full Vimshottari dasha timeline with narrative interpretation for each Mahadasha period, key life events timing, and dasha-transit correlations. Use for 'generate dasha report', 'full dasha analysis', 'lifetime dasha report'.",
		Title:       "Dasha Analysis Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/dasha-analysis", in.toQuery())
	})
}

type ReportHoroscopeInput struct {
	Rashi string `json:"rashi" jsonschema:"Moon sign / rashi: Aries, Taurus, ..., Pisces (case-insensitive)"`
	Date  string `json:"date,omitempty" jsonschema:"optional date YYYY-MM-DD; defaults to today"`
	Lang  string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerReportHoroscopeDaily(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_horoscope_daily",
		Description: "Generate a detailed daily horoscope report for a Moon sign — structured report with separate sections for career, love, finance, health, and daily guidance. More comprehensive than the quick daily reading. Use for 'detailed daily horoscope report', 'full rashifal report'.",
		Title:       "Daily Horoscope Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportHoroscopeInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("rashi", strings.ToLower(strings.TrimSpace(in.Rashi)))
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/reports/horoscope/daily", q)
	})
}

type ReportHoroscopeMonthInput struct {
	Rashi string `json:"rashi" jsonschema:"Moon sign / rashi: Aries, Taurus, ..., Pisces (case-insensitive)"`
	Year  int    `json:"year" jsonschema:"year e.g. 2026"`
	Month int    `json:"month" jsonschema:"month 1–12"`
	Lang  string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerReportHoroscopeMonthly(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_horoscope_monthly",
		Description: "Generate a detailed monthly horoscope report for a Moon sign — structured report covering the month's major transit themes, weekly breakdown, and guidance for career, relationships, and finance for the full month. Use for 'monthly horoscope report', 'full month rashifal', 'monthly prediction for Gemini'.",
		Title:       "Monthly Horoscope Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportHoroscopeMonthInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("rashi", strings.ToLower(strings.TrimSpace(in.Rashi)))
		q.Set("year", strconv.Itoa(in.Year))
		q.Set("month", strconv.Itoa(in.Month))
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/reports/horoscope/monthly", q)
	})
}

func registerReportHoroscopeWeekly(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_horoscope_weekly",
		Description: "Generate a detailed weekly horoscope report for a Moon sign — structured report covering the week's planetary movement, day-by-day highlights, and focus areas. Use for 'weekly horoscope report', 'full week rashifal', 'detailed weekly prediction'.",
		Title:       "Weekly Horoscope Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportHoroscopeInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("rashi", strings.ToLower(strings.TrimSpace(in.Rashi)))
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/reports/horoscope/weekly", q)
	})
}

func registerReportKundliBrihad(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_kundli_brihad",
		Description: "Generate the Brihad (comprehensive) Kundli report — the most detailed birth chart report covering chart analysis, all divisional charts (D1–D60), planetary strengths (Shadbala), ashtakavarga, all dashas, yogas, doshas, remedies, and lifetime predictions. Use for 'complete kundli report', 'brihad kundli', 'detailed horoscope report'.",
		Title:       "Brihad Kundli Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/kundli/brihad", in.toQuery())
	})
}

func registerReportKundliDetailed(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_kundli_detailed",
		Description: "Generate a Detailed Kundli report — comprehensive birth chart report with D1/D9/D10 charts, planetary strengths, major yogas, doshas, Vimshottari dasha timeline, and predictions by life area. Less exhaustive than Brihad but complete for most purposes. Use for 'detailed kundli report', 'full birth chart report', 'janm kundli report'.",
		Title:       "Detailed Kundli Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/kundli/detailed", in.toQuery())
	})
}

func registerReportKundliLite(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_kundli_lite",
		Description: "Generate a Lite Kundli report — a concise birth chart summary with lagna, Moon sign, nakshatra, key planets, current dasha, and a brief personality reading. Ideal for a quick reference or introductory reading. Use for 'quick kundli', 'lite birth chart', 'basic horoscope report'.",
		Title:       "Lite Kundli Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/kundli/lite", in.toQuery())
	})
}

func registerReportMangalDosha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_mangal_dosha",
		Description: "Generate a comprehensive Mangal Dosha report — full analysis of Mangal Dosha presence, sub-rules, all cancellation conditions checked, severity, impact on marriage, and detailed remedies with specific mantras, pujas, and gemstone guidance. Use for 'Mangal Dosha report', 'Manglik report', 'complete Manglik analysis'.",
		Title:       "Mangal Dosha Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/mangal-dosha", in.toQuery())
	})
}

type ReportMatchMakingInput struct {
	AshtakootaInput
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerReportMatchMaking(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_match_making",
		Description: "Generate a comprehensive Match Making (Kundli Milan) report for a couple — full Ashtakoota analysis with all 8 kootas scored, Mangal Dosha compatibility, Dasha synchrony, Navamsa compatibility, and an overall written recommendation. Use for 'kundli matching report', 'full compatibility report', 'marriage compatibility report'.",
		Title:       "Match Making Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportMatchMakingInput) (*mcp.CallToolResult, any, error) {
		q := boyGirlQuery(in.AshtakootaInput)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/reports/match-making", q)
	})
}

type ReportNumerologyInput struct {
	Name string `json:"name" jsonschema:"full name as written on official documents"`
	Dob  string `json:"dob" jsonschema:"date of birth in YYYY-MM-DD format"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerReportNumerology(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_numerology",
		Description: "Generate a comprehensive Numerology report — calculates all key numbers (Driver/Mulank, Conductor/Bhagyank, Name number, Soul Urge, Personality, Destiny) from name and date of birth, with complete interpretation of each number's meaning and their combined influence. Use for 'numerology report', 'complete numerology analysis', 'full number reading'.",
		Title:       "Numerology Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportNumerologyInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("name", in.Name)
		q.Set("dob", in.Dob)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/reports/numerology", q)
	})
}

func registerReportSadeSati(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_sade_sati",
		Description: "Generate a comprehensive Sade Sati report — covers current Sade Sati status, all past and upcoming cycles from birth to age 90, phase-by-phase narrative for the current or next cycle, and detailed remedies (mantras, donations, Shani puja guidance). Use for 'Sade Sati report', 'complete Saturn transit report', 'full Shani report'.",
		Title:       "Sade Sati Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/sade-sati", in.toQuery())
	})
}

func registerReportVarshaphal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "report_varshaphal",
		Description: "Generate a comprehensive Varshaphal (Solar Return) annual report — full annual chart with varsha lord, muntha, annual lagna, Tajika yogas, Mudda dasha timeline, and predictions for all 12 life areas for the current solar-return year. Use for 'varshaphal report', 'annual horoscope report', 'solar return year report'.",
		Title:       "Varshaphal Annual Report",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ReportBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/reports/varshaphal", in.toQuery())
	})
}

// =====================================================================
// Vedic Remedies (1 tool)
// =====================================================================

type VedicRemediesInput struct {
	BirthInput
	Category string `json:"category,omitempty" jsonschema:"optional remedy category: mantra, gemstone, pooja, vrat, yantra, color, number; if omitted returns all categories"`
	Lang     string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerVedicRemedies(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "vedic_remedies",
		Description: "Get Vedic astrological remedies for a birth chart — personalised remedies based on weak planets, doshas, and challenging dashas. Returns mantras (with count/timing), gemstone recommendations, vratas, poojas, yantras, lucky colors, and lucky numbers. Use for 'what are my remedies', 'remedy for weak Saturn', 'gemstone recommendation', 'puja for my chart'.",
		Title:       "Vedic Remedies",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VedicRemediesInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.Category != "" {
			q.Set("category", in.Category)
		}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies", q)
	})
}

// =====================================================================
// Western — missing endpoints (15 tools)
// =====================================================================

// WesternBirthInput and its toQuery() method are declared in tools_v03.go.

// TwoPersonInput (personA.*/personB.* params) is declared in tools_v03.go.
// The composite/davison tools below reuse it directly.
// WesternSynastryInput was removed — it sent person1_*/person2_* params
// which the API never reads (handler uses personA.*/personB.* prefixes).

type WesternAstrocartographyInput struct {
	WesternBirthInput
	ObserverLat float64 `json:"observer_lat" jsonschema:"observer latitude for local space calculations"`
	ObserverLon float64 `json:"observer_lon" jsonschema:"observer longitude for local space calculations"`
}

func registerWesternAstrocartographyLocalSpace(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_astrocartography_local_space",
		Description: "Compute Local Space astrocartography — the directions from an observer's location in which each natal planet rises, sets, culminates, or anti-culminates. Used to find beneficial directions for travel, relocation, or orienting a home. Use for 'local space directions', 'which direction is Jupiter from my home', 'astrocartography local space'.",
		Title:       "Astrocartography Local Space",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternAstrocartographyInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("observer_lat", strconv.FormatFloat(in.ObserverLat, 'f', -1, 64))
		q.Set("observer_lon", strconv.FormatFloat(in.ObserverLon, 'f', -1, 64))
		return callPassthrough(ctx, c, "/v1/western/astrocartography/local-space", q)
	})
}

func registerWesternCompositeHouses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_composite_houses",
		Description: "Compute the house cusps of the Composite chart for two people — the midpoint chart combining both natal charts to represent the relationship itself. Returns the composite Ascendant, Midheaven, and all 12 house cusps. Use for 'composite chart houses', 'relationship chart ascendant', 'composite house cusps'.",
		Title:       "Western Composite Houses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/composite/houses", in.toQuery())
	})
}

func registerWesternCompositePlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_composite_planets",
		Description: "Compute the planet placements in the Composite chart for two people — midpoint positions for all planets in the relationship chart. Returns sign, degree, house, and aspect pattern in the composite. Use for 'composite chart planets', 'relationship chart placements', 'composite Venus sign'.",
		Title:       "Western Composite Planets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/composite/planets", in.toQuery())
	})
}

func registerWesternDavisonHouses(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_davison_houses",
		Description: "Compute the house cusps of the Davison Relationship chart — the chart for the midpoint in time and space between two people's births (an actual chart for a real moment, unlike the composite). Returns Davison Ascendant and house cusps. Use for 'Davison chart houses', 'Davison ascendant', 'relationship chart Davison method'.",
		Title:       "Western Davison Houses",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/davison/houses", in.toQuery())
	})
}

func registerWesternDavisonPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_davison_planets",
		Description: "Compute the planet placements in the Davison Relationship chart — all planet positions in the Davison chart's tropical zodiac. Use for 'Davison chart planets', 'relationship chart planet placements Davison'.",
		Title:       "Western Davison Planets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TwoPersonInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/davison/planets", in.toQuery())
	})
}

func registerWesternDignitiesReceptions(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_dignities_receptions",
		Description: "Compute essential dignities and mutual receptions for a Western birth chart — each planet's essential dignity state (domicile, exaltation, triplicity, term, face, detriment, fall), any mutual reception pairs, and the overall dignity score. Use for 'planet dignities', 'mutual receptions', 'essential dignities western chart', 'which planets are strong in my chart'.",
		Title:       "Western Essential Dignities & Receptions",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/dignities/receptions", in.toQuery())
	})
}

// WesternNarrativeLordInput is used by the static lord-narrative endpoints
// that take only a planet name (lord) and optional language — no birth chart.
type WesternNarrativeLordInput struct {
	Lord string `json:"lord" jsonschema:"Firdaria time-lord planet name: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default)"`
}

func registerWesternNarrativeFirdaria(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_narrative_firdaria",
		Description: "Get the Firdaria time-lord narrative for a given planet — the Hellenistic planetary period system (similar to Vimshottari but Ptolemaic). Pass the current Firdaria lord (e.g. 'Saturn') to get a narrative of what that period brings. Use for 'firdaria period meaning', 'Hellenistic time lord interpretation', 'Saturn firdaria narrative'.",
		Title:       "Western Firdaria Time Lord Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternNarrativeLordInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lord", in.Lord)
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/western/narrative/timelord/firdaria", q)
	})
}

func registerWesternNatalAscendant(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_ascendant",
		Description: "Compute the Western (tropical) natal Ascendant — rising sign, degree, and Midheaven in the tropical zodiac. Distinct from the Vedic sidereal ascendant (typically ~23° earlier). Use for 'western ascendant', 'tropical rising sign', 'western natal rising', 'what is my western lagna'.",
		Title:       "Western Natal Ascendant (Tropical)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/ascendant", in.toQuery())
	})
}

func registerWesternNatalMC(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_mc",
		Description: "Compute the Midheaven (Medium Coeli, MC) and its tropical sign for a Western birth chart — the career and public reputation point, along with its aspects to natal planets. Use for 'what is my midheaven', 'western MC sign', 'tropical Midheaven', 'career point western astrology'.",
		Title:       "Western Natal Midheaven (MC)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/mc", in.toQuery())
	})
}

func registerWesternNatalOutOfBounds(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_out_of_bounds",
		Description: "Check which planets are Out of Bounds (OOB) in a Western natal chart — planets with declination beyond ±23°26′ (the Sun's maximum), considered to operate outside normal rules and with heightened, unconventional energy. Use for 'out of bounds planets', 'OOB planets', 'which planets are outside solar bounds'.",
		Title:       "Western Natal Out-of-Bounds Planets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/out-of-bounds", in.toQuery())
	})
}

func registerWesternNatalPlanets(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_planets",
		Description: "Get all planet placements in the Western (tropical) birth chart — sign, degree, house, retrograde flag, and speed for all 10 planets (Sun through Pluto) plus Chiron. Distinct from Vedic chart_planets which uses sidereal zodiac. Use for 'western birth chart planets', 'tropical planet placements', 'western natal chart'.",
		Title:       "Western Natal Planets (Tropical)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternBirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/western/natal/planets", in.toQuery())
	})
}

type WesternSolarReturnJDInput struct {
	WesternBirthInput
	Year int `json:"year" jsonschema:"year for the solar return (the year the return occurs, e.g. 2026)"`
}

func registerWesternSolarReturnExactJD(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_solar_return_exact_jd",
		Description: "Get the exact Julian Day (JD) and calendar datetime of the Western (tropical) solar return for a given year — the precise moment the Sun returns to its natal tropical longitude. Foundation for accurate solar-return chart casting. Use for 'exact tropical solar return time', 'western solar return JD', 'true birthday moment western'.",
		Title:       "Western Solar Return Exact JD",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternSolarReturnJDInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/western/returns/solar/exact-jd", q)
	})
}

// western_transit_summary is registered in tools_v04.go (registerWesternTransitSummary).
// That version correctly sends birth.lat/birth.lon/birth.date/birth.time/birth.tz
// (resolveWesternNatalContext reads "birth." prefix) plus transit date/time/tz and
// verbosity. A duplicate registration here was removed — it sent plain lat/lon which
// the handler never reads, causing 400 missing_param: birth.lat on every call.
//
// western_horoscope_monthly and western_horoscope_weekly are already
// registered in registerV04Tools (tools_v04.go) — no duplicate needed.
