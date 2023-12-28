import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
import type { ApiInterfaces } from '@shared/interfaces/apiInterfaces';

import { ServerRoutes } from '../../../shared/routes';
import { socket } from './socket';

export namespace Api {
  export async function handleClickAction(clickAction: AppInterfaces.ClickAction) {
    switch (clickAction.type) {
      case 'redirect':
        window.location.href = clickAction.url;
        break;
      case 'open':
        window.open(clickAction.url);
        break;
      case 'fetch':
        return makeServerRequest(clickAction.url, 'GET');
      case 'none':
      default:
        break;
    }

    return null;
  }

  async function apiRequest<T, U>(route: ServerRoutes, data: T, backup?: (data: T) => Promise<U>): Promise<U> {
    if (socket.connected) {
      const response: U = await socket.apiRequest<T, U>(route, data);

      return response;
    } else {
      if (backup) {
        return backup(data);
      }

      throw new Error('Socket not connected and no backup function provided');
    }
  }

  export async function getConfig(): Promise<AppInterfaces.AppConfig> {
    return await apiRequest<{}, AppInterfaces.AppConfig>(ServerRoutes.CONF, {}, async () => {
      const config = await fetch(ServerRoutes.CONF);

      return await config.json();
    });
  }

  export async function pingHost(host: AppInterfaces.Host): Promise<ApiInterfaces.PingHost.Response> {
    const data: ApiInterfaces.PingHost.Request = { ip: host.ip };

    return await apiRequest<ApiInterfaces.PingHost.Request, ApiInterfaces.PingHost.Response>(
      ServerRoutes.PING_HOST,
      data,
      async () => {
        const response = await fetch(ServerRoutes.PING_HOST, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data),
        });

        return await response.json();
      }
    );
  }

  export async function makeServerRequest<R>(
    url: string,
    method: string,
    body?: string
  ): Promise<ApiInterfaces.MakeRequest.Response<R>> {
    const data: ApiInterfaces.MakeRequest.Request = { url, method, body };

    return await apiRequest<ApiInterfaces.MakeRequest.Request, ApiInterfaces.MakeRequest.Response<any>>(
      ServerRoutes.MAKE_REQUEST,
      data,
      async () => {
        const response = await fetch(ServerRoutes.MAKE_REQUEST, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(data),
        });

        return await response.json();
      }
    );
  }
}
