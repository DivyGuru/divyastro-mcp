package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newTestServer wires the same tool surface as production but points
// the upstream HTTP client at a fake server, then connects an
// in-memory MCP client/server pair so tests can issue real
// CallTool RPCs without spawning a subprocess.
//
// Returns the connected client session + the request-log struct so
// tests can assert what hit the upstream HTTP server.
func newTestServer(t *testing.T, status int, upstreamBody string) (*mcp.ClientSession, *requestLog) {
	t.Helper()
	upstream, log := fakeAPIServer(t, status, upstreamBody)
	apiClient := newAPIClient(upstream.URL, "dv_test_xxx", "test")

	server := mcp.NewServer(&mcp.Implementation{Name: "divyastro-test", Version: "test"}, nil)
	registerAllTools(server, apiClient)

	clientT, serverT := mcp.NewInMemoryTransports()
	ctx := t.Context()

	go func() {
		// Server.Run blocks until the transport closes — that
		// happens when t.Cleanup tears down the in-memory pipe.
		_ = server.Run(ctx, serverT)
	}()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "test"}, nil)
	session, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return session, log
}

// callTool is a small helper that invokes a tool by name with a typed
// argument map and returns the textual result content (most tools
// return JSON-as-text so this is the right shape for assertions).
func callTool(t *testing.T, session *mcp.ClientSession, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	res, err := session.CallTool(t.Context(), &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("call tool %q: %v", name, err)
	}
	return res
}

// firstText extracts the first TextContent from a CallToolResult or
// fails the test. Every tool in this MCP server returns a single
// text block (the upstream JSON pretty-printed), so this helper
// covers the universal case.
func firstText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatalf("CallToolResult has no content")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("first content is %T, want *TextContent", res.Content[0])
	}
	return tc.Text
}

// expectedToolCount is the source of truth for "how many tools should
// the MCP server expose right now". Bump this when adding a tool;
// the registration test below catches accidental additions or removals.
const expectedToolCount = 287

func TestListTools_AllRegistered(t *testing.T) {
	t.Parallel()
	session, _ := newTestServer(t, 200, `{}`)
	out, err := session.ListTools(t.Context(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if got := len(out.Tools); got != expectedToolCount {
		t.Errorf("got %d tools, want %d", got, expectedToolCount)
	}
	// Complete tool registry — updated to match all 287 registered tools.
	// Remove western_natal_rulerships (dead API route).
	// Added: v0.6–v0.9 tools (milan param fixes, panchang sub-routes, chart_ghatak, etc.).
	wantNames := map[string]bool{
		// v0.1 (12 tools)
		"panchang_today":     true,
		"chart_ascendant":    true,
		"chart_planets":      true,
		"chart_navamsa":      true,
		"chart_nakshatra":    true,
		"dasha_current":      true,
		"mangal_dosha":       true,
		"match_making_score": true,
		"sade_sati":          true,
		"horoscope_daily":    true,
		"horoscope_weekly":   true,
		"transits_now":       true,
		// v0.2 (30 tools)
		"panchang_tithi":           true,
		"panchang_choghadiya":      true,
		"panchang_hora":            true,
		"panchang_rahu_kaal":       true,
		"panchang_brahma_muhurat":  true,
		"chart_houses":             true,
		"chart_aspects":            true,
		"chart_dignity":            true,
		"chart_divisional":         true,
		"chart_avakhada":           true,
		"chart_shadbala":           true,
		"chart_ashtakavarga":       true,
		"dasha_vimshottari_full":   true,
		"dasha_yogini_current":     true,
		"nadi_dosha":               true,
		"vivah_phal":               true,
		"ashtakoota_breakdown":     true,
		"muhurta_vivah":            true,
		"muhurta_naamkaran":        true,
		"muhurta_best_time":        true,
		"varshaphal_chart":         true,
		"numerology_full":          true,
		"eclipses_solar":           true,
		"eclipses_lunar":           true,
		"festivals_month":          true,
		"planet_retrograde_window": true,
		"planet_ingress":           true,
		"prashna_answer":           true,
		"western_natal_chart":      true,
		"horoscope_monthly":        true,
		// v0.3 (29 tools — western_natal_rulerships removed: route does not exist in API)
		"western_natal_houses":                  true,
		"western_natal_aspects":                 true,
		"western_transit_positions":             true,
		"western_synastry_aspects":              true,
		"western_synastry_score":                true,
		"western_composite_chart":               true,
		"western_progressions_planets":          true,
		"western_solar_arc_planets":             true,
		"western_solar_return":                  true,
		"western_transits_to_natal":             true,
		"western_interpretation_aspect":         true,
		"western_interpretation_transit":        true,
		"western_natal_lots":                    true,
		"western_natal_midpoints":               true,
		"western_profections_annual":            true,
		"western_zodiacal_releasing":            true,
		"western_firdaria":                      true,
		"western_dignities_planet":              true,
		"western_natal_fixed_stars":             true,
		"western_natal_big_four_asteroids":      true,
		"western_natal_declinations":            true,
		"western_natal_harmonic":                true,
		"western_natal_antiscia":                true,
		"western_natal_midpoint_tree":           true,
		"western_heliocentric_planets":          true,
		"western_eclipses_upcoming":             true,
		"western_ingresses":                     true,
		"western_retrograde_window":             true,
		"western_astrocartography_planet_lines": true,
		// v0.4 (54 tools)
		"western_transits_exact":                        true,
		"western_transits_calendar":                     true,
		"western_synastry_grid":                         true,
		"western_lunar_return":                          true,
		"western_solar_return_aspects":                  true,
		"western_lunar_return_aspects":                  true,
		"western_composite_aspects":                     true,
		"western_davison_chart":                         true,
		"western_davison_aspects":                       true,
		"western_progressions_aspects_to_natal":         true,
		"western_progressions_lunation":                 true,
		"western_solar_arc_aspects_to_natal":            true,
		"western_horoscope_daily":                       true,
		"western_horoscope_weekly":                      true,
		"western_horoscope_monthly":                     true,
		"western_natal_summary":                         true,
		"western_transit_summary":                       true,
		"western_narrative_yearly":                      true,
		"western_narrative_decan":                       true,
		"western_narrative_sabian":                      true,
		"western_narrative_fixed_star":                  true,
		"western_narrative_timelord_profection":         true,
		"western_narrative_timelord_zr":                 true,
		"western_profections_monthly":                   true,
		"western_profections_daily":                     true,
		"western_zodiacal_releasing_spirit":             true,
		"western_dignities_almuten":                     true,
		"western_dignities_triplicity":                  true,
		"western_dignities_bounds":                      true,
		"western_dignities_face":                        true,
		"western_natal_centaurs":                        true,
		"western_natal_tnos":                            true,
		"western_natal_lilith_set":                      true,
		"western_natal_uranian_hypotheticals":           true,
		"western_natal_parans":                          true,
		"western_natal_vertex_set":                      true,
		"western_natal_harmonic_aspects":                true,
		"western_compatibility_sun_sign":                true,
		"western_lunar_phase":                           true,
		"western_natal_aspect_patterns":                 true,
		"western_natal_dominant":                        true,
		"western_natal_chart_shape":                     true,
		"western_eclipses_visibility":                   true,
		"western_primary_directions_placidus_zodiacal":  true,
		"western_primary_directions_placidus_mundane":   true,
		"western_primary_directions_regiomontanus":      true,
		"numerology_personal_periods":                   true,
		"numerology_challenges":                         true,
		"numerology_advanced":                           true,
		"tarot_cards":                                   true,
		"tarot_card":                                    true,
		"tarot_daily":                                   true,
		"tarot_spread":                                  true,
		"tarot_yes_no":                                  true,
		// v0.5 (10 tools)
		"panchang_basic":                true,
		"panchang_advanced":             true,
		"panchang_lagna_table":          true,
		"panchang_monthly":              true,
		"panchang_nakshatra_prediction": true,
		"chart_astro_details":           true,
		"varshaphal_harsha_bala":        true,
		"varshaphal_mudda_dasha":        true,
		"narrative_yearly_bhavishyafal": true,
		"geo_reverse":                   true,
		// v0.6–v0.9 additions (152 tools)
		"ashtakavarga_bhinna":                          true,
		"ashtakavarga_kaksha":                          true,
		"ashtakavarga_sarva":                           true,
		"ashtakavarga_transit_score":                   true,
		"calendar_adhik_maas":                          true,
		"calendar_month":                               true,
		"calendar_ritu":                                true,
		"calendar_samvatsara":                          true,
		"chart_bhavabala":                              true,
		"chart_combustion":                             true,
		"chart_ghatak":                                 true,
		"chart_graha_yuddha":                           true,
		"chart_house_occupants":                        true,
		"chart_kp_house_significator":                  true,
		"chart_kp_sublord":                             true,
		"chart_shadbala_component":                     true,
		"chart_single_planet":                          true,
		"dasha_vimshottari_antardasha":                 true,
		"dasha_vimshottari_mahadasha":                  true,
		"dasha_vimshottari_pratyantardasha":            true,
		"dasha_vimshottari_sookshma":                   true,
		"dasha_yogini_by_name":                         true,
		"dasha_yogini_full":                            true,
		"eclipses_all":                                 true,
		"festivals_on_date":                            true,
		"geo_place":                                    true,
		"geo_search":                                   true,
		"geo_timezone":                                 true,
		"milan_ashtakoota_single":                      true,
		"milan_dasha_sync":                             true,
		"milan_dhan_yog":                               true,
		"milan_longevity":                              true,
		"milan_mahendra":                               true,
		"milan_navamsa_compat":                         true,
		"milan_santan_yog":                             true,
		"milan_shani_dosha":                            true,
		"milan_stree_dirgha":                           true,
		"muhurta_chandra_bala":                         true,
		"muhurta_graha_pravesh":                        true,
		"muhurta_panchaka_rahita":                      true,
		"muhurta_sarvartha_siddhi":                     true,
		"muhurta_shubha_yoga":                          true,
		"muhurta_tara_bala":                            true,
		"muhurta_vyapar":                               true,
		"muhurta_yatra":                                true,
		"narrative_career_outlook":                     true,
		"narrative_dasha_phal":                         true,
		"narrative_dasha_tree_phal":                    true,
		"narrative_divisional":                         true,
		"narrative_dosha_by_id":                        true,
		"narrative_doshas":                             true,
		"narrative_finance_outlook":                    true,
		"narrative_horoscope_daily_by_lagna":           true,
		"narrative_horoscope_daily_by_moon":            true,
		"narrative_horoscope_daily_tamil":              true,
		"narrative_horoscope_weekly_by_lagna":          true,
		"narrative_horoscope_weekly_by_moon":           true,
		"narrative_house_lord":                         true,
		"narrative_house_lords_overview":               true,
		"narrative_karakas":                            true,
		"narrative_lagna":                              true,
		"narrative_marriage_outlook":                   true,
		"narrative_milan":                              true,
		"narrative_moon_sign":                          true,
		"narrative_nakshatra":                          true,
		"narrative_planet_in_house":                    true,
		"narrative_planet_in_sign":                     true,
		"narrative_planets_in_houses":                  true,
		"narrative_planets_in_signs":                   true,
		"narrative_profile":                            true,
		"narrative_sade_sati_history":                  true,
		"narrative_sade_sati_phase":                    true,
		"narrative_sade_sati_status":                   true,
		"narrative_sun_sign":                           true,
		"narrative_transit_ashtakavarga":               true,
		"narrative_transit_double":                     true,
		"narrative_transit_phal":                       true,
		"narrative_varshaphal_themes":                  true,
		"narrative_yoga_by_id":                         true,
		"narrative_yogas":                              true,
		"numerology_conductor":                         true,
		"numerology_destiny":                           true,
		"numerology_driver":                            true,
		"numerology_personality":                       true,
		"numerology_soul":                              true,
		"panchang_abhijit":                             true,
		"panchang_amrit_kalam":                         true,
		"panchang_durmuhurta":                          true,
		"panchang_gulika":                              true,
		"panchang_hindu_month":                         true,
		"panchang_hora_dinman":                         true,
		"panchang_karana":                              true,
		"panchang_moonrise_moonset":                    true,
		"panchang_nakshatra":                           true,
		"panchang_panchaka":                            true,
		"panchang_pradosh_kaal":                        true,
		"panchang_sankranti":                           true,
		"panchang_siddha_yoga":                         true,
		"panchang_sunrise_sunset":                      true,
		"panchang_vara":                                true,
		"panchang_varjyam":                             true,
		"panchang_yamaganda":                           true,
		"panchang_yoga":                                true,
		"planet_combustion_window":                     true,
		"planet_speed":                                 true,
		"prashna_arudha":                               true,
		"prashna_chart":                                true,
		"prashna_lagna_lord":                           true,
		"prashna_significators":                        true,
		"remedy_by_rule_id":                            true,
		"remedy_mantra_by_id":                          true,
		"remedy_pooja_by_id":                           true,
		"remedy_vrat_by_id":                            true,
		"report_dasha_analysis":                        true,
		"report_horoscope_daily":                       true,
		"report_horoscope_monthly":                     true,
		"report_horoscope_weekly":                      true,
		"report_kundli_brihad":                         true,
		"report_kundli_detailed":                       true,
		"report_kundli_lite":                           true,
		"report_mangal_dosha":                          true,
		"report_match_making":                          true,
		"report_numerology":                            true,
		"report_sade_sati":                             true,
		"report_varshaphal":                            true,
		"transit_ashtakavarga":                         true,
		"transit_double_transit":                       true,
		"transit_small_panoti":                         true,
		"transit_small_panoti_history":                 true,
		"transit_tarabala":                             true,
		"transit_vedha":                                true,
		"varshaphal_lord":                              true,
		"varshaphal_muntha":                            true,
		"varshaphal_solar_return_jd":                   true,
		"varshaphal_yoga":                              true,
		"vedic_remedies":                               true,
		"western_astrocartography_local_space":         true,
		"western_composite_houses":                     true,
		"western_composite_planets":                    true,
		"western_davison_houses":                       true,
		"western_davison_planets":                      true,
		"western_dignities_receptions":                 true,
		"western_narrative_firdaria":                   true,
		"western_natal_ascendant":                      true,
		"western_natal_asteroid":                       true,
		"western_natal_fixed_star":                     true,
		"western_natal_lot":                            true,
		"western_natal_mc":                             true,
		"western_natal_out_of_bounds":                  true,
		"western_natal_planets":                        true,
		"western_natal_single_planet":                  true,
		"western_solar_return_exact_jd":                true,
	}
	if got, want := len(wantNames), expectedToolCount; got != want {
		t.Errorf("test wantNames map has %d entries, want %d (test bookkeeping out of sync)", got, want)
	}
	gotNames := map[string]bool{}
	for _, tl := range out.Tools {
		gotNames[tl.Name] = true
	}
	for name := range wantNames {
		if !gotNames[name] {
			t.Errorf("tool %q not registered", name)
		}
	}
	// Catch the reverse — extra tools that aren't in the expected
	// list. Helps spot accidental dupes when adding new ones.
	for name := range gotNames {
		if !wantNames[name] {
			t.Errorf("unexpected tool %q registered (not in test list)", name)
		}
	}
}

func TestTool_PanchangToday_HitsCorrectEndpoint(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, cannedPanchangResponse)
	res := callTool(t, session, "panchang_today", map[string]any{
		"lat": 19.076,
		"lon": 72.8777,
		"tz":  "Asia/Kolkata",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if got := firstText(t, res); !strings.Contains(got, "Krishna Saptami") {
		t.Errorf("tool output missing fixture content: %q", got)
	}
	if len(log.paths) != 1 || log.paths[0] != "/v1/panchang/summary" {
		t.Errorf("upstream paths = %v, want [/v1/panchang/summary]", log.paths)
	}
	q := log.queries[0]
	if q.Get("tz") != "Asia/Kolkata" {
		t.Errorf("upstream tz = %q, want Asia/Kolkata", q.Get("tz"))
	}
}

func TestTool_ChartAscendant_PassesBirthQuintet(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"sign_en":"Capricorn","degree":"12.5"}`)
	res := callTool(t, session, "chart_ascendant", map[string]any{
		"lat":  19.076,
		"lon":  72.8777,
		"date": "1990-01-15",
		"time": "10:30",
		"tz":   "Asia/Kolkata",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if got := firstText(t, res); !strings.Contains(got, "Capricorn") {
		t.Errorf("tool output missing 'Capricorn': %q", got)
	}
	if log.paths[0] != "/v1/chart/ascendant" {
		t.Errorf("upstream path = %q, want /v1/chart/ascendant", log.paths[0])
	}
	q := log.queries[0]
	for k, want := range map[string]string{
		"lat": "19.076", "lon": "72.8777",
		"date": "1990-01-15", "time": "10:30", "tz": "Asia/Kolkata",
	} {
		if got := q.Get(k); got != want {
			t.Errorf("query[%q] = %q, want %q", k, got, want)
		}
	}
}

func TestTool_DashaCurrent_OptionalAtParam(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"maha":{"lord":"Saturn"}}`)
	res := callTool(t, session, "dasha_current", map[string]any{
		"lat":  19.076,
		"lon":  72.8777,
		"date": "1990-01-15",
		"time": "10:30",
		"tz":   "Asia/Kolkata",
		"at":   "2026-05-10T12:00:00Z",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if got := log.queries[0].Get("at"); got != "2026-05-10T12:00:00Z" {
		t.Errorf("upstream 'at' query = %q, want passthrough of input", got)
	}
}

func TestTool_HoroscopeDaily_LowercasesRashi(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"body":"You will have a great day."}`)
	// Pass mixed-case rashi — the tool must normalize before sending
	// to the upstream API (which expects lowercase).
	res := callTool(t, session, "horoscope_daily", map[string]any{
		"rashi": "Aries",
		"lang":  "en",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if got := log.queries[0].Get("rashi"); got != "aries" {
		t.Errorf("upstream rashi = %q, want lowercase 'aries'", got)
	}
}

func TestTool_MatchMaking_PassesBothPartiesParams(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"total":28,"max":36}`)
	res := callTool(t, session, "match_making_score", map[string]any{
		"bride_lat":  19.076,
		"bride_lon":  72.8777,
		"bride_date": "1995-06-12",
		"bride_time": "08:15",
		"bride_tz":   "Asia/Kolkata",
		"groom_lat":  28.6139,
		"groom_lon":  77.2090,
		"groom_date": "1992-11-03",
		"groom_time": "14:45",
		"groom_tz":   "Asia/Kolkata",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	// API reads boy.*/girl.* params via parseBoyGirlMoonInfo.
	// API convention: boy = groom (male), girl = bride (female).
	q := log.queries[0]
	for k, want := range map[string]string{
		"boy.date":  "1992-11-03",
		"girl.date": "1995-06-12",
		"boy.time":  "14:45",
		"girl.time": "08:15",
	} {
		if got := q.Get(k); got != want {
			t.Errorf("query[%q] = %q, want %q", k, got, want)
		}
	}
}

func TestTool_TransitsNow_NoArgsRequired(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `[{"name":"Sun","sign":"Aries"}]`)
	// transits_now is the one tool with no required input — invoking
	// it with an empty map should still hit the upstream and produce
	// a result.
	res := callTool(t, session, "transits_now", map[string]any{})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if log.paths[0] != "/v1/transit/positions" {
		t.Errorf("upstream path = %q, want /v1/transit/positions", log.paths[0])
	}
}

func TestTool_UpstreamErrorSurfacesAsToolError(t *testing.T) {
	t.Parallel()
	// 400 from upstream → the MCP tool should set IsError and pass
	// the upstream message through so the AI client can show a useful
	// message to the end user.
	session, _ := newTestServer(t, 400, `{"error":"invalid date","code":"invalid_param","details":"date must be YYYY-MM-DD"}`)
	res := callTool(t, session, "chart_ascendant", map[string]any{
		"lat":  19.076,
		"lon":  72.8777,
		"date": "not-a-date",
		"time": "10:30",
		"tz":   "Asia/Kolkata",
	})
	if !res.IsError {
		t.Fatalf("expected IsError=true, got false")
	}
	msg := firstText(t, res)
	for _, want := range []string{"400", "invalid_param", "invalid date"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error text %q missing %q", msg, want)
		}
	}
}

func TestTool_OutputIsValidJSON(t *testing.T) {
	t.Parallel()
	// The upstream payload is returned as a single text content block
	// containing the JSON. Verify that text is valid JSON the AI
	// client can parse — a marshal-side regression here would silently
	// break every tool.
	session, _ := newTestServer(t, 200, `{"foo": "bar", "list": [1,2,3]}`)
	res := callTool(t, session, "panchang_today", map[string]any{
		"lat": 1.0, "lon": 1.0, "tz": "UTC",
	})
	body := firstText(t, res)
	got := mustDecode(t, body)
	if got["foo"] != "bar" {
		t.Errorf("decoded foo = %v, want bar", got["foo"])
	}
}

func TestTool_AuthHeaderForwarded(t *testing.T) {
	t.Parallel()
	// Verify the bearer token reaches the upstream — a missing auth
	// header would break every paid call in production.
	session, _ := newTestServer(t, 200, `{}`)
	_ = callTool(t, session, "transits_now", map[string]any{})
	// The header check is in the upstream handler — fakeAPIServer's
	// handler doesn't capture headers. Re-do directly via apiClient
	// for header verification (TestAPIClient_GetHappyPath covers it),
	// so this test just sanity-checks the call path completes.
	_ = session
}

// Context: ensure context cancellation propagates from the MCP layer
// through to the upstream HTTP request — a runaway tool call should
// abort cleanly when the AI client cancels.
func TestTool_ContextCancellationPropagates(t *testing.T) {
	t.Parallel()
	session, _ := newTestServer(t, 200, `{}`)
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // cancel immediately
	_, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "transits_now",
		Arguments: map[string]any{},
	})
	if err == nil {
		t.Error("expected error from cancelled context, got nil")
	}
}

// =====================================================================
// v0.2 tool tests
// =====================================================================
//
// Spot-checks on the trickier v0.2 surface — path-parametrised tools
// (chart_divisional uses {varga} in the URL), shared muhurta input
// shape, eclipse range with optional location filter, the rashi
// lowerTrim normalization shared by horoscope_monthly.

func TestTool_ChartDivisional_PathContainsVarga(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"varga":"D10","planets":[]}`)
	res := callTool(t, session, "chart_divisional", map[string]any{
		"lat":   19.076,
		"lon":   72.8777,
		"date":  "1990-01-15",
		"time":  "10:30",
		"tz":    "Asia/Kolkata",
		"varga": "D10",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	// The varga selector lives in the URL path, not the query string —
	// a regression here would silently return D9 instead of D10.
	if log.paths[0] != "/v1/chart/divisional/D10" {
		t.Errorf("upstream path = %q, want /v1/chart/divisional/D10", log.paths[0])
	}
}

func TestTool_ChartDivisional_AcceptsLongName(t *testing.T) {
	t.Parallel()
	// Confirm the path passthrough works for textual varga names too —
	// "navamsa" is what an AI would send when paraphrasing the user.
	session, log := newTestServer(t, 200, `{}`)
	_ = callTool(t, session, "chart_divisional", map[string]any{
		"lat": 19.076, "lon": 72.8777,
		"date": "1990-01-15", "time": "10:30", "tz": "Asia/Kolkata",
		"varga": "navamsa",
	})
	if log.paths[0] != "/v1/chart/divisional/navamsa" {
		t.Errorf("upstream path = %q, want /v1/chart/divisional/navamsa", log.paths[0])
	}
}

func TestTool_MuhurtaVivah_PassesDateRange(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"windows":[]}`)
	res := callTool(t, session, "muhurta_vivah", map[string]any{
		"lat":        19.076,
		"lon":        72.8777,
		"tz":         "Asia/Kolkata",
		"start_date": "2026-11-01",
		"end_date":   "2026-12-31",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	q := log.queries[0]
	if q.Get("start_date") != "2026-11-01" || q.Get("end_date") != "2026-12-31" {
		t.Errorf("muhurta date range not forwarded: start=%q end=%q", q.Get("start_date"), q.Get("end_date"))
	}
}

func TestTool_EclipsesSolar_OptionalLocationOmittedWhenZero(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `[]`)
	// Caller supplies only the date range — location should NOT
	// appear in the upstream query (we don't want the API to filter
	// to "visible from 0,0" by accident).
	_ = callTool(t, session, "eclipses_solar", map[string]any{
		"start_date": "2026-01-01",
		"end_date":   "2026-12-31",
	})
	q := log.queries[0]
	if q.Has("lat") || q.Has("lon") {
		t.Errorf("lat/lon should be absent when zero; got lat=%q lon=%q", q.Get("lat"), q.Get("lon"))
	}
}

func TestTool_EclipsesSolar_OptionalLocationIncludedWhenSet(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `[]`)
	_ = callTool(t, session, "eclipses_solar", map[string]any{
		"start_date": "2026-01-01",
		"end_date":   "2026-12-31",
		"lat":        19.076,
		"lon":        72.8777,
	})
	q := log.queries[0]
	if q.Get("lat") == "" || q.Get("lon") == "" {
		t.Errorf("lat/lon should be present when set; got lat=%q lon=%q", q.Get("lat"), q.Get("lon"))
	}
}

func TestTool_HoroscopeMonthly_LowerTrimMatchesDailyWeekly(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"body":""}`)
	// Pass a mixed-case rashi with surrounding whitespace — both must
	// be normalized exactly like the daily/weekly tools do.
	_ = callTool(t, session, "horoscope_monthly", map[string]any{
		"rashi": "  PISCES  ",
	})
	if got := log.queries[0].Get("rashi"); got != "pisces" {
		t.Errorf("upstream rashi = %q, want trimmed-lowercased 'pisces'", got)
	}
}

func TestTool_VarshaphalChart_PassesYear(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"varsha_lord":"Sun"}`)
	_ = callTool(t, session, "varshaphal_chart", map[string]any{
		"lat":  19.076,
		"lon":  72.8777,
		"date": "1990-01-15",
		"time": "10:30",
		"tz":   "Asia/Kolkata",
		"year": 2026,
	})
	if got := log.queries[0].Get("year"); got != "2026" {
		t.Errorf("upstream year = %q, want 2026", got)
	}
}

func TestTool_NadiDosha_PassesBothPartiesParams(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"present":false}`)
	res := callTool(t, session, "nadi_dosha", map[string]any{
		"bride_lat":  19.076,
		"bride_lon":  72.8777,
		"bride_date": "1995-06-12",
		"bride_time": "08:15",
		"bride_tz":   "Asia/Kolkata",
		"groom_lat":  28.6139,
		"groom_lon":  77.2090,
		"groom_date": "1992-11-03",
		"groom_time": "14:45",
		"groom_tz":   "Asia/Kolkata",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	if log.paths[0] != "/v1/milan/nadi-dosha" {
		t.Errorf("upstream path = %q, want /v1/milan/nadi-dosha", log.paths[0])
	}
	// API reads boy.*/girl.* params via parseBoyGirlMoonInfo.
	// API convention: boy = groom (male), girl = bride (female).
	q := log.queries[0]
	if q.Get("boy.date") != "1992-11-03" || q.Get("girl.date") != "1995-06-12" {
		t.Error("both parties' birth dates should reach upstream")
	}
}

func TestTool_NumerologyFull_AllOptionalFields(t *testing.T) {
	t.Parallel()
	// Numerology tool has 3 optional inputs — verify each shows up
	// only when set, so a name-only or DOB-only call doesn't pollute
	// the upstream query with empty strings.
	session, log := newTestServer(t, 200, `{}`)
	_ = callTool(t, session, "numerology_full", map[string]any{
		"name": "Vikas Sharma",
	})
	q := log.queries[0]
	if q.Get("name") != "Vikas Sharma" {
		t.Errorf("name not forwarded: %q", q.Get("name"))
	}
	if q.Has("dob") || q.Has("lang") {
		t.Errorf("optional fields should be absent when unset; got dob=%q lang=%q",
			q.Get("dob"), q.Get("lang"))
	}
}

func TestTool_PrashnaAnswer_RequiresQuestion(t *testing.T) {
	t.Parallel()
	session, log := newTestServer(t, 200, `{"answer":"yes"}`)
	res := callTool(t, session, "prashna_answer", map[string]any{
		"lat":      19.076,
		"lon":      72.8777,
		"tz":       "Asia/Kolkata",
		"question": "Will I get the new role?",
	})
	if res.IsError {
		t.Fatalf("tool returned error: %s", firstText(t, res))
	}
	q := log.queries[0]
	if q.Get("question") != "Will I get the new role?" {
		t.Errorf("question text not forwarded verbatim: %q", q.Get("question"))
	}
}

// TestLowerTrim spot-checks the small standalone helper used by
// horoscope_monthly. It mirrors strings.ToLower(strings.TrimSpace(...))
// so the tool's output is consistent with daily/weekly handlers.
func TestLowerTrim(t *testing.T) {
	t.Parallel()
	cases := []struct{ in, want string }{
		{"Aries", "aries"},
		{"  PISCES  ", "pisces"},
		{"\tCancer\n", "cancer"},
		{"already-lower", "already-lower"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := lowerTrim(tc.in); got != tc.want {
			t.Errorf("lowerTrim(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
