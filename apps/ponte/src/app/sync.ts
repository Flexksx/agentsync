import {
  type Config,
  createDefaultConfig,
  err,
  getEnabledVendors,
  getStaleLinkPaths,
  ok,
  parseVendorNames,
  type Result,
  type VendorName,
  type VendorPlan,
} from "@ponte/core";
import { readConfig, writeConfig, writePrompt } from "../infra/config-file";
import { fileExists, writeText } from "../infra/filesystem";
import { applyPlan, readSymlinks } from "../infra/links";
import {
  configDirectoryPath,
  overridePromptPath,
  promptFilePath,
} from "../infra/paths";
import { buildVendorPlans } from "./resolve";

export type SyncRequest = {
  readonly promptOverride: string | undefined;
  readonly requestedVendors: readonly string[];
};

export type Bootstrap = {
  readonly configDirectory: string;
  readonly systemPromptFile: string;
};

export type SyncReport = {
  readonly vendors: readonly VendorName[];
  readonly stale: number;
  readonly bootstrap: Bootstrap | null;
};

type PendingSync = {
  readonly vendors: readonly VendorName[];
  readonly plans: Readonly<Record<VendorName, VendorPlan>>;
  readonly stale: Readonly<Record<string, readonly string[]>>;
  readonly bootstrap: Bootstrap | null;
};

export const findConfig = (): Promise<Config | null> => readConfig();

const bootstrapConfig = async (): Promise<{
  config: Config;
  bootstrap: Bootstrap;
}> => {
  const config = createDefaultConfig();
  await writeConfig(config);
  await writePrompt(config.systemPromptFile, "");
  return {
    config,
    bootstrap: {
      configDirectory: configDirectoryPath(),
      systemPromptFile: config.systemPromptFile,
    },
  };
};

const configuredPromptPath = async (
  config: Config,
): Promise<Result<string, string>> => {
  const path = promptFilePath(config.systemPromptFile);
  if (!(await fileExists(path))) {
    return err(`system prompt not found: ${config.systemPromptFile}`);
  }
  return ok(path);
};

const materializedOverridePath = async (override: string): Promise<string> => {
  if (await fileExists(override)) return override;
  const path = overridePromptPath();
  await writeText(path, override);
  return path;
};

const pendingSync = async (
  request: SyncRequest,
): Promise<Result<PendingSync, string>> => {
  const existing = await readConfig();
  const { config, bootstrap } =
    existing === null
      ? await bootstrapConfig()
      : { config: existing, bootstrap: null };

  const vendors =
    request.requestedVendors.length > 0
      ? parseVendorNames(request.requestedVendors)
      : getEnabledVendors(config);
  if (vendors.length === 0) {
    return err("no agents enabled in config - run with -a to specify agents");
  }

  const promptResult =
    request.promptOverride === undefined
      ? await configuredPromptPath(config)
      : ok(await materializedOverridePath(request.promptOverride));
  if (!promptResult.ok) return promptResult;

  const plans = await buildVendorPlans(config, promptResult.value);
  const stale: Record<string, readonly string[]> = {};
  for (const vendor of vendors) {
    stale[vendor] = getStaleLinkPaths(
      plans[vendor],
      await readSymlinks(plans[vendor]),
    );
  }
  return ok({ vendors, plans, stale, bootstrap });
};

const countStale = (pending: PendingSync): number =>
  pending.vendors.reduce(
    (total, vendor) => total + (pending.stale[vendor]?.length ?? 0),
    0,
  );

export const syncVendors = async (
  request: SyncRequest,
  apply: boolean,
): Promise<Result<SyncReport, string>> => {
  const result = await pendingSync(request);
  if (!result.ok) return result;
  const pending = result.value;
  if (apply) {
    for (const vendor of pending.vendors) {
      await applyPlan(pending.plans[vendor], pending.stale[vendor] ?? []);
    }
  }
  return ok({
    vendors: pending.vendors,
    stale: countStale(pending),
    bootstrap: pending.bootstrap,
  });
};
