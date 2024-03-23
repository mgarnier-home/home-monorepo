import type { Request, Response } from 'express';

import type { ApiInterfaces as Api } from '@shared/interfaces/apiInterfaces';
import { Handlers } from './handlers';

export namespace RestControllers {
  export const getConf = async (req: Request, res: Response) => {
    res.send(await Handlers.getAppConf());
  };

  export const postPingHost = async (
    req: Request<{}, {}, Api.PingHost.Request>,
    res: Response<Api.PingHost.Response>
  ) => {
    res.send(await Handlers.pingHost(req.body.ip));
  };

  export const postMakeRequest = async (
    req: Request<{}, {}, Api.MakeRequest.Request>,
    res: Response<Api.MakeRequest.Response<any>>
  ) => {
    res.send(await Handlers.makeRequest(req.body.url, req.body.method, req.body.body));
  };

  export const postStatusChecks = async (
    req: Request<{}, {}, Api.StatusChecks.Request>,
    res: Response<Api.StatusChecks.Response>
  ) => {
    res.send(await Handlers.getStatusChecks(req.body.statusChecks));
  };
}
