import { z } from 'zod';

export const SuccessCodeSchema = z.object({
  code: z.number(),
  color: z.string().optional(),
});

export type SuccessCode = z.infer<typeof SuccessCodeSchema>;
