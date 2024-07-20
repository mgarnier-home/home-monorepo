import { Cpu } from './cpu';
import { Gpu } from './gpu';
import { Network } from './network';
import { Ram } from './ram';
import { Storage } from './storage';

export namespace HwInfo {
  export const staticInfo = async () => {
    const data = await Promise.all([
      Storage.staticInfo,
      Network.staticInfo,
      Ram.staticInfo,
      Gpu.staticInfo,
      Cpu.staticInfo,
    ]);

    return {
      storage: data[0],
      network: data[1],
      ram: data[2],
      gpu: data[3],
      cpu: data[4],
    };
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
