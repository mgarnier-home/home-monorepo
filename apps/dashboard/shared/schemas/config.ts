import { z } from 'zod';

import { ClickActionSchema, ClickActionTypeSchema } from './clickAction';
import { SuccessCodeSchema } from './successCode';

// Config

export namespace Config {
  export const ConfigClickActionSchema = z.union([ClickActionTypeSchema, ClickActionSchema]);

  // export const GlobalConfigSchema = z.object({
  //   statusCheckInterval: z.number(),
  //   pingInterval: z.number(),
  //   statsApiUrl: z.string(),
  // });

  export const StatusCheckSchema = z.object({
    success: z.number(),
    url: z.string().optional(),
    name: z.string().optional(),
    codes: z.array(SuccessCodeSchema).optional(),
    clickAction: ConfigClickActionSchema.optional(),
  });

  export const ServiceSchema = z.object({
    name: z.string(),
    url: z.string(),
    icon: z.string().optional(),
    order: z.number().optional(),
    clickAction: ConfigClickActionSchema.optional(),
    statusChecks: z.array(StatusCheckSchema).optional(),
  });

  export const HostSchema = z.object({
    name: z.string(),
    id: z.string().optional(),
    icon: z.string().optional(),
    ip: z.string().ip().optional(),
    enablePing: z.boolean().default(false),
    nodesightUrl: z.string().optional(),
    order: z.number().optional(),
    services: z.array(ServiceSchema).optional(),
  });

  export type Host = z.infer<typeof HostSchema>;
  export type Service = z.infer<typeof ServiceSchema>;
  export type StatusCheck = z.infer<typeof StatusCheckSchema>;
}
