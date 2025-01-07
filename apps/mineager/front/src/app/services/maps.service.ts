import { Injectable } from '@angular/core';
import { CreateMapRequest, Map } from '../interfaces/map.interfaces';
import { ApiService } from './api.service';

@Injectable({
  providedIn: 'root',
})
export class MapsService extends ApiService {
  constructor() {
    super();
  }

  public getMaps(): Promise<Map[]> {
    return this.get(`${this.getApiUrl()}/map`);
  }

  public getMap(name: string): Promise<Map> {
    return this.get(`${this.getApiUrl()}/map/${name}`);
  }

  public deleteMap(name: string): Promise<void> {
    return this.delete(`${this.getApiUrl()}/map/${name}`);
  }

  public createMap(mapDto: CreateMapRequest): Promise<Map> {
    return this.postJson(`${this.getApiUrl()}/map`, mapDto);
  }
}
