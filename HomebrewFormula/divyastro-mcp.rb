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
  version "0.3.1"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-darwin-arm64"
      sha256 "0c831042d448251755a2966b4aeaac79e42c33cc5ae8a77d8a7a52ac96fae893"
    end
    on_intel do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-darwin-amd64"
      sha256 "08c70ddc503507f0c17cdbcadba25fbd85e6964c2f86316de9c390bcbdd1d322"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-linux-arm64"
      sha256 "07d35ccc9a4647e98d1fc3fe217f4983527b936a2c82c92c4b989c191608cb2a"
    end
    on_intel do
      url "https://github.com/DivyGuru/divyastro-mcp/releases/download/v#{version}/divyastro-mcp-linux-amd64"
      sha256 "8ef1118a5fbb6b936b444d34dd7412fb384988de47fd98cf4d9b8055f4612de9"
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
