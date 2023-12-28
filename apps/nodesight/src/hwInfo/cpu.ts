import { HwCpu } from 'nodesight-types';
import * as si from 'systeminformation';

import { config } from '../utils/config.js';

const normalizeCpuModel = (cpuModel: string) => {
  return cpuModel
    .replace(/Processor/g, '')
    .replace(/[A-Za-z0-9]*-Core/g, '')
    .trim();
};

const getStaticInfo = async (): Promise<HwCpu.Static> => {
  const cpuInfo = await si.cpu();

  return {
    brand: cpuInfo.manufacturer,
    model: normalizeCpuModel(cpuInfo.brand),
    cores: cpuInfo.physicalCores,
    ecores: cpuInfo.efficiencyCores ?? 0,
    pcores: cpuInfo.performanceCores ?? 0,
    threads: cpuInfo.cores,
    frequency: cpuInfo.speed,
  };
};

export class Cpu {
  public static staticInfo = getStaticInfo();

  public static current = async (): Promise<HwCpu.Load> => {
    const cpu = await Cpu.staticInfo;

    const cpuLoad = (await si.currentLoad()).cpus;

    let temps: si.Systeminformation.CpuTemperatureData['cores'] = [];
    let mainTemp = 0;

    if (config.disableCpuTemps === false) {
      const cpuTemp = await si.cpuTemperature();
      const threadsPerCore = (cpu.threads - cpu.ecores) / cpu.pcores;
      temps = cpuTemp.cores.flatMap((temp, i) => (i < cpu.pcores ? Array(threadsPerCore).fill(temp) : temp));
      mainTemp = cpuTemp.main; // AVG temp of all cores, in case no per-core data is found
    }

    return {
      cores: cpuLoad.map(({ load }, i) => ({
        load,
        temp: temps[i] ?? mainTemp,
        core: i,
      })),
    };
  };
}
