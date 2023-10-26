export namespace HwRam {
  export type Static = {
    size: number;
    layout: {
      brand: string;
      type: string;
      frequency: number;
      size: number;
    }[];
  };
  export type Load = { load: number };

  export namespace History {
    export type Value = Load & {
      timestamp: number;
    };
  }
}
