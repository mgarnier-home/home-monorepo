import type { WidgetInterfaces } from './widgetInterfaces';

export namespace AppInterfaces {
  export interface AppConfig {
    pageConfig: PageConfig;
    globalConfig: GlobalConfig;
    hosts: Host[];
  }

  export interface GlobalConfig {
    statusCheckInterval: number;
    pingInterval: number;
    statsApiUrl: string;
  }

  export interface PageConfig {
    background: string;
    pageTitle: string;
    pageIcon: string;
  }

  export interface Host {
    name: string;
    id: string;
    icon: string;
    ip: string;
    enablePing: boolean;
    nodesightUrl?: string;
    order?: number;
    services: HostService[];
    widgets: WidgetInterfaces.Widget[];
  }

  export type ClickActionType = 'redirect' | 'open' | 'fetch' | 'none';

  export interface ClickAction {
    type: ClickActionType;
    url: string;
  }

  export interface HostService {
    name: string;
    icon: string;
    url: string;
    order?: number;
    clickAction: ClickAction | ClickActionType;
    statusChecks: HostServiceStatusCheck[];
  }

  export type HostServiceStatusCheck =
    | {
        type: 'singleCode';
        name: string;
        url: string;
        clickAction: ClickAction | ClickActionType;
        success: number;
        color: string;
      }
    | {
        type: 'multipleCodes';
        name: string;
        url: string;
        clickAction: ClickAction | ClickActionType;
        codes: { code: number; color: string }[];
      };
}
