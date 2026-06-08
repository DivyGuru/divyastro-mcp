# divyastro-mcp

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/DivyGuru/divyastro-mcp.svg)](https://pkg.go.dev/github.com/DivyGuru/divyastro-mcp)
[![Release](https://img.shields.io/github/v/release/DivyGuru/divyastro-mcp)](https://github.com/DivyGuru/divyastro-mcp/releases)

**Model Context Protocol (MCP) server for [DivyAstroAPI](https://divyastroapi.com)** — exposes 288 Vedic & Western astrology tools to AI assistants like Claude Desktop, Cursor, Continue.dev, and Codex CLI.

End users don't write code, copy curl commands, or read API docs. They just chat with their AI assistant naturally and the assistant invokes the right tool when an astrology question comes up:

> **User:** Mera lagna kya hai? Birth: 15 Jan 1990, 10:30 AM, Mumbai.
>
> **Claude:** *(calls `chart_ascendant`)* Aapka lagna **Capricorn (Makara)** hai, 12°34′ par. Lagna lord Saturn 10th house mein sthit hai...

## Architecture

```
┌─────────────────────────┐
│  USER LAPTOP            │
│  Claude Desktop / Cursor│
│        ↓ stdio          │
│  divyastro-mcp (this)   │
│        ↓ HTTPS          │
└────────│────────────────┘
         ↓ public internet
  api.divyastroapi.com
```

`divyastro-mcp` runs as a **subprocess of the AI client on your laptop** — not on any server. It forwards tool calls as authenticated HTTPS requests to `api.divyastroapi.com`. Every call costs the same number of credits as a direct curl; there's no MCP surcharge.

## Installation

### macOS / Linux — Homebrew (recommended)

```bash
brew tap DivyGuru/divyastro-mcp https://github.com/DivyGuru/divyastro-mcp
brew install divyastro-mcp

# Verify
divyastro-mcp 2>&1 | head -1
# → "divyastro-mcp: DIVYASTRO_API_KEY env var is required."
```

### Any platform — pre-built binary

1. Visit the [Releases page](https://github.com/DivyGuru/divyastro-mcp/releases)
2. Download the `divyastro-mcp-<os>-<arch>` binary matching your machine
3. Make it executable and put it on your `PATH`:

   ```bash
   # macOS / Linux
   chmod +x divyastro-mcp-darwin-arm64
   sudo mv divyastro-mcp-darwin-arm64 /usr/local/bin/divyastro-mcp

   # Windows: rename to divyastro-mcp.exe and add the folder to PATH
   ```

### Go developers — `go install`

```bash
go install github.com/DivyGuru/divyastro-mcp/cmd/divyastro-mcp@latest
```

The binary lands in `$GOBIN` (or `$GOPATH/bin`).

### Build from source

```bash
git clone https://github.com/DivyGuru/divyastro-mcp.git
cd divyastro-mcp
make build
# Binary: ./divyastro-mcp
```

## Get an API key

Sign up at [divyastroapi.com](https://divyastroapi.com) and create a key from your dashboard. The MCP server uses the same `dv_live_<hex>` keys as direct API access — there's no separate "MCP key".

Free tier (5,000 calls/month, 500/day) is enough to evaluate every tool.

## Configuring your AI client

The MCP server reads its credentials from environment variables that the AI client injects at subprocess spawn. Add the config block below to your client's MCP config file.

### Claude Desktop

Config file location:
- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

Add (or merge with existing `mcpServers`):

```json
{
  "mcpServers": {
    "divyastro": {
      "command": "/usr/local/bin/divyastro-mcp",
      "env": {
        "DIVYASTRO_API_KEY": "dv_live_REPLACE_WITH_YOUR_KEY"
      }
    }
  }
}
```

Quit and reopen Claude Desktop. The tools will appear in the tool list at the bottom of the chat window.

### Cursor

Cursor reads `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "divyastro": {
      "command": "/usr/local/bin/divyastro-mcp",
      "env": {
        "DIVYASTRO_API_KEY": "dv_live_REPLACE_WITH_YOUR_KEY"
      }
    }
  }
}
```

Restart Cursor (`Cmd-Shift-P → Reload Window`).

### OpenAI Codex CLI

Codex CLI reads `~/.codex/config.toml`:

```toml
[mcp_servers.divyastro]
command = "/usr/local/bin/divyastro-mcp"
env = { DIVYASTRO_API_KEY = "dv_live_REPLACE_WITH_YOUR_KEY" }
```

Restart Codex and the tools become available in your session.

### Continue.dev / Windsurf / other MCP clients

Any MCP-aware client supports the same shape — `command` + `env` with the absolute path to the binary and the API key. Refer to your client's docs for the config-file location.

## Tools

The MCP server exposes 288 tools across Vedic and Western astrology. See [docs/TOOLS.md](docs/TOOLS.md) for the full catalog.

**High-level coverage:**

| Domain | Count | Examples |
|---|---|---|
| Panchang | 26 | `panchang_today`, `panchang_basic`, `panchang_advanced`, `panchang_nakshatra`, `panchang_yoga`, `panchang_karana`, `panchang_sunrise_sunset`, `panchang_abhijit`, `panchang_siddha_yoga` |
| Birth Chart (Vedic) | 20 | `chart_ascendant`, `chart_planets`, `chart_houses`, `chart_divisional`, `chart_shadbala`, `chart_bhavabala`, `chart_kp_sublord`, `chart_ghatak` |
| Dasha | 5 | `dasha_current`, `dasha_vimshottari_full`, `dasha_yogini_current`, `dasha_yogini_full` |
| Compatibility (Milan) | 14 | `match_making_score`, `mangal_dosha`, `milan_dasha_sync`, `milan_navamsa_compat`, `milan_stree_dirgha` |
| Muhurta | 11 | `muhurta_vivah`, `muhurta_graha_pravesh`, `muhurta_vyapar`, `muhurta_yatra`, `muhurta_sarvartha_siddhi` |
| Transit | 9 | `transits_now`, `sade_sati`, `transit_double_transit`, `transit_tarabala`, `transit_vedha`, `transit_small_panoti` |
| Ashtakavarga | 3 | `ashtakavarga_sarva`, `ashtakavarga_kaksha`, `ashtakavarga_transit_score` |
| Calendar | 4 | `calendar_adhik_maas`, `calendar_month`, `calendar_ritu`, `calendar_samvatsara` |
| Sky Events | 4 | `eclipses_solar`, `eclipses_lunar`, `festivals_month`, `festivals_on_date` |
| Planet Moments | 4 | `planet_retrograde_window`, `planet_ingress`, `planet_combustion_window`, `planet_speed` |
| Horary (Prashna) | 5 | `prashna_answer`, `prashna_chart`, `prashna_arudha`, `prashna_significators` |
| Numerology | 7 | `numerology_full`, `numerology_driver`, `numerology_conductor`, `numerology_soul`, `numerology_destiny` |
| Varshaphal | 8 | `varshaphal_chart`, `varshaphal_lord`, `varshaphal_muntha`, `varshaphal_yoga`, `varshaphal_harsha_bala` |
| Geo | 3 | `geo_reverse`, `geo_search`, `geo_timezone` |
| Vedic Narrative | 32 | `narrative_profile`, `narrative_career_outlook`, `narrative_marriage_outlook`, `narrative_yogas`, `narrative_doshas` |
| Reports | 12 | `report_kundli_brihad`, `report_kundli_detailed`, `report_match_making`, `report_varshaphal` |
| Vedic Remedies | 1 | `vedic_remedies` |
| Horoscope | 5 | `horoscope_daily`, `horoscope_weekly`, `narrative_horoscope_daily_by_lagna`, `narrative_horoscope_weekly_by_moon` |
| Western Astrology | 46 | `western_natal_planets`, `western_synastry`, `western_composite_planets`, `western_davison_houses`, `western_dignities_receptions`, `western_narrative_firdaria` |

## Troubleshooting

**The tools don't appear in Claude Desktop.**
Check the Claude Desktop logs at `~/Library/Logs/Claude/mcp*.log`. The most common errors:

- `command not found` — the path in `command` is wrong. Use `which divyastro-mcp` to find the absolute path and paste that.
- `DIVYASTRO_API_KEY env var is required` — the `env` block is missing or your API key has a typo.

**A tool call returns `API 401: Unauthorized`.**
Your API key is invalid, revoked, or expired. Generate a new one in the dashboard.

**A tool call returns `API 429`.**
You've hit your plan's rate limit or daily cap. Either upgrade or wait for the next reset window (free tier resets at midnight UTC).

**A tool call returns `context deadline exceeded`.**
The upstream API took longer than 15s — usually a network blip. Retry once. If it persists, check [api.divyastroapi.com/healthz](https://api.divyastroapi.com/healthz).

## Privacy

`divyastro-mcp` runs entirely on your local machine. It only sends data to `api.divyastroapi.com` (HTTPS) when an AI tool calls one of the registered tools. It does not phone home, collect telemetry, or log requests anywhere. The binary is statically linked Go — you can audit the source in this repository.

## Development

```bash
# Run tests
make test

# Format + vet
make fmt vet

# Build for current platform
make build

# Cross-compile for all release platforms
make build-all VERSION=0.3.0
```

## Roadmap

- **v0.3.x (current):** Local stdio MCP server, 288 tools
- **v1.0 (planned):** Stable tool surface; semantic versioning kicks in
- **Remote HTTP MCP (separate project):** Hosted MCP endpoint at `mcp.divyastroapi.com` for ChatGPT and web/mobile MCP clients — tracked separately from this repo

## License

[MIT](LICENSE) — Copyright (c) 2026 DivyGuru. You may use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of this software.

The MIT license applies to the MCP client code in this repository only. The DivyAstroAPI service it talks to (`api.divyastroapi.com`) is a separate commercial product governed by its own [Terms of Service](https://divyastroapi.com/terms).

## Support

- 🐛 Bug reports & feature requests: [GitHub Issues](https://github.com/DivyGuru/divyastro-mcp/issues)
- 📖 API documentation: [api.divyastroapi.com/docs](https://api.divyastroapi.com/docs)
- 💬 Questions: divyguru108@gmail.com
