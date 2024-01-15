import type { DockerMessage, SyslogMessage } from '@shared/interfaces';

import fs from 'fs';
import path from 'path';
import { SimpleCache, Utils } from 'utils';

import { config } from './config';
import { Syslog } from './syslog';

import type { WriteStream } from './interfaces';

const syslog = new Syslog();
const fileStreams: { [key: string]: WriteStream } = {}; // key = full path of the file
const filePathCache = new SimpleCache<string>(10 * 1); // key = messageKey => value = full path of the file

const fileWatcher = setInterval(() => {
  console.log('Checking file streams');

  for (const key in fileStreams) {
    const fileStream = fileStreams[key];

    if (Date.now() - fileStream!.lastWrite > 30 * 100) {
      console.log(`Closing file stream for ${key}`);

      fileStream!.stream.end();
      delete fileStreams[key];
    }
  }
}, 30 * 100);

const getDockerMessage = (msg: string): DockerMessage => {
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

const handleSyslogMessage = (msg: SyslogMessage) => {
  // console.time('handleSyslogMessage');
  const isDockerMessage = msg.message.includes('DCMSG:');

  if (isDockerMessage) {
    const dockerMessage = getDockerMessage(msg.message);

    msg.dockerMessage = dockerMessage;
  }

  // console.timeEnd('handleSyslogMessage');

  // console.log(getLogFilePath(msg));

  writeMessage(msg);
};

const getMessageKey = (msg: SyslogMessage): string => {
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

const getLogFilePath = (msg: SyslogMessage): string => {
  const folderPath = getMessageKey(msg);
  const folderFullPath = path.join(config.storagePath, folderPath);

  let fullPath = filePathCache.get(folderPath);
  let fileNb = 0;

  if (!fullPath) {
    // console.log('Cache miss');

    // si le dossier existe, sinon fileNb sera 0
    if (fs.existsSync(folderFullPath)) {
      // on récupère tous les fichiers du dossier
      const fileNames = fs.readdirSync(folderFullPath);

      const filesDatas = fileNames
        .filter((fileName) => fileName.endsWith('.log') && fileName.includes(msg.dockerMessage?.containerId ?? 'no_id')) // on ne garde que les fichiers de log pour le container actuel
        .map((fileName) => ({
          // on récupère les stats de chaque fichier
          stats: fs.statSync(path.join(folderFullPath, fileName)),
          fileName,
        }));

      const filesDatasFiltered = filesDatas.filter(
        //
        (file) => file.stats.size < config.maxLogFileSize
      ); // on ne garde que les fichiers qui ne sont pas trop gros

      if (filesDatasFiltered.length > 0) {
        const file = filesDatasFiltered[filesDatasFiltered.length - 1]!; // on prend le dernier fichier

        const [, nbStr] = path.basename(file.fileName, '.log').split('_');

        fileNb = Number(nbStr);
      } else {
        fileNb = filesDatas.length;
      }
    }

    fullPath = path.join(folderFullPath, `${msg.dockerMessage?.containerId ?? 'no_id'}_${fileNb}.log`);

    // console.log(fullPath);

    filePathCache.set(folderPath, fullPath);
  } else {
    // console.log('Cache hit');
  }

  return fullPath;
};

const getFileStream = (msg: SyslogMessage): WriteStream => {
  const logFilePath = getLogFilePath(msg);

  if (!fileStreams[logFilePath]) {
    console.log(`Creating file stream for ${logFilePath}`);

    const logDir = path.dirname(logFilePath);

    if (!fs.existsSync(logDir)) {
      fs.mkdirSync(logDir, { recursive: true });
    }

    const stream = fs.createWriteStream(logFilePath, { flags: 'a' });

    fileStreams[logFilePath] = {
      stream,
      lastWrite: 0,
      initialFileSize: fs.existsSync(logFilePath) ? fs.statSync(logFilePath).size : 0,
    };
  }

  return fileStreams[logFilePath]!;
};

const writeMessage = (msg: SyslogMessage) => {
  const fileStream = getFileStream(msg);

  const msgStr = JSON.stringify(msg);

  fileStream.stream.write(msgStr + '\n');
  fileStream.lastWrite = Date.now();

  if (fileStream.initialFileSize + fileStream.stream.bytesWritten > config.maxLogFileSize) {
    filePathCache.invalidate(getMessageKey(msg));
  }
};

const main = async () => {
  console.log('Starting syslog server');
  console.log(config);

  console.log(new Date().toLocaleString());

  syslog.on('start', () => {
    console.log('Syslog server started');
  });

  syslog.on('message', (msg: SyslogMessage) => {
    // console.log('Message received', msg);

    console.time('handleSyslogMessage');
    handleSyslogMessage(msg);
    console.timeEnd('handleSyslogMessage');
  });

  syslog.on('error', (err) => {
    console.log('Error', err);
  });

  syslog.on('close', () => {
    console.log('Syslog server stopped');
  });

  syslog.start({ port: config.syslogPort });
};

main();
