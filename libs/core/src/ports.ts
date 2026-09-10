import type { Config } from "./domain/config";
import type { VendorPlan } from "./domain/plan";
import type { Platform } from "./domain/platform";
import type { ProjectConfig, ProjectLayout, ProjectLock } from "./domain/project";
import type { SkillSource } from "./domain/source";

export type ResolvedSource = {
  readonly directory: string;
  readonly commit: string | null;
};

export type Filesystem = {
  readonly fileExists: (path: string) => Promise<boolean>;
  readonly directoryExists: (path: string) => Promise<boolean>;
  readonly writeText: (path: string, content: string) => Promise<void>;
  readonly listFiles: (directory: string) => Promise<string[]>;
  readonly removeDirectory: (path: string) => Promise<void>;
  readonly copyDirectoryWithoutGit: (from: string, to: string) => Promise<void>;
  readonly directoriesDiffer: (left: string, right: string) => Promise<boolean>;
};

export type ConfigRepository = {
  readonly readConfig: () => Promise<Config | null>;
  readonly writeConfig: (config: Config) => Promise<void>;
  readonly readPrompt: (filename: string) => Promise<string | null>;
  readonly writePrompt: (filename: string, content: string) => Promise<void>;
  readonly resolveContent: (fileOrLiteral: string) => Promise<string>;
};

export type ProjectRepository = {
  readonly findProjectRoot: (start: string) => Promise<string | null>;
  readonly readProjectConfig: (root: string) => Promise<ProjectConfig>;
  readonly readProjectLock: (layout: ProjectLayout) => Promise<ProjectLock>;
  readonly writeProjectLock: (layout: ProjectLayout, lock: ProjectLock) => Promise<void>;
};

export type SourceResolver = {
  readonly resolve: (source: SkillSource, cacheDirectory: string) => Promise<string>;
  readonly resolveDetails: (source: SkillSource, cacheDirectory: string) => Promise<ResolvedSource>;
};

export type LinkManager = {
  readonly readSymlinks: (plan: VendorPlan) => Promise<Map<string, string>>;
  readonly applyPlan: (plan: VendorPlan, stale: readonly string[]) => Promise<void>;
};

export type Environment = {
  readonly platform: () => Platform;
  readonly cwd: () => string;
  readonly home: () => string;
  readonly configDirectory: () => string;
  readonly dataDirectory: () => string;
  readonly overridePromptPath: () => string;
  readonly gitCacheDirectory: () => string;
  readonly configFile: () => string;
  readonly promptFile: (filename: string) => string;
};
