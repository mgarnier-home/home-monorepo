interface Port {
  PrivatePort: number;
  PublicPort?: number;
  Type: string;
}

interface Mount {
  Name: string;
  Source: string;
  Destination: string;
  Driver: string;
  Mode: string;
  RW: boolean;
  Propagation: string;
}

interface Network {
  IPAMConfig?: any;
  Links?: any;
  Aliases?: any;
  NetworkID: string;
  EndpointID: string;
  Gateway: string;
  IPAddress: string;
  IPPrefixLen: number;
  IPv6Gateway?: string;
  GlobalIPv6Address?: string;
  GlobalIPv6PrefixLen?: number;
  MacAddress: string;
}

interface NetworkSettings {
  Networks: {
    [key: string]: Network;
  };
}

interface HostConfig {
  NetworkMode: string;
}

export interface Container {
  Id: string;
  Names: string[];
  Image: string;
  ImageID: string;
  Command: string;
  Created: number;
  State: string;
  Status: string;
  Ports: Port[];
  Labels: {
    [key: string]: string;
  };
  SizeRw?: number;
  SizeRootFs?: number;
  HostConfig: HostConfig;
  NetworkSettings: NetworkSettings;
  Mounts: Mount[];
}
