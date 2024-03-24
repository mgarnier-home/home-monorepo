import { SocketService } from 'src/app/services/socket.service';

import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [],
  templateUrl: './root.component.html',
  styleUrl: './root.component.scss',
})
export class RootComponent {
  constructor(private socketService: SocketService) {
    console.log('Test');
    socketService.connect();
  }
}
