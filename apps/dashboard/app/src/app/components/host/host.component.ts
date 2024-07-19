import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
  selector: 'app-host',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './host.component.html',
  styleUrl: './host.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HostComponent {}
