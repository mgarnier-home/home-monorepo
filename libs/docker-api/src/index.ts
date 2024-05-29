import { Container } from './interfaces';

export class Docker {
  private host: string;
  private port: number;

  constructor(host: string, port: number) {
    this.host = host;
    this.port = port;
  }

  async listContainers(): Promise<Container[]> {
    const response = await fetch(`http://${this.host}:${this.port}/containers/json`);

    return await response.json();
  }
}

export * from './interfaces';
