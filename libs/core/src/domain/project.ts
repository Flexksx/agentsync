import { join } from "node:path";
import { absoluteSources, type SourceEntry, type VendorConfig } from "./config";
import type { Platform } from "./platform";
import { VENDORS, type VendorName, vendorLayouts } from "./vendor";

export type ProjectConfig = {
  readonly vendors?: Readonly<Partial<Record<VendorName, VendorConfig>>>;
  readonly skills: Readonly<Record<string, SourceEntry>>;
};

export type LockEntry = { readonly commit: string };

export type ProjectLock = { readonly skills: Readonly<Record<string, LockEntry>> };

export type ProjectLayout = {
  readonly root: string;
  readonly skills: string;
  readonly vendorSkillDirectories: readonly string[];
  readonly sources: string;
  readonly lockFile: string;
};

export type ProjectSkillTarget = { readonly name: string; readonly directory: string };

export const PROJECT_CONFIG_FILE = "ponte.toml";

export const PROJECT_SOURCES_DIRECTORY = join(".ponte", "sources");

export const PROJECT_SKILLS_DIRECTORY = join(".agents", "skills");

const PROJECT_LOCK_FILE = join(".ponte", "lock.toml");

export const projectEnabledVendors = (config: ProjectConfig): VendorName[] | undefined => {
  if (config.vendors === undefined) return undefined;
  return VENDORS.filter(name => config.vendors?.[name]?.enabled === true);
};

export const projectLayout = (
  root: string,
  platform: Platform = "posix",
  enabledVendors?: readonly VendorName[],
): ProjectLayout => {
  const layouts = vendorLayouts(root, platform);
  const vendorDirs =
    enabledVendors !== undefined
      ? enabledVendors.map(name => layouts[name].skills)
      : Object.values(layouts).map(l => l.skills);
  return {
    root,
    skills: join(root, PROJECT_SKILLS_DIRECTORY),
    vendorSkillDirectories: vendorDirs,
    sources: join(root, PROJECT_SOURCES_DIRECTORY),
    lockFile: join(root, PROJECT_LOCK_FILE),
  };
};

export const vendoredSkillPath = (layout: ProjectLayout, name: string): string =>
  join(layout.sources, name);

export const normalizeProjectConfig = (config: ProjectConfig, root: string): ProjectConfig => ({
  ...(config.vendors !== undefined && { vendors: config.vendors }),
  skills: absoluteSources(config.skills, root),
});
