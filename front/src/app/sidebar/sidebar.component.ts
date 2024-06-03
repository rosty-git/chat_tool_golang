import { Component } from '@angular/core';

import { DirectChannelsComponent } from '../directChannels/directChannels.component';
import { OpenChannelsComponent } from '../openChannels/openChannels.component';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [OpenChannelsComponent, DirectChannelsComponent],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss',
})
export class SidebarComponent {}
