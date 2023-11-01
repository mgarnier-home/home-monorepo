import dotenv from "dotenv";
import Fs from "node:fs";
import Path from "node:path";
import { fileURLToPath } from "node:url";

import { Config } from "./interfaces.js";

dotenv.config();

export const __filename = fileURLToPath(import.meta.url);

export const __dirname = Path.dirname(__filename);

const configFilePath = process.env.CONFIG_FILE || Path.resolve(__dirname, "../../config.json");

const loadConfigFromFile = (): Config => {
  const config = Fs.readFileSync(configFilePath, "utf-8");

  return JSON.parse(config) as Config;
};

const loadConfigFromEnv = (): Config => {
  const config: Config = {
    serverPort: parseInt(process.env.SERVER_PORT || "3000"),
    dbHost: process.env.DB_HOST || "",
    dbOrg: process.env.DB_ORG || "",
    dbBucket: process.env.DB_BUCKET || "",
    dbToken: process.env.DB_TOKEN || "",
  };

  return config;
};

export const config = (process.env.CONFIG_FILE ? loadConfigFromFile() : loadConfigFromEnv()) as Config;
