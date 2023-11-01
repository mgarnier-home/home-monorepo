import dotenv from "dotenv";
import express from "express";
import http from "http";
import ViteExpress from "vite-express";

import { ServerRoutes } from "@shared/routes";

import { RestControllers } from "./restControllers";
import { config } from "./utils/config";
import { bindSocketIOServer } from "./wsControllers";

const log = (...args: any[]) => {
  console.log(`[API]`, ...args);
};

dotenv.config();

if (process.env.NODE_ENV === "production") ViteExpress.config({ mode: "production" });

const expressApp = express();
const httpServer = http.createServer(expressApp);

expressApp.use("/", (req, res, next) => {
  log(`${req.method} ${req.url}`);
  next();
});

expressApp.use(express.json());
expressApp.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", "*");
  next();
});
expressApp.use(express.urlencoded({ extended: true }));

expressApp.use((err: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error(err.stack);

  res.status(500).send(err.message);
});

expressApp.use(express.static(config.iconsPath));

expressApp.get(ServerRoutes.CONF, RestControllers.getConf);
expressApp.post(ServerRoutes.PING_HOST, RestControllers.postPingHost);
expressApp.post(ServerRoutes.MAKE_REQUEST, RestControllers.postMakeRequest);

ViteExpress.bind(expressApp, httpServer);

bindSocketIOServer(httpServer);

httpServer.listen(config.serverPort, () => {
  console.log(`Server listening on port ${config.serverPort}`);
});