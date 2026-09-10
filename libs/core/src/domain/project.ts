import { join } from "node:path";
import type { Platform } from "./platform";
import { buildVendorLayouts, type VendorName } from "./vendor";

export type ProjectLayout = {
  readonly root: string;
  readonly skills: string;
  readonly vendorSkillDirectories: readonly string[];
  readonly sources: string;
  readonly lockFile: string;
};

export type ProjectSkillTarget = { readonly name: string; readonly directory: string };

export const PROJECT_SOURCES_DIRECTORY = join(".ponte", "sources");

export const PROJECT_SKILLS_DIRECTORY = join(".agents", "skills");

const PROJECT_LOCK_FILE = join(".ponte", "lock.toml");

export const projectLayout = (
  root: string,
  platform: Platform = "posix",
  enabledVendors?: readonly VendorName[],
): ProjectLayout => {
  const layouts = buildVendorLayouts(root, platform);
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
