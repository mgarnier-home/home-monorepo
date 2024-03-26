import type { MakeRequestResponse } from './utils';
import type { Widget } from './widget';

export namespace App {
  export type ClickActionType = 'redirect' | 'open' | 'fetch' | 'none';

  export type ClickAction = {
    type: ClickActionType;
    url: string;
  };

  export type Ping = {
    ping: boolean;
    duration: number;
    ms: number;
  };

  export namespace State {
    export type SuccessCode = {
      color: string;
      code: number;
    };

    export type StatusCheck = {
      id: string;
      name: string;
      successCodes: SuccessCode[];
      lastRequest: MakeRequestResponse<void>;
      clickAction?: ClickAction;
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
      ping: Ping | null | undefined;
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
      order: number;
      services: Service[];
      // widgets: Widget[];
    };

    export type Service = {
      name: string;
      icon: string;
      url: string;
      order: number;
      clickAction?: ClickAction;
      statusChecks: StatusCheck[];
    };

    export type StatusCheck = {
      id: string;
      name: string;
      url: string;
      clickAction?: ClickAction;
      successCodes: { code: number; color: string }[];
    };
  }
}
