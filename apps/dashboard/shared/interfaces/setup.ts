import type { Widget } from './widget';

export namespace Setup {
  export type App = {
    global: Setup.Global;
    hosts: Host[];
  };

  export type Global = {
    statusCheckInterval: number;
    pingInterval: number;
    statsApiUrl: string;
  };

  export type Host = {
    name: string;
    id: string;
    icon: string;
    ip: string;
    enablePing: boolean;
    nodesightUrl?: string;
    order?: number;
    services: HostService[];
    widgets: Widget[];
  };

  export type HostService = {
    name: string;
    icon: string;
    url: string;
    order?: number;
    clickAction?: ClickAction | ClickActionType;
    statusChecks: HostServiceStatusCheck[];
  };

  export type ClickActionType = 'redirect' | 'open' | 'fetch' | 'none';

  export type ClickAction = {
    type: ClickActionType;
    url: string;
  };

  export type HostServiceStatusCheck =
    | {
        type: 'singleCode';
        name: string;
        url: string;
        clickAction?: ClickAction | ClickActionType;
        success: number;
        color: string;
      }
    | {
        type: 'multipleCodes';
        name: string;
        url: string;
        clickAction?: ClickAction | ClickActionType;
        codes: { code: number; color: string }[];
      };
}
