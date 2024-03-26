import { z } from 'zod';

export const ClickActionTypeSchema = z.union([
  z.literal('redirect'),
  z.literal('open'),
  z.literal('fetch'),
  z.literal('none'),
]);

export const ClickActionSchema = z.object({
  type: ClickActionTypeSchema,
  url: z.string(),
});

export type ClickActionType = z.infer<typeof ClickActionTypeSchema>;
export type ClickAction = z.infer<typeof ClickActionSchema>;
