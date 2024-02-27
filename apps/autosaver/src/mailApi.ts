import { logger } from 'logger';
import nodemailer from 'nodemailer';

import { config } from './utils/config';

export namespace MailApi {
  const getTransporter = () =>
    nodemailer.createTransport({
      host: config.backupConfig.mail?.host,
      port: config.backupConfig.mail?.port,
      secure: config.backupConfig.mail?.secure,
      auth: {
        user: config.backupConfig.mail?.login,
        pass: config.backupConfig.mail?.password,
      },
    });

  const sendMail = async (subject: string, text: string, to?: string) => {
    try {
      if (!config.backupConfig.mail) throw new Error('Mail is not enabled');

      const transporter = getTransporter();

      const response = await transporter.sendMail({
        from: config.backupConfig.mail.login,
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

  export const sendFileInfos = async (archivePassword: string, fileUrl: string, fileName: string) => {
    return await sendMail(
      `Infos for : ${fileName}`,
      `Archive password : ${archivePassword}\nUrl : ${fileUrl}`,
      config.backupConfig.mail?.infoTo
    );
  };

  export const sendError = async (error: string) => {
    return await sendMail(`MySaveCLI Error`, error, config.backupConfig.mail?.errorTo);
  };
}
