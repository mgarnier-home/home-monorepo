import { HwRam } from 'nodesight-types';
import * as si from 'systeminformation';
import { Utils } from 'utils';

const getStaticInfo = async (): Promise<HwRam.Static> => {
  const [memInfo, memLayout] = await Promise.all([si.mem(), si.memLayout()]);

  return {
    size: Utils.convert(memInfo.total, 'B', 'MB'),
    layout: memLayout.map(({ manufacturer, type, clockSpeed, size }) => ({
      brand: manufacturer ?? '',
      type: type ?? '',
      frequency: clockSpeed ?? -1,
      size: Utils.convert(size, 'B', 'MB'),
    })),
  };
};

export class Ram {
  public static staticInfo = getStaticInfo();

  public static current = async (): Promise<HwRam.Load> => {
    const memInfo = Utils.convert((await si.mem()).active, 'B', 'MB');

    return { load: memInfo };
  };
}
