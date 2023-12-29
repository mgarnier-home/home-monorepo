import { config as cfg } from './config';

export class NtfyUtils {
  private static config = cfg;

  public static sendNotification = async (title: string, message: string, tags: string) => {
    if (!NtfyUtils.config?.ntfyTopic || !NtfyUtils.config?.ntfyServer) {
      console.log('Cant send a notification : no ntfy topic or server defined');

      return;
    }

    let ntfyUrl = `${NtfyUtils.config.ntfyProtocol || 'http'}://${
      //
      NtfyUtils.config.ntfyServer.endsWith('/') ? NtfyUtils.config.ntfyServer.slice(0, -1) : NtfyUtils.config.ntfyServer
    }/${
      //
      NtfyUtils.config.ntfyTopic.startsWith('/') ? NtfyUtils.config.ntfyTopic.slice(1) : NtfyUtils.config.ntfyTopic
    }`;

    try {
      await fetch(ntfyUrl, {
        method: 'POST',
        body: message,
        headers: {
          Title: title + ' - ' + new Date().toLocaleString('fr-FR', { timeZone: 'Europe/Paris' }),
          Tags: tags,
        },
      });
    } catch (error) {
      console.log('url', ntfyUrl);

      console.error('Error while sending a NTFY notification', error);
    }
  };
}
