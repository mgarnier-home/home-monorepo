import * as si from "systeminformation";

import { HwRam } from "@mgarnier11/nodesight-types";

import { Utils } from "../utils/utils.js";

const getStaticInfo = async (): Promise<HwRam.Static> => {
  const [memInfo, memLayout] = await Promise.all([si.mem(), si.memLayout()]);

  return {
    size: Utils.bytesToMb(memInfo.total),
    layout: memLayout.map(({ manufacturer, type, clockSpeed, size }) => ({
      brand: manufacturer ?? "",
      type: type ?? "",
      frequency: clockSpeed ?? -1,
      size: Utils.bytesToMb(size),
    })),
  };
};

const staticInfo = await getStaticInfo();

export class Ram {
  public static staticInfo = staticInfo;

  public static current = async (): Promise<HwRam.Load> => {
    const memInfo = Utils.bytesToMb((await si.mem()).active);

    return { load: memInfo };
  };
}
