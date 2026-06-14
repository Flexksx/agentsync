from abc import ABC, abstractmethod

from agentsync.agent_vendor.models import AgentVendorConfiguration, AgentVendorName


class VendorConfigurationNotFoundError(Exception):
    pass


class AgentVendorConfigurationPort(ABC):
    @abstractmethod
    def get_configuration(self, vendor_name: AgentVendorName) -> AgentVendorConfiguration: ...
