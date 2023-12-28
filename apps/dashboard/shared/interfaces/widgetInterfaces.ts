export namespace WidgetInterfaces {
  export enum Type {
    Stats = 'stats',
  }

  export namespace Stats {
    export type Options = {
      type: OptionsType;
      time: OptionsTime;
      history?: OptionsTimeHistory;
      gpuId?: number;
      gpuMaxMemory?: number;
      ramMaxMemory?: number;
    };

    export enum OptionsType {
      Cpu = 'cpu',
      Gpu = 'gpu',
      Ram = 'ram',
      Network = 'network',
      Storage = 'storage',
    }

    export enum OptionsTime {
      Live = 'live',
      History = 'history',
    }

    export enum OptionsTimeHistory {
      Last15Min = 'last15Min',
      LastHour = 'lastHour',
      LastDay = 'lastDay',
      LastWeek = 'lastWeek',
    }

    export type Widget = {
      name: string;
      type: Type.Stats;
      options: Options;
    };
  }

  export type Widget = Stats.Widget;
}
