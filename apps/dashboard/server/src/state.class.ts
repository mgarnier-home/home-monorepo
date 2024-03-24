import fs from 'fs';
import jsYaml from 'js-yaml';
import { logger } from 'logger';

import { config } from './utils/config';

import type { App } from '@shared/interfaces/app';

class State {
  constructor() {
    this.reloadSetup = this.reloadSetup.bind(this);
  }

  private _hosts: App.Setup.Host[] = [];
  private _globalConfig: App.Setup.GlobalConfig = {
    statusCheckInterval: 30000,
    pingInterval: 30000,
    statsApiUrl: '',
  };

  public get globalConfig(): App.Setup.GlobalConfig {
    return structuredClone(this._globalConfig);
  }

  public get hosts(): App.Setup.Host[] {
    return structuredClone(this._hosts);
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

  private static sanitizeClickAction(
    serviceUrl: string,
    clickAction?: App.ClickAction | App.ClickActionType
  ): App.ClickAction | undefined {
    if (!clickAction) {
      return undefined;
    }

    return typeof clickAction === 'string'
      ? { type: clickAction, url: serviceUrl }
      : { ...clickAction, url: clickAction.url || serviceUrl };
  }

  private static sanitizeStatusCheck(
    statusCheck: App.Setup.HostServiceStatusCheck,
    serviceUrl: string
  ): App.Setup.HostServiceStatusCheck {
    const { clickAction } = statusCheck;

    const type: any = Array.isArray((statusCheck as any).codes) ? 'multipleCodes' : 'singleCode';

    return {
      ...statusCheck,
      type,
      clickAction: State.sanitizeClickAction(serviceUrl, clickAction),
    };
  }

  private static sanitizeService(service: App.Setup.HostService): App.Setup.HostService {
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
      statusChecks: statusChecks.map((statusCheck) => State.sanitizeStatusCheck(statusCheck, service.url)),
    };
  }

  private static sanitizeHost(host: App.Setup.Host): App.Setup.Host {
    const nodesightUrl = host.nodesightUrl;

    return {
      ...host,
      nodesightUrl: nodesightUrl ? (nodesightUrl.endsWith('/') ? nodesightUrl.slice(0, -1) : nodesightUrl) : undefined,
      services: host.services.map((service) => State.sanitizeService(service)),
    };
  }
}

const state = new State();

export { state };
