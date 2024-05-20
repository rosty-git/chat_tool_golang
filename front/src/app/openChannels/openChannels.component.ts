import { NgClass } from '@angular/common';
import { Component, inject, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { type Channel } from '../sidebar/sidebar.component';
import { AppStore } from '../store/app.store';

@Component({
  selector: 'app-open-channels',
  standalone: true,
  imports: [NgIconComponent, NgClass],
  templateUrl: './openChannels.component.html',
  styleUrl: './openChannels.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class OpenChannelsComponent {
  readonly store = inject(AppStore);

  collapsed = true;

  active = '';

  @Input() channels: Channel[] = [];

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: string) {
    this.active = id;

    this.store.setIsChannelsActive();
    this.store.setActiveChannel(id);
  }
}
