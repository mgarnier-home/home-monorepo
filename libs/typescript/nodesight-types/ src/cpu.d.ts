export namespace HwCpu {
  export type Static = {
    brand: string;
    model: string;
    cores: number;
    ecores: number;
    pcores: number;
    threads: number;
    frequency: number;
  };
  export type CoreLoad = {
    load: number;
    temp: number;
    core: number;
  };
  export type Load = {
    cores: CoreLoad[];
  };

  export namespace History {
    export type Value = CoreLoad & {
      timestamp: number;
    };
  }
}
