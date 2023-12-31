export interface UploadApi {
  cleanOldFolders(backupFolder: string): Promise<void>;

  uploadFile(filePath: string, destFolder: string): Promise<string>;

  protectFile(fileLocation: string): Promise<string>;
}
