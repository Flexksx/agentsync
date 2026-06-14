from dataclasses import dataclass
from enum import Enum


class AgentVendorName(Enum):
    CLAUDE_CODE = "claude-code"
    CODEX = "codex"
    GEMINI_CLI = "gemini-cli"
    CURSOR_AGENT = "cursor-agent"


@dataclass
class AgentVendorConfiguration:
    vendor_name: AgentVendorName
    package_name: str
    skills_directory_path: str
    subagents_directory_path: str
