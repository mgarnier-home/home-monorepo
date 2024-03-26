import { z } from 'zod';

import { ClickActionSchema } from './clickAction';
import { SuccessCodeSchema } from './successCode';

export namespace State {
  export const StatusCheckSchema = z.object({
    id: z.string(),
    name: z.string(),
    successCodes: z.array(SuccessCodeSchema),
    lastRequest: z.object({ code: z.number(), duration: z.number() }),
    clickAction: ClickActionSchema.optional(),
  });

  export const ServiceSchema = z.object({
    id: z.string(),
    name: z.string(),
    icon: z.string(),
    order: z.number(),
    clickAction: ClickActionSchema.optional(),
    statusChecks: z.array(StatusCheckSchema),
  });

  export const HostSchema = z.object({
    id: z.string(),
    name: z.string(),
    icon: z.string(),
    ping: z.object({ ping: z.boolean(), duration: z.number(), ms: z.number() }).nullable(),
    services: z.array(ServiceSchema),
  });

  export type Host = z.infer<typeof HostSchema>;
  export type Service = z.infer<typeof ServiceSchema>;
  export type StatusCheck = z.infer<typeof StatusCheckSchema>;
}
