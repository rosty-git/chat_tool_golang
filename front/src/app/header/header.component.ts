import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { NgIconComponent } from '@ng-icons/core';

import { DataService } from '../data.service';
import { UserAvatarComponent } from '../user-avatar/user-avatar.component';
import { UserStatusComponent } from '../user-status/user-status.component';

@Component({
  selector: 'app-header',
  standalone: true,
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
  imports: [
    NgIconComponent,
    UserAvatarComponent,
    NgClass,
    UserStatusComponent,
    ReactiveFormsModule,
  ],
})
export class HeaderComponent {
  userId: string = '';

  userName: string = '';

  isOpen = false;

  status = 'online';

  isSearchFocused = false;

  isSearchHovered = false;

  searchForm = new FormGroup({
    text: new FormControl(''),
  });

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

  updateStatus(options: { status: string; manual: boolean }) {
    this.dataService.updateStatus(options);
  }

  onFocus() {
    if (!this.searchForm.value.text) {
      this.isSearchFocused = true;
    }
  }

  onBlur() {
    if (!this.searchForm.value.text) {
      this.isSearchFocused = false;
    }
  }

  onMouseEnter() {
    if (!this.searchForm.value.text) {
      this.isSearchHovered = true;
    }
  }

  onMouseLeave() {
    if (!this.searchForm.value.text) {
      this.isSearchHovered = false;
    }
  }

  searchKeyDown() {
    if (this.searchForm.value.text) {
      console.log(this.searchForm.value.text);

      this.dataService.searchPosts(this.searchForm.value.text);
    }
  }
}
