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
}

export interface DockerHost {
  label: string;
  ip: string;
  dockerPort: number;
}

export interface AutoscalerConfig {
  autoscalerHosts: DockerHost[];
}
