import { Component, inject } from '@angular/core';

import { SocketService } from '../../services/socket.service';

@Component({
    selector: 'app-root',
    imports: [],
    templateUrl: './root.component.html',
    styleUrl: './root.component.scss'
})
export class RootComponent {
  private socketService: SocketService = inject(SocketService);

  public dashboardConfig = this.socketService.dashboardConfig;

  constructor() {
    this.socketService.connect();
  }
}
