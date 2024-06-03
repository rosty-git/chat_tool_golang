import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';
import { heroGlobeAlt } from '@ng-icons/heroicons/outline';

import { Channel, ChannelsState, DataService } from '../data.service';

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
  collapsed = true;

  active = '';

  mouseOn = '';

  channels$: Channel[] = [];

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  constructor(private dataService: DataService) {
    this.dataService.openChannels$.subscribe((value) => {
      this.channels$ = value;
    });

    this.dataService.channelsActive$.subscribe((value) => {
      this.channelsState$ = value;
    });
  }

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: string) {
    this.active = id;

    this.dataService.setOpenActive(id);
  }

  mouseover(channelId: string) {
    this.mouseOn = channelId;
  }

  mouseout() {
    this.mouseOn = '';
  }
}
