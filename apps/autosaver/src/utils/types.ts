export enum DirectoryType {
  local = 'local',
  cifs = 'cifs',
}

export enum ArchiveApiType {
  TAR = 'tar',
  ZIP = 'zip',
}

export type LocalDirectory = {
  type: DirectoryType.local;
  path: string;
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
  backupPath: string;
  backupDest: Directory;
  cifsDirectories: CifsDirectory[] | undefined;
  rsync: Directory | undefined;
  mail:
    | {
        host: string;
        port: number;
        secure: boolean;
        infoTo: string;
        errorTo: string;
        login: string;
        password: string;
      }
    | undefined;
  deleteFiles: boolean;
};

export type Config = {
  serverPort: number;
  backupConfigPath: string;
  archiveApiType: ArchiveApiType;
  keepAliveUrl: string;
};

export type DirectoryToBackup = {
  name: string;
  path: string;
  success?: boolean;
  filesNb?: number;
  size?: number;
  duration?: number;
};
