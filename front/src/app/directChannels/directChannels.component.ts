import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { Channel, ChannelsState, DataService } from '../data.service';
import { DirectChannelComponent } from '../direct-channel/direct-channel.component';
import { UserAvatarComponent } from '../user-avatar/user-avatar.component';

@Component({
  selector: 'app-direct-channels',
  standalone: true,
  imports: [
    NgIconComponent,
    NgClass,
    UserAvatarComponent,
    DirectChannelComponent,
  ],
  templateUrl: './directChannels.component.html',
  styleUrl: './directChannels.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class DirectChannelsComponent {
  channels$: Channel[] = [];

  collapsed = false;

  active = '';

  mouseOn = '';

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  constructor(private dataService: DataService) {
    this.dataService.channelsState$.subscribe((value) => {
      const directChannels: Channel[] = [];
      // eslint-disable-next-line no-restricted-syntax, guard-for-in
      for (const channelId in value.channels) {
        if (value.channels[channelId].type === 'D') {
          directChannels.push({
            id: value.channels[channelId].id,
            name: value.channels[channelId].name,
            unread: value.channels[channelId].unread,
            membersIds: value.channels[channelId].membersIds,
            index: value.channels[channelId].index,
          });
        }
      }

      this.channels$ = directChannels.sort((a, b) => a.index - b.index);

      this.channelsState$ = value;
    });

    this.dataService.directChannelCollapse$.subscribe((value) => {
      this.collapsed = value;
    });
  }

  onClick(): void {
    this.dataService.setDirectChannelCollapse(!this.collapsed);
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
