export namespace HwGpu {
  export type GpuInfo = {
    brand: string;
    model: string;
    memoryTotal: number;
    index: number;
  };

  export type Static = {
    gpus: GpuInfo[];
  };

  export type GpuLoad = {
    model: string;
    index: number;
    load: number;
    memoryUsed: number;
    powerDraw: number;
    temp: number;
  };

  export type Load = {
    gpus: GpuLoad[];
  };

  export namespace History {
    export type Value = Load & {
      timestamp: number;
    };
  }
}
