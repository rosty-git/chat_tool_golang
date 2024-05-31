import { NgClass } from '@angular/common';
import { Component, inject, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { type Channel } from '../sidebar/sidebar.component';
import { ChannelsStore } from '../store/channels.store';

@Component({
  selector: 'app-direct-channels',
  standalone: true,
  imports: [NgIconComponent, NgClass],
  templateUrl: './directChannels.component.html',
  styleUrl: './directChannels.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class DirectChannelsComponent {
  readonly store = inject(ChannelsStore);

  @Input() channels: Channel[] = [];

  collapsed = true;

  active = '';

  mouseOn = '';

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: string) {
    this.active = id;

    this.store.setIsContactsActive();
    this.store.setActiveChannel(id);
  }

  mouseover(channelId: string) {
    this.mouseOn = channelId;
  }

  mouseout() {
    this.mouseOn = '';
  }
}
