# divyastro-mcp

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/DivyGuru/divyastro-mcp.svg)](https://pkg.go.dev/github.com/DivyGuru/divyastro-mcp)
[![Release](https://img.shields.io/github/v/release/DivyGuru/divyastro-mcp)](https://github.com/DivyGuru/divyastro-mcp/releases)

**Model Context Protocol (MCP) server for [DivyAstroAPI](https://divyastroapi.com)** — exposes 72 Vedic & Western astrology tools to AI assistants like Claude Desktop, Cursor, Continue.dev, and Codex CLI.

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

The MCP server exposes 72 tools across Vedic and Western astrology. See [docs/TOOLS.md](docs/TOOLS.md) for the full catalog.

**High-level coverage:**

| Domain | Count | Examples |
|---|---|---|
| Panchang | 6 | `panchang_today`, `panchang_tithi`, `panchang_choghadiya`, `panchang_rahu_kaal` |
| Birth Chart (Vedic) | 11 | `chart_ascendant`, `chart_planets`, `chart_houses`, `chart_divisional`, `chart_shadbala` |
| Dasha | 3 | `dasha_current`, `dasha_vimshottari_full`, `dasha_yogini_current` |
| Compatibility (Milan) | 5 | `match_making_score`, `ashtakoota_breakdown`, `mangal_dosha`, `nadi_dosha`, `vivah_phal` |
| Muhurta | 3 | `muhurta_vivah`, `muhurta_naamkaran`, `muhurta_best_time` |
| Transit | 2 | `transits_now`, `sade_sati` |
| Sky Events | 3 | `eclipses_solar`, `eclipses_lunar`, `festivals_month` |
| Planet Moments | 2 | `planet_retrograde_window`, `planet_ingress` |
| Horary / Numerology / Varshaphal | 3 | `prashna_answer`, `numerology_full`, `varshaphal_chart` |
| Horoscope | 3 | `horoscope_daily`, `horoscope_weekly`, `horoscope_monthly` |
| Western Astrology | 31 | `western_natal_chart`, `western_synastry`, `western_progressions`, `western_solar_return` |

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

- **v0.3.x (current):** Local stdio MCP server, 72 tools
- **v1.0 (planned):** Stable tool surface; semantic versioning kicks in
- **Remote HTTP MCP (separate project):** Hosted MCP endpoint at `mcp.divyastroapi.com` for ChatGPT and web/mobile MCP clients — tracked separately from this repo

## License

[MIT](LICENSE) — Copyright (c) 2026 DivyGuru. You may use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of this software.

The MIT license applies to the MCP client code in this repository only. The DivyAstroAPI service it talks to (`api.divyastroapi.com`) is a separate commercial product governed by its own [Terms of Service](https://divyastroapi.com/terms).

## Support

- 🐛 Bug reports & feature requests: [GitHub Issues](https://github.com/DivyGuru/divyastro-mcp/issues)
- 📖 API documentation: [api.divyastroapi.com/docs](https://api.divyastroapi.com/docs)
- 💬 Questions: divyguru108@gmail.com
