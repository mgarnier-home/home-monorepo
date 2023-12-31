export interface UptoboxFolder {
  fld_id: number;
  usr_id: number;
  fld_parent_id: number;
  fld_desc: string;
  fld_name: string;
  fld_password: string;
  subFolderLimitReached: boolean;
  name: string;
  hash: string;
  fileCount: number;
  totalFileSize: number;
}

export interface UptoboxFile {
  file_code: string;
  file_created: string;
  file_descr: string;
  file_downloads: number;
  file_last_download: string;
  file_name: string;
  file_password: string;
  file_public: boolean;
  file_size: number;
  id: number;
  last_stream: string;
  nb_stream: number;
  transcoded: number;
}

export enum ArchiveStrategy {
  ZIP = 'ZIP',
  TAR = 'TAR',
  ARCHIVER = 'ARCHIVER',
}

export enum UploadStrategy {
  UPTOBOX = 'UPTOBOX',
  CP = 'CP',
}

export interface Config {
  archiveStrategy: ArchiveStrategy;
  uploadStrategy: UploadStrategy;
  enableMail: boolean;
  enableRsync: boolean;
  mailHost: string;
  mailPort: number;
  mailSecure: boolean;
  infoMailTo: string;
  errorMailTo: string;
  uptoboxToken: string;
  mailLogin: string;
  mailPassword: string;
  cronSchedule: string;
  serverPort: number;
  folderToBackup: string;
  uploadDestFolder: string;
  rsyncDestFolder: string;
  deleteFiles: boolean;
  smbConfig: SmbConfig;
}

export interface FolderToBackup {
  name: string;
  path: string;
  success?: boolean;
  filesNb?: number;
  size?: number;
}

export interface SmbConfig {
  hosts: {
    ip: string;
    user: string;
    password: string;
    mountRoot: string;
    mounts: {
      dest: string;
      src: string;
    }[];
  }[];
}
