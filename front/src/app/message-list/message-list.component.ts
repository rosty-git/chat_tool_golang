import { Component } from '@angular/core';

import { ChannelsState, DataService } from '../data.service';
import { MessageBoxComponent } from '../message-box/message-box.component';
import { MessageInputComponent } from '../message-input/message-input.component';

@Component({
  selector: 'app-message-list',
  standalone: true,
  templateUrl: './message-list.component.html',
  styleUrl: './message-list.component.scss',
  imports: [MessageBoxComponent, MessageInputComponent],
})
export class MessageListComponent {
  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  constructor(private dataService: DataService) {
    this.dataService.channelsActive$.subscribe((value) => {
      this.channelsState$ = value;
    });
  }
}
