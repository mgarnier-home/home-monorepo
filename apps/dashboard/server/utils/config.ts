import fs from "node:fs";
import path from "node:path";

import { ServerConfig } from "./serverConfig";

const configFilePath = process.env.CONFIG_FILE || path.resolve(__dirname, "../../config.json");

const loadConfigFromFile = (): ServerConfig => {
  const config = fs.readFileSync(configFilePath, "utf-8");

  return JSON.parse(config) as ServerConfig;
};

const loadConfigFromEnv = (): ServerConfig => {
  const appConfPath = process.env.APP_CONF_PATH || "conf.yml";
  const appConfPathResolved = appConfPath.startsWith("/")
    ? appConfPath
    : path.resolve(__dirname, "../../", appConfPath);

  const iconsPath = process.env.ICONS_PATH || "./icons";
  const iconsPathResolved = iconsPath.startsWith("/") ? iconsPath : path.resolve(__dirname, "../../", iconsPath);

  const config: ServerConfig = {
    appConfPath: appConfPathResolved,
    iconsPath: iconsPathResolved,
    serverPort: Number(process.env.SERVER_PORT) || 3000,
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as ServerConfig;
