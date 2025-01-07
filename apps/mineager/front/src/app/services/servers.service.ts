import { Injectable } from '@angular/core';
import { CreateServerRequest, Server } from '../interfaces/server.interfaces';
import { ApiService } from './api.service';

@Injectable({
  providedIn: 'root',
})
export class ServersService extends ApiService {
  constructor() {
    super();
  }

  public getServers(hostName: string): Promise<Server[]> {
    return this.get(`${this.getApiUrl()}/${hostName}/server`);
  }

  public getServer(hostName: string, serverName: string): Promise<Server> {
    return this.get(`${this.getApiUrl()}/${hostName}/server/${serverName}`);
  }

  public deleteServer(hostName: string, serverName: string): Promise<void> {
    return this.delete(`${this.getApiUrl()}/${hostName}/server/${serverName}`);
  }

  public createServer(hostName: string, serverDto: CreateServerRequest): Promise<Server> {
    return this.postJson(`${this.getApiUrl()}/${hostName}/server`, serverDto);
  }

  public startServer(hostName: string, serverName: string): Promise<void> {
    return this.postJson(`${this.getApiUrl()}/${hostName}/server/${serverName}/start`, {});
  }

  public stopServer(hostName: string, serverName: string): Promise<void> {
    return this.postJson(`${this.getApiUrl()}/${hostName}/server/${serverName}/stop`, {});
  }
}
