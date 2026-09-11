import {
  type Config,
  type ReadConfig,
  type Result,
  resolvePromptPath,
} from "@ponte/core";
import {
  readConfig,
  readPrompt,
  resolveContent,
  writeConfig,
  writePrompt,
} from "../infra/config-file";
import {
  copyDirectoryWithoutGit,
  directoriesDiffer,
  directoryExists,
  fileExists,
  listFiles,
  removeDirectory,
  writeText,
} from "../infra/filesystem";
import {
  resolveSource as gitResolveSource,
  resolveSourceDetails as gitResolveSourceDetails,
} from "../infra/git";
import { applyPlan, readSymlinks } from "../infra/links";
import {
  configDirectoryPath,
  currentDirectory,
  currentPlatform,
  gitCacheDirectoryPath,
  homeDirectory,
  overridePromptPath,
} from "../infra/paths";
import {
  findProjectRoot,
  readProjectConfig,
  readProjectLock,
  writeProjectLock,
} from "../infra/project-file";
import type {
  FindProject,
  GetProjectStatusReport,
  ListProjectSkills,
  Project,
} from "./project";
import {
  createCopyVendorSkill,
  createFindProject,
  createGetProjectStatusReport,
  createListProjectSkills,
  createResolveProjectSkills,
} from "./project";
import type { ProjectSyncReport } from "./project-sync";
import { createSyncProject } from "./project-sync";
import type { ProjectUpdateReport } from "./project-update";
import { createRunProjectUpdate } from "./project-update";
import { createBuildVendorPlans } from "./resolve";
import type { StatusReport } from "./status";
import { createGetStatusReport } from "./status";
import type { SyncReport, SyncRequest } from "./sync";
import { createSyncVendors } from "./sync";
import { createReadSystemPrompt, createSetSystemPrompt } from "./sysprompt";

export type App = {
  findProject: FindProject;
  readConfig: ReadConfig;
  readSystemPrompt: (config: Config) => Promise<string | null>;
  setSystemPrompt: (config: Config, fileOrLiteral: string) => Promise<void>;
  listProjectSkills: ListProjectSkills;
  getProjectStatusReport: GetProjectStatusReport;
  syncProject: (project: Project, apply: boolean) => Promise<ProjectSyncReport>;
  getStatusReport: (config: Config) => Promise<Result<StatusReport, string>>;
  syncVendors: (
    request: SyncRequest,
    apply: boolean,
  ) => Promise<Result<SyncReport, string>>;
  runProjectUpdate: (
    project: Project,
    name: string | undefined,
    force: boolean,
  ) => Promise<Result<ProjectUpdateReport, string>>;
};

export const bootstrap = (): App => {
  const platform = currentPlatform();
  const cwd = currentDirectory();
  const home = homeDirectory();
  const configDir = configDirectoryPath();
  const overridePath = overridePromptPath();
  const cacheDir = gitCacheDirectoryPath();

  const resolveSource = (source: Parameters<typeof gitResolveSource>[0]) =>
    gitResolveSource(source, cacheDir);
  const resolveSourceDetails = (
    source: Parameters<typeof gitResolveSourceDetails>[0],
  ) => gitResolveSourceDetails(source, cacheDir);
  const promptPath = (filename: string) =>
    resolvePromptPath(configDir, filename);

  const buildVendorPlans = createBuildVendorPlans(
    resolveSource,
    listFiles,
    home,
    platform,
  );
  const copyVendorSkill = createCopyVendorSkill(
    resolveSourceDetails,
    copyDirectoryWithoutGit,
  );
  const resolveProjectSkills = createResolveProjectSkills(
    readProjectLock,
    resolveSource,
    directoryExists,
    copyVendorSkill,
  );

  return {
    findProject: createFindProject(
      findProjectRoot,
      readProjectConfig,
      cwd,
      platform,
    ),
    readConfig,
    readSystemPrompt: createReadSystemPrompt(readPrompt),
    setSystemPrompt: createSetSystemPrompt(writePrompt, resolveContent),
    listProjectSkills: createListProjectSkills(readProjectLock),
    getProjectStatusReport: createGetProjectStatusReport(
      resolveProjectSkills,
      readSymlinks,
    ),
    syncProject: createSyncProject(
      resolveProjectSkills,
      readSymlinks,
      applyPlan,
      writeProjectLock,
    ),
    getStatusReport: createGetStatusReport(
      fileExists,
      readSymlinks,
      buildVendorPlans,
      promptPath,
    ),
    syncVendors: createSyncVendors({
      readConfig,
      writeConfig,
      writePrompt,
      fileExists,
      writeText,
      readSymlinks,
      applyPlan,
      buildVendorPlans,
      configDirectory: configDir,
      overridePromptPath: overridePath,
      resolvePromptPath: promptPath,
    }),
    runProjectUpdate: createRunProjectUpdate({
      readProjectLock,
      writeProjectLock,
      resolveSource,
      copyVendorSkill,
      directoryExists,
      removeDirectory,
      directoriesDiffer,
    }),
  };
};
