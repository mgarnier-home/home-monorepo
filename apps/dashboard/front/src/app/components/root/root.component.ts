import { Component, computed, inject } from '@angular/core';

import { SocketService } from '../../services/socket.service';
import { HostComponent } from '../host/host.component';

@Component({
  selector: 'app-root',
  imports: [HostComponent],
  templateUrl: './root.component.html',
  styleUrl: './root.component.scss',
})
export class RootComponent {
  private socketService: SocketService = inject(SocketService);

  public hosts = computed(() => {
    return this.socketService.dashboardConfig()?.hosts || [];
  });

  constructor() {
    // this.socketService.connect();
  }
}
