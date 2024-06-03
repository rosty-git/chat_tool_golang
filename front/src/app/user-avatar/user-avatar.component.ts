import { NgClass } from '@angular/common';
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import * as CRC32 from 'crc-32';

import { DataService } from '../data.service';
import { UserStatusComponent } from '../user-status/user-status.component';

const stringToNumberInRange = (
  str: string,
  min: number,
  max: number,
): number => {
  const hash = CRC32.str(str);

  const positiveHash = Math.abs(hash);

  const range = max - min + 1;

  const result = (positiveHash % range) + min;

  return result;
};

const colors: string[] = [
  '#FF5733', // Bright Red
  '#33FF57', // Bright Green
  '#3357FF', // Bright Blue
  '#FF33A1', // Bright Pink
  '#FF8C33', // Orange
  '#FFD700', // Gold
  '#ADFF2F', // Green Yellow
  '#00CED1', // Dark Turquoise
  '#9400D3', // Dark Violet
  '#FF4500', // Orange Red
];

@Component({
  selector: 'app-user-avatar',
  standalone: true,
  imports: [UserStatusComponent, NgClass],
  templateUrl: './user-avatar.component.html',
  styleUrl: './user-avatar.component.scss',
})
export class UserAvatarComponent implements OnChanges {
  @Input() userName: string = '';

  @Input() userId: string = '';

  @Input() size: string = 'small';

  colorNumber = 0;

  userColor: string = '';

  status$ = 'offline';

  constructor(private dataService: DataService) {
    this.dataService.statuses$.subscribe((value) => {
      if (this.userId !== '' && Object.hasOwn(value, this.userId)) {
        this.status$ = value[this.userId];
      }
    });
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['userId']) {
      this.colorNumber = stringToNumberInRange(this.userName, 0, 9);

      this.userColor = colors[this.colorNumber];
    }
  }
}
