import type { Widget } from './widget';

export namespace App {
  export type ClickActionType = 'redirect' | 'open' | 'fetch' | 'none';

  export type ClickAction = {
    type: ClickActionType;
    url: string;
  };

  export namespace State {
    export type SuccessCode = {
      color: string;
      code: number;
    };

    export type StatusCheck = {
      name: string;
      successCodes: SuccessCode[];
      lastCheck: number;
      clickAction: ClickAction;
    };

    export type Service = {
      name: string;
      icon: string;
      clickAction: ClickAction;
      statusChecks: StatusCheck[];
    };

    export type Host = {
      id: string;
      name: string;
      icon: string;
      ip: string;
      ping: number | null;
      services: Service[];
    };
  }

  export namespace Setup {
    export type GlobalConfig = {
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
}
