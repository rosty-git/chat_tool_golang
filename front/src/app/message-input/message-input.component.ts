import { NgClass } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';

import { ApiService } from '../api.service';
import { AppStore } from '../store/app.store';

@Component({
  selector: 'app-message-input',
  standalone: true,
  imports: [ReactiveFormsModule, NgClass],
  templateUrl: './message-input.component.html',
  styleUrl: './message-input.component.scss',
})
export class MessageInputComponent {
  readonly store = inject(AppStore);

  constructor(private api: ApiService) {}

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  sendMessage() {
    if (this.messageForm.value.message) {
      this.api
        .post('/v1/api/posts', {
          message: this.messageForm.value.message,
          channel: this.store.activeChannel(),
        })
        .subscribe({
          next: (response: unknown) => {
            console.log('auth ok', response);
            this.messageForm.reset();
          },

          error: (err) => {
            console.error('auth error', err);
          },
        });
    }
  }
}
