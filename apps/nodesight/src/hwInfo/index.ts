import { HwCpu, HwGpu, HwNetwork, HwRam, HwStorage } from "utils";

import { Cpu } from "./cpu.js";
import { Gpu } from "./gpu.js";
import { Network } from "./network.js";
import { Ram } from "./ram.js";
import { Storage } from "./storage.js";

export namespace HwInfo {
  export const staticInfo = {
    storage: Storage.staticInfo,
    network: Network.staticInfo,
    ram: Ram.staticInfo,
    gpu: Gpu.staticInfo,
    cpu: Cpu.staticInfo,
  };

  export const current = async () => {
    const data = await Promise.all([Storage.current(), Network.current(), Ram.current(), Gpu.current(), Cpu.current()]);

    return {
      storage: data[0],
      network: data[1],
      ram: data[2],
      gpu: data[3],
      cpu: data[4],
    };
  };
}
