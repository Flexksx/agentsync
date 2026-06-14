import shutil

from agentsync.agent_vendor.models import AgentVendorName
from agentsync.agent_vendor.ports.configuration import AgentVendorConfigurationPort


def is_installed(
    vendor_name: AgentVendorName, configuration_port: AgentVendorConfigurationPort
) -> bool:
    configuration = configuration_port.get_configuration(vendor_name)
    return shutil.which(configuration.package_name) is not None
