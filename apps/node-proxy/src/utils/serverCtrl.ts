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
          console.log(`[${serverIp}] Executing command : ${command}`);

          connection.exec(command, (error, channel) => {
            console.log(`[${serverIp}] Command executed`, error);

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

              console.log(`[${serverIp}] Command output : ${chunk.toString()}`);
            });

            channel.on('exit', () => {
              console.log(`[${serverIp}] Command exited`);

              clearTimeout(to);
              channel.destroy();
              resolve(data);
            });

            channel.on('close', () => {
              console.log(`[${serverIp}] Command closed`);

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

  static async getServicesFromSSH(serverIp: string, sshUsername: string, sshPassword: string): Promise<Service[]> {
    const str = await ServerControl.executeCommand(
      serverIp,
      sshUsername,
      sshPassword,
      `docker ps --format "{{.Names}}|{{.ID}}|{{.Ports}}"`
    );

    const lines = str.split('\n');

    const services: Service[] = [];

    for (const line of lines) {
      const parts = line.split('|');

      if (parts.length === 3) {
        const name = parts[0] ?? '';
        const forwardedPorts = parts[2]?.split(',').filter((p) => p.includes('->'));

        for (const forwardedPort of forwardedPorts ?? []) {
          const protocol = forwardedPort.includes('tcp') ? Protocol.TCP : Protocol.UDP;
          const forwardedPortInfos = forwardedPort.split('->');

          const port = forwardedPortInfos[0]?.split(':')[1];

          services.push({
            name,
            proxyPort: parseInt(port ?? ''),
            protocol,
          });
        }
      }
    }

    return services;
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

  static async stopServer(serverHost: string, sshUsername: string, sshPassword: string) {
    await ServerControl.executeCommand(serverHost, sshUsername, sshPassword, 'sudo pm-suspend &');
  }
}