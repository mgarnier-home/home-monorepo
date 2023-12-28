import { createContext } from 'react';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';

export const ConfigContext = createContext<AppInterfaces.AppConfig>({} as AppInterfaces.AppConfig);
