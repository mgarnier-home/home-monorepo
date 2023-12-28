import { createContext } from 'react';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';

export const WidgetContext = createContext<AppInterfaces.Host>({} as AppInterfaces.Host);
