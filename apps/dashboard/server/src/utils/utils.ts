import ping from 'ping';

import { logger } from '@libs/logger';
import { Utils } from '@libs/utils';

import type { App } from '@shared/interfaces/app';
import type { MakeRequestResponse } from '@shared/interfaces/utils';
export const pingHost = async (host: App.Setup.Host): Promise<App.Ping> => {
  const startTime = Date.now();

  const pingResult = await ping.promise.probe(host.ip, { timeout: 3 });

  const duration = Date.now() - startTime;

  logger.info(`Pinged host ${host.ip} in ${duration}ms, got ${pingResult.alive}: ${pingResult.time}ms`);

  const ms = Math.floor(Number(pingResult.time));

  return { ping: pingResult.alive, duration, ms };
};

export const makeRequest = async <Data>(
  url: string,
  method: string,
  body?: string
): Promise<MakeRequestResponse<Data>> => {
  // const cached = requestsCache.get(url);

  const startTime = Date.now();

  let code = 0;
  let data = undefined;

  // if (!cached) {
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
  // } else {
  //   code = cached.code;
  //   data = cached.data;
  //   log(`[MakeRequest] using cached response for ${url}`);
  // }

  const duration = Date.now() - startTime;

  logger.info(`Request to ${method} ${url} in ${duration}ms, got ${code}`);

  // requestsCache.set(url, { code, data });

  return { code, duration, data };
};
