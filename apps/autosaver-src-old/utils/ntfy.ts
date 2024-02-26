import { NtfyUtils } from 'ntfy-utils';

import { FolderToBackup } from './interfaces';

export const sendBackupRecap = async (folders: FolderToBackup[]) => {
  const success = folders.every((f) => f.success);
  const failure = folders.every((f) => !f.success);

  const tags = success ? 'partying_face' : failure ? 'no_entry' : 'warning';

  const title = 'Autosaver';

  const recaptable = folders
    .map((f) => ({
      success: f.success,
      name: f.name,
      filesNb: f.filesNb,
      size: f.size !== undefined ? (f.size / 1024 / 1024).toFixed(2) + ' MB' : '',
    }))
    .map((r) => `${r.success ? '🟢' : '🔴'} ${r.name} : ${r.filesNb} files, ${r.size}`)
    .join('\n');

  const message = success
    ? `Backup done
    
    ${recaptable}`
    : failure
    ? `Backup failed`
    : `Backup partially failed
    
    ${recaptable}`;

  NtfyUtils.sendNotification(title, message, tags);
};

export const sendError = async (error: any) => {
  const title = 'Autosaver';
  const tags = 'bomb';
  const message = `An error occured : ${error}`;

  NtfyUtils.sendNotification(title, message, tags);
};
