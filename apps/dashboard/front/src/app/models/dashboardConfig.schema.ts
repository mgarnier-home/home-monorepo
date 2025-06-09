import { z } from 'zod';

export const responseColorSchema = z.object({
  color: z.string(),
  code: z.number(),
});

export const actionTypeEnum = z.enum(['open', 'fetch', 'request']);

export const actionSchema = z.discriminatedUnion('type', [
  z.object({
    name: z.string(),
    type: z.literal(actionTypeEnum.enum.open),
    url: z.string().url(),
  }),
  z.object({
    name: z.string(),
    type: z.literal(actionTypeEnum.enum.fetch),
    url: z.string().url(),
  }),
  z.object({
    name: z.string(),
    type: z.literal(actionTypeEnum.enum.request),
    method: z.enum([
      'GET',
      'POST',
      'PUT',
      'DELETE',
      'PATCH',
      'HEAD',
      'OPTIONS',
    ]),
    url: z.string().url(),
  }),
]);

export const healthCheckSchema = z.object({
  name: z.string(),
  url: z.string(),
  action: actionSchema,
  responsesColors: z.array(responseColorSchema),
});

export type HealthCheck = z.infer<typeof healthCheckSchema>;

export const serviceSchema = z.object({
  name: z.string(),
  icon: z.string(),
  dockerName: z.string().optional(),
  healthCheck: healthCheckSchema.optional(),
  healthChecks: z.array(healthCheckSchema).optional(),
});

export type Service = z.infer<typeof serviceSchema>;

export const hostSchema = z.object({
  name: z.string(),
  nodesight: z.string().url(),
  ip: z.string().ip(),
  icon: z.string(),
  services: z.array(serviceSchema),
});

export type Host = z.infer<typeof hostSchema>;

export const dashboardConfigSchema = z.object({
  hosts: z.array(hostSchema).optional(),
  statsApiUrl: z.string().url(),
});

export type DashboardConfig = z.infer<typeof dashboardConfigSchema>;
