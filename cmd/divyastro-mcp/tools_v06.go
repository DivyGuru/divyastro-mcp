package main

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV06Tools wires all Vedic compute endpoints not yet covered by
// v0.1–v0.5. Grouped by domain:
//
//   Ashtakavarga (3):  sarva, kaksha, transit-score
//   Calendar (4):      adhik-maas, month, ritu, samvatsara
//   Chart (6):         bhavabala, combustion, graha-yuddha, house-occupants,
//                      kp-house-significator, kp-sublord
//   Dasha (1):         yogini full list
//   Festivals (1):     on-date
//   Geo (2):           search, timezone
//   Milan (8):         dasha-sync, dhan-yog, longevity, mahendra,
//                      navamsa-compat, santan-yog, shani-dosha, stree-dirgha
//   Muhurta (8):       chandra-bala, graha-pravesh, panchaka-rahita,
//                      sarvartha-siddhi, shubha-yoga, tara-bala, vyapar, yatra
//   Numerology (5):    conductor, destiny, driver, personality, soul
//   Panchang (2):      hora-dinman, sankranti
//   Planet-moments (2):combustion-window, speed
//   Prashna (4):       arudha, chart, lagna-lord, significators
//   Transit (6):       ashtakavarga, double-transit, small-panoti,
//                      small-panoti-history, tarabala, vedha
//   Varshaphal (4):    lord, muntha, solar-return-jd, yoga
//
// Total new: 56 tools.

func registerV06Tools(s *mcp.Server, c *apiClient) {
	// Ashtakavarga
	registerAshtakavargaSarva(s, c)
	registerAshtakavargaKaksha(s, c)
	registerAshtakavargaTransitScore(s, c)

	// Calendar
	registerCalendarAdhikMaas(s, c)
	registerCalendarMonth(s, c)
	registerCalendarRitu(s, c)
	registerCalendarSamvatsara(s, c)

	// Chart
	registerChartBhavabala(s, c)
	registerChartCombustion(s, c)
	registerChartGrahaYuddha(s, c)
	registerChartHouseOccupants(s, c)
	registerChartKPHouseSignificator(s, c)
	registerChartKPSublord(s, c)

	// Dasha
	registerDashaYoginiFull(s, c)

	// Festivals
	registerFestivalsOnDate(s, c)

	// Geo
	registerGeoSearch(s, c)
	registerGeoTimezone(s, c)

	// Milan
	registerMilanDashaSync(s, c)
	registerMilanDhanYog(s, c)
	registerMilanLongevity(s, c)
	registerMilanMahendra(s, c)
	registerMilanNavamsaCompat(s, c)
	registerMilanSantanYog(s, c)
	registerMilanShaniDosha(s, c)
	registerMilanStreeDirgha(s, c)

	// Muhurta
	registerMuhurtaChandraBala(s, c)
	registerMuhurtaGrahaPravesh(s, c)
	registerMuhurtaPanchakaRahita(s, c)
	registerMuhurtaSarvarthaSiddhi(s, c)
	registerMuhurtaShubhaYoga(s, c)
	registerMuhurtaTaraBala(s, c)
	registerMuhurtaVyapar(s, c)
	registerMuhurtaYatra(s, c)

	// Numerology
	registerNumerologyConductor(s, c)
	registerNumerologyDestiny(s, c)
	registerNumerologyDriver(s, c)
	registerNumerologyPersonality(s, c)
	registerNumerologySoul(s, c)

	// Panchang
	registerPanchangHoraDinman(s, c)
	registerPanchangSankranti(s, c)

	// Planet-moments
	registerPlanetCombustionWindow(s, c)
	registerPlanetSpeed(s, c)

	// Prashna
	registerPrashnaArudha(s, c)
	registerPrashnaChart(s, c)
	registerPrashnaLagnaLord(s, c)
	registerPrashnaSignificators(s, c)

	// Transit
	registerTransitAshtakavarga(s, c)
	registerTransitDoubleTransit(s, c)
	registerTransitSmallPanoti(s, c)
	registerTransitSmallPanotiHistory(s, c)
	registerTransitTarabala(s, c)
	registerTransitVedha(s, c)

	// Varshaphal
	registerVarshaphalLord(s, c)
	registerVarshaphalMuntha(s, c)
	registerVarshaphalSolarReturnJD(s, c)
	registerVarshaphalYoga(s, c)
}

// =====================================================================
// Ashtakavarga (3 tools)
// =====================================================================

func registerAshtakavargaSarva(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ashtakavarga_sarva",
		Description: "Get the Sarvashtakavarga table — the combined 12-sign benefic-point grid from all 9 grahas. Each sign shows total benefic points (max 56). High-scoring signs indicate strong houses for transit timing. Use for 'how many ashtakavarga points does my 10th house have', 'sarva ashtakavarga', 'which house is strongest by SAV'.",
		Title:       "Sarvashtakavarga Table",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/ashtakavarga/sarva", in.toQuery())
	})
}

func registerAshtakavargaKaksha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ashtakavarga_kaksha",
		Description: "Get Kaksha (sub-division) level Ashtakavarga analysis — divides each sign into 8 equal 3°45′ segments, each ruled by a graha. Used for precise transit timing within a sign. Use for 'kaksha ashtakavarga', 'which kaksha is Saturn transiting', 'precise transit timing'.",
		Title:       "Ashtakavarga Kaksha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/ashtakavarga/kaksha", in.toQuery())
	})
}

type AshtakavargaTransitScoreInput struct {
	BirthInput
	TransitDate string `json:"transit_date,omitempty" jsonschema:"optional date to evaluate transit positions in YYYY-MM-DD; defaults to today"`
}

func registerAshtakavargaTransitScore(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ashtakavarga_transit_score",
		Description: "Score current (or a given date's) planetary transits against a natal chart using Ashtakavarga benefic points. Returns a score per transiting planet indicating how favourable the transit is for that natal chart. Use for 'how good is Saturn's transit for me now', 'ashtakavarga transit score today'.",
		Title:       "Ashtakavarga Transit Score",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakavargaTransitScoreInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.TransitDate != "" {
			q.Set("transit_date", in.TransitDate)
		}
		return callPassthrough(ctx, c, "/v1/ashtakavarga/transit-score", q)
	})
}

// =====================================================================
// Calendar (4 tools)
// =====================================================================

type CalendarYearInput struct {
	Year int    `json:"year" jsonschema:"Gregorian year, e.g. 2026"`
	Tz   string `json:"tz,omitempty" jsonschema:"IANA timezone for date boundaries (e.g. Asia/Kolkata); defaults to UTC"`
}

func registerCalendarAdhikMaas(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "calendar_adhik_maas",
		Description: "Find the Adhik Maas (intercalary/leap month) in a given year — the extra lunar month inserted to sync the Hindu lunar calendar with the solar year. Returns the month name, start date, and end date. Adhik Maas occurs roughly every 2.5–3 years. Use for 'is there adhik maas in 2026', 'leap month this year', 'mal maas dates'.",
		Title:       "Adhik Maas (Leap Month)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CalendarYearInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/calendar/adhik-maas", q)
	})
}

type CalendarMonthInput struct {
	Year  int    `json:"year" jsonschema:"Gregorian year, e.g. 2026"`
	Month int    `json:"month" jsonschema:"Gregorian month 1–12"`
	Tz    string `json:"tz,omitempty" jsonschema:"IANA timezone (e.g. Asia/Kolkata); defaults to UTC"`
}

func registerCalendarMonth(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "calendar_month",
		Description: "Get the Hindu calendar details for a Gregorian month — maps each day to its Hindu month (Chaitra/Vaishakha/…), paksha (Shukla/Krishna), tithi, and Vikram Samvat year. Use for 'what is the Hindu date today', 'which Hindu month is it', 'tithi calendar for July 2026'.",
		Title:       "Hindu Calendar Month",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CalendarMonthInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		q.Set("month", strconv.Itoa(in.Month))
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/calendar/month", q)
	})
}

func registerCalendarRitu(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "calendar_ritu",
		Description: "Get the current Hindu Ritu (season) for a date — one of 6 seasons: Vasanta (spring), Grishma (summer), Varsha (monsoon), Sharad (autumn), Hemanta (pre-winter), Shishira (winter). Each Ritu spans 2 solar months. Use for 'what ritu is it now', 'which season in Hindu calendar', 'ritu for June'.",
		Title:       "Hindu Ritu (Season)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CalendarYearInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/calendar/ritu", q)
	})
}

func registerCalendarSamvatsara(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "calendar_samvatsara",
		Description: "Get the Samvatsara (Jupiter year) name for a given year — the 60-year Jupiter cycle used in Hindu astrology (Prabhava, Vibhava, Shukla, … Akshaya). Returns Vikram Samvat, Shaka Samvat, and Samvatsara name with its ruling deity. Use for 'what samvatsara is 2026', 'Hindu year name', 'Jupiter year'.",
		Title:       "Samvatsara (Hindu Year Name)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CalendarYearInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/calendar/samvatsara", q)
	})
}

// =====================================================================
// Chart — missing compute tools (6)
// =====================================================================

func registerChartBhavabala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_bhavabala",
		Description: "Compute Bhava Bala — the six-component strength of each of the 12 houses: Bhavadhipati Bala (lord strength), Bhava Digbala, Bhava Drishti Bala, and more. Identifies which houses are strong/weak in the chart. Use for 'which is my strongest house', 'bhava bala', '10th house strength'.",
		Title:       "Bhava Bala (House Strength)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/bhavabala", in.toQuery())
	})
}

func registerChartCombustion(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_combustion",
		Description: "Check which planets are combust (too close to the Sun) in a birth chart — returns each planet's angular distance from Sun, whether it is combust by classical limits, and combust severity. Use for 'is my Mercury combust', 'combust planets in my chart', 'asta graha'.",
		Title:       "Combust Planets",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/combustion", in.toQuery())
	})
}

func registerChartGrahaYuddha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_graha_yuddha",
		Description: "Identify Graha Yuddha (planetary war) in a birth chart — occurs when two visible planets are within 1° of each other. Returns the warring pair, winner (the planet with higher ecliptic latitude), loser, and degrees of separation. Use for 'graha yuddha in my chart', 'planetary war', 'which planets are in war'.",
		Title:       "Graha Yuddha (Planetary War)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/graha-yuddha", in.toQuery())
	})
}

func registerChartHouseOccupants(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_house_occupants",
		Description: "List all planets occupying each of the 12 houses in a birth chart — a house-indexed view of the chart (vs planet-indexed from chart_planets). Use for 'what planets are in my 5th house', 'house occupants', 'planets in each bhava'.",
		Title:       "House Occupants",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/house-occupants", in.toQuery())
	})
}

type KPInput struct {
	BirthInput
	Cusp int `json:"cusp" jsonschema:"house cusp number 1–12 for KP analysis"`
}

func registerChartKPHouseSignificator(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_kp_house_significator",
		Description: "Get KP (Krishnamurti Paddhati) house significators for a given cusp — lists planets that significate the house through occupancy, lordship, and sub-lord chain. Used in KP astrology for event timing and prediction. Use for 'KP significators for 7th house', 'KP house analysis', 'krishnamurti paddhati'.",
		Title:       "KP House Significator",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in KPInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("cusp", strconv.Itoa(in.Cusp))
		return callPassthrough(ctx, c, "/v1/chart/kp-house-significator", q)
	})
}

type KPSublordInput struct {
	BirthInput
	Planet string `json:"planet" jsonschema:"planet name for KP sublord analysis: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu"`
}

func registerChartKPSublord(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_kp_sublord",
		Description: "Get the KP sublord of a planet — the sub-lord of the nakshatra sub-division the planet occupies (KP star lord → sub lord → sub-sub lord chain). Central to KP timing. Use for 'KP sublord of Moon', 'star lord and sub lord', 'KP planet analysis'.",
		Title:       "KP Sublord",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in KPSublordInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("planet", in.Planet)
		return callPassthrough(ctx, c, "/v1/chart/kp-sublord", q)
	})
}

// =====================================================================
// Dasha (1 tool)
// =====================================================================

func registerDashaYoginiFull(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_yogini_full",
		Description: "Get the full Yogini Dasha sequence for a birth chart — all 8 Yogini lords (Mangala, Pingala, Dhanya, Bhramari, Bhadrika, Ulka, Siddha, Sankata) with start/end dates across 36 years. Use for 'full yogini dasha timeline', 'all yogini periods', 'yogini dasha from birth'.",
		Title:       "Yogini Dasha Full Timeline",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/dasha/yogini", in.toQuery())
	})
}

// =====================================================================
// Festivals (1 tool)
// =====================================================================

type FestivalsOnDateInput struct {
	Date   string `json:"date" jsonschema:"date to query in YYYY-MM-DD format"`
	Tz     string `json:"tz,omitempty" jsonschema:"IANA timezone for date boundaries (e.g. Asia/Kolkata); defaults to UTC"`
	Region string `json:"region,omitempty" jsonschema:"optional region: 'north_india', 'south_india', 'maharashtra', 'gujarat', 'bengal'"`
}

func registerFestivalsOnDate(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "festivals_on_date",
		Description: "List Hindu festivals and observances falling on a specific date — Ekadashi, Pradosh, Purnima, Amavasya, major festivals, vratas, and regional observances. Use for 'what festival is today', 'which vrat is on 15 July 2026', 'is there any festival tomorrow'.",
		Title:       "Festivals on a Specific Date",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in FestivalsOnDateInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("date", in.Date)
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		if in.Region != "" {
			q.Set("region", in.Region)
		}
		return callPassthrough(ctx, c, "/v1/festivals/on-date", q)
	})
}

// =====================================================================
// Geo (2 tools)
// =====================================================================

type GeoSearchInput struct {
	Q      string `json:"q" jsonschema:"city or place name to search (e.g. 'Mumbai', 'New Delhi', 'Varanasi')"`
	Limit  int    `json:"limit,omitempty" jsonschema:"max results to return (1–20, default 10)"`
	Cc     string `json:"cc,omitempty" jsonschema:"optional ISO-3166-1 alpha-2 country code to narrow results (e.g. 'IN', 'US')"`
}

func registerGeoSearch(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "geo_search",
		Description: "Search for a city or place by name — returns matching locations with latitude, longitude, state/province, country, and IANA timezone. Use when the user mentions a city name and you need coordinates for chart or panchang computations. Use for 'where is Nashik', 'coordinates of Jaipur', 'find city lat lon'.",
		Title:       "Geo Search (City → Coordinates)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in GeoSearchInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("q", in.Q)
		if in.Limit > 0 {
			q.Set("limit", strconv.Itoa(in.Limit))
		}
		if in.Cc != "" {
			q.Set("cc", in.Cc)
		}
		return callPassthrough(ctx, c, "/v1/geo/search", q)
	})
}

type GeoTimezoneInput struct {
	Lat float64 `json:"lat" jsonschema:"latitude in decimal degrees"`
	Lon float64 `json:"lon" jsonschema:"longitude in decimal degrees"`
}

func registerGeoTimezone(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "geo_timezone",
		Description: "Get the IANA timezone for a latitude/longitude — returns timezone name (e.g. Asia/Kolkata), UTC offset, and current DST status. Use before calling any chart/panchang endpoint when you only have coordinates without an explicit timezone. Use for 'what timezone is at these coordinates', 'timezone for lat/lon'.",
		Title:       "Geo Timezone Lookup",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in GeoTimezoneInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		return callPassthrough(ctx, c, "/v1/geo/timezone", q)
	})
}

// =====================================================================
// Milan — remaining compatibility tools (8)
// =====================================================================

func registerMilanDashaSync(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_dasha_sync",
		Description: "Check Dasha synchrony between two partners — are they running compatible mahadasha lords at the same time? Good dasha synchrony indicates the couple will flourish together during major life events. Use for 'dasha compatibility', 'are our dashas compatible', 'dasha sync for marriage'.",
		Title:       "Milan Dasha Sync",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/dasha-sync", boyGirlQuery(in))
	})
}

func registerMilanDhanYog(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_dhan_yog",
		Description: "Check Dhan Yog (wealth yoga) for a single person's chart — assesses 2nd/11th house lord strength and Jupiter placement for prosperity. Use for 'dhan yog in my chart', 'wealth yoga', 'financial yoga analysis'.",
		Title:       "Milan Dhan Yog",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/dhan-yog", birthPrefixQuery(in))
	})
}

func registerMilanLongevity(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_longevity",
		Description: "Assess longevity for a single person's chart — classical analysis of lagna, 8th house, and Moon sign longevity indicators. Use for 'longevity in my chart', 'ayush analysis', 'lifespan indicators'.",
		Title:       "Milan Longevity",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/longevity", birthPrefixQuery(in))
	})
}

func registerMilanMahendra(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_mahendra",
		Description: "Check Mahendra Koota — the 4th additional compatibility factor beyond Ashtakoota. Counts from bride's nakshatra to groom's in multiples of 4, 13, 22 etc. A positive result indicates the groom will protect and provide for the bride. Use for 'mahendra koota', 'mahendra compatibility'.",
		Title:       "Milan Mahendra Koota",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/mahendra", boyGirlQuery(in))
	})
}

func registerMilanNavamsaCompat(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_navamsa_compat",
		Description: "Check Navamsa compatibility — compares the D9 (Navamsa) chart placements of both partners for deeper marriage harmony beyond the D1 chart. Assesses 7th house navamsa lord synergy and mutual aspect in the navamsa. Use for 'navamsa compatibility', 'D9 chart matching', 'deeper marriage compatibility'.",
		Title:       "Milan Navamsa Compatibility",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/navamsa-compat", boyGirlQuery(in))
	})
}

func registerMilanSantanYog(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_santan_yog",
		Description: "Check Santan Yog for a single person's chart — assesses 5th house lord and Jupiter placement to indicate prospects for children. Use for 'santan yog in my chart', 'children prospects', 'progeny yoga'.",
		Title:       "Milan Santan Yog (Progeny)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/santan-yog", birthPrefixQuery(in))
	})
}

func registerMilanShaniDosha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_shani_dosha",
		Description: "Check Shani Dosha for a single person's chart — Saturn in 1st, 4th, 7th, 10th, or 12th house from lagna. Returns dosha presence, which house, and cancellation rules. Use for 'do I have shani dosha', 'Saturn affliction in my chart', 'is my shani dosha cancelled'.",
		Title:       "Milan Shani Dosha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/shani-dosha", birthPrefixQuery(in))
	})
}

func registerMilanStreeDirgha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_stree_dirgha",
		Description: "Check Stree Dirgha — the distance from bride's nakshatra to groom's must be more than 9 nakshatras for a beneficial match (indicates the groom will lead a longer, more prosperous life than the bride expected). Use for 'stree dirgha', 'nakshatra distance check', 'bride groom nakshatra compatibility'.",
		Title:       "Milan Stree Dirgha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/stree-dirgha", boyGirlQuery(in))
	})
}

// =====================================================================
// Muhurta — remaining 8 tools
// =====================================================================

type ChandraBalaInput struct {
	Lat           float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon           float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	Date          string  `json:"date" jsonschema:"date to evaluate in YYYY-MM-DD"`
	Time          string  `json:"time" jsonschema:"time in HH:MM 24-hour format"`
	Tz            string  `json:"tz" jsonschema:"IANA timezone (e.g. Asia/Kolkata)"`
	BirthMoonSign int     `json:"birth_moon_sign" jsonschema:"natal Moon sign index 0-11 (0=Aries, 1=Taurus … 11=Pisces)"`
}

func registerMuhurtaChandraBala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_chandra_bala",
		Description: "Compute Chandra Bala — the Moon's strength relative to a person's natal Moon sign on a given date. Score 1–6 based on the Moon's position (1st=Janma, 2nd=Sampat, etc.). Used for personal muhurta selection. Use for 'is today good for my Moon sign', 'chandra bala today', 'Moon strength for my rashi'.",
		Title:       "Muhurta Chandra Bala",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ChandraBalaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("date", in.Date)
		q.Set("time", in.Time)
		q.Set("tz", in.Tz)
		q.Set("birth_moon_sign", strconv.Itoa(in.BirthMoonSign))
		return callPassthrough(ctx, c, "/v1/muhurta/chandra-bala", q)
	})
}

func registerMuhurtaGrahaPravesh(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_graha_pravesh",
		Description: "Find auspicious Griha Pravesh (housewarming) muhurta windows in a date range — checks tithi, nakshatra, vara, Sun position, and lagna for the ceremony. Use for 'housewarming date', 'griha pravesh muhurta', 'auspicious date to move into new home'.",
		Title:       "Griha Pravesh Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/graha-pravesh", in.toQuery())
	})
}

func registerMuhurtaPanchakaRahita(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_panchaka_rahita",
		Description: "Find Panchaka-free (Panchaka Rahita) windows in a date range — Panchaka occurs when Moon is in Aquarius or Pisces and creates obstacles for travel, construction, funerals, and certain ceremonies. Returns safe windows clear of Panchaka. Use for 'panchaka free time', 'when is panchaka', 'avoid panchaka days'.",
		Title:       "Panchaka Rahita (Panchaka-Free Windows)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/panchaka-rahita", in.toQuery())
	})
}

func registerMuhurtaSarvarthaSiddhi(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_sarvartha_siddhi",
		Description: "Find Sarvartha Siddhi Yoga windows — one of the most auspicious panchang yogas formed by specific vara+nakshatra combinations. Said to confer success in all endeavours. Use for 'sarvartha siddhi dates', 'best yoga for starting a business', 'auspicious yoga today'.",
		Title:       "Sarvartha Siddhi Yoga",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/sarvartha-siddhi", in.toQuery())
	})
}

func registerMuhurtaShubhaYoga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_shubha_yoga",
		Description: "Find Shubha Yoga (auspicious panchang yoga) windows in a date range — checks for Amrit Siddhi, Sarvartha Siddhi, Ravi Yoga, and other classical panchang yogas in combination. Returns windows when multiple auspicious yogas align. Use for 'shubha yoga dates', 'multiple auspicious yogas', 'best panchang day'.",
		Title:       "Shubha Yoga Windows",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/shubha-yoga", in.toQuery())
	})
}

type TaraBalaInput struct {
	Lat     float64 `json:"lat" jsonschema:"observer latitude in decimal degrees"`
	Lon     float64 `json:"lon" jsonschema:"observer longitude in decimal degrees"`
	Date    string  `json:"date" jsonschema:"date to evaluate in YYYY-MM-DD"`
	Time    string  `json:"time" jsonschema:"time in HH:MM 24-hour format"`
	Tz      string  `json:"tz" jsonschema:"IANA timezone (e.g. Asia/Kolkata)"`
	BirthNak int    `json:"birth_nak" jsonschema:"natal nakshatra index 0-26 (0=Ashwini, 1=Bharani … 26=Revati)"`
}

func registerMuhurtaTaraBala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_tara_bala",
		Description: "Compute Tara Bala score for a given date relative to a person's birth nakshatra — rates the day's Moon nakshatra as Janma/Sampat/Vipat/Kshema/Pratyak/Sadhana/Vadha/Mitra/Param-Mitra. Used for personalized muhurta selection. Use for 'is today good for my nakshatra', 'tara bala score', 'which day is good for my birth nakshatra'.",
		Title:       "Muhurta Tara Bala",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in TaraBalaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("date", in.Date)
		q.Set("time", in.Time)
		q.Set("tz", in.Tz)
		q.Set("birth_nak", strconv.Itoa(in.BirthNak))
		return callPassthrough(ctx, c, "/v1/muhurta/tara-bala", q)
	})
}

func registerMuhurtaVyapar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_vyapar",
		Description: "Find auspicious Vyapar (business/trade) muhurta windows — suitable for starting a business, signing contracts, launching a product, or opening a shop. Checks Mercury strength, Jupiter aspect, Labh choghadiya, and avoidance of inauspicious tithis. Use for 'business launch date', 'shop opening muhurta', 'vyapar muhurta', 'when to start business'.",
		Title:       "Vyapar (Business) Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/vyapar", in.toQuery())
	})
}

func registerMuhurtaYatra(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "muhurta_yatra",
		Description: "Find auspicious Yatra (travel/journey) muhurta windows — avoids Disha Shool (directional inauspiciousness), Rahu Kaal, Nak Shool, and favours Pushya/Hasta/Ashwini/Anuradha nakshatras for journeys. Use for 'good time to travel', 'yatra muhurta', 'auspicious day for journey', 'when to start travel'.",
		Title:       "Yatra (Travel) Muhurta",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MuhurtaWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/muhurta/yatra", in.toQuery())
	})
}

// =====================================================================
// Numerology — individual numbers (5 tools)
// =====================================================================

type NumerologyNameInput struct {
	Name string `json:"name" jsonschema:"full name as written on official documents"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

type NumerologyDobInput struct {
	Dob  string `json:"dob" jsonschema:"date of birth in YYYY-MM-DD format"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func numerologyNameQuery(in NumerologyNameInput) url.Values {
	q := url.Values{}
	q.Set("name", in.Name)
	if in.Lang != "" {
		q.Set("lang", in.Lang)
	}
	return q
}

func numerologyDobQuery(in NumerologyDobInput) url.Values {
	q := url.Values{}
	q.Set("dob", in.Dob)
	if in.Lang != "" {
		q.Set("lang", in.Lang)
	}
	return q
}

func registerNumerologyDriver(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_driver",
		Description: "Compute the Mulank / Driver number — the single-digit number derived from the day of birth (1–31, reduced to 1–9, preserving 11/22/33). The Driver represents personality, instinct, and primary nature. Use for 'what is my driver number', 'mulank', 'date of birth number'.",
		Title:       "Numerology Driver (Mulank)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyDobInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/driver", numerologyDobQuery(in))
	})
}

func registerNumerologyConductor(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_conductor",
		Description: "Compute the Bhagyank / Conductor number — derived from the full date of birth (DD+MM+YYYY reduced to a single digit). The Conductor governs destiny, life path, and fortune. Use for 'what is my bhagyank', 'life path number', 'conductor number', 'destiny from date of birth'.",
		Title:       "Numerology Conductor (Bhagyank)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyDobInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/conductor", numerologyDobQuery(in))
	})
}

func registerNumerologySoul(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_soul",
		Description: "Compute the Soul Urge number — derived from the vowels of the full name. Reveals inner motivation, deepest desires, and what truly drives a person at soul level. Use for 'soul urge number', 'heart's desire', 'inner motivation number', 'vowel number from name'.",
		Title:       "Numerology Soul Urge",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyNameInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/soul", numerologyNameQuery(in))
	})
}

func registerNumerologyPersonality(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_personality",
		Description: "Compute the Personality number — derived from the consonants of the full name. Represents the outer personality, how others perceive you, and the mask you present to the world. Use for 'personality number', 'consonant number', 'how others see me numerology'.",
		Title:       "Numerology Personality",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyNameInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/personality", numerologyNameQuery(in))
	})
}

func registerNumerologyDestiny(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "numerology_destiny",
		Description: "Compute the Destiny number — derived from all letters of the full name (vowels + consonants). Represents life purpose, talents, and the path you are meant to fulfil. Use for 'destiny number', 'expression number', 'full name number', 'life purpose number'.",
		Title:       "Numerology Destiny",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NumerologyNameInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/numerology/destiny", numerologyNameQuery(in))
	})
}

// =====================================================================
// Panchang — remaining 2 tools
// =====================================================================

func registerPanchangHoraDinman(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_hora_dinman",
		Description: "Get the Dinman (day length) and Ratriман (night length) for a date — the exact duration of daytime from sunrise to sunset and nighttime from sunset to next sunrise, in ghatis (1 ghati = 24 minutes). Used in muhurta calculations. Use for 'day length today', 'dinman ratriман', 'how long is the day', 'ghati duration'.",
		Title:       "Hora Dinman (Day/Night Duration)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/hora-dinman", in.toQuery())
	})
}

func registerPanchangSankranti(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_sankranti",
		Description: "Get solar Sankranti (Sun's sign ingress) dates for a year — the exact moment the Sun enters each of the 12 sidereal signs, including Makar Sankranti (Capricorn), Mesh Sankranti (Aries/New Year), and all others. Use for 'when is Makar Sankranti', 'solar ingress dates', 'sun sign change 2026', 'sankranti calendar'.",
		Title:       "Sankranti (Solar Ingress) Calendar",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in CalendarYearInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("year", strconv.Itoa(in.Year))
		if in.Tz != "" {
			q.Set("tz", in.Tz)
		}
		return callPassthrough(ctx, c, "/v1/panchang/sankranti", q)
	})
}

// =====================================================================
// Planet-moments — remaining 2 tools
// =====================================================================

type PlanetDateInput struct {
	Planet string `json:"planet" jsonschema:"planet name: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu"`
	Date   string `json:"date,omitempty" jsonschema:"optional date YYYY-MM-DD; defaults to today"`
}

func registerPlanetCombustionWindow(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "planet_combustion_window",
		Description: "Find the combustion window for a planet — the date range when the planet is too close to the Sun (within classical combust limits) and loses strength. Mercury combustion is very frequent; Venus, Mars, Jupiter, Saturn combust periodically. Use for 'when is Mercury combust', 'combustion dates for Venus', 'when does Mars come out of combustion'.",
		Title:       "Planet Combustion Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PlanetWindowInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/planet-moments/combustion-window", in.toQuery())
	})
}

func registerPlanetSpeed(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "planet_speed",
		Description: "Get the current speed and motion state of a planet — degrees per day, whether it is direct/retrograde/stationary, and its speed relative to mean motion. Use for 'how fast is Jupiter moving', 'is Saturn stationary', 'planet speed today', 'graha gati'.",
		Title:       "Planet Speed & Motion",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PlanetDateInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("planet", in.Planet)
		if in.Date != "" {
			q.Set("date", in.Date)
		}
		return callPassthrough(ctx, c, "/v1/planet-moments/speed", q)
	})
}

// =====================================================================
// Prashna — remaining 4 tools
// =====================================================================

func registerPrashnaArudha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "prashna_arudha",
		Description: "Compute the Arudha Lagna of a Prashna chart — the perceived image/manifestation of the question's subject. Used in KP and Jaimini Prashna for nuanced prediction beyond the basic yes/no. Use for 'prashna arudha lagna', 'horary chart arudha', 'prashna image lagna'.",
		Title:       "Prashna Arudha Lagna",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PrashnaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		q.Set("question", in.Question)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/prashna/arudha", q)
	})
}

func registerPrashnaChart(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "prashna_chart",
		Description: "Cast the full Prashna (horary) chart for the moment a question is asked — returns lagna, all planet placements, houses, and panchang at the moment of the question. Use when you need the complete chart data for detailed horary analysis, not just the yes/no answer. Use for 'prashna kundli', 'horary chart', 'prashna chart details'.",
		Title:       "Prashna Chart",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PrashnaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		q.Set("question", in.Question)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/prashna/chart", q)
	})
}

func registerPrashnaLagnaLord(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "prashna_lagna_lord",
		Description: "Get the lagna lord analysis for a Prashna chart — where the lagna lord is placed, its strength, and what it indicates about the querent's situation and chances of success. Use for 'prashna lagna lord', 'horary lagna lord analysis', 'prashna ascendant lord'.",
		Title:       "Prashna Lagna Lord",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PrashnaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		q.Set("question", in.Question)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/prashna/lagna-lord", q)
	})
}

func registerPrashnaSignificators(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "prashna_significators",
		Description: "Get the KP significators for a Prashna chart — lists planets significating the houses relevant to the question (e.g. houses 2, 6, 10, 11 for job questions). Used for KP horary timing and prediction. Use for 'prashna significators', 'KP horary significators', 'which planets significate my question'.",
		Title:       "Prashna KP Significators",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PrashnaInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		q.Set("lat", strconv.FormatFloat(in.Lat, 'f', -1, 64))
		q.Set("lon", strconv.FormatFloat(in.Lon, 'f', -1, 64))
		q.Set("tz", in.Tz)
		q.Set("question", in.Question)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/prashna/significators", q)
	})
}

// =====================================================================
// Transit — remaining 6 tools
// =====================================================================

func registerTransitAshtakavarga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_ashtakavarga",
		Description: "Score current planetary transits against a natal chart using Ashtakavarga — shows which houses transiting planets are moving through and how many benefic points those signs have in the natal chart. High-point transits are favourable; low-point are challenging. Use for 'ashtakavarga transit today', 'is this transit good for me', 'transit ashtakavarga score'.",
		Title:       "Transit Ashtakavarga",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/ashtakavarga", birthPrefixQuery(in))
	})
}

func registerTransitDoubleTransit(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_double_transit",
		Description: "Compute Double Transit (Jupiter-Saturn-Rahu joint transit) analysis — the classical technique where an event manifests only when Jupiter AND Saturn both transit houses connected to a natal significator simultaneously. Used for timing major life events (marriage, career change, relocation). Use for 'double transit analysis', 'Jupiter Saturn transit together', 'when will I get married by double transit'.",
		Title:       "Double Transit Analysis",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/double-transit", in.toQuery())
	})
}

type SmallPanotiInput struct {
	BirthInput
	At string `json:"at,omitempty" jsonschema:"optional date-time to evaluate in RFC3339; defaults to now"`
}

func registerTransitSmallPanoti(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_small_panoti",
		Description: "Check whether a person is currently in Small Panoti (Ashtama Shani) — the 2.5-year Saturn transit through the 8th house from natal Moon, considered challenging. Returns status, start/end dates, and intensity. Use for 'is Saturn in my 8th house', 'ashtama shani', 'small panoti status', 'kantaka shani'.",
		Title:       "Small Panoti (Ashtama Shani)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in SmallPanotiInput) (*mcp.CallToolResult, any, error) {
		q := birthPrefixQuery(in.BirthInput)
		if in.At != "" {
			q.Set("at", in.At)
		}
		return callPassthrough(ctx, c, "/v1/transit/small-panoti", q)
	})
}

func registerTransitSmallPanotiHistory(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_small_panoti_history",
		Description: "Get the complete history and future windows of Small Panoti (Ashtama Shani) for a natal chart — all past and upcoming 2.5-year periods when Saturn transits the 8th house from natal Moon. Use for 'all my ashtama shani periods', 'small panoti history', 'when was my last small panoti'.",
		Title:       "Small Panoti History & Future",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/small-panoti/history", birthPrefixQuery(in))
	})
}

func registerTransitTarabala(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_tarabala",
		Description: "Compute Tara Bala for today's transiting Moon relative to a natal chart — rates the day as Janma/Sampat/Vipat/Kshema/Pratyak/Sadhana/Vadha/Mitra/Param-Mitra tara. Used to personalize daily muhurta. Use for 'tara bala today for my chart', 'is today's Moon good for me', 'tarabala transit'.",
		Title:       "Transit Tara Bala",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/tarabala", birthPrefixQuery(in))
	})
}

func registerTransitVedha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "transit_vedha",
		Description: "Check Vedha (transit obstruction) for a natal chart — certain house transits are blocked (vedha'd) by another planet transiting a specific opposing house simultaneously. Returns which transits are obstructed by Vedha right now. Use for 'vedha today', 'transit obstruction', 'which planet transit is blocked by vedha'.",
		Title:       "Transit Vedha (Obstruction)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/transit/vedha", birthPrefixQuery(in))
	})
}

// =====================================================================
// Varshaphal — remaining 4 tools
// =====================================================================

func registerVarshaphalLord(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_lord",
		Description: "Get the Varsha Lord (year ruler) for an annual chart — the planet that rules the Varshaphal year based on the Tajika system. The Varsha Lord strongly influences the themes of the solar-return year. Use for 'varsha lord 2026', 'year ruler annual chart', 'who rules my varshaphal year'.",
		Title:       "Varshaphal Lord (Year Ruler)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/lord", q)
	})
}

func registerVarshaphalMuntha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_muntha",
		Description: "Get Muntha placement for a Varshaphal year — Muntha moves one sign per year from birth lagna and its house position in the annual chart is a key indicator of the year's themes. Muntha in 1st/5th/9th/10th is favourable; 6th/8th/12th is challenging. Use for 'muntha 2026', 'varshaphal muntha placement', 'muntha in which house this year'.",
		Title:       "Varshaphal Muntha",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/muntha", q)
	})
}

func registerVarshaphalSolarReturnJD(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_solar_return_jd",
		Description: "Get the precise Julian Day (JD) and calendar datetime of the solar return for a given year — the exact moment the Sun returns to its natal longitude. This is the true astrological birthday and the foundation of all Varshaphal calculations. Use for 'exact solar return time 2026', 'varshaphal solar return JD', 'true birthday moment'.",
		Title:       "Varshaphal Solar Return JD",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/solar-return-jd", q)
	})
}

func registerVarshaphalYoga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "varshaphal_yoga",
		Description: "Get Tajika Yogas in the Varshaphal annual chart — Ikkaval, Induvara, Ithasala, Ishrafa, Nakta, Yamaya, Manau, Khallasara, Duphali-Kuttha, and other Tajika aspect-based yogas that indicate the tone of the solar-return year. Use for 'varshaphal yogas 2026', 'tajika yogas annual chart', 'annual chart yoga analysis'.",
		Title:       "Varshaphal Yogas (Tajika)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in VarshaphalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		q.Set("year", strconv.Itoa(in.Year))
		return callPassthrough(ctx, c, "/v1/varshaphal/yoga", q)
	})
}
