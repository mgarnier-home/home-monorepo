export interface Config {
  webhookSecret: string;
  serverPort: number;
  nodeEnv: string;
  arm64Host: string;
  amd64Host: string;
  smeeUrl: string;
}
