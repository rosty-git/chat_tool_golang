import { NgClass } from '@angular/common';
import { Component, Input } from '@angular/core';
import { NgIconComponent, provideIcons } from '@ng-icons/core';
import {
  heroChevronDownMini,
  heroChevronRightMini,
} from '@ng-icons/heroicons/mini';

import { type Contact } from '../sidebar/sidebar.component';

@Component({
  selector: 'app-contacts',
  standalone: true,
  imports: [NgIconComponent, NgClass],
  templateUrl: './contacts.component.html',
  styleUrl: './contacts.component.scss',
  viewProviders: [provideIcons({ heroChevronRightMini, heroChevronDownMini })],
})
export class ContactsComponent {
  @Input() contacts: Contact[] = [];

  collapsed = true;

  active = 0;

  onClick(): void {
    this.collapsed = !this.collapsed;
  }

  setActive(id: number) {
    this.active = id;
  }
}
