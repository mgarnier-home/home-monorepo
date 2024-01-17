import fs from 'fs';

export interface ServerConfig {
  storagePath: string;
  hostsMap: Record<string, string>;
  maxLogFileSize: number;
  serverPort: number;
  syslogPort: number;
  devMode: boolean;
}

export interface WriteStream {
  stream: fs.WriteStream;
  lastWrite: number;
  initialFileSize: number;
}
