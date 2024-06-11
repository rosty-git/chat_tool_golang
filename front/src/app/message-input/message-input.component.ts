import { NgClass } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
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
      const message = this.messageForm.value.message as string;

      this.dataService
        .sendPost({
          message,
          channelId: this.channelsState$.active,
        })
        .catch((err: HttpErrorResponse) => {
          if (err.status === 0) {
            this.dataService.addOfflineMessage({
              message,
              channelId: this.channelsState$.active,
            });
          }
        });

      this.messageForm.reset();
    }
  }

  // eslint-disable-next-line class-methods-use-this
  handleEnterKey(event: Event) {
    const keyboardEvent = event as KeyboardEvent;

    if (keyboardEvent.key === 'Enter' && !keyboardEvent.shiftKey) {
      event.preventDefault();
    }
  }
}
