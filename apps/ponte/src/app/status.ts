import {
  type Config,
  err,
  getVendorState,
  ok,
  type Result,
  VENDORS,
  type VendorName,
  type VendorState,
} from "@ponte/core";
import { fileExists } from "../infra/filesystem";
import { readSymlinks } from "../infra/links";
import { promptFilePath } from "../infra/paths";
import { buildVendorPlans } from "./resolve";

export type VendorStatus = {
  readonly name: VendorName;
  readonly enabled: boolean;
  readonly linkCount: number;
  readonly state: VendorState;
};

export type StatusReport = {
  readonly promptFile: string;
  readonly vendors: readonly VendorStatus[];
};

export const getStatusReport = async (config: Config): Promise<Result<StatusReport, string>> => {
  const promptPath = promptFilePath(config.systemPromptFile);
  if (!(await fileExists(promptPath))) {
    return err(`system prompt not found: ${config.systemPromptFile}`);
  }

  const plans = await buildVendorPlans(config, promptPath);
  const vendors: VendorStatus[] = [];
  for (const name of VENDORS) {
    const enabled = config.vendors[name]?.enabled === true;
    const actual = await readSymlinks(plans[name]);
    vendors.push({
      name,
      enabled,
      linkCount: actual.size,
      state: enabled ? getVendorState(plans[name], actual) : "disabled",
    });
  }
  return ok({ promptFile: config.systemPromptFile, vendors });
};
