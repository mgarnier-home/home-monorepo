import { z } from 'zod';

export const responseColorSchema = z.object({
  color: z.string(),
  code: z.number(),
});

export const actionSchema = z.discriminatedUnion('type', [
  z.object({
    name: z.string(),
    type: z.literal('open'),
    url: z.string().url(),
  }),
  z.object({
    name: z.string(),
    type: z.literal('fetch'),
    url: z.string().url(),
  }),
  z.object({
    name: z.string(),
    type: z.literal('request'),
    method: z.string(),
    url: z.string().url(),
    body: z.string().optional(),
  }),
]);

export const serviceStatusSchema = z.object({
  name: z.string(),
  url: z.string(),
  action: actionSchema,
  responsesColors: z.array(responseColorSchema),
});

export const serviceSchema = z.object({
  name: z.string(),
  icon: z.string(),
  dockerName: z.string().optional(),
  primary: serviceStatusSchema.optional(),
  secondaries: z.array(serviceStatusSchema).optional(),
});

export const hostSchema = z.object({
  name: z.string(),
  nodesight: z.string().url(),
  ip: z.string().ip(),
  icon: z.string(),
  services: z.array(serviceSchema),
});

export const dashboardConfigSchema = z.object({
  hosts: z.array(hostSchema).optional(),
  statsApiUrl: z.string().url(),
});
