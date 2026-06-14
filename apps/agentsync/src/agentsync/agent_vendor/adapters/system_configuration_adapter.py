import platform
from pathlib import Path

from agentsync.agent_vendor.models import AgentVendorConfiguration, AgentVendorName
from agentsync.agent_vendor.ports.configuration import (
    AgentVendorConfigurationPort,
    VendorConfigurationNotFoundError,
)


class UnsupportedPlatformError(Exception):
    pass


_POSIX_CONFIGURATIONS: dict[AgentVendorName, AgentVendorConfiguration] = {
    AgentVendorName.CLAUDE_CODE: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CLAUDE_CODE,
        package_name="claude",
        global_instruction_file_path=str(Path.home() / ".claude" / "CLAUDE.md"),
        skills_directory_path=str(Path.home() / ".claude" / "skills"),
        subagents_directory_path=str(Path.home() / ".claude" / "subagents"),
    ),
    AgentVendorName.CODEX: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CODEX,
        package_name="codex",
        global_instruction_file_path=str(Path.home() / ".codex" / "instructions.md"),
        skills_directory_path=str(Path.home() / ".codex" / "skills"),
        subagents_directory_path=str(Path.home() / ".codex" / "subagents"),
    ),
    AgentVendorName.GEMINI_CLI: AgentVendorConfiguration(
        vendor_name=AgentVendorName.GEMINI_CLI,
        package_name="gemini",
        global_instruction_file_path=str(Path.home() / ".gemini" / "GEMINI.md"),
        skills_directory_path=str(Path.home() / ".gemini" / "skills"),
        subagents_directory_path=str(Path.home() / ".gemini" / "subagents"),
    ),
    AgentVendorName.CURSOR_AGENT: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CURSOR_AGENT,
        package_name="cursor",
        global_instruction_file_path=str(Path.home() / ".cursor" / "rules" / "global.mdc"),
        skills_directory_path=str(Path.home() / ".cursor" / "skills"),
        subagents_directory_path=str(Path.home() / ".cursor" / "subagents"),
    ),
}

_WINDOWS_CONFIGURATIONS: dict[AgentVendorName, AgentVendorConfiguration] = {
    AgentVendorName.CLAUDE_CODE: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CLAUDE_CODE,
        package_name="claude",
        global_instruction_file_path=str(
            Path.home() / "AppData" / "Roaming" / "Claude" / "CLAUDE.md"
        ),
        skills_directory_path=str(Path.home() / "AppData" / "Roaming" / "Claude" / "skills"),
        subagents_directory_path=str(Path.home() / "AppData" / "Roaming" / "Claude" / "subagents"),
    ),
    AgentVendorName.CODEX: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CODEX,
        package_name="codex",
        global_instruction_file_path=str(
            Path.home() / "AppData" / "Roaming" / "Codex" / "instructions.md"
        ),
        skills_directory_path=str(Path.home() / "AppData" / "Roaming" / "Codex" / "skills"),
        subagents_directory_path=str(Path.home() / "AppData" / "Roaming" / "Codex" / "subagents"),
    ),
    AgentVendorName.GEMINI_CLI: AgentVendorConfiguration(
        vendor_name=AgentVendorName.GEMINI_CLI,
        package_name="gemini",
        global_instruction_file_path=str(
            Path.home() / "AppData" / "Roaming" / "Gemini" / "GEMINI.md"
        ),
        skills_directory_path=str(Path.home() / "AppData" / "Roaming" / "Gemini" / "skills"),
        subagents_directory_path=str(Path.home() / "AppData" / "Roaming" / "Gemini" / "subagents"),
    ),
    AgentVendorName.CURSOR_AGENT: AgentVendorConfiguration(
        vendor_name=AgentVendorName.CURSOR_AGENT,
        package_name="cursor",
        global_instruction_file_path=str(
            Path.home() / "AppData" / "Roaming" / "Cursor" / "rules" / "global.mdc"
        ),
        skills_directory_path=str(Path.home() / "AppData" / "Roaming" / "Cursor" / "skills"),
        subagents_directory_path=str(Path.home() / "AppData" / "Roaming" / "Cursor" / "subagents"),
    ),
}


class SystemAgentVendorConfigurationAdapter(AgentVendorConfigurationPort):
    def get_configuration(self, vendor_name: AgentVendorName) -> AgentVendorConfiguration:
        configurations = self._platform_configurations()
        if vendor_name not in configurations:
            raise VendorConfigurationNotFoundError(vendor_name)
        return configurations[vendor_name]

    def _platform_configurations(self) -> dict[AgentVendorName, AgentVendorConfiguration]:
        system = platform.system()
        if system in ("Linux", "Darwin"):
            return _POSIX_CONFIGURATIONS
        if system == "Windows":
            return _WINDOWS_CONFIGURATIONS
        raise UnsupportedPlatformError(system)
