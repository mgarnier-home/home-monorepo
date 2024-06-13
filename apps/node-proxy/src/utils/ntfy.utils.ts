import { NtfyUtils } from 'ntfy-utils';

export const sendStartingServer = async (hostname: string, incomingIp: string, service: string) => {
  NtfyUtils.sendNotification(
    'Node proxy',
    `Host ${hostname} is starting\nRequest coming from ${incomingIp} proxied to ${service}`,
    ''
  );
};

export const sendStoppingServer = async (hostname: string) => {
  NtfyUtils.sendNotification('Node proxy', `Host ${hostname} is stopping`, '');
};
