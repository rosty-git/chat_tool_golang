import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { NgIconComponent } from '@ng-icons/core';

import { DataService } from '../data.service';
import { UserAvatarComponent } from '../user-avatar/user-avatar.component';
import { UserStatusComponent } from '../user-status/user-status.component';

@Component({
  selector: 'app-header',
  standalone: true,
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
  imports: [NgIconComponent, UserAvatarComponent, NgClass, UserStatusComponent],
})
export class HeaderComponent {
  userId: string = '';

  userName: string = '';

  isOpen = false;

  status = 'online';

  toggleDropdown() {
    this.isOpen = !this.isOpen;
  }

  constructor(private dataService: DataService) {
    this.dataService.userId$.subscribe((value) => {
      this.userId = value;
    });

    this.dataService.userName$.subscribe((value) => {
      this.userName = value;
    });

    this.dataService.userStatus$.subscribe((value) => {
      this.status = value;
    });
  }

  updateStatus(options: {status: string, manual: boolean}) {
    this.dataService.updateStatus(options)
  }
}
