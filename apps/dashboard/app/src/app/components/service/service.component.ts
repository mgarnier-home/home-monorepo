import { Component, inject, Input, OnInit, signal } from '@angular/core';
import { Service, serviceSchema } from '@shared/schemas/dashboard-config.schema';
import { z } from 'zod';
import { StateService } from '../../services/state.service';
import { ServiceState } from '@shared/schemas/dashboard-state.schema';

@Component({
  selector: 'app-service',
  imports: [],
  templateUrl: './service.component.html',
  styleUrl: './service.component.scss',
})
export class ServiceComponent implements OnInit {
  private stateService = inject(StateService);

  @Input() public service: Service | null = null;
  @Input() public serviceId: string = '';

  public serviceState = signal<ServiceState | null>(null);

  ngOnInit() {
    console.log('ServiceComponent ngOnInit');
    console.log('ServiceComponent service', this.service);

    this.serviceState = this.stateService.getServiceState(this.serviceId);
  }
}
