import { Component, inject, Input, OnInit, signal } from '@angular/core';
import { z } from 'zod';
import { StateService } from '../../services/state.service';
import { Service } from '../../models/dashboardConfig.schema';
import { ServiceState } from '../../models/dashboardState.schema';

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
