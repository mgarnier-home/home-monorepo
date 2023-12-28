import { HwNetwork } from 'nodesight-types';
import * as si from 'systeminformation';
import { Utils } from 'utils';

const getStaticInfo = async (): Promise<HwNetwork.Static> => {
  const networkInfo = await si.networkInterfaces();

  const defaultNet = Array.isArray(networkInfo) ? networkInfo.find((net) => net.default)! : networkInfo;

  return {
    type: Utils.capFirst(defaultNet.type),
    interfaceSpeed: defaultNet.speed ?? 0,
    publicIp: defaultNet.ip4 ?? '',
    speed: defaultNet.speed ?? 0,
  };
};
(async () => await si.networkStats())(); // get correct values from the first api call

export class Network {
  public static staticInfo = getStaticInfo();

  public static current = async (): Promise<HwNetwork.Load> => {
    const networkStats = (await si.networkStats())[0];

    return {
      up: Math.round(networkStats?.tx_sec ?? 0),
      down: Math.round(networkStats?.rx_sec ?? 0),
    };
  };
}
