import type {
  Config,
  ReadPrompt,
  ResolveContent,
  WritePrompt,
} from "@ponte/core";

export const createReadSystemPrompt =
  (readPrompt: ReadPrompt) =>
  (config: Config): Promise<string | null> =>
    readPrompt(config.systemPromptFile);

export const createSetSystemPrompt =
  (writePrompt: WritePrompt, resolveContent: ResolveContent) =>
  async (config: Config, fileOrLiteral: string): Promise<void> => {
    await writePrompt(
      config.systemPromptFile,
      await resolveContent(fileOrLiteral),
    );
  };
