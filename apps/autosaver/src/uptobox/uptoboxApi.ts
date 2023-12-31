import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import { blobFromSync } from 'fetch-blob/from.js';
import fs from 'fs';
import path from 'path';

import { config } from '../utils/config.js';
import { UptoboxFile, UptoboxFolder } from '../utils/interfaces.js';
import { ProgressBlob } from '../utils/progressBlob.js';
import { Utils } from '../utils/utils.js';

class UptoboxApiError extends Error {
  public response: AxiosResponse | Response;

  constructor(message: string, response: AxiosResponse | Response) {
    super(message);

    this.response = response;
    this.name = 'UptoboxApiError';

    console.error('UptoboxApiError', response);
  }
}

export class UptoboxApi {
  static getUptoboxCode(uptoboxFileUrl: string): string {
    const parts = uptoboxFileUrl.split('/');

    return parts[parts.length - 1] as string;
  }

  static getFolderUptoboxPath(folderPath: string): string {
    if (folderPath.startsWith('/')) return folderPath.replace(/^\/+/, '//');
    else return `//${folderPath}`;
  }

  static async getUploadUrl(): Promise<string> {
    const requestConfig: AxiosRequestConfig = {
      method: 'GET',
      url: 'https://uptobox.com/api/upload',
      data: {
        token: config.uptoboxToken,
      },
    };

    const response = await axios(requestConfig);

    if (response?.data?.data?.uploadLink) {
      return `https:${response.data.data.uploadLink}`;
    } else {
      throw new UptoboxApiError('getUploadUrl : empty response', response);
    }
  }

  static async uploadFileWithRetry(filePath: string, nbRetries = 0, retry: number = 0): Promise<string> {
    try {
      const uploadedFile = await UptoboxApi.uploadFile(filePath);

      return uploadedFile;
    } catch (error) {
      if (retry < nbRetries) {
        console.log('caught error  : ', error);
        console.log(`retrying upload in 10s (${retry + 1}/${nbRetries})`);

        await Utils.timeout(10000);

        return UptoboxApi.uploadFileWithRetry(filePath, nbRetries, retry + 1);
      } else {
        throw error;
      }
    }
  }

  static async uploadFile(filePath: string): Promise<string> {
    const uploadUrl = await UptoboxApi.getUploadUrl();

    const fileExists = fs.existsSync(filePath);

    if (fileExists) {
      await Utils.timeout(5000);

      let intPercentage = 0;

      const blob = new ProgressBlob([blobFromSync(filePath)], {}, (streamActual, streamTotal) => {
        const actualPercentage = Math.floor((streamActual / streamTotal) * 100);
        const streamTotalMb = Math.floor(streamTotal / 1024 / 1024);
        const streamActualMb = Math.floor(streamActual / 1024 / 1024);

        if (actualPercentage > intPercentage) {
          intPercentage = actualPercentage;

          console.log(`uploading file : ${streamActualMb}Mb / ${streamTotalMb}Mb (${intPercentage}%)`);
        }
      });

      const formData = new FormData();

      formData.append('files', blob, path.basename(filePath));

      console.log('uploadFile : starting upload');

      const response = await fetch(uploadUrl, {
        method: 'POST',
        body: formData,
      });

      if (response.status === 200) {
        const data = (await response.json()) as any;

        console.log('uploadFile : response', data, path.basename(filePath));

        if (data.files.length > 0) {
          const uploadedFile = data.files[0];

          return uploadedFile.url;
        } else {
          throw new UptoboxApiError(`uploadFile : 0 file in the response`, response);
        }
      } else {
        throw new UptoboxApiError(`uploadFile : status is not 200 : ${response.status}`, response);
      }
    } else {
      throw new Error(`uploadFile : file does not exist ${filePath}`);
    }
  }

  static async updateFile(
    uptoboxFileUrl: string,
    data: { newPassword?: string; newName?: string; newDescription?: string; newPublic?: boolean }
  ) {
    const fileCode = UptoboxApi.getUptoboxCode(uptoboxFileUrl);

    const requestConfig: AxiosRequestConfig = {
      url: `https://uptobox.com/api/user/files`,
      method: 'PATCH',
      data: {
        token: config.uptoboxToken,
        file_code: fileCode,
      },
    };

    if (data.newPassword !== undefined) requestConfig.data.password = data.newPassword;
    if (data.newName !== undefined) requestConfig.data.new_name = data.newName;
    if (data.newDescription !== undefined) requestConfig.data.description = data.newDescription;
    if (data.newPublic !== undefined) requestConfig.data.public = data.newPublic ? '1' : '0';

    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return true;
      } else {
        throw new UptoboxApiError(
          `updateFile : api error : ${response.data.data} ${response.data.statusCode}`,
          response
        );
      }
    } else {
      throw new UptoboxApiError(`updateFile : status is not 200 : ${response.status}`, response);
    }
  }

  static async protectFile(uptoboxFileUrl: string, usePassword: boolean = false): Promise<string> {
    const password = Utils.generatePassword();

    if (usePassword) {
      await UptoboxApi.updateFile(uptoboxFileUrl, { newPassword: password, newPublic: false });
    } else {
      await UptoboxApi.updateFile(uptoboxFileUrl, { newPublic: false });
    }

    return password;
  }

  static async getFolder(folderPath: string): Promise<UptoboxFolder> {
    const uptoboxFolderPath = UptoboxApi.getFolderUptoboxPath(folderPath);

    const requestConfig: AxiosRequestConfig = {
      url: 'https://uptobox.com/api/user/files',
      method: 'GET',
      params: {
        token: config.uptoboxToken,
        path: uptoboxFolderPath,
        limit: 1,
      },
    };

    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return response.data.data.currentFolder;
      } else {
        if (response.data.data === 'Chemin introuvable' || response.data.data === 'Could not find current path') {
          return UptoboxApi.createFolder(folderPath);
        } else {
          throw new UptoboxApiError(
            `getFolder : api error : ${response.data.data} ${response.data.statusCode}`,
            response
          );
        }
      }
    } else {
      throw new UptoboxApiError(`getFolder : status is not 200 : ${response.status}`, response);
    }
  }

  static async createFolder(folderPath: string): Promise<UptoboxFolder> {
    const createFolderRecur = async (paths: string[], index: number): Promise<boolean> => {
      const pathsToAdd = [...paths].reverse().slice(paths.length - index);

      const folderPath = `//${pathsToAdd.reverse().join('/')}`;

      const requestConfig: AxiosRequestConfig = {
        method: 'PUT',
        url: 'https://uptobox.com/api/user/files',
        data: {
          token: config.uptoboxToken,
          path: folderPath,
          name: paths[index],
        },
      };

      const response = await axios(requestConfig);

      if (response.status === 200) {
        if (index === paths.length - 1) return true;
        else return createFolderRecur(paths, index + 1);
      } else {
        throw new UptoboxApiError(`createFolderRecur : status is not 200 : ${response.status}`, response);
      }
    };

    const uptoboxFolderPath = UptoboxApi.getFolderUptoboxPath(folderPath);

    const folderPaths = uptoboxFolderPath.split('/');

    const res = await createFolderRecur(folderPaths.slice(2), 0);

    if (res) {
      return UptoboxApi.getFolder(folderPath);
    } else {
      throw new Error(`createFolder : createFolderRecur error`);
    }
  }

  static async moveFile(uptoboxFileUrl: string, destFolder: UptoboxFolder) {
    const requestConfig: AxiosRequestConfig = {
      method: 'PATCH',
      url: 'https://uptobox.com/api/user/files',
      data: {
        token: config.uptoboxToken,
        file_codes: `${UptoboxApi.getUptoboxCode(uptoboxFileUrl)}`,
        destination_fld_id: destFolder.fld_id,
        action: 'move',
      },
    };
    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return true;
      } else {
        throw new UptoboxApiError(`moveFile : api error : ${response.data.data} ${response.data.statusCose}`, response);
      }
    } else {
      throw new UptoboxApiError(`moveFile : status is not 200 : ${response.status}`, response);
    }
  }

  static async listContent(folder: UptoboxFolder): Promise<{ files: UptoboxFile[]; folders: UptoboxFolder[] }> {
    const requestConfig: AxiosRequestConfig = {
      method: 'GET',
      url: 'https://uptobox.com/api/user/files',
      data: {
        token: config.uptoboxToken,
        path: folder.fld_name,
      },
    };

    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return {
          files: response.data.data.files,
          folders: response.data.data.folders,
        };
      } else {
        throw new UptoboxApiError(
          `listContent : api error : ${response.data.data} ${response.data.statusCose}`,
          response
        );
      }
    } else {
      throw new UptoboxApiError(`listContent : status is not 200 : ${response.status}`, response);
    }
  }

  static async deleteFiles(files: UptoboxFile[]) {
    const requestConfig: AxiosRequestConfig = {
      method: 'DELETE',
      url: 'https://uptobox.com/api/user/files',
      data: {
        token: config.uptoboxToken,
        file_codes: files.map((file) => file.file_code).join(','),
      },
    };

    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return true;
      } else {
        throw new UptoboxApiError(
          `deleteFiles : api error : ${response.data.data} ${response.data.statusCose}`,
          response
        );
      }
    } else {
      throw new UptoboxApiError(`deleteFiles : status is not 200 : ${response.status}`, response);
    }
  }

  static async deleteFolder(folder: UptoboxFolder) {
    const folderContent = await UptoboxApi.listContent(folder);

    if (folderContent.files.length > 0) {
      await UptoboxApi.deleteFiles(folderContent.files);
    }

    if (folderContent.folders.length > 0) {
      for (const folder of folderContent.folders) {
        await UptoboxApi.deleteFolder(folder);
      }
    }

    const requestConfig: AxiosRequestConfig = {
      method: 'DELETE',
      url: 'https://uptobox.com/api/user/files',
      data: {
        token: config.uptoboxToken,
        fld_id: folder.fld_id,
      },
    };

    const response = await axios(requestConfig);

    if (response.status === 200) {
      if (response.data.message === 'Success' && response.data.statusCode === 0) {
        return true;
      } else {
        throw new UptoboxApiError(
          `deleteFolder : api error : ${response.data.data} ${response.data.statusCose}`,
          response
        );
      }
    } else {
      throw new UptoboxApiError(`deleteFolder : status is not 200 : ${response.status}`, response);
    }
  }
}
