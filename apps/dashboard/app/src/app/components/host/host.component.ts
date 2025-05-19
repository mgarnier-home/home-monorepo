import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject, Input, OnInit, signal } from '@angular/core';
import { hostSchema } from '@shared/schemas/dashboard-config.schema';
import { z } from 'zod';
import { ServiceComponent } from '../service/service.component';
import { StateService } from '../../services/state.service';
import { HostState } from '@shared/schemas/dashboard-state.schema';

@Component({
  selector: 'app-host',
  imports: [CommonModule, ServiceComponent],
  templateUrl: './host.component.html',
  styleUrl: './host.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HostComponent implements OnInit {
  private stateService = inject(StateService);

  @Input() public host: z.infer<typeof hostSchema> | null = null;
  @Input() public hostId: string = '';

  public hostState = signal<HostState | null>(null);

  ngOnInit() {
    console.log('HostComponent ngOnInit');
    console.log('HostComponent host', this.host);

    this.hostState = this.stateService.getHostState(this.hostId);
  }
}
