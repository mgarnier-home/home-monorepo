import { Injectable } from '@angular/core';
import { ApiService } from './api.service';

@Injectable({
  providedIn: 'root',
})
export class HostsService extends ApiService {
  constructor() {
    super();
  }

  getHosts() {
    return this.get('hosts');
  }
}
