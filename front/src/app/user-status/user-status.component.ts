import { NgStyle } from '@angular/common';
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';

@Component({
  selector: 'app-user-status',
  standalone: true,
  imports: [NgStyle],
  templateUrl: './user-status.component.html',
  styleUrl: './user-status.component.scss',
})
export class UserStatusComponent implements OnChanges {
  @Input() userId: string = '';

  @Input() status: string = 'offline';

  @Input() color: string = '#1e325c';

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['userId']) {
      console.log(this.userId);
    }
  }
}
