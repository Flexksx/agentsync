import type { SkillSource } from "./source";

export type ResolvedSource = {
  readonly directory: string;
  readonly commit: string | null;
};

export type ResolveSource = (source: SkillSource) => Promise<string>;
export type ResolveSourceDetails = (
  source: SkillSource,
) => Promise<ResolvedSource>;
