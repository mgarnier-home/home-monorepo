import { HwCpu, HwGpu, HwNetwork, HwRam, HwStorage } from '@libs/nodesight-types';

export interface Config {
  hostname: string;
  serverPort: number;
  updateInterval: number;
  statsApiUrl: string;
  enableStatsApi: boolean;

  disableCpuTemps: boolean;
}

export interface Current {
  cpu: HwCpu.Load;
  ram: HwRam.Load;
  gpu: HwGpu.Load;
  network: HwNetwork.Load;
  storage: HwStorage.Load;
}
