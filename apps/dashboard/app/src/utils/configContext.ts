import { createContext } from "react";

import AppInterfaces from "@shared/interfaces/appInterfaces";

export const ConfigContext = createContext<AppInterfaces.AppConfig>({} as AppInterfaces.AppConfig);
