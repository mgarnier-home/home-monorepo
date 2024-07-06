import fs from 'fs';
import jsYaml from 'js-yaml';
import { logger } from 'logger';
import { SimpleCache, Utils } from 'utils';

import { config } from './utils/config';
import { makeRequest, pingHost } from './utils/utils';

import type { App } from '@shared/interfaces/app';
import type { MakeRequestResponse } from '@shared/interfaces/utils';
class State {
  constructor() {
    this.reloadSetup = this.reloadSetup.bind(this);
    this._refreshHostPing = this._refreshHostPing.bind(this);
  }

  private _hosts: App.Setup.Host[] = [];
  private _hostPings: Map<string, { lastPing: App.Ping | null; interval: NodeJS.Timeout }> = new Map();
  private _statusChecksCodesCache = new SimpleCache<MakeRequestResponse<void>>(10);
  private _globalConfig: App.Setup.GlobalConfig = {
    statusCheckInterval: 30000,
    pingInterval: 30000,
    statsApiUrl: '',
  };

  public async reloadSetup() {
    if (fs.existsSync(config.appSetupPath)) {
      const appSetupContentStr = await fs.promises.readFile(config.appSetupPath, 'utf-8');

      const { globalConfig, hosts } = jsYaml.load(appSetupContentStr) as {
        globalConfig: App.Setup.GlobalConfig;
        hosts: App.Setup.Host[];
      };
      const { statsApiUrl } = globalConfig;

      this._globalConfig = {
        statusCheckInterval: globalConfig.statusCheckInterval ?? 30000,
        pingInterval: globalConfig.pingInterval ?? 30000,
        statsApiUrl: statsApiUrl.endsWith('/') ? statsApiUrl.slice(0, -1) : statsApiUrl,
      };

      this._hosts = hosts.map((host) => State.sanitizeHost(host));

      logger.info('App setup reloaded');

      logger.debug(this._hosts);
      logger.debug(this._globalConfig);
    } else {
      logger.error('Config file not found');
    }
  }

  public startPinging() {
    this._hosts.forEach(this._refreshHostPing);
  }

  private _refreshHostPing(host: App.Setup.Host) {
    const hostPing = this._hostPings.get(host.id);

    if (host.enablePing && !hostPing) {
      const interval = setInterval(async () => {
        const ping = await pingHost(host);

        this._hostPings.set(host.id, { lastPing: ping, interval });
      }, this._globalConfig.pingInterval);

      this._hostPings.set(host.id, { lastPing: null, interval });
    }

    if (!host.enablePing && hostPing) {
      clearInterval(hostPing.interval);
      this._hostPings.delete(host.id);
    }
  }

  public stopPinging() {
    this._hostPings.forEach(({ interval }) => clearInterval(interval));
    this._hostPings.clear();
  }

  private async _getStatusChecks(): Promise<App.State.StatusCheck[]> {
    const setupStatusChecks = this._hosts.flatMap((host) => host.services.flatMap((service) => service.statusChecks));

    const statusChecks = await Promise.all(
      setupStatusChecks.map<Promise<App.State.StatusCheck>>(async (statusCheck) => {
        if (this._statusChecksCodesCache.get(statusCheck.id) === null) {
          const response = await makeRequest<void>(statusCheck.url, 'GET');

          this._statusChecksCodesCache.set(statusCheck.id, response);
        }

        const cachedCode = this._statusChecksCodesCache.get(statusCheck.id)!;

        return {
          id: statusCheck.id,
          name: statusCheck.name,
          successCodes:
            statusCheck.type === 'singleCode'
              ? [{ code: statusCheck.success, color: statusCheck.color }]
              : statusCheck.codes,
          clickAction: statusCheck.clickAction,
          lastRequest: cachedCode,
        };
      })
    );

    return statusChecks;
  }

  private static sanitizeClickAction(
    serviceUrl: string,
    clickAction?: App.ClickAction | App.ClickActionType
  ): App.ClickAction | undefined {
    if (!clickAction) {
      return undefined;
    }

    const sanitizedClickAction =
      typeof clickAction === 'string'
        ? { type: clickAction, url: serviceUrl }
        : { ...clickAction, url: clickAction.url || serviceUrl };

    return sanitizedClickAction;
  }

  private static sanitizeService(service: App.Setup.Service, hostName: string): App.Setup.Service {
    const { clickAction, statusChecks } = service;

    const sanitizedService = {
      ...service,
      clickAction: State.sanitizeClickAction(service.url, clickAction),
      statusChecks: statusChecks.map((statusCheck, index) => {
        const { clickAction } = statusCheck;

        const type: any = Array.isArray((statusCheck as any).codes) ? 'multipleCodes' : 'singleCode';

        const sanitizedStatusCheck = {
          ...statusCheck,
          name: statusCheck.name || index === 1 ? 'Service' : '',
          id: Utils.hashCode(hostName + service.name + service.url + statusCheck.url + type),
          type,
          clickAction: State.sanitizeClickAction(service.url, clickAction),
          url: statusCheck.url || service.url,
        };

        return sanitizedStatusCheck;
      }),
    };

    return sanitizedService;
  }

  private static sanitizeHost(host: App.Setup.Host): App.Setup.Host {
    const nodesightUrl = host.nodesightUrl;

    const sanitizedHost = {
      ...host,
      nodesightUrl: nodesightUrl ? (nodesightUrl.endsWith('/') ? nodesightUrl.slice(0, -1) : nodesightUrl) : undefined,
      services: host.services.map((service) => State.sanitizeService(service, host.name)),
    };

    return sanitizedHost;
  }
}

const state = new State();

export { state };
