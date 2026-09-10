import { getStaleLinkPaths } from "@ponte/core";
import { applyPlan, readSymlinks } from "../infra/links";
import { writeProjectLock } from "../infra/project-file";
import { type Project, resolveProjectSkills } from "./project";

export type ProjectSyncReport = {
  readonly root: string;
  readonly vendored: number;
  readonly linked: number;
  readonly stale: number;
};

export const syncProject = async (
  project: Project,
  apply: boolean,
): Promise<ProjectSyncReport> => {
  const resolution = await resolveProjectSkills(project, apply);
  const stale = getStaleLinkPaths(
    resolution.plan,
    await readSymlinks(resolution.plan),
  );
  if (apply) {
    await applyPlan(resolution.plan, stale);
    if (resolution.vendored.length > 0)
      await writeProjectLock(project.layout, resolution.lock);
  }
  return {
    root: project.layout.root,
    vendored: resolution.vendored.length,
    linked: resolution.plan.links.length,
    stale: stale.length,
  };
};
