import { ArchiveStrategy } from '../utils/interfaces.js';
import { ArchiveApi } from './archiveApi';
import { ArchiveApiArchiver } from './archiveApiArchiver.js';
import { ArchiveApiTar } from './archiveApiTar.js';
import { ArchiveApiZip } from './archiveApiZip.js';

export const getArchiveApi = (strategy: ArchiveStrategy): ArchiveApi => {
  switch (strategy) {
    case ArchiveStrategy.ZIP:
      return new ArchiveApiZip();
    case ArchiveStrategy.TAR:
      return new ArchiveApiTar();
    case ArchiveStrategy.ARCHIVER:
      return new ArchiveApiArchiver();
    default:
      throw new Error(`Unknown archive strategy : ${strategy}`);
  }
};
