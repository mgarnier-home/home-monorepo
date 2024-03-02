import { logger } from 'logger';
import nodemailer from 'nodemailer';

import { getBackupConfig } from './utils/config';
import { BackupConfig } from './utils/types';

const mail = {
  withBackupConfig: (backupConfig: BackupConfig) => {
    const getTransporter = () =>
      nodemailer.createTransport({
        host: backupConfig.mail?.host,
        port: backupConfig.mail?.port,
        secure: backupConfig.mail?.secure,
        auth: {
          user: backupConfig.mail?.login,
          pass: backupConfig.mail?.password,
        },
      });

    const sendMail = async (subject: string, text: string, to?: string) => {
      try {
        if (!backupConfig.mail) throw new Error('Mail is not enabled');

        const transporter = getTransporter();

        const response = await transporter.sendMail({
          from: backupConfig.mail.login,
          to,
          subject,
          text,
        });

        transporter.close();

        return response;
      } catch (error) {
        logger.error('Cannot send mail : ', error);
      }
    };

    return {
      sendFileInfos: async (archivePassword: string, fileUrl: string, fileName: string) => {
        return await sendMail(
          `Infos for : ${fileName}`,
          `Archive password : ${archivePassword}\nUrl : ${fileUrl}`,
          backupConfig.mail?.infoTo
        );
      },

      sendError: async (error: string) => {
        return await sendMail(`MySaveCLI Error`, error, backupConfig.mail?.errorTo);
      },
    };
  },
};

export const mailApi = mail.withBackupConfig(getBackupConfig());
