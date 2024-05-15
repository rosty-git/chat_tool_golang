import { NgClass } from '@angular/common';
import { Component, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { type Channel } from '../sidebar/sidebar.component';

@Component({
  selector: 'app-channels',
  standalone: true,
  imports: [NgIconComponent, NgClass],
  templateUrl: './channels.component.html',
  styleUrl: './channels.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class ChannelsComponent {
  collapsed = true;

  active = 0;

  @Input() channels: Channel[] = [];

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: number) {
    this.active = id;
  }
}
