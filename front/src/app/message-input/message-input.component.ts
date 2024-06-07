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
      this.dataService
        .sendPost({
          message: this.messageForm.value.message,
          channelId: this.channelsState$.active,
        })
        .catch((err: HttpErrorResponse) => {
          if (err.status === 0) {
            this.dataService.addOfflineMessage({
              message: this.messageForm.value.message as string,
              channelId: this.channelsState$.active,
            });
          }
        })
        .finally(() => {
          this.messageForm.reset();
        });
    }
  }
}
