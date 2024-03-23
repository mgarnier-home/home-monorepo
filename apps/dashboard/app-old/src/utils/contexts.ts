import { createContext } from 'react';

import { ApiInterfaces } from '@shared/interfaces/apiInterfaces';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
export const ConfigContext = createContext<AppInterfaces.AppConfig>({} as AppInterfaces.AppConfig);

export const WidgetContext = createContext<AppInterfaces.Host>({} as AppInterfaces.Host);

export const StatusChecksContext = createContext<Record<string, ApiInterfaces.StatusChecks.ResponseData>>({});
