import fs from 'fs';
import jsYaml from 'js-yaml';
import { logger } from 'logger';
import pLimit from 'p-limit';
import { Utils } from 'utils';

import { config } from './utils/config';
import { makeRequest, pingHost } from './utils/utils';

import type { App } from '@shared/interfaces/app';
class State {
  constructor() {
    this.reloadSetup = this.reloadSetup.bind(this);
    this._refreshHostPing = this._refreshHostPing.bind(this);
  }

  private _hosts: App.Setup.Host[] = [];
  private _hostPings: Map<string, { lastPing: App.Ping | null; interval: NodeJS.Timeout }> = new Map();
  private _statusChecks: Map<string, { lastCheck: number | null; interval: NodeJS.Timeout }> = new Map();
  private _globalConfig: App.Setup.GlobalConfig = {
    statusCheckInterval: 30000,
    pingInterval: 30000,
    statsApiUrl: '',
  };

  private _getStatusChecksList(): App.Setup.StatusCheck[] {
    return this._hosts.flatMap((host) => host.services.flatMap((service) => service.statusChecks));
  }

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

  public startTracking() {
    // this._hosts.forEach(this._refreshHostPing);
    // this._refreshStatusChecks();
    this._test();
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

  private _refreshStatusChecks() {
    const statusCheckList = this._getStatusChecksList();

    const delayBetweenChecks = this._globalConfig.statusCheckInterval / statusCheckList.length;
    let actualDelay = 0;

    for (const statusCheck of statusCheckList) {
      const statusCheckData = this._statusChecks.get(statusCheck.id);

      if (!statusCheckData) {
        const timeout = setTimeout(() => {
          const interval = setInterval(() => {}, this._globalConfig.statusCheckInterval);
        }, actualDelay);

        actualDelay += delayBetweenChecks;

        this._statusChecks.set(statusCheck.id, { lastCheck: null, interval: timeout });
      }
    }
  }

  private async _test() {
    const statusCheckList = [
      ...this._getStatusChecksList(),
      // ...this._getStatusChecksList(),
      // ...this._getStatusChecksList(),
    ];

    const limit = pLimit(10);

    const statusCheckListMap = statusCheckList.map((statusCheck) => {
      return limit(() => makeRequest(statusCheck.url, 'GET'));
    });

    console.time('test');

    const results = await Promise.all(statusCheckListMap);

    console.timeEnd('test');

    // console.log(results);
  }

  public stopTracking() {
    this._hostPings.forEach(({ interval }) => clearInterval(interval));
    this._hostPings.clear();
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

  private static sanitizeStatusCheck(
    statusCheck: App.Setup.StatusCheck,
    hostName: string,
    serviceName: string,
    serviceUrl: string
  ): App.Setup.StatusCheck {
    const { clickAction } = statusCheck;

    const type: any = Array.isArray((statusCheck as any).codes) ? 'multipleCodes' : 'singleCode';

    const sanitizedStatusCheck = {
      ...statusCheck,
      id: Utils.hashCode(hostName + serviceName + serviceUrl + statusCheck.url + type),
      type,
      clickAction: State.sanitizeClickAction(serviceUrl, clickAction),
      url: statusCheck.url || serviceUrl,
    };

    return sanitizedStatusCheck;
  }

  private static sanitizeService(service: App.Setup.Service, hostName: string): App.Setup.Service {
    const { clickAction, statusChecks } = service;

    const sanitizedService = {
      ...service,
      clickAction: State.sanitizeClickAction(service.url, clickAction),
      statusChecks: statusChecks.map((statusCheck) =>
        State.sanitizeStatusCheck(statusCheck, hostName, service.name, service.url)
      ),
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
