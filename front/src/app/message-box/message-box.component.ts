import { HttpParams } from '@angular/common/http';
import { Component, effect, inject } from '@angular/core';
import { getState } from '@ngrx/signals';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';

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

export type GetPostsResp = {
  posts: PostItem[];
};

@Component({
  selector: 'app-message-box',
  standalone: true,
  templateUrl: './message-box.component.html',
  styleUrl: './message-box.component.scss',
  imports: [MessageItemComponent, InfiniteScrollModule],
})
export class MessageBoxComponent {
  readonly channelsStore = inject(ChannelsStore);

  posts$ = this.dataService.posts$;

  postItems: PostItem[] = [];

  throttle = 50;

  scrollDistance = 2;

  scrollUpDistance = 2;

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

          const posts = (response as GetPostsResp).posts.sort(
            (a, b) => new Date(a.created_at).getTime()
              - new Date(b.created_at).getTime(),
          );

          this.dataService.setPosts(posts);
        },
        error: (err: unknown) => {
          console.error('error', err);
        },
      });
    });
  }

  onUp() {
    console.log('scrolled up!', this);
  }

  onScrollDown() {
    console.log('scrolled down!', this);
  }
}
