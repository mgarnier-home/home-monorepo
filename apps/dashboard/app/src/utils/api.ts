import type { AppInterfaces } from '@shared/interfaces/appInterfaces';
import type { ApiInterfaces } from '@shared/interfaces/apiInterfaces';

import { SERVER_ROUTES, SERVER_ROUTES_METHODS, ServerRoute } from '../../../shared/routes';
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

  async function apiRequest<T, U>(route: ServerRoute, data: T): Promise<U> {
    if (socket.connected) {
      const response: U = await socket.apiRequest<T, U>(route, data);

      return response;
    } else {
      const method = SERVER_ROUTES_METHODS[route];
      const response = await fetch(route, {
        method,
        ...(method === 'POST'
          ? {
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify(data),
            }
          : {}),
      });

      return await response.json();
    }
  }

  export async function getConfig(): Promise<AppInterfaces.AppConfig> {
    return await apiRequest<{}, AppInterfaces.AppConfig>(SERVER_ROUTES.CONF, {});
  }

  export async function pingHost(host: AppInterfaces.Host): Promise<ApiInterfaces.PingHost.Response> {
    const data: ApiInterfaces.PingHost.Request = { ip: host.ip };

    return await apiRequest<ApiInterfaces.PingHost.Request, ApiInterfaces.PingHost.Response>(
      SERVER_ROUTES.PING_HOST,
      data
    );
  }

  export async function makeServerRequest<R>(
    url: string,
    method: string,
    body?: string
  ): Promise<ApiInterfaces.MakeRequest.Response<R>> {
    const data: ApiInterfaces.MakeRequest.Request = { url, method, body };

    return await apiRequest<ApiInterfaces.MakeRequest.Request, ApiInterfaces.MakeRequest.Response<any>>(
      SERVER_ROUTES.MAKE_REQUEST,
      data
    );
  }

  export async function getStatusChecks(
    statusChecks: ApiInterfaces.StatusChecks.Request
  ): Promise<ApiInterfaces.StatusChecks.Response> {
    return await apiRequest<ApiInterfaces.StatusChecks.Request, ApiInterfaces.StatusChecks.Response>(
      SERVER_ROUTES.STATUS_CHECKS,
      statusChecks
    );
  }
}
