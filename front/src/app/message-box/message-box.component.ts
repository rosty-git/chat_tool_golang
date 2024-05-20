import { HttpParams } from '@angular/common/http';
import { Component, effect, inject } from '@angular/core';
import { getState } from '@ngrx/signals';

import { ApiService } from '../api.service';
import { MessageItemComponent } from '../message-item/message-item.component';
import { AppStore } from '../store/app.store';

type PostItem = {
  id: string;
  message: string;
  created_at: string;
  user: {
    name: string;
  };
};

type GetPostsResp = {
  posts: PostItem[];
};

@Component({
  selector: 'app-message-box',
  standalone: true,
  templateUrl: './message-box.component.html',
  styleUrl: './message-box.component.scss',
  imports: [MessageItemComponent],
})
export class MessageBoxComponent {
  readonly store = inject(AppStore);

  postItems: PostItem[] = [];

  constructor(private api: ApiService) {
    effect(() => {
      const state = getState(this.store);
      console.log('books state changed', state.activeChannel);

      const params = new HttpParams().append('limit', 20);

      this.api.get(`/v1/api/posts/${state.activeChannel}`, params).subscribe({
        next: (response) => {
          console.log('Get Channels', response);

          this.postItems = (response as GetPostsResp).posts;
        },

        error: (err: unknown) => {
          console.error('error', err);
        },
      });
    });
  }
}
