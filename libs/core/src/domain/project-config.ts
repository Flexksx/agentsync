import { resolveSourcePaths, type SourceEntry, type VendorConfig } from "./config";
import { VENDORS, type VendorName } from "./vendor";

export type ProjectConfig = {
  readonly vendors?: Readonly<Partial<Record<VendorName, VendorConfig>>>;
  readonly skills: Readonly<Record<string, SourceEntry>>;
};

export type LockEntry = { readonly commit: string };

export type ProjectLock = { readonly skills: Readonly<Record<string, LockEntry>> };

export const PROJECT_CONFIG_FILE = "ponte.toml";

export const getProjectEnabledVendors = (config: ProjectConfig): VendorName[] | undefined => {
  if (config.vendors === undefined) return undefined;
  return VENDORS.filter(name => config.vendors?.[name]?.enabled === true);
};

export const resolveProjectConfigPaths = (config: ProjectConfig, root: string): ProjectConfig => ({
  ...(config.vendors !== undefined && { vendors: config.vendors }),
  skills: resolveSourcePaths(config.skills, root),
});
