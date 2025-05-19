import { z } from 'zod';

export const statusEnum = z.enum(['ok', 'warning', 'error']);

export const healthCheckStateSchema = z.object({
  status: statusEnum,
  name: z.string(),
});

export const serviceStateSchema = z.object({
  healthCheck: statusEnum,
  dockerStatus: statusEnum,
  healthChecks: z.array(healthCheckStateSchema).optional(),
});

export type ServiceState = z.infer<typeof serviceStateSchema>;

export const hostStateSchema = z.object({
  ping: z.number(),
  status: statusEnum,
  dockerStatus: statusEnum,
});

export type HostState = z.infer<typeof hostStateSchema>;
