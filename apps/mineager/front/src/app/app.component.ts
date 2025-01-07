import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MapsService } from './services/maps.service';
import { ServersService } from './services/servers.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  standalone: true,
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'front';

  constructor(private mapsService: MapsService, private serversService: ServersService) {}

  async test() {
    try {
      console.log(await this.mapsService.getMaps());
      const servers = await this.serversService.getServers('athena');

      console.log(await this.serversService.stopServer('athena', servers[0].name));
    } catch (error) {
      console.error(error);
    }
  }
}
