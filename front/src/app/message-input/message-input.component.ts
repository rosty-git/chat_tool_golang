import { NgClass } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { DataService } from '../data.service';
import { ChannelsStore } from '../store/channels.store';

@Component({
  selector: 'app-message-input',
  standalone: true,
  imports: [ReactiveFormsModule, NgClass],
  templateUrl: './message-input.component.html',
  styleUrl: './message-input.component.scss',
})
export class MessageInputComponent {
  readonly store = inject(ChannelsStore);

  constructor(
    private dataService: DataService,
  ) {}

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  sendMessage() {
    if (this.messageForm.value.message) {
      this.dataService.sendPost({
        message: this.messageForm.value.message,
        channelId: this.store.active(),
      });

      this.messageForm.reset();
    }
  }
}
