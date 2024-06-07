import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { Channel, ChannelsState, DataService } from '../data.service';
import { UserAvatarComponent } from '../user-avatar/user-avatar.component';

@Component({
  selector: 'app-direct-channels',
  standalone: true,
  imports: [NgIconComponent, NgClass, UserAvatarComponent],
  templateUrl: './directChannels.component.html',
  styleUrl: './directChannels.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class DirectChannelsComponent {
  channels$: Channel[] = [];

  collapsed = true;

  active = '';

  mouseOn = '';

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  constructor(private dataService: DataService) {
    this.dataService.channelsActive$.subscribe((value) => {
      const directChannels: Channel[] = [];
      // eslint-disable-next-line no-restricted-syntax, guard-for-in
      for (const channelId in value.channels) {
        if (value.channels[channelId].type === 'D') {
          directChannels.push({
            id: value.channels[channelId].id,
            name: value.channels[channelId].name,
            unread: value.channels[channelId].unread,
            membersIds: value.channels[channelId].membersIds,
          });
        }
      }
      this.channels$ = directChannels;

      this.channelsState$ = value;
    });
  }

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: string) {
    this.active = id;

    this.dataService.setDirectActive(id);
  }

  mouseover(channelId: string) {
    this.mouseOn = channelId;
  }

  mouseout() {
    this.mouseOn = '';
  }
}
