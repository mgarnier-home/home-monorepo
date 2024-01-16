import { Utils } from 'utils';

import type { DockerMessage, SyslogMessage } from './interfaces';

export const getMessageKey = (msg: SyslogMessage): string => {
  const date = `${
    msg.date.getFullYear() //
  }-${
    Utils.padStart(msg.date.getMonth() + 1, 2) //
  }-${
    Utils.padStart(msg.date.getDate(), 2) //
  }`;
  if (msg.dockerMessage) {
    return `${msg.host}/${msg.dockerMessage.containerName}_${date}`;
  }

  return `${msg.host}/${date}`;
};

export const getDockerMessage = (msg: string): DockerMessage => {
  const [datePart, msgPart1] = msg.split('DCMSG:') as [string, string];

  const [, dateStr] = datePart.split('>') as [string, string];

  const [containerNameAndId, msgPart2] = msgPart1.split('[') as [string, string];

  const [containerName, containerId] = containerNameAndId.split(':') as [string, string];

  const [, message] = msgPart2.split(']:') as [string, string];

  return {
    date: new Date(`${new Date().getFullYear()} ${dateStr} UTC`),
    message: message.trim(),
    containerName,
    containerId,
  };
};
