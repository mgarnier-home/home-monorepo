import { z } from 'zod';

export const statusEnum = z.enum(['ok', 'warning', 'error', 'unknown']);

export const healthCheckStateSchema = z.object({
  id: z.string(),
  name: z.string(),
  code: z.number(),
  responseTime: z.number(),
  response: z.string(),
});

export type HealthCheckState = z.infer<typeof healthCheckStateSchema>;

export const serviceStateSchema = z.object({
  id: z.string(),
  healthCheck: healthCheckStateSchema.optional(),
  dockerStatus: statusEnum,
  healthChecks: z.array(healthCheckStateSchema).optional(),
});

export type ServiceState = z.infer<typeof serviceStateSchema>;

export const hostStateSchema = z.object({
  id: z.string(),
  ping: z.number(),
  status: statusEnum,
  // dockerStatus: statusEnum,
});

export type HostState = z.infer<typeof hostStateSchema>;
