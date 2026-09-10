import type { Config } from "@ponte/core";
import { readPrompt, resolveContent, writePrompt } from "../infra/config-file";

export const readSystemPrompt = async (
  config: Config,
): Promise<string | null> => readPrompt(config.systemPromptFile);

export const setSystemPrompt = async (
  config: Config,
  fileOrLiteral: string,
): Promise<void> => {
  await writePrompt(
    config.systemPromptFile,
    await resolveContent(fileOrLiteral),
  );
};
