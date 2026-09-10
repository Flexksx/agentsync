import {
  buildProjectPlan,
  getProjectEnabledVendors,
  getVendorState,
  isGitSource,
  type LockEntry,
  type ProjectConfig,
  type ProjectLayout,
  type ProjectLock,
  parseSource,
  projectLayout,
  type SourceEntry,
  type VendorPlan,
  type VendorState,
  vendoredSkillPath,
} from "@ponte/core";
import { copyDirectoryWithoutGit, directoryExists } from "../infra/filesystem";
import { resolveSource, resolveSourceDetails } from "../infra/git";
import { readSymlinks } from "../infra/links";
import {
  currentDirectory,
  currentPlatform,
  gitCacheDirectoryPath,
} from "../infra/paths";
import {
  findProjectRoot,
  readProjectConfig,
  readProjectLock,
} from "../infra/project-file";

export type Project = {
  readonly layout: ProjectLayout;
  readonly config: ProjectConfig;
};

export type ProjectSkill = {
  readonly name: string;
  readonly directory: string;
  readonly vendored: boolean;
  readonly commit: string | null;
};

export type ProjectResolution = {
  readonly skills: readonly ProjectSkill[];
  readonly plan: VendorPlan;
  readonly lock: ProjectLock;
  readonly vendored: readonly string[];
};

export type ProjectSkillRow = {
  readonly name: string;
  readonly entry: SourceEntry;
  readonly vendored: boolean;
  readonly commit: string | null;
};

export type ProjectStatusReport = {
  readonly root: string;
  readonly skillsDirectory: string;
  readonly linkCount: number;
  readonly state: VendorState;
};

export const findProject = async (): Promise<Project | null> => {
  const root = await findProjectRoot(currentDirectory());
  if (root === null) return null;
  const config = await readProjectConfig(root);
  const enabled = getProjectEnabledVendors(config);
  return { layout: projectLayout(root, currentPlatform(), enabled), config };
};

export const copyVendorSkill = async (
  layout: ProjectLayout,
  name: string,
  entry: SourceEntry,
): Promise<string | null> => {
  const resolved = await resolveSourceDetails(
    parseSource(entry.source, entry.ref, entry.subdir),
    gitCacheDirectoryPath(),
  );
  await copyDirectoryWithoutGit(
    resolved.directory,
    vendoredSkillPath(layout, name),
  );
  return resolved.commit;
};

const localSkillDirectory = (entry: SourceEntry): Promise<string> =>
  resolveSource(
    parseSource(entry.source, entry.ref, entry.subdir),
    gitCacheDirectoryPath(),
  );

export const resolveProjectSkills = async (
  project: Project,
  materialize: boolean,
): Promise<ProjectResolution> => {
  const locked: Record<string, LockEntry> = {
    ...(await readProjectLock(project.layout)).skills,
  };
  const skills: ProjectSkill[] = [];
  const vendored: string[] = [];
  for (const [name, entry] of Object.entries(project.config.skills)) {
    if (!isGitSource(entry.source)) {
      const directory = await localSkillDirectory(entry);
      skills.push({ name, directory, vendored: false, commit: null });
      continue;
    }
    const directory = vendoredSkillPath(project.layout, name);
    if (!(await directoryExists(directory))) {
      vendored.push(name);
      if (materialize) {
        const commit = await copyVendorSkill(project.layout, name, entry);
        if (commit !== null) locked[name] = { commit };
      }
    }
    skills.push({
      name,
      directory,
      vendored: true,
      commit: locked[name]?.commit ?? null,
    });
  }
  return {
    skills,
    plan: buildProjectPlan(project.layout, skills),
    lock: { skills: locked },
    vendored,
  };
};

export const listProjectSkills = async (
  project: Project,
): Promise<ProjectSkillRow[]> => {
  const lock = await readProjectLock(project.layout);
  return Object.entries(project.config.skills).map(([name, entry]) => ({
    name,
    entry,
    vendored: isGitSource(entry.source),
    commit: lock.skills[name]?.commit ?? null,
  }));
};

export const getProjectStatusReport = async (
  project: Project,
): Promise<ProjectStatusReport> => {
  const { plan } = await resolveProjectSkills(project, false);
  const actual = await readSymlinks(plan);
  return {
    root: project.layout.root,
    skillsDirectory: project.layout.skills,
    linkCount: actual.size,
    state: getVendorState(plan, actual),
  };
};
