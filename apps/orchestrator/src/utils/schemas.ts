import { z } from 'zod';

export interface Config {
  serverPort: number;
  composeEnvFilesPaths: string[];
  composeFolderPath: string;
  stackFilePath: string;
}

export const hostSchema = z.object({
  name: z.string(),
  ip: z.string(),
  username: z.string(),
});

export const stackSchema = z.object({
  stacks: z.array(z.string()),
  hosts: z.array(hostSchema),
});

export type Host = z.infer<typeof hostSchema>;
export type Stack = z.infer<typeof stackSchema>;
