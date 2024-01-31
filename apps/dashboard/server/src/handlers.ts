import jsYaml from 'js-yaml';
import { logger } from 'logger';
import fs from 'node:fs';
import ping from 'ping';
import { Utils } from 'utils';

import { config } from './utils/config';

import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
const log = (...args: any[]) => {
  logger.info(`[Handlers]`, ...args);
};

export namespace Handlers {
  const sanitizeAppConf = (appConf: AppInterfaces.AppConfig): AppInterfaces.AppConfig => {
    const sanitizeClickActions = (host: AppInterfaces.Host): AppInterfaces.Host => {
      if (host.services === undefined) {
        host.services = [];
      }

      if (host.nodesightUrl && host.nodesightUrl.endsWith('/')) {
        host.nodesightUrl = host.nodesightUrl.slice(0, -1);
      }

      host.services.forEach((service) => {
        if (typeof service.clickAction === 'string') {
          service.clickAction = { type: service.clickAction, url: service.url } as AppInterfaces.ClickAction;
        }

        if (service.clickAction && !service.clickAction.url) {
          service.clickAction.url = service.url;
        }

        service.statusChecks.forEach((statusCheck) => {
          if (Array.isArray((statusCheck as any).codes)) {
            statusCheck.type = 'multipleCodes';
          } else {
            statusCheck.type = 'singleCode';
          }

          if (typeof statusCheck.clickAction === 'string') {
            statusCheck.clickAction = {
              type: statusCheck.clickAction,
              url: statusCheck.url,
            } as AppInterfaces.ClickAction;
          }

          if (statusCheck.clickAction && !statusCheck.clickAction.url) {
            statusCheck.clickAction.url = statusCheck.url || service.url;
          }
        });
      });

      return host;
    };

    appConf.globalConfig = {
      ...{ statusCheckInterval: 10000, pingInterval: 10000, statsApiUrl: '' },
      ...appConf.globalConfig,
    };

    if (appConf.globalConfig.statsApiUrl.endsWith('/')) {
      appConf.globalConfig.statsApiUrl = appConf.globalConfig.statsApiUrl.slice(0, -1);
    }

    appConf.hosts.forEach((host) => {
      host = sanitizeClickActions(host);
    });

    return appConf;
  };

  export const getAppConf = async (): Promise<string> => {
    const isYml = config.appConfPath.endsWith('.yml');

    if (fs.existsSync(config.appConfPath)) {
      const appConfContentStr = await fs.promises.readFile(config.appConfPath, 'utf-8');

      let appConfContent: AppInterfaces.AppConfig;

      if (isYml) {
        appConfContent = jsYaml.load(appConfContentStr) as AppInterfaces.AppConfig;
      } else {
        appConfContent = JSON.parse(appConfContentStr) as AppInterfaces.AppConfig;
      }

      appConfContent = sanitizeAppConf(appConfContent);

      return JSON.stringify(appConfContent);
    } else {
      throw new Error('Config file not found');
    }
  };

  export const pingHost = async (ip: string): Promise<{ ping: boolean; duration: number; ms: number }> => {
    const startTime = Date.now();

    const pingResult = await ping.promise.probe(ip, { timeout: 3 });

    const duration = Date.now() - startTime;

    log(`[PingHost] pinged host ${ip} in ${duration}ms, got ${pingResult.alive}: ${pingResult.time}ms`);

    const ms = Math.floor(Number(pingResult.time));

    return { ping: pingResult.alive, duration, ms };
  };

  export const makeRequest = async <Data>(
    url: string,
    method: string,
    body?: string
  ): Promise<{ code: number; duration: number; data?: Data }> => {
    const startTime = Date.now();

    let code = 0;
    let data = undefined;

    try {
      const response = await Utils.fetchWithTimeout(url, 10000, {
        method: method,
        headers: {
          Status: 'true',
        },
        body: body,
      });

      data = await response.text();

      try {
        data = JSON.parse(data);
      } catch (error) {}

      code = response.status;
    } catch (error) {
      logger.error(error);

      code = 500;
    }

    const duration = Date.now() - startTime;

    log(`[MakeRequest] to ${url} in ${duration}ms, got ${code}`);

    return { code, duration, data };
  };
}
