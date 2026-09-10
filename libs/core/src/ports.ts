import type { Config } from "./domain/config";
import type { VendorPlan } from "./domain/plan";
import type { ProjectLayout } from "./domain/project";
import type { ProjectConfig, ProjectLock } from "./domain/project-config";
import type { SkillSource } from "./domain/source";

export type ResolvedSource = {
  readonly directory: string;
  readonly commit: string | null;
};

// Source resolution
export type ResolveSource = (source: SkillSource) => Promise<string>;
export type ResolveSourceDetails = (
  source: SkillSource,
) => Promise<ResolvedSource>;

// Symlink management
export type ReadSymlinks = (plan: VendorPlan) => Promise<Map<string, string>>;
export type ApplyPlan = (
  plan: VendorPlan,
  stale: readonly string[],
) => Promise<void>;

// Global config
export type ReadConfig = () => Promise<Config | null>;
export type WriteConfig = (config: Config) => Promise<void>;
export type ReadPrompt = (filename: string) => Promise<string | null>;
export type WritePrompt = (filename: string, content: string) => Promise<void>;
export type ResolveContent = (fileOrLiteral: string) => Promise<string>;

// Project
export type FindProjectRoot = (start: string) => Promise<string | null>;
export type ReadProjectConfig = (root: string) => Promise<ProjectConfig>;
export type ReadProjectLock = (layout: ProjectLayout) => Promise<ProjectLock>;
export type WriteProjectLock = (
  layout: ProjectLayout,
  lock: ProjectLock,
) => Promise<void>;

// Filesystem
export type FileExists = (path: string) => Promise<boolean>;
export type DirectoryExists = (path: string) => Promise<boolean>;
export type WriteText = (path: string, content: string) => Promise<void>;
export type ListFiles = (directory: string) => Promise<string[]>;
export type RemoveDirectory = (path: string) => Promise<void>;
export type CopyDirectoryWithoutGit = (
  from: string,
  to: string,
) => Promise<void>;
export type DirectoriesDiffer = (
  left: string,
  right: string,
) => Promise<boolean>;
