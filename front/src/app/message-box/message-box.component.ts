import { HttpParams } from '@angular/common/http';
import { Component, effect, inject } from '@angular/core';
import { getState } from '@ngrx/signals';

import { ApiService } from '../api.service';
import { DataService } from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';
import { ChannelsStore } from '../store/channels.store';

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
  readonly channelsStore = inject(ChannelsStore);

  posts$ = this.dataService.posts$;

  postItems: PostItem[] = [];

  constructor(
    private api: ApiService,
    private dataService: DataService,
  ) {
    effect(() => {
      this.posts$.subscribe((posts) => {
        this.postItems = posts;
      });

      const state = getState(this.channelsStore);

      const params = new HttpParams().append('limit', 20);

      this.api.get(`/v1/api/posts/${state.active}`, params).subscribe({
        next: (response) => {
          console.log('Get Channels', response);

          // this.postItems = (response as GetPostsResp).posts;

          this.dataService.setPosts((response as GetPostsResp).posts);
        },
        error: (err: unknown) => {
          console.error('error', err);
        },
      });
    });
  }
}
