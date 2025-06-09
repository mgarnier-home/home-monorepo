import { Injectable, signal, Signal, WritableSignal } from '@angular/core';
import { z } from 'zod';
import { HostState, ServiceState } from '../models/dashboardState.schema';

@Injectable({
  providedIn: 'root',
})
export class StateService {
  private hosts: Map<string, WritableSignal<HostState | null>> = new Map();
  private services: Map<string, WritableSignal<ServiceState | null>> =
    new Map();

  public getHostState(hostId: string): WritableSignal<HostState | null> {
    if (!this.hosts.has(hostId)) {
      this.hosts.set(hostId, signal<HostState | null>(null));
    }
    return this.hosts.get(hostId)!;
  }

  public updateHostState(hostId: string, state: HostState) {
    const hostState = this.getHostState(hostId);
    hostState.set(state);
  }

  public getServiceState(
    serviceId: string
  ): WritableSignal<ServiceState | null> {
    if (!this.services.has(serviceId)) {
      this.services.set(serviceId, signal<ServiceState | null>(null));
    }
    return this.services.get(serviceId)!;
  }

  public updateServiceState(serviceId: string, state: ServiceState) {
    const serviceState = this.getServiceState(serviceId);
    serviceState.set(state);
  }
}
