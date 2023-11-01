import * as si from "systeminformation";
import { HwStorage } from "utils";

import { Utils } from "../utils/utils.js";

const getPartitionsWindows = (
  disks: HwStorage.Disk[],
  blocks: si.Systeminformation.BlockDevicesData[]
): HwStorage.Partition[] => {
  const localBlock = blocks.filter((block) => block.physical === "Local");

  return localBlock.map((block) => ({
    disk: disks.find((disk) => disk.device === block.device)!,
    identifier: block.uuid,
    uuid: block.uuid,
    label: block.label,
    size: block.size,
    mount: block.mount,
  }));
};

const getPartitionsLinux = (
  disks: HwStorage.Disk[],
  blocks: si.Systeminformation.BlockDevicesData[]
): HwStorage.Partition[] => {
  const partBlocks = blocks.filter((block) => block.type === "part");

  return partBlocks.map((block) => ({
    disk: disks.find((disk) => block.name.startsWith(disk.device.replace("/dev/", "")))!,
    identifier: block.uuid,
    uuid: block.uuid,
    label: block.label,
    size: Utils.bytesToMb(block.size),
    mount: block.mount,
  }));
};

const getStaticInfo = async (): Promise<HwStorage.Static> => {
  const [disksLayout, blocks] = await Promise.all([si.diskLayout(), si.blockDevices()]);

  const disks: HwStorage.Disk[] = disksLayout.map((disk) => ({
    device: disk.device,
    name: disk.name,
    brand: disk.vendor,
    type: disk.type,
    size: Utils.bytesToMb(disk.size),
  }));

  const totalSize = disksLayout.reduce((acc, disk) => acc + disk.size, 0);

  const partitions = Utils.platformIsWindows(process.platform)
    ? getPartitionsWindows(disks, blocks)
    : getPartitionsLinux(disks, blocks);

  return {
    size: Utils.bytesToMb(totalSize),
    disks,
    partitions: partitions,
  };
};

const staticInfo = await getStaticInfo();

export class Storage {
  public static staticInfo = staticInfo;

  public static current = async (): Promise<HwStorage.Load> => {
    const [fsSize] = await Promise.all([si.fsSize()]);

    return fsSize
      .map((fs) => ({
        size: Utils.bytesToMb(fs.size),
        used: Utils.bytesToMb(fs.used),
        available: Utils.bytesToMb(fs.available),
        use: fs.use,
        partition: staticInfo.partitions.find((partition) => partition.mount === fs.mount)!,
      }))
      .filter((fs) => fs.partition !== undefined);
  };
}
