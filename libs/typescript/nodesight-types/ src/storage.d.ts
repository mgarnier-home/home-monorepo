export namespace HwStorage {
  export type Disk = {
    device: string;
    name: string;
    brand: string;
    type: string;
    size: number;
  };
  export type Partition = {
    disk: Disk;
    identifier: string;
    uuid: string;
    label: string;
    size: number;
    mount: string;
  };
  export type Static = {
    size: number;
    disks: Disk[];
    partitions: Partition[];
  };
  export type Load = {
    partition: Partition;
    size: number;
    used: number;
    available: number;
    use: number;
  }[];
}
