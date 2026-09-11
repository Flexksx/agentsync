import type { VendorPlan } from "./plan";

export type ReadSymlinks = (plan: VendorPlan) => Promise<Map<string, string>>;
export type ApplyPlan = (
  plan: VendorPlan,
  stale: readonly string[],
) => Promise<void>;
