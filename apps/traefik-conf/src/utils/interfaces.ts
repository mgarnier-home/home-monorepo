export interface Config {
  traefikConfDirectory: string;
  composeDirectory: string;
  saveDataFile: string;
  serverPort: number;
  stacksToIgnore: string[];
  redirectionName: string;
  fbxAppToken: string;
  fbxAppId: string;
  fbxApiDomain: string;
  fbxHttpsPort: number;
  fbxApiBaseUrl: string;
  fbxApiVersion: string;
}

export interface Proxy {
  sourceIP: string;
  destIP: string;
  activated: boolean;
}

export interface Host {
  name: string;
  ip: string;
}

export interface AppData {
  proxies: Proxy[];
  hosts: Host[];
}

export interface StackInfos {
  path: string;
  host: Host;
  stack: string;
}

export interface TraefikService {
  host: Host;
  stackName: string;
  serviceName: string;
  portVariable: string;
}
