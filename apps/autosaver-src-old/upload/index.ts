import path from 'path';

import { __dirname } from '../utils/config.js';
import { UploadStrategy } from '../utils/interfaces.js';
import { UploadApi } from './uploadApi.js';
import { UploadApiCp } from './uploadApiCp.js';
import { UploadApiUptobox } from './uploadApiUptobox.js';

export const getUploadApi = (
  strategy: UploadStrategy,
  uploadDestFolder: string
): { uploadApi: UploadApi; uploadDestFolder: string } => {
  switch (strategy) {
    case UploadStrategy.UPTOBOX:
      return { uploadApi: new UploadApiUptobox(), uploadDestFolder: uploadDestFolder.replace(/^[^a-zA-Z0-9]+/, '') };
    case UploadStrategy.CP:
      return {
        uploadApi: new UploadApiCp(),
        uploadDestFolder: uploadDestFolder.startsWith('/')
          ? uploadDestFolder
          : path.join(__dirname, '../../', uploadDestFolder),
      };
    default:
      throw new Error(`Unknown upload strategy : ${strategy}`);
  }
};
