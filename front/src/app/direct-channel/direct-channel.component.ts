import { Component, Input } from '@angular/core';

import { UserAvatarComponent } from '../user-avatar/user-avatar.component';

@Component({
  selector: 'app-direct-channel',
  standalone: true,
  imports: [UserAvatarComponent],
  templateUrl: './direct-channel.component.html',
  styleUrl: './direct-channel.component.scss',
})
export class DirectChannelComponent {
  @Input() userId: string = '';

  @Input() userName: string = '';

  @Input() unread: number = 0;
}
