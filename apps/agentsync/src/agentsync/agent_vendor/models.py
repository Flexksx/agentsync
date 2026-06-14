from dataclasses import dataclass
from enum import Enum
from pathlib import Path


class AgentVendorName(Enum):
    CLAUDE_CODE = "claude-code"
    CODEX = "codex"
    GEMINI_CLI = "gemini-cli"
    CURSOR_AGENT = "cursor-agent"


@dataclass
class AgentVendorConfiguration:
    vendor_name: AgentVendorName
    package_name: str
    global_instruction_file_path: Path
    skills_directory_path: Path
    subagents_directory_path: Path
