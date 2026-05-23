# Homebrew formula for divyastro-mcp.
#
# Two distribution paths from this single file:
#
#  1. Tap from THIS repo (immediate):
#       brew tap DivyGuru/divyastro-mcp https://github.com/DivyGuru/divyastro-mcp
#       brew install divyastro-mcp
#
#  2. Promote to the official homebrew-core later for `brew install divyastro-mcp`
#     without a tap. Requires a stable release history first (per Homebrew's
#     audit policy: 30 days + GitHub stars threshold).
#
# Maintainer release flow when cutting `v0.3.1`:
#   1. git tag v0.3.1 && git push origin v0.3.1
#   2. Wait for the GitHub Actions release workflow to publish binaries
#      + checksums.txt at the Release URL.
#   3. Update `version` below.
#   4. Update each `sha256` from dist/divyastro-mcp-checksums.txt.
#   5. Commit + push. Users `brew upgrade divyastro-mcp` picks it up.
#
# Why a binary formula (not source build)?
#  - The MCP server is pure Go and the release pipeline already
#    cross-compiles to every supported triple. Distributing pre-built
#    bottles avoids requiring Go on the user's machine and is faster
#    (sub-second `brew install`).
#  - Homebrew binary formulae are first-class — both `cask`-style
#    (apps) and `formula`-style (CLIs) support shipping pre-built
#    artifacts via `url` + `sha256`.

class DivyastroMcp < Formula
  desc "DivyAstroAPI MCP server — 72 Vedic + Western astrology tools for Claude Desktop, Cursor, Continue"
  homepage "https://github.com/DivyGuru/divyastro-mcp"
  version "0.3.0"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-darwin-arm64"
      sha256 "010267986f3ed065638f4219bfdf9cca7abc71ddb59f3125b3721955eabe99e8"
    end
    on_intel do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-darwin-amd64"
      sha256 "76a40df87f8c86368a85cadbd743add0d9111838167894494833e181140b848e"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-linux-arm64"
      sha256 "4295cde39a1508a136f382fd4eab459a311b689f4cace55b26669e3ad7b209e2"
    end
    on_intel do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-linux-amd64"
      sha256 "ed98f5e1168ef35fffa9012b0c1766d4e10eaacfa702a146fc47315ccd533827"
    end
  end

  def install
    # The downloaded artifact is the bare binary (not a tarball).
    # Rename to the canonical name and install into Homebrew's bin.
    src = "divyastro-mcp-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}"
    bin.install src => "divyastro-mcp"
  end

  def caveats
    <<~EOS
      To use divyastro-mcp with Claude Desktop, add it to your config:

        macOS:   ~/Library/Application Support/Claude/claude_desktop_config.json
        Windows: %APPDATA%\\Claude\\claude_desktop_config.json

        {
          "mcpServers": {
            "divyastro": {
              "command": "#{HOMEBREW_PREFIX}/bin/divyastro-mcp",
              "env": { "DIVYASTRO_API_KEY": "dv_live_..." }
            }
          }
        }

      Get an API key at https://divyastroapi.com/dashboard/keys
      Cursor users: edit ~/.cursor/mcp.json with the same shape.
    EOS
  end

  test do
    # Smoke test — running without an API key must fail fast with the
    # documented error message. This catches build regressions where
    # the auth gate is silently skipped.
    output = shell_output("#{bin}/divyastro-mcp 2>&1", 1)
    assert_match "DIVYASTRO_API_KEY", output
  end
end
