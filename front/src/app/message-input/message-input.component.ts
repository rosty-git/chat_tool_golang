import { NgClass } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { ApiService } from '../api.service';
import { DataService } from '../data.service';
import { ChannelsStore } from '../store/channels.store';

type PostItem = {
  id: string;
  message: string;
  created_at: string;
  user: {
    name: string;
  };
};

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
    private api: ApiService,
    private dataService: DataService,
  ) {}

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  sendMessage() {
    if (this.messageForm.value.message) {
      this.api
        .post('/v1/api/posts', {
          message: this.messageForm.value.message,
          channel: this.store.active(),
        })
        .subscribe({
          next: (response: unknown) => {
            console.log('post created', response);

            this.dataService.addPost((response as { post: PostItem }).post);

            this.messageForm.reset();
          },

          error: (err) => {
            console.error('auth error', err);
          },
        });
    }
  }
}
