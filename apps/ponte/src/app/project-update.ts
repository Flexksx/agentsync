import {
  err,
  isGitSource,
  type LockEntry,
  ok,
  type ProjectLayout,
  type ProjectLock,
  parseSource,
  type Result,
  type SourceEntry,
  vendoredSkillPath,
} from "@ponte/core";
import { directoriesDiffer, directoryExists, removeDirectory } from "../infra/filesystem";
import { resolveSource } from "../infra/git";
import { gitCacheDirectoryPath } from "../infra/paths";
import { readProjectLock, writeProjectLock } from "../infra/project-file";
import { copyVendorSkill, type Project } from "./project";

export type UpdatedSkill = { readonly name: string; readonly commit: string | null };

export type ProjectUpdateReport = {
  readonly root: string;
  readonly updated: readonly UpdatedSkill[];
};

type Target = readonly [string, SourceEntry];

const namedTarget = (project: Project, name: string): Result<Target, string> => {
  const entry = project.config.skills[name];
  if (entry === undefined) return err(`unknown project skill: ${name}`);
  if (!isGitSource(entry.source)) {
    return err(`${name} is a local skill, so there is nothing to update`);
  }
  return ok([name, entry]);
};

const updateTargets = (
  project: Project,
  name: string | undefined,
): Result<readonly Target[], string> => {
  if (name === undefined) {
    return ok(
      Object.entries(project.config.skills).filter(([, entry]) => isGitSource(entry.source)),
    );
  }
  const result = namedTarget(project, name);
  if (!result.ok) return result;
  return ok([result.value]);
};

const isDirty = async (
  layout: ProjectLayout,
  lock: ProjectLock,
  [name, entry]: Target,
): Promise<boolean> => {
  const directory = vendoredSkillPath(layout, name);
  if (!(await directoryExists(directory))) return false;
  const commit = lock.skills[name]?.commit;
  if (commit === undefined) return true;
  const pristine = await resolveSource(
    parseSource(entry.source, commit, entry.subdir),
    gitCacheDirectoryPath(),
  );
  return directoriesDiffer(pristine, directory);
};

const dirtyTargets = async (
  layout: ProjectLayout,
  lock: ProjectLock,
  targets: readonly Target[],
): Promise<string[]> => {
  const dirty: string[] = [];
  for (const target of targets) {
    if (await isDirty(layout, lock, target)) dirty.push(target[0]);
  }
  return dirty;
};

export const runProjectUpdate = async (
  project: Project,
  name: string | undefined,
  force: boolean,
): Promise<Result<ProjectUpdateReport, string>> => {
  const targetsResult = updateTargets(project, name);
  if (!targetsResult.ok) return targetsResult;
  const targets = targetsResult.value;

  const lock = await readProjectLock(project.layout);
  if (!force) {
    const dirty = await dirtyTargets(project.layout, lock, targets);
    if (dirty.length > 0) {
      return err(
        `${dirty.join(", ")}: the vendored copy differs from its locked commit, or the lock entry is missing - commit the copy, or pass --force to overwrite it`,
      );
    }
  }
  const locked: Record<string, LockEntry> = { ...lock.skills };
  const updated: UpdatedSkill[] = [];
  for (const [skill, entry] of targets) {
    await removeDirectory(vendoredSkillPath(project.layout, skill));
    const commit = await copyVendorSkill(project.layout, skill, entry);
    if (commit !== null) locked[skill] = { commit };
    updated.push({ name: skill, commit });
  }
  await writeProjectLock(project.layout, { skills: locked });
  return ok({ root: project.layout.root, updated });
};
