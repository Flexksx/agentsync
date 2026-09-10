import {
  buildVendorLayouts,
  buildVendorPlan,
  type Config,
  type ListFiles,
  type Platform,
  parseSource,
  type ResolvedEntry,
  type ResolveSource,
  type SourceEntry,
  type VendorName,
  type VendorPlan,
} from "@ponte/core";

export type BuildVendorPlans = (
  config: Config,
  promptPath: string,
) => Promise<Record<VendorName, VendorPlan>>;

export const createBuildVendorPlans =
  (
    resolveSource: ResolveSource,
    listFiles: ListFiles,
    home: string,
    platform: Platform,
  ): BuildVendorPlans =>
  async (config, promptPath) => {
    const resolveDirectory = (entry: SourceEntry): Promise<string> =>
      resolveSource(parseSource(entry.source, entry.ref, entry.subdir));

    const skills = await Promise.all(
      Object.entries(config.skills).map(
        async ([name, entry]): Promise<ResolvedEntry> => ({
          name,
          sourceDirectory: await resolveDirectory(entry),
          files: [],
        }),
      ),
    );

    const subagents = await Promise.all(
      Object.entries(config.subagents).map(async ([name, entry]) => {
        const sourceDirectory = await resolveDirectory(entry);
        return {
          name,
          sourceDirectory,
          files: await listFiles(sourceDirectory),
        };
      }),
    );

    const layouts = buildVendorLayouts(home, platform);
    const plans = {} as Record<VendorName, VendorPlan>;
    for (const [name, layout] of Object.entries(layouts)) {
      plans[name as VendorName] = buildVendorPlan(
        layout,
        promptPath,
        skills,
        subagents,
      );
    }
    return plans;
  };
