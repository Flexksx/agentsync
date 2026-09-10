import {
  buildProjectPlan,
  type CopyDirectoryWithoutGit,
  type DirectoryExists,
  type FindProjectRoot,
  getProjectEnabledVendors,
  getVendorState,
  isGitSource,
  type LockEntry,
  type Platform,
  type ProjectConfig,
  type ProjectLayout,
  type ProjectLock,
  parseSource,
  projectLayout,
  type ReadProjectConfig,
  type ReadProjectLock,
  type ReadSymlinks,
  type ResolveSource,
  type ResolveSourceDetails,
  type SourceEntry,
  type VendorPlan,
  type VendorState,
  vendoredSkillPath,
} from "@ponte/core";

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

export type CopyVendorSkill = (
  layout: ProjectLayout,
  name: string,
  entry: SourceEntry,
) => Promise<string | null>;

export type ResolveProjectSkills = (
  project: Project,
  materialize: boolean,
) => Promise<ProjectResolution>;

export type FindProject = () => Promise<Project | null>;

export type ListProjectSkills = (
  project: Project,
) => Promise<ProjectSkillRow[]>;

export type GetProjectStatusReport = (
  project: Project,
) => Promise<ProjectStatusReport>;

export const createFindProject =
  (
    findProjectRoot: FindProjectRoot,
    readProjectConfig: ReadProjectConfig,
    cwd: string,
    platform: Platform,
  ): FindProject =>
  async () => {
    const root = await findProjectRoot(cwd);
    if (root === null) return null;
    const config = await readProjectConfig(root);
    const enabled = getProjectEnabledVendors(config);
    return { layout: projectLayout(root, platform, enabled), config };
  };

export const createCopyVendorSkill =
  (
    resolveSourceDetails: ResolveSourceDetails,
    copyDirectoryWithoutGit: CopyDirectoryWithoutGit,
  ): CopyVendorSkill =>
  async (layout, name, entry) => {
    const resolved = await resolveSourceDetails(
      parseSource(entry.source, entry.ref, entry.subdir),
    );
    await copyDirectoryWithoutGit(
      resolved.directory,
      vendoredSkillPath(layout, name),
    );
    return resolved.commit;
  };

export const createResolveProjectSkills =
  (
    readProjectLock: ReadProjectLock,
    resolveSource: ResolveSource,
    directoryExists: DirectoryExists,
    copyVendorSkill: CopyVendorSkill,
  ): ResolveProjectSkills =>
  async (project, materialize) => {
    const locked: Record<string, LockEntry> = {
      ...(await readProjectLock(project.layout)).skills,
    };
    const skills: ProjectSkill[] = [];
    const vendored: string[] = [];
    for (const [name, entry] of Object.entries(project.config.skills)) {
      if (!isGitSource(entry.source)) {
        const directory = await resolveSource(
          parseSource(entry.source, entry.ref, entry.subdir),
        );
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

export const createListProjectSkills =
  (readProjectLock: ReadProjectLock): ListProjectSkills =>
  async project => {
    const lock = await readProjectLock(project.layout);
    return Object.entries(project.config.skills).map(([name, entry]) => ({
      name,
      entry,
      vendored: isGitSource(entry.source),
      commit: lock.skills[name]?.commit ?? null,
    }));
  };

export const createGetProjectStatusReport =
  (
    resolveProjectSkills: ResolveProjectSkills,
    readSymlinks: ReadSymlinks,
  ): GetProjectStatusReport =>
  async project => {
    const { plan } = await resolveProjectSkills(project, false);
    const actual = await readSymlinks(plan);
    return {
      root: project.layout.root,
      skillsDirectory: project.layout.skills,
      linkCount: actual.size,
      state: getVendorState(plan, actual),
    };
  };
