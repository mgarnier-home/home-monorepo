import { Express, Request, Response } from "express";
import { body, ValidationChain } from "express-validator";

import { Database } from "../db/database.js";
import { ApiUtils } from "./utils.js";

export const setupApiRoutes = <T, L>(
  app: Express,
  endpoint: string,
  postValidators: ValidationChain[],
  database: Database<T, L>
) => {
  app.post(`/api/:hostname/${endpoint}`, postValidators, ApiUtils.checkValidators, postEndpoint<T, L>(database));

  app.get(`/api/:hostname/${endpoint}`, get15Min<T, L>(database));
  app.get(`/api/:hostname/${endpoint}/hour`, get1Hour<T, L>(database));
  app.get(`/api/:hostname/${endpoint}/6hour`, get6Hour<T, L>(database));
  app.get(`/api/:hostname/${endpoint}/day`, get1Day<T, L>(database));
  app.get(`/api/:hostname/${endpoint}/week`, get1Week<T, L>(database));
};

const postEndpoint =
  <T, L>(database: Database<T, L>) =>
  async (req: Request<{}, {}, L>, res: Response<any, ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    console.log(`Received load for ${hostname}`);

    await database.logLoad(hostname.toLowerCase(), req.body);

    res.sendStatus(200);
  };

const get15Min =
  <T, L>(database: Database<T, L>) =>
  async (req: Request, res: Response<T[], ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    const result = await database.getDataLast15Min(hostname);

    res.send(result);
  };

const get1Hour =
  <T, L>(database: Database<T, L>) =>
  async (req: Request, res: Response<T[], ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    const result = await database.getDataLast1Hour(hostname);

    res.send(result);
  };

const get6Hour =
  <T, L>(database: Database<T, L>) =>
  async (req: Request, res: Response<T[], ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    const result = await database.getDataLast6Hour(hostname);

    res.send(result);
  };

const get1Day =
  <T, L>(database: Database<T, L>) =>
  async (req: Request, res: Response<T[], ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    const result = await database.getDataLast1Day(hostname);

    res.send(result);
  };

const get1Week =
  <T, L>(database: Database<T, L>) =>
  async (req: Request, res: Response<T[], ApiUtils.Hostname>) => {
    const { hostname } = res.locals;

    const result = await database.getDataLast1Week(hostname);

    res.send(result);
  };
