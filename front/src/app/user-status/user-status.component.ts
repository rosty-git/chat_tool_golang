import { NgStyle } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-user-status',
  standalone: true,
  imports: [NgStyle],
  templateUrl: './user-status.component.html',
  styleUrl: './user-status.component.scss',
})
export class UserStatusComponent {
  @Input() userId: string = '';

  @Input() status: string = 'offline';

  @Input() color: string = '#1e325c';
}
