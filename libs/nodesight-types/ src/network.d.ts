export namespace HwNetwork {
  export type Static = {
    interfaceSpeed: number;
    speed: number;
    type: string;
    publicIp: string;
  };
  export type Load = {
    up: number;
    down: number;
  };

  export namespace History {
    export type Value = Load & {
      timestamp: number;
    };
  }
}
