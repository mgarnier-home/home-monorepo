import { logger } from 'logger';
import { NodeSSH } from 'node-ssh';
import Ping from 'ping';
import { Client } from 'ssh2';
import Wol from 'wol';

import { Protocol, Service } from './interfaces.js';

export class ServerControl {
  public static getServerStatus(hostIP: string) {
    return new Promise<boolean>(async (resolve, reject) => {
      Ping.sys.probe(
        hostIP,
        (isAlive) => {
          resolve(!!isAlive);
        },
        { timeout: 0.5 }
      );
    });
  }

  public static async startServer(serverMac: string) {
    await Wol.wake(serverMac);
  }

  public static async startServerAndWaitForPing(serverMac: string, host: string) {
    return new Promise<boolean>(async (resolve, reject) => {
      const timeout = setTimeout(() => {
        resolve(false);
      }, 15000);

      await ServerControl.startServer(serverMac);

      const checkInterval = setInterval(async () => {
        const status = await ServerControl.getServerStatus(host);

        if (status) {
          clearInterval(checkInterval);
          clearTimeout(timeout);
          resolve(true);
        }
      }, 1000);
    });
  }

  private static getSshConnection<T>(
    serverIp: string,
    sshUsername: string,
    sshPassword: string,
    fn: (connection: Client) => Promise<T>,
    timeout: number = 1000
  ): Promise<T> {
    return new Promise<T>(async (resolve, reject) => {
      const sshCon = new NodeSSH();

      const dispose = () => {
        sshCon.dispose();
      };

      const res = (r: T) => {
        dispose();
        resolve(r);
      };

      const rej = (err: any) => {
        dispose();
        reject(err);
      };

      try {
        await sshCon.connect({
          host: serverIp,
          username: sshUsername,
          password: sshPassword,
        });

        const connection = sshCon.connection;

        if (!connection) {
          return rej('No connection');
        } else {
          const to = setTimeout(() => {
            connection.destroy();
            rej('Timeout');
          }, timeout);

          fn(connection)
            .then((r) => {
              clearTimeout(to);
              connection.destroy();
              res(r);
            })
            .catch(rej);
        }
      } catch (err) {
        return rej(err);
      }
    });
  }

  private static executeCommand(
    serverIp: string,
    sshUsername: string,
    sshPassword: string,
    command: string,
    maxExecutionTime: number = 1000
  ): Promise<string> {
    return ServerControl.getSshConnection<string>(
      serverIp,
      sshUsername,
      sshPassword,
      (connection) => {
        return new Promise<string>((resolve, reject) => {
          logger.info(`[${serverIp}] Executing command : ${command}`);

          connection.exec(command, (error, channel) => {
            logger.info(`[${serverIp}] Command executed`, error);

            if (error) {
              return reject(error);
            }

            const to = setTimeout(() => {
              channel.destroy();
              reject('Timeout');
            }, maxExecutionTime);

            let data = '';

            channel.on('data', (chunk: any) => {
              data += chunk.toString();

              logger.info(`[${serverIp}] Command output : ${chunk.toString()}`);
            });

            channel.on('exit', () => {
              logger.info(`[${serverIp}] Command exited`);

              clearTimeout(to);
              channel.destroy();
              resolve(data);
            });

            channel.on('close', () => {
              logger.info(`[${serverIp}] Command closed`);

              clearTimeout(to);
              channel.destroy();
              resolve(data);
            });
          });
        });
      },
      maxExecutionTime + 500
    );
  }

  static getServicesFromEnv(hostname: string): Service[] {
    const services: Service[] = [];

    for (const key in process.env) {
      if (key.endsWith('_PORT')) {
        const parts = key.split('_');

        const host = parts[0] ?? '';
        const service = parts.slice(1, parts.length - 1).join('-');
        const port = parseInt(process.env[key] ?? '');

        if (host.toUpperCase() === hostname.toUpperCase() && !Number.isNaN(port)) {
          services.push({
            name: service.toLowerCase(),
            proxyPort: port,
            protocol: Protocol.TCP,
          });
        }
      }
    }
    return services;
  }

  static async getServicesFromDocker(serverIp: string, dockerPort: number): Promise<Service[]> {
    return [];
  }

  static async stopServer(serverHost: string, sshUsername: string, sshPassword: string) {
    await ServerControl.executeCommand(serverHost, sshUsername, sshPassword, 'sudo pm-suspend &');
  }
}
