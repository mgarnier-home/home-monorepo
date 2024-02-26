import { Table } from 'console-table-printer';
import generator from 'generate-password';

import { FolderToBackup } from './interfaces';

export class Utils {
  static generatePassword(length: number = 12) {
    return generator.generate({ length, numbers: true });
  }

  static printRecapTable(foldersToBackup: FolderToBackup[]) {
    const table = new Table({
      columns: [
        { name: 'name', title: 'Folder name' },
        { name: 'filesNb', title: 'Nb files' },
        { name: 'size', title: 'Zip size' },
      ],
    });

    for (const folderToBackup of foldersToBackup) {
      table.addRow(
        {
          name: folderToBackup.name,
          filesNb: folderToBackup.filesNb,
          size: folderToBackup.size !== undefined ? (folderToBackup.size / 1024 / 1024).toFixed(2) + ' MB' : '',
        },
        {
          color: folderToBackup.success ? 'green' : 'red',
        }
      );
    }

    table.printTable();
  }
}
