export interface Config {
  serverPort: number;
  dataFilePath: string;
}

export interface Cron {
  name: string;
  schedule: string;
  command: string;
}

export interface CronConfig {
  crons: Cron[];
}

export interface CronExecution {
  cronName: string;
  date: Date;
  result: string;
  success: boolean;
}
