import type { Widget } from '@shared/interfaces/widget';
import jsYaml from 'js-yaml';

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

class Setup {
  private _hosts: Host[] = [];
  private _globalConfig: GlobalConfig = {
    statusCheckInterval: 30000,
    pingInterval: 30000,
    statsApiUrl: '',
  };

  public get globalConfig(): GlobalConfig {
    return structuredClone(this._globalConfig);
  }

  public get hosts(): Host[] {
    return structuredClone(this._hosts);
  }

  public reloadAppSetup(appSetupContentStr: string) {
    const { globalConfig, hosts } = jsYaml.load(appSetupContentStr) as { globalConfig: GlobalConfig; hosts: Host[] };
    const { statsApiUrl } = globalConfig;

    this._globalConfig = {
      statusCheckInterval: globalConfig.statusCheckInterval ?? 30000,
      pingInterval: globalConfig.pingInterval ?? 30000,
      statsApiUrl: statsApiUrl.endsWith('/') ? statsApiUrl.slice(0, -1) : statsApiUrl,
    };

    this._hosts = hosts.map((host) => Setup.sanitizeHost(host));
  }

  private static sanitizeClickAction(
    serviceUrl: string,
    clickAction?: ClickAction | ClickActionType
  ): ClickAction | undefined {
    if (!clickAction) {
      return undefined;
    }

    return typeof clickAction === 'string'
      ? { type: clickAction, url: serviceUrl }
      : { ...clickAction, url: clickAction.url || serviceUrl };
  }

  private static sanitizeStatusCheck(statusCheck: HostServiceStatusCheck, serviceUrl: string): HostServiceStatusCheck {
    const { clickAction } = statusCheck;

    const type: any = Array.isArray((statusCheck as any).codes) ? 'multipleCodes' : 'singleCode';

    return {
      ...statusCheck,
      type,
      clickAction: Setup.sanitizeClickAction(serviceUrl, clickAction),
    };
  }

  private static sanitizeService(service: HostService): HostService {
    const { clickAction, statusChecks } = service;

    let sanitizedClickAction;
    if (clickAction) {
      sanitizedClickAction =
        typeof clickAction === 'string'
          ? { type: clickAction, url: service.url }
          : { ...clickAction, url: clickAction.url || service.url };
    }

    return {
      ...service,
      clickAction: sanitizedClickAction,
      statusChecks: statusChecks.map((statusCheck) => Setup.sanitizeStatusCheck(statusCheck, service.url)),
    };
  }

  private static sanitizeHost(host: Host): Host {
    const nodesightUrl = host.nodesightUrl;

    return {
      ...host,
      nodesightUrl: nodesightUrl ? (nodesightUrl.endsWith('/') ? nodesightUrl.slice(0, -1) : nodesightUrl) : undefined,
      services: host.services.map((service) => Setup.sanitizeService(service)),
    };
  }
}

const setup = new Setup();

export { Setup, setup };
