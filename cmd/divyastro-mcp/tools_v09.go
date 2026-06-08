package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV09Tools covers panchang sub-endpoints and chart/ghatak that are
// registered via chi r.Route() in the API (not via registerRoute), so they
// were missed in the initial coverage scan.
//
//   Panchang (16): sunrise-sunset, moonrise-moonset, nakshatra, yoga,
//     karana, vara, yamaganda, gulika, abhijit, pradosh-kaal,
//     durmuhurta, hindu-month, panchaka, varjyam, amrit-kalam,
//     siddha-yoga
//   Chart (1): ghatak
//
// Total new: 17 tools.

func registerV09Tools(s *mcp.Server, c *apiClient) {
	// Panchang sub-endpoints
	registerPanchangSunriseSunset(s, c)
	registerPanchangMoonriseMoonset(s, c)
	registerPanchangNakshatra(s, c)
	registerPanchangYoga(s, c)
	registerPanchangKarana(s, c)
	registerPanchangVara(s, c)
	registerPanchangYamaganda(s, c)
	registerPanchangGulika(s, c)
	registerPanchangAbhijit(s, c)
	registerPanchangPradoshKaal(s, c)
	registerPanchangDurmuhurta(s, c)
	registerPanchangHinduMonth(s, c)
	registerPanchangPanchaka(s, c)
	registerPanchangVarjyam(s, c)
	registerPanchangAmritKalam(s, c)
	registerPanchangSiddhaYoga(s, c)

	// Chart
	registerChartGhatak(s, c)
}

// =====================================================================
// Panchang sub-endpoints (16 tools)
// All use PanchangAtInput (lat/lon/tz/date) declared in tools_v02.go.
// =====================================================================

func registerPanchangSunriseSunset(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_sunrise_sunset",
		Description: "Get the precise sunrise and sunset times for a date and location — returns solar event timestamps, duration of daylight, and the nakshatra/tithi at sunrise. Use for 'what time is sunrise today', 'sunset time in Mumbai', 'dawn and dusk times'.",
		Title:       "Sunrise & Sunset Times",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/sunrise-sunset", in.toQuery())
	})
}

func registerPanchangMoonriseMoonset(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_moonrise_moonset",
		Description: "Get the moonrise and moonset times for a date and location — returns timestamps for when the Moon rises and sets, its current phase, and nakshatra. Use for 'what time does the Moon rise tonight', 'moonrise time', 'chandrodaya samay'.",
		Title:       "Moonrise & Moonset Times",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/moonrise-moonset", in.toQuery())
	})
}

func registerPanchangNakshatra(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_nakshatra",
		Description: "Get the Moon's nakshatra (lunar mansion) for a date and location — returns the current nakshatra name, pada, ruling planet, start/end timestamps, and the next nakshatra. Use for 'what nakshatra is today', 'which nakshatra is the Moon in', 'today's chandra nakshatra'.",
		Title:       "Daily Nakshatra",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/nakshatra", in.toQuery())
	})
}

func registerPanchangYoga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_yoga",
		Description: "Get the panchang yoga for a date — one of the 27 yogas (Vishkumbha through Vaidhriti) determined by the combined longitude of Sun and Moon. Some yogas are auspicious (Siddha, Shubha, Amrita) and others inauspicious (Vishkumbha, Atiganda, Vyaghata). Use for 'what panchang yoga is today', 'is today's yoga auspicious', 'which of the 27 yogas is active'.",
		Title:       "Panchang Yoga",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/yoga", in.toQuery())
	})
}

func registerPanchangKarana(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_karana",
		Description: "Get the panchang karana for a date and location — one of the 11 karanas (half-tithis): Bava, Balava, Kaulava, Taitila, Gara, Vanija, Vishti (Bhadra), Shakuni, Chatushpada, Naga, Kimstughna. Vishti/Bhadra karana is inauspicious for starting activities. Use for 'what karana is today', 'is this Bhadra', 'karana for this tithi'.",
		Title:       "Panchang Karana",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/karana", in.toQuery())
	})
}

func registerPanchangVara(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_vara",
		Description: "Get the Hindu vara (weekday) details for a date — returns the vara name (Ravivara/Somavara/…), its ruling planet, the vara's general auspiciousness for different activities, and the vara lord's strength. Use for 'what vara is today', 'which planet rules Wednesday', 'is today's vara good for travel'.",
		Title:       "Panchang Vara (Weekday)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/vara", in.toQuery())
	})
}

func registerPanchangYamaganda(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_yamaganda",
		Description: "Get the Yamaganda muhurta window for a date and location — a daily inauspicious 90-minute period (varies by weekday) ruled by Yama; activities started during Yamaganda are considered ill-starred. Returns start and end timestamps. Use for 'what time is Yamaganda today', 'avoid Yamaganda for this activity', 'Yama Kanda time'.",
		Title:       "Yamaganda (Inauspicious Window)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/yamaganda", in.toQuery())
	})
}

func registerPanchangGulika(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_gulika",
		Description: "Get the Gulika Kaal muhurta window for a date and location — a daily inauspicious period associated with Saturn's son Gulika (Mandi). Considered inauspicious for starting new ventures. Returns start and end timestamps. Use for 'what time is Gulika Kaal today', 'when is Mandi Kaal', 'avoid Gulika for muhurta'.",
		Title:       "Gulika Kaal Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/gulika", in.toQuery())
	})
}

func registerPanchangAbhijit(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_abhijit",
		Description: "Get the Abhijit muhurta window for a date and location — the auspicious 48-minute period around solar noon, considered one of the most powerful muhurtas for new beginnings (except Wednesdays when it is weak). Returns start and end timestamps. Use for 'when is Abhijit muhurta today', 'noon muhurta', 'best time today for starting something'.",
		Title:       "Abhijit Muhurta Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/abhijit", in.toQuery())
	})
}

func registerPanchangPradoshKaal(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_pradosh_kaal",
		Description: "Get the Pradosh Kaal window for a date and location — the auspicious 1.5-hour twilight period after sunset especially sacred to Lord Shiva. Falls on Trayodashi tithi (13th lunar day). Returns the window timestamps and whether it is Pradosh Vrat day. Use for 'when is Pradosh Kaal today', 'Shiva puja time this evening', 'Trayodashi pradosh window'.",
		Title:       "Pradosh Kaal",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/pradosh-kaal", in.toQuery())
	})
}

func registerPanchangDurmuhurta(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_durmuhurta",
		Description: "Get the Durmuhurta windows for a date and location — two inauspicious 48-minute periods in the day determined by the weekday (different for each vara). Starting anything important during Durmuhurta is avoided. Returns start/end timestamps for all Durmuhurta windows of the day. Use for 'when is Durmuhurta today', 'avoid Durmuhurta for this meeting', 'inauspicious times today'.",
		Title:       "Durmuhurta Windows",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/durmuhurta", in.toQuery())
	})
}

func registerPanchangHinduMonth(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_hindu_month",
		Description: "Get the current Hindu lunar month details for a date — returns Purnimanta and Amanta month name (Chaitra/Vaishakha/…), paksha (Shukla/Krishna), Vikram Samvat, Shaka Samvat, and whether it is an Adhik (intercalary) month. Use for 'what Hindu month is today', 'which Shukla paksha day is it', 'current Vikram Samvat date'.",
		Title:       "Hindu Month Details",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/hindu-month", in.toQuery())
	})
}

func registerPanchangPanchaka(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_panchaka",
		Description: "Check if the current time is in Panchaka — the inauspicious period when the Moon transits Aquarius or Pisces (creates obstacles for travel, construction, funerals, and starting new activities). Returns whether Panchaka is active, which of the 5 types it is (Roga/Mrityu/Agni/Raja/Chora), and the end time. Use for 'is it Panchaka right now', 'panchaka today', 'should I avoid this activity due to panchaka'.",
		Title:       "Panchaka Status",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/panchaka", in.toQuery())
	})
}

func registerPanchangVarjyam(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_varjyam",
		Description: "Get the Varjyam window for a date and location — a daily inauspicious 1h36m period within a nakshatra determined by the Moon's position; all activities including travel, medical procedures, and ceremonies are avoided. Returns start/end timestamps. Use for 'when is Varjyam today', 'varjyam window to avoid', 'inauspicious nakshatra period today'.",
		Title:       "Varjyam Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/varjyam", in.toQuery())
	})
}

func registerPanchangAmritKalam(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_amrit_kalam",
		Description: "Get the Amrit Kalam window for a date and location — an auspicious 1h36m period within the current nakshatra, considered especially beneficial for starting important activities and medical treatments. The counterpart to Varjyam within the same nakshatra. Returns start/end timestamps. Use for 'when is Amrit Kalam today', 'amrit kaal timing', 'best nakshatra window today'.",
		Title:       "Amrit Kalam Window",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/amrit-kalam", in.toQuery())
	})
}

func registerPanchangSiddhaYoga(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "panchang_siddha_yoga",
		Description: "Check if the current date has Siddha Yoga — an auspicious panchang yoga formed by specific vara+nakshatra combinations (e.g. Sunday+Hasta, Monday+Mrigashira, Saturday+Rohini, etc.). Siddha Yoga is highly auspicious for starting new ventures and ceremonies. Returns whether it is active and which vara-nakshatra pair created it. Use for 'is today Siddha Yoga', 'Siddha Yoga dates this month', 'best yoga for beginning something'.",
		Title:       "Siddha Yoga Check",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in PanchangAtInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/panchang/siddha-yoga", in.toQuery())
	})
}

// =====================================================================
// Chart — Ghatak Chakra
// =====================================================================

func registerChartGhatak(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_ghatak",
		Description: "Get the Ghatak Chakra (inauspicious indicators) for a birth chart — identifies the inauspicious weekday, tithi, nakshatra, lagna, and month specific to this chart's Moon sign. These are the conditions most harmful to the native and should be avoided for muhurta. Use for 'what is my ghatak nakshatra', 'ghatak chakra for my chart', 'which day is inauspicious for my rashi'.",
		Title:       "Ghatak Chakra",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in BirthInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/chart/ghatak", in.toQuery())
	})
}

// PanchangAtInput is declared in tools_v02.go — referenced here for all 16 panchang handlers.
// BirthInput is declared in tools.go — used by registerChartGhatak.
