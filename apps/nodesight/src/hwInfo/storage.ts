import { HwStorage } from 'nodesight-types';
import * as si from 'systeminformation';
import { Utils } from 'utils';

const getPartitionsWindows = (
  disks: HwStorage.Disk[],
  blocks: si.Systeminformation.BlockDevicesData[]
): HwStorage.Partition[] => {
  const localBlock = blocks.filter((block) => block.physical === 'Local');

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
  const partBlocks = blocks.filter((block) => block.type === 'part');

  return partBlocks.map((block) => ({
    disk: disks.find((disk) => block.name.startsWith(disk.device.replace('/dev/', '')))!,
    identifier: block.uuid,
    uuid: block.uuid,
    label: block.label,
    size: Utils.convert(block.size, 'B', 'MB'),
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
    size: Utils.convert(disk.size, 'B', 'MB'),
  }));

  const totalSize = disksLayout.reduce((acc, disk) => acc + disk.size, 0);

  const partitions = Utils.platformIsWindows(process.platform)
    ? getPartitionsWindows(disks, blocks)
    : getPartitionsLinux(disks, blocks);

  return {
    size: Utils.convert(totalSize, 'B', 'MB'),
    disks,
    partitions: partitions,
  };
};

export class Storage {
  public static staticInfo = getStaticInfo();

  public static current = async (): Promise<HwStorage.Load> => {
    const staticInfo = await Storage.staticInfo;

    const [fsSize] = await Promise.all([si.fsSize()]);

    return fsSize
      .map((fs) => ({
        size: Utils.convert(fs.size, 'B', 'MB'),
        used: Utils.convert(fs.used, 'B', 'MB'),
        available: Utils.convert(fs.available, 'B', 'MB'),
        use: fs.use,
        partition: staticInfo.partitions.find((partition) => partition.mount === fs.mount)!,
      }))
      .filter((fs) => fs.partition !== undefined);
  };
}
