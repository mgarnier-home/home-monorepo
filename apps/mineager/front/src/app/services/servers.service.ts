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
    return this.get(`/hosts/${hostName}/servers`);
  }

  public getServer(hostName: string, serverName: string): Promise<Server> {
    return this.get(`/hosts/${hostName}/servers/${serverName}`);
  }

  public deleteServer(hostName: string, serverName: string): Promise<void> {
    return this.delete(`/hosts/${hostName}/servers/${serverName}`);
  }

  public createServer(hostName: string, serverDto: CreateServerRequest): Promise<Server> {
    return this.postJson(`/hosts/${hostName}/servers`, serverDto);
  }

  public startServer(hostName: string, serverName: string): Promise<void> {
    return this.postJson(`/hosts/${hostName}/servers/${serverName}/start`, {});
  }

  public stopServer(hostName: string, serverName: string): Promise<void> {
    return this.postJson(`/hosts/${hostName}/servers/${serverName}/stop`, {});
  }
}
