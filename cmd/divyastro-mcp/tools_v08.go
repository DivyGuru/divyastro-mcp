package main

import (
	"context"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerV08Tools covers all remaining API routes not yet in the MCP server:
//
//   Ashtakavarga (1):  bhinna/{planet}
//   Chart (4):         planet/{name}, divisional/{varga}, shadbala/{component},
//                      ashtakvarga/{planet}
//   Dasha (5):         vimshottari/{md}, /{md}/{ad}, /{md}/{ad}/{pd},
//                      /{md}/{ad}/{pd}/{sd}, yogini/{name}
//   Eclipses (1):      all (solar+lunar combined)
//   Geo (1):           place/{id}
//   Meta (7):          endpoints, ayanamsas, version, vedic/yogas,
//                      vedic/doshas, vedic/locales, vedic/narrative-axes
//   Milan (1):         ashtakoota/{koota}
//   Tarot (1):         card/{id}
//   Vedic Narrative (3): divisional/{varga}, yoga/{id}, dosha/{id}
//   Vedic Remedies (4): {rule_id}, mantra/{id}, pooja/{id}, vrat/{id}
//   Western (6):       natal/planet/{name}, natal/lot/{name},
//                      natal/asteroid/{number}, natal/fixed-stars/{name},
//                      natal/harmonic/{n}, retrograde-window/{planet}
//   Western Narrative (1): fixed-star/{name}
//
// Total new: 35 tools  (plus /v1/docs skipped — HTML page, not JSON).

func registerV08Tools(s *mcp.Server, c *apiClient) {
	// Ashtakavarga
	registerAshtakavargaBhinna(s, c)

	// Chart — parameterised
	registerChartSinglePlanet(s, c)
	// registerChartDivisional already in registerV02Tools
	registerChartShadbalaComponent(s, c)
	// registerChartAshtakvargaPlanet removed — duplicate of chart_ashtakvarga (registerV02Tools)

	// Dasha — vimshottari sub-periods + yogini single
	registerDashaVimshottariMaha(s, c)
	registerDashaVimshottariAntar(s, c)
	registerDashaVimshottariPratyantar(s, c)
	registerDashaVimshottariSookshma(s, c)
	registerDashaYoginiFindByName(s, c)

	// Eclipses combined
	registerEclipsesAll(s, c)

	// Geo place lookup
	registerGeoPlace(s, c)

	// Meta tools excluded — static developer reference data, not useful for AI assistant queries.

	// Milan — single koota
	registerMilanSingleKoota(s, c)

	// Tarot — single card (registerTarotCard already in registerV04Tools)

	// Vedic Narrative — parameterised static
	registerNarrativeDivisional(s, c)
	registerNarrativeYogaByID(s, c)
	registerNarrativeDoshaByID(s, c)

	// Vedic Remedies — parameterised static
	registerRemedyByRuleID(s, c)
	registerRemedyMantraByID(s, c)
	registerRemedyPoojaByID(s, c)
	registerRemedyVratByID(s, c)

	// Western — parameterised natal + mundane
	registerWesternNatalSinglePlanet(s, c)
	registerWesternNatalLot(s, c)
	registerWesternNatalAsteroid(s, c)
	registerWesternNatalFixedStarByName(s, c)
	// registerWesternNatalHarmonic already in registerV03Tools
	// registerWesternRetrogradeWindow already in registerV03Tools
	// registerWesternNarrativeFixedStar already in registerV04Tools
}

// =====================================================================
// Ashtakavarga — Bhinna (planet-specific table)
// =====================================================================

type AshtakavargaBhinnaInput struct {
	BirthInput
	Planet string `json:"planet" jsonschema:"planet name: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu"`
}

func registerAshtakavargaBhinna(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "ashtakavarga_bhinna",
		Description: "Get the Bhinna (individual planet) Ashtakavarga table for one planet — the 12-sign grid showing how many of the 8 reference points (7 planets + lagna) contribute benefic points to each sign for that planet's transit. Max 8 per sign. Use for 'Jupiter ashtakavarga table', 'Saturn bhinna ashtakavarga', 'how many points does Mars have in each sign'.",
		Title:       "Bhinna Ashtakavarga (Single Planet)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in AshtakavargaBhinnaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/ashtakavarga/bhinna/"+in.Planet, in.BirthInput.toQuery())
	})
}

// =====================================================================
// Chart — parameterised endpoints
// =====================================================================

type ChartSinglePlanetInput struct {
	BirthInput
	Name string `json:"name" jsonschema:"planet name: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu"`
}

func registerChartSinglePlanet(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_single_planet",
		Description: "Get the detailed placement of a single planet in a Vedic birth chart — sign, degree, nakshatra, pada, house, speed, retrograde flag, dignity, and aspects received. More concise than chart_planets when you only need one graha. Use for 'where is Saturn in my chart', 'Jupiter placement details', 'Mars degree and nakshatra'.",
		Title:       "Single Planet Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ChartSinglePlanetInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/chart/planet/"+in.Name, q)
	})
}

type ChartDivisionalInput struct {
	BirthInput
	Varga string `json:"varga" jsonschema:"divisional chart code: D1, D2, D3, D4, D5, D6, D7, D8, D9, D10, D11, D12, D16, D20, D24, D27, D30, D40, D45, D60"`
}

// registerChartDivisional is declared in tools_v02.go — ChartDivisionalInput
// type is still needed here for the parameterised path builder below.

type ChartShadbalaComponentInput struct {
	BirthInput
	Component string `json:"component" jsonschema:"Shadbala component: sthana (positional), dig (directional), kala (temporal), cheshta (motional), naisargika (natural), drik (aspectual)"`
}

func registerChartShadbalaComponent(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "chart_shadbala_component",
		Description: "Get a specific Shadbala component for all planets — one of the 6 strength types: sthana bala (positional), dig bala (directional), kala bala (temporal), cheshta bala (motional/retrogression), naisargika bala (natural), drik bala (aspectual). Use when you need to understand why a planet is strong or weak in a specific dimension. Use for 'dig bala of all planets', 'temporal strength shadbala', 'cheshta bala'.",
		Title:       "Shadbala Component Breakdown",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in ChartShadbalaComponentInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/chart/shadbala/"+in.Component, q)
	})
}

// =====================================================================
// Dasha — Vimshottari sub-period drill-down + Yogini single
// =====================================================================

type DashaVimshottariMahaInput struct {
	BirthInput
	Md string `json:"md" jsonschema:"Mahadasha lord abbreviation: sun, moon, mars, mer, jup, ven, sat, rah, ket"`
}

func registerDashaVimshottariMaha(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_vimshottari_mahadasha",
		Description: "Get all Antardasha sub-periods within a specific Vimshottari Mahadasha — returns the 9 antardasha periods (with start/end dates) nested under the requested mahadasha lord. Use when you want 'all antardasha periods inside my Saturn mahadasha', 'sub-periods in Jupiter dasha', 'when does Rahu-Sun antardasha start in my Rahu mahadasha'.",
		Title:       "Vimshottari Antardasha List (within a Mahadasha)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaVimshottariMahaInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari/"+in.Md, q)
	})
}

type DashaVimshottariAntarInput struct {
	BirthInput
	Md string `json:"md" jsonschema:"Mahadasha lord abbreviation: sun, moon, mars, mer, jup, ven, sat, rah, ket"`
	Ad string `json:"ad" jsonschema:"Antardasha lord abbreviation: sun, moon, mars, mer, jup, ven, sat, rah, ket"`
}

func registerDashaVimshottariAntar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_vimshottari_antardasha",
		Description: "Get all Pratyantardasha sub-sub-periods within a specific Vimshottari Mahadasha + Antardasha combination. Use for 'all pratyantardasha inside my Saturn-Jupiter antardasha', 'sub-sub-periods in Rahu-Moon', 'drill down into a specific dasha pair'.",
		Title:       "Vimshottari Pratyantardasha List",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaVimshottariAntarInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari/"+in.Md+"/"+in.Ad, q)
	})
}

type DashaVimshottariPratyantarInput struct {
	BirthInput
	Md string `json:"md" jsonschema:"Mahadasha lord abbreviation: sun, moon, mars, mer, jup, ven, sat, rah, ket"`
	Ad string `json:"ad" jsonschema:"Antardasha lord abbreviation"`
	Pd string `json:"pd" jsonschema:"Pratyantardasha lord abbreviation"`
}

func registerDashaVimshottariPratyantar(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_vimshottari_pratyantardasha",
		Description: "Get all Sookshma sub-sub-sub-periods within a specific 3-level Vimshottari combination (Maha + Antar + Pratyantar). Use for 'sookshma periods inside Saturn-Jupiter-Mercury', 'drill 3 levels into dasha timeline'.",
		Title:       "Vimshottari Sookshma List",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaVimshottariPratyantarInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari/"+in.Md+"/"+in.Ad+"/"+in.Pd, q)
	})
}

type DashaVimshottariSookshmaInput struct {
	BirthInput
	Md string `json:"md" jsonschema:"Mahadasha lord abbreviation: sun, moon, mars, mer, jup, ven, sat, rah, ket"`
	Ad string `json:"ad" jsonschema:"Antardasha lord abbreviation"`
	Pd string `json:"pd" jsonschema:"Pratyantardasha lord abbreviation"`
	Sd string `json:"sd" jsonschema:"Sookshma lord abbreviation"`
}

func registerDashaVimshottariSookshma(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_vimshottari_sookshma",
		Description: "Get all Prana (5th level) sub-sub-sub-sub-periods within a specific 4-level Vimshottari combination. The finest dasha granularity — each Prana period spans a few days. Use for 'prana periods in Saturn-Jupiter-Mercury-Venus sookshma', '5th level vimshottari dasha'.",
		Title:       "Vimshottari Prana List (5th Level)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaVimshottariSookshmaInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/dasha/vimshottari/"+in.Md+"/"+in.Ad+"/"+in.Pd+"/"+in.Sd, q)
	})
}

type DashaYoginiByNameInput struct {
	BirthInput
	Name string `json:"name" jsonschema:"Yogini name: mangala, pingala, dhanya, bhramari, bhadrika, ulka, siddha, sankata"`
}

func registerDashaYoginiFindByName(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dasha_yogini_by_name",
		Description: "Get the Yogini Dasha sub-periods for a specific Yogini lord — when does a particular Yogini period (Mangala, Pingala, Dhanya, Bhramari, Bhadrika, Ulka, Siddha, or Sankata) run and what are its sub-periods. Use for 'when is my Siddha yogini period', 'Ulka yogini sub-periods', 'find Bhadrika in my timeline'.",
		Title:       "Yogini Dasha by Name",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in DashaYoginiByNameInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/dasha/yogini/"+in.Name, q)
	})
}

// =====================================================================
// Eclipses — combined solar + lunar
// =====================================================================

func registerEclipsesAll(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "eclipses_all",
		Description: "Get all upcoming eclipses (both solar and lunar) in a single call — returns the combined list of solar and lunar eclipses in chronological order with type, date, visibility path, and magnitude. Use for 'all eclipses this year', 'next eclipse solar or lunar', 'eclipse calendar 2026'.",
		Title:       "All Eclipses (Solar + Lunar Combined)",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in EclipseRangeInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/eclipses/all", in.toQuery())
	})
}

// =====================================================================
// Geo — place lookup by ID
// =====================================================================

type GeoPlaceInput struct {
	ID string `json:"id" jsonschema:"place ID returned by geo_search (e.g. 'IN:MH:Mumbai' or a numeric geoname ID)"`
}

func registerGeoPlace(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "geo_place",
		Description: "Get full details for a specific place by its geo ID — resolves a place ID returned by geo_search into complete details (name, state, country, lat, lon, elevation, IANA timezone, population). Use when you have a place ID from a previous search and need to resolve its full data. Use for 'get place details by ID', 'resolve geo ID'.",
		Title:       "Geo Place Lookup by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in GeoPlaceInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		return callPassthrough(ctx, c, "/v1/geo/place/"+in.ID, q)
	})
}

// =====================================================================
// Milan — single Ashtakoota koota breakdown
// =====================================================================

type MilanSingleKootaInput struct {
	AshtakootaInput
	Koota string `json:"koota" jsonschema:"koota name: varna, vashya, tara, yoni, graha_maitri, gana, bhakoot, nadi"`
}

func registerMilanSingleKoota(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "milan_ashtakoota_single",
		Description: "Get the score and detailed explanation for a single Ashtakoota koota — one of: varna, vashya, tara, yoni, graha_maitri, gana, bhakoot, or nadi. Returns the score for that koota with the classical reasoning. Use for 'explain our nadi koota', 'what is our yoni koota score', 'bhakoot dosha details', 'why did we score 0 on nadi'.",
		Title:       "Single Ashtakoota Koota Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in MilanSingleKootaInput) (*mcp.CallToolResult, any, error) {
		return callPassthrough(ctx, c, "/v1/milan/ashtakoota/"+in.Koota, boyGirlQuery(in.AshtakootaInput))
	})
}

// registerTarotCard is declared in tools_v04.go (integer ID, 0–77).

// =====================================================================
// Vedic Narrative — parameterised static content
// =====================================================================

type NarrativeDivisionalInput struct {
	BirthInput
	Varga string `json:"varga" jsonschema:"divisional chart code: D1, D2, D3, D4, D5, D6, D7, D8, D9, D10, D11, D12, D16, D20, D24, D27, D30, D40, D45, D60"`
	Lang  string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func registerNarrativeDivisional(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_divisional",
		Description: "Get a narrative interpretation of a specific divisional (varga) chart — a written analysis of the lagna, key planet placements, and the life area governed by that divisional chart (e.g. D10 for career, D7 for children, D9 for marriage). Use for 'interpret my D10 chart', 'D9 navamsa reading', 'what does my career chart say'.",
		Title:       "Divisional Chart Narrative",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeDivisionalInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/divisional/"+in.Varga, q)
	})
}

type NarrativeYogaByIDInput struct {
	BirthInput
	ID   string `json:"id" jsonschema:"yoga ID from meta_vedic_yogas (e.g. 'hamsa', 'malavya', 'raja_yoga_1_5', 'kemadruma')"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func registerNarrativeYogaByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_yoga_by_id",
		Description: "Get a narrative interpretation of a specific yoga in a birth chart — detailed explanation of what the yoga means, how strongly it is formed in this chart, its classical effects, and what the native can expect from it. Use meta_vedic_yogas to find yoga IDs. Use for 'explain my Hamsa yoga', 'how strong is my Raja Yoga', 'interpret Kemadruma in my chart'.",
		Title:       "Yoga Narrative by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeYogaByIDInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/yoga/"+in.ID, q)
	})
}

type NarrativeDoshaByIDInput struct {
	BirthInput
	ID   string `json:"id" jsonschema:"dosha ID from meta_vedic_doshas (e.g. 'mangal_dosha', 'kaal_sarp', 'pitru_dosha', 'grahan_yoga')"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr, bn, kn, ta, te, gu"`
}

func registerNarrativeDoshaByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "narrative_dosha_by_id",
		Description: "Get a narrative interpretation of a specific dosha in a birth chart — whether it is present, its severity, classical cancellation conditions that apply, its effects on the native's life, and remedial suggestions. Use for 'explain my Kaal Sarp Dosha', 'is my Pitru Dosha cancelled', 'Mangal Dosha details from chart'.",
		Title:       "Dosha Narrative by ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in NarrativeDoshaByIDInput) (*mcp.CallToolResult, any, error) {
		q := in.BirthInput.toQuery()
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/narrative/dosha/"+in.ID, q)
	})
}

// =====================================================================
// Vedic Remedies — parameterised static content
// =====================================================================

type RemedyByRuleIDInput struct {
	ID   string `json:"id" jsonschema:"rule ID for the remedy (e.g. a planet-in-sign or yoga rule ID that has an associated remedy)"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerRemedyByRuleID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_by_rule_id",
		Description: "Get the remedy associated with a specific astrological rule ID — returns the complete remedy package (mantras, poojas, vratas, gemstone, yantra, color, number) for a rule that has a remedy attached. Use when you have a specific rule ID and want its remedy. Use for 'remedy for rule saturn_in_aries', 'what remedy is associated with this rule'.",
		Title:       "Remedy by Rule ID",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByRuleIDInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/"+in.ID, q)
	})
}

type RemedyByIDInput struct {
	ID   string `json:"id" jsonschema:"remedy item ID — use meta_vedic_yogas or vedic_remedies to discover valid IDs"`
	Lang string `json:"lang,omitempty" jsonschema:"response language: en (default), hi, mr"`
}

func registerRemedyMantraByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_mantra_by_id",
		Description: "Get the full details of a specific mantra remedy by ID — Sanskrit text, transliteration, meaning, recitation count, timing (day/nakshatra), and deity. Use for 'get details of mantra remedy X', 'how to recite this mantra', 'mantra text and meaning'.",
		Title:       "Mantra Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/mantra/"+in.ID, q)
	})
}

func registerRemedyPoojaByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_pooja_by_id",
		Description: "Get the full details of a specific pooja (ritual worship) remedy by ID — puja name, deity, procedure summary, recommended day/nakshatra, materials needed, and expected benefit. Use for 'how to perform this pooja', 'pooja details for Saturn', 'Shani puja procedure'.",
		Title:       "Pooja Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/pooja/"+in.ID, q)
	})
}

func registerRemedyVratByID(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "remedy_vrat_by_id",
		Description: "Get the full details of a specific vrat (fasting/observance) remedy by ID — vrat name, deity, day of observance, procedure, food restrictions, duration, and expected benefit. Use for 'how to observe this vrat', 'Shanivar vrat details', 'fasting remedy procedure'.",
		Title:       "Vrat Remedy Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RemedyByIDInput) (*mcp.CallToolResult, any, error) {
		q := url.Values{}
		if in.Lang != "" {
			q.Set("lang", in.Lang)
		}
		return callPassthrough(ctx, c, "/v1/vedic/remedies/vrat/"+in.ID, q)
	})
}

// =====================================================================
// Western — parameterised natal + mundane (6 tools)
// =====================================================================

type WesternSinglePlanetInput struct {
	WesternBirthInput
	Name string `json:"name" jsonschema:"planet name: Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn, Uranus, Neptune, Pluto, Chiron"`
}

func registerWesternNatalSinglePlanet(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_single_planet",
		Description: "Get detailed data for a single planet in a Western (tropical) birth chart — sign, degree, house, speed, retrograde flag, dignity, aspects to other planets, and declination. Use for 'where is Pluto in my western chart', 'Venus tropical sign', 'Chiron placement western'.",
		Title:       "Western Natal Single Planet Detail",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternSinglePlanetInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/western/natal/planet/"+in.Name, q)
	})
}

type WesternLotInput struct {
	WesternBirthInput
	Name string `json:"name" jsonschema:"Arabic lot name: fortune, spirit, eros, necessity, courage, victory, nemesis, basis"`
}

func registerWesternNatalLot(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_lot",
		Description: "Compute a specific Arabic Lot (Part) in a Western birth chart — Lot of Fortune, Spirit, Eros, Necessity, Courage, Victory, Nemesis, or Basis. Returns the lot's zodiac position, sign, house, and brief interpretation. Use for 'Part of Fortune in my chart', 'Lot of Spirit', 'Arabic lot calculation'.",
		Title:       "Western Arabic Lot",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternLotInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/western/natal/lot/"+in.Name, q)
	})
}

type WesternAsteroidInput struct {
	WesternBirthInput
	Number int `json:"number" jsonschema:"asteroid number (e.g. 1 for Ceres, 2 for Pallas, 3 for Juno, 4 for Vesta, 2060 for Chiron, 433 for Eros)"`
}

func registerWesternNatalAsteroid(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_asteroid",
		Description: "Get the position of a specific asteroid in a Western birth chart by its MPC catalog number — sign, degree, house, retrograde flag. Supports all numbered asteroids in the Swiss Ephemeris database. Use for 'where is Ceres in my chart (asteroid 1)', 'Vesta position (asteroid 4)', 'asteroid 433 Eros placement'.",
		Title:       "Western Natal Asteroid",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternAsteroidInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/western/natal/asteroid/"+strconv.Itoa(in.Number), q)
	})
}

type WesternFixedStarInput struct {
	WesternBirthInput
	Name string `json:"name" jsonschema:"fixed star name (e.g. Algol, Regulus, Spica, Antares, Aldebaran, Fomalhaut, Achernar, Sirius, Vega, Arcturus)"`
}

func registerWesternNatalFixedStarByName(s *mcp.Server, c *apiClient) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "western_natal_fixed_star",
		Description: "Check if a specific fixed star is conjunct any natal planet or angle in a Western birth chart — returns the star's tropical longitude, which natal point it aspects (if any), orb in degrees, and classical interpretation. Use for 'is Algol conjunct my chart', 'Regulus on my Ascendant', 'fixed star Spica placement'.",
		Title:       "Western Natal Fixed Star Conjunction",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in WesternFixedStarInput) (*mcp.CallToolResult, any, error) {
		q := in.WesternBirthInput.toQuery()
		return callPassthrough(ctx, c, "/v1/western/natal/fixed-stars/"+in.Name, q)
	})
}

type WesternHarmonicInput struct {
	WesternBirthInput
	N int `json:"n" jsonschema:"harmonic number (2–36); common ones: 4 (square series), 5 (quintile), 7 (septile), 8 (semi-square series), 9 (novile)"`
}

// registerWesternNatalHarmonic, registerWesternRetrogradeWindow, and
// registerWesternNarrativeFixedStar are declared in tools_v03.go / tools_v04.go.
// WesternHarmonicInput type kept here — it is referenced by tools_v03.go
// only if declared there; check below shows it is NOT declared there,
// so we keep the type but skip the function body.
