import fs from 'fs';
import jsYaml from 'js-yaml';

import { config } from './utils/config';

import type { Setup } from '@shared/interfaces/setup';

const sanitizeClickAction = (
  serviceUrl: string,
  clickAction?: Setup.ClickAction | Setup.ClickActionType
): Setup.ClickAction | undefined => {
  if (!clickAction) {
    return undefined;
  }

  return typeof clickAction === 'string'
    ? { type: clickAction, url: serviceUrl }
    : { ...clickAction, url: clickAction.url || serviceUrl };
};

const sanitizeStatusCheck = (
  statusCheck: Setup.HostServiceStatusCheck,
  serviceUrl: string
): Setup.HostServiceStatusCheck => {
  const { clickAction } = statusCheck;

  const type: any = Array.isArray((statusCheck as any).codes) ? 'multipleCodes' : 'singleCode';

  return {
    ...statusCheck,
    type,
    clickAction: sanitizeClickAction(serviceUrl, clickAction),
  };
};

const sanitizeService = (service: Setup.HostService): Setup.HostService => {
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
    statusChecks: statusChecks.map((statusCheck) => sanitizeStatusCheck(statusCheck, service.url)),
  };
};

const sanitizeHost = (host: Setup.Host): Setup.Host => {
  const nodesightUrl = host.nodesightUrl;

  return {
    ...host,
    nodesightUrl: nodesightUrl ? (nodesightUrl.endsWith('/') ? nodesightUrl.slice(0, -1) : nodesightUrl) : undefined,
    services: host.services.map((service) => sanitizeService(service)),
  };
};

const sanitizeAppSetup = (appSetup: Setup.App): Setup.App => {
  const statsApiUrl = appSetup.global.statsApiUrl;

  return {
    global: {
      statusCheckInterval: appSetup.global.statusCheckInterval ?? 30000,
      pingInterval: appSetup.global.pingInterval ?? 30000,
      statsApiUrl: statsApiUrl.endsWith('/') ? statsApiUrl.slice(0, -1) : statsApiUrl,
    },
    hosts: appSetup.hosts.map((host) => sanitizeHost(host)),
  };
};

export const getAppSetup = async (): Promise<Setup.App> => {
  if (fs.existsSync(config.appSetupPath)) {
    const appSetupContentStr = await fs.promises.readFile(config.appSetupPath, 'utf-8');

    const appSetupContent = jsYaml.load(appSetupContentStr) as Setup.App;

    return sanitizeAppSetup(appSetupContent);
  } else {
    throw new Error('Config file not found');
  }
};
