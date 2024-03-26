import { z } from 'zod';

import { ClickActionSchema } from './clickAction';
import { SuccessCodeSchema } from './successCode';

export namespace SanitizedConfig {
  export const StatusCheckSchema = z.object({
    id: z.string(),
    name: z.string(),
    url: z.string(),
    successCodes: z.array(SuccessCodeSchema),
    clickAction: ClickActionSchema.optional(),
  });

  export const ServiceSchema = z.object({
    id: z.string(),
    name: z.string(),
    icon: z.string(),
    url: z.string(),
    order: z.number(),
    clickAction: ClickActionSchema.optional(),
    statusChecks: z.array(StatusCheckSchema),
  });

  export const HostSchema = z.discriminatedUnion('enablePing', [
    z.object({
      name: z.string(),
      id: z.string(),
      icon: z.string(),
      enablePing: z.literal(false),
      ip: z.string().optional(),
      services: z.array(ServiceSchema),
    }),
    z.object({
      name: z.string(),
      id: z.string(),
      icon: z.string(),
      enablePing: z.literal(true),
      ip: z.string(),
      services: z.array(ServiceSchema),
    }),
  ]);

  export type Host = z.infer<typeof HostSchema>;
  export type Service = z.infer<typeof ServiceSchema>;
  export type StatusCheck = z.infer<typeof StatusCheckSchema>;
}
