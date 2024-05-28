export interface Config {
  dataFilePath: string;
  serverPort: number;
  stacksToIgnore: string[];
}

export interface Proxy {
  sourceIP: string;
  destIP: string;
  activated: boolean;
}

export interface Host {
  name: string;
  ip: string;
  apiPort: number;
}

export interface AppData {
  proxies: Proxy[];
  hosts: Host[];
  domainName: string;
  defaultEntrypoints: string;

  defaultMiddlewares: string;
  defaultTls: {
    defaultCertResolver: string;
    defaultOptions: string;
  };
}

export interface StackInfos {
  path: string;
  host: Host;
  stack: string;
}

export interface TraefikService {
  host: Host;
  serviceName: string;
  portVariable: string;
  entryPoints?: string;
  middlewares?: string;
  rule?: string;
  tlsResolver?: string;
}
