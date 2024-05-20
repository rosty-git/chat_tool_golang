import { Component, inject } from '@angular/core';

import { MessageBoxComponent } from '../message-box/message-box.component';
import { MessageInputComponent } from '../message-input/message-input.component';
import { AppStore } from '../store/app.store';

@Component({
  selector: 'app-message-list',
  standalone: true,
  templateUrl: './message-list.component.html',
  styleUrl: './message-list.component.scss',
  imports: [MessageBoxComponent, MessageInputComponent],
})
export class MessageListComponent {
  readonly store = inject(AppStore);
}
