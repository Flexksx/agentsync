# ponte

> **ponte** — Portuguese for *bridge*. Pronounced **pon-chee** (`/ˈpõ.tʃi/`).

Sync AI agent instructions, skills, and subagents across vendors —
Claude Code, Codex, Gemini CLI, Cursor — from a single config.

## What it is

Every AI coding agent keeps its own config in its own place: Claude Code
reads `~/.claude/`, Codex reads `~/.codex/`, Gemini CLI and Cursor each
have their own dotfiles. The same system prompt, the same skills, the
same subagent definitions end up copy-pasted and drifting across four
trees.

ponte is the bridge between them. You declare your system prompt,
skills, and subagents **once** in `~/.config/ponte/`. `ponte sync`
resolves every source, builds an immutable, content-addressed **store
generation**, and activates it by symlinking each vendor's config
directory into that store. One source of truth fans out to every tool.

Because the store is immutable and content-addressed:

- Editing a source file has **no effect on running agents** until the
  next `ponte sync` — no accidental mid-session changes.
- Identical inputs reuse the same generation (same hash, no rebuild).
- Agents cannot modify their own configuration — the store is read-only.
- Every past generation stays addressable until you `ponte gc` it.

## Install

```sh
# Nix
nix profile install github:flexksx/ponte

# Go
go install github.com/flexksx/ponte/apps/ponte@latest
```

A [home-manager module](#nix--home-manager-reference) is also available.

## Showcase

```sh
# First run bootstraps ~/.config/ponte/config.toml and an empty AGENTS.md
ponte sync

# Point ponte at your system prompt, then activate it everywhere
ponte sysprompt set ~/prompts/my-prompt.md
ponte sync

# See where every vendor stands — active generation and drift
ponte status

# List what's declared
ponte skills
ponte subagents

# Sync to a single vendor only
ponte sync -a claude-code

# Preview a change without touching the store or any vendor
ponte sync --dry-run

# Reclaim store generations no vendor points to anymore
ponte gc

# Read the full manual
ponte manual
```

To add a skill, declare it in `~/.config/ponte/config.toml` and run
`ponte sync` — it lands at `~/.claude/skills/<name>`,
`~/.codex/skills/<name>`, and every other enabled vendor at once.

```toml
[[skills]]
name = "my-skill"
[skills.source]
type = "local"
path = "skills/my-skill"
```

See [MANUAL.md](MANUAL.md) for the full CLI reference and usage guide.

## Configuration reference

ponte is configured two ways: directly via `config.toml`, or
declaratively via the home-manager module (which generates that same
`config.toml`). Both expose the same shape.

### TOML reference

All configuration lives in `~/.config/ponte/config.toml`. The first
`ponte sync` creates it with defaults.

```toml
# Path to the system prompt file. A bare filename resolves relative to
# ~/.config/ponte/; an absolute path is read as-is, letting an external
# repo own the prompt. Defaults to AGENTS.md when omitted.
system_prompt_file = "AGENTS.md"

# Per-vendor toggles. Omitting a vendor defaults to enabled = true.
[agents.claude-code]
enabled = true

[agents.codex]
enabled = true

[agents.gemini-cli]
enabled = true

[agents.cursor-agent]
enabled = false

# Skills — zero or more entries. Each is a directory containing a
# SKILL.md plus supporting files, synced to every enabled vendor.

[[skills]]
name = "software-engineering"
[skills.source]
type = "local"
path = "skills/software-engineering"   # relative to ~/.config/ponte/

[[skills]]
name = "ast-grep"
[skills.source]
type = "git"
url    = "https://github.com/example/ast-grep-skill"
ref    = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"   # full SHA recommended
subdir = ""   # optional: subdirectory inside the repo holding the skill

# Subagents — zero or more entries. Each source resolves to a directory
# of agent definition files, flattened into every enabled vendor's
# agents directory. Same source schema as skills.

[[subagents]]
name = "claude"
[subagents.source]
type = "local"
path = "subagents/claude"   # relative to ~/.config/ponte/
```

| Key | Type | Default | Meaning |
|-----|------|---------|---------|
| `system_prompt_file` | string | `AGENTS.md` | System prompt path. Bare name → relative to `~/.config/ponte/`; absolute → read as-is. |
| `[agents.<vendor>].enabled` | bool | `true` | Whether sync targets that vendor. Vendors: `claude-code`, `codex`, `gemini-cli`, `cursor-agent`. |
| `[[skills]].name` | string | — | Skill name; becomes the synced directory name. |
| `[[skills]].source` | table | — | Where the skill is fetched from (see below). |
| `[[subagents]].name` | string | — | Subagent group name. |
| `[[subagents]].source` | table | — | Directory of agent files; same schema as a skill source. |

**Source tables** (`[skills.source]` / `[subagents.source]`):

| Field | Used by | Meaning |
|-------|---------|---------|
| `type` | both | `"local"` or `"git"`. |
| `path` | `local` | Directory path. Relative resolves against `~/.config/ponte/`; absolute used as-is. |
| `url` | `git` | Remote URL. Cloned into `~/.cache/ponte/sources/`. |
| `ref` | `git` | Branch, tag, or commit. **Prefer full commit SHAs** — a moving ref changes the store hash and forces rebuilds. |
| `subdir` | `git` | Optional subdirectory inside the repo. Omit for the repo root. |

### Nix / home-manager reference

The module installs the binary and generates `config.toml`
declaratively. The system prompt file (`AGENTS.md`) is intentionally
left **unmanaged** so `ponte sysprompt set` keeps working, and `ponte
sync` is **never run automatically** — run it yourself after a rebuild.

```nix
# flake.nix
{
  inputs.ponte.url = "github:flexksx/ponte";
}
```

```nix
# home configuration
{ inputs, ... }: {
  imports = [ inputs.ponte.homeManagerModules.ponte ];

  programs.ponte = {
    enable = true;

    # Read the system prompt from an absolute path (e.g. a config repo)
    # instead of ~/.config/ponte/AGENTS.md. A bare filename stays
    # relative to that dir.
    systemPromptFile = "/home/me/config/ai_agents/AGENTS.md";

    # Toggle individual vendors; unset vendors default to enabled.
    agents."gemini-cli".enable = false;

    skills = [
      {
        name = "my-skill";
        source = {
          type = "git";
          url = "https://github.com/me/skills";
          ref = "main";
          subdir = "my-skill";
        };
      }
    ];

    # Subagents: each source resolves to a directory of agent files that
    # are flattened into every enabled vendor's agents directory on sync.
    subagents = [
      {
        name = "claude";
        source = {
          type = "local";
          path = "/home/me/config/ai_agents/subagents/claude";
        };
      }
    ];
  };
}
```

| Option | Type | Default | Meaning |
|--------|------|---------|---------|
| `enable` | bool | `false` | Install ponte and generate `config.toml`. |
| `package` | package | flake's default | The ponte package to install. |
| `systemPromptFile` | string | `"AGENTS.md"` | Maps to `system_prompt_file`. Bare name → relative to `~/.config/ponte/`; absolute → read as-is. |
| `agents.<vendor>.enable` | bool | `true` | Per-vendor toggle. Vendors: `claude-code`, `codex`, `gemini-cli`, `cursor-agent`. |
| `skills` | list of `{ name; source; }` | `[]` | Skill declarations; `source` matches the TOML source schema. |
| `subagents` | list of `{ name; source; }` | `[]` | Subagent declarations; same `source` schema as skills. |
| `settings` | TOML attrset | `{}` | Escape hatch for keys the module doesn't model; merged into `config.toml` and **takes precedence** over generated values. |

`source` submodule fields: `type` (`"local"`/`"git"`), `path`, `url`,
`ref`, `subdir` — same meanings as the [TOML source table](#toml-reference).

## Documentation

See [MANUAL.md](MANUAL.md) for the full configuration reference and usage guide.
