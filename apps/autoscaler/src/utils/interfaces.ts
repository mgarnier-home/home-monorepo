export interface Config {
  webhookSecret: string;
  serverPort: number;
  nodeEnv: string;
  smeeUrl: string;
  autoscalerConfigPath: string;
  runnerImage: string;
  runnerOrgName: string;
  runnerRepoUrl: string;
  runnerAccessToken: string;
  sshUser: string;
  sshPrivateKey: string;
}

export interface DockerHost {
  label: string;
  ip: string;
  sshPort: number;
}

export interface AutoscalerConfig {
  autoscalerHosts: DockerHost[];
}
