import { z } from 'zod';

export const socketEvents = z.enum(['dashboardConfig', 'hostStateUpdate', 'serviceStateUpdate']);
