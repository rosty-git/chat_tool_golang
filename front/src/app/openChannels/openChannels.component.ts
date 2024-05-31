import { NgClass } from '@angular/common';
import { Component, inject, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';
import { heroGlobeAlt } from '@ng-icons/heroicons/outline';

import { type Channel } from '../sidebar/sidebar.component';
import { ChannelsStore } from '../store/channels.store';

@Component({
  selector: 'app-open-channels',
  standalone: true,
  imports: [NgIconComponent, NgClass],
  templateUrl: './openChannels.component.html',
  styleUrl: './openChannels.component.scss',
  viewProviders: [
    provideIcons({ heroChevronRightMini, heroChevronDownMini, heroGlobeAlt }),
  ],
})
export class OpenChannelsComponent {
  readonly store = inject(ChannelsStore);

  collapsed = true;

  active = '';

  mouseOn = '';

  @Input() channels: Channel[] = [];

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: string) {
    this.active = id;

    this.store.setIsChannelsActive();
    this.store.setActiveChannel(id);
  }

  mouseover(channelId: string) {
    this.mouseOn = channelId;
  }

  mouseout() {
    this.mouseOn = ''
  }
}
