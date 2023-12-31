import nodemailer from 'nodemailer';

import { config } from '../utils/config.js';

export class MailApi {
  private static getTransporter() {
    return nodemailer.createTransport({
      host: config.mailHost,
      port: config.mailPort,
      secure: config.mailSecure,
      auth: {
        user: config.mailLogin,
        pass: config.mailPassword,
      },
    });
  }

  private static async sendMail(to: string, subject: string, text: string) {
    try {
      if (!config.enableMail) throw new Error('Mail is not enabled');

      const transporter = MailApi.getTransporter();

      const response = await transporter.sendMail({
        from: config.mailLogin,
        to,
        subject,
        text,
      });

      transporter.close();

      return response;
    } catch (error) {
      console.error('Cannot send mail : ', error);
    }
  }

  static async sendFileInfos(uptoboxPassword: string, archivePassword: string, fileUrl: string, fileName: string) {
    return await MailApi.sendMail(
      config.infoMailTo,
      `Infos for : ${fileName}`,
      `Uptobox password: ${uptoboxPassword}\nArchive password : ${archivePassword}\nUrl : ${fileUrl}`
    );
  }

  static async sendError(error: string) {
    return await MailApi.sendMail(config.errorMailTo, `MySaveCLI Error`, error);
  }
}
