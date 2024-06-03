import { NgClass } from '@angular/common';
import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { ChannelsState, DataService } from '../data.service';

@Component({
  selector: 'app-message-input',
  standalone: true,
  imports: [ReactiveFormsModule, NgClass],
  templateUrl: './message-input.component.html',
  styleUrl: './message-input.component.scss',
})
export class MessageInputComponent {
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

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  sendMessage() {
    if (this.messageForm.value.message) {
      this.dataService.sendPost({
        message: this.messageForm.value.message,
        channelId: this.channelsState$.active,
      });

      this.messageForm.reset();
    }
  }
}
