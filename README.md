# ponte

Sync AI agent instructions and skills across vendors — Claude Code, Codex, Gemini CLI, Cursor — from a single config.

Instead of managing separate dotfiles per tool, declare your system prompt and skills once. `ponte sync` builds an immutable artifact in a content-addressed store and activates it via symlinks, so edits to your source never silently affect a running agent.

## Install

```sh
# Nix
nix profile install github:flexksx/ponte

# Go
go install github.com/flexksx/ponte/apps/ponte@latest
```

## Quick start

```sh
# First run creates ~/.config/ponte/config.toml and AGENTS.md
ponte sync

# Set your system prompt
ponte sysprompt set ~/my-prompt.md

# Declare a skill in config.toml, then sync
ponte sync

# Read the full manual
ponte manual
```

## Documentation

See [MANUAL.md](MANUAL.md) for the full configuration reference and usage guide.
