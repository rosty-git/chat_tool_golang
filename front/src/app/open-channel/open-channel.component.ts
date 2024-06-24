import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-open-channel',
  standalone: true,
  imports: [],
  templateUrl: './open-channel.component.html',
  styleUrl: './open-channel.component.scss',
})
export class OpenChannelComponent {
  @Input() name: string = '';
}
