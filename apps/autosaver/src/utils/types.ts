export enum DirectoryType {
  local = 'local',
  cifs = 'cifs',
}

export type LocalDirectory = {
  type: DirectoryType.local;
  srcPath: string;
  destPath: string;
};

export type CifsDirectory = {
  type: DirectoryType.cifs;
  hostPath: string;
  mountPath: string;
  user: string;
  password: string;
  ip: string;
  port: number;
};

export type Directory = LocalDirectory | CifsDirectory;

export type BackupConfig = {
  directories: Directory[];
  rsync: Directory | undefined;
  mail:
    | {
        mailHost: string;
        mailPort: number;
        mailSecure: boolean;
        infoMailTo: string;
        errorMailTo: string;
        mailLogin: string;
        mailPassword: string;
      }
    | undefined;
  cronSchedule: string;
  deleteFiles: boolean;
};

export type Config = {
  serverPort: number;
  backupConfig: BackupConfig;
};
