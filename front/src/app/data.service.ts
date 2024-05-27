import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { ApiService } from './api.service';

type PostItem = {
  id: string;
  message: string;
  created_at: string;
  user: {
    name: string;
  };
};

const USER_UPDATE_STATUS_INTERVAL = 60_000;

const getLastCreatedAt = (posts: PostItem[]): string => {
  if (posts.length === 0) {
    return '';
  }

  return posts.reduce(
    (latest, post) => (post.created_at > latest ? post.created_at : latest),
    posts[0].created_at,
  );
};

@Injectable({
  providedIn: 'root',
})
export class DataService {
  constructor(private api: ApiService) {}

  private posts = new BehaviorSubject<PostItem[]>([]);

  private lastCreatedAt = new BehaviorSubject<string>('');

  posts$ = this.posts.asObservable();

  lastCreatedAt$ = this.lastCreatedAt.asObservable();

  userStatus = 'online';

  userStatusLastUpdate = new Date().getTime() - 100_000;

  setPosts(newPosts: PostItem[]) {
    this.posts.next(newPosts);

    this.lastCreatedAt.next(getLastCreatedAt(newPosts));
  }

  addPost(newPost: PostItem) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, newPost];
    this.posts.next(updatedPosts);

    this.lastCreatedAt.next(newPost.created_at);
  }

  addPosts(newPosts: PostItem[]) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, ...newPosts];
    this.posts.next(updatedPosts);

    this.lastCreatedAt.next(getLastCreatedAt(newPosts));
  }

  updateOnlineStatus() {
    if (new Date().getTime() - this.userStatusLastUpdate > USER_UPDATE_STATUS_INTERVAL) {
      console.log('Update status');

      console.log(
        'Last update',
        this.userStatusLastUpdate - new Date().getTime(),
      );

      this.api
        .put('/v1/api/statuses', {
          status: 'online',
          manual: false,
          dnd_end_time: '2000-10-31T00:00:00.000-00:00',
        })
        .subscribe({
          next: (response: unknown) => {
            console.log('post created', response);
          },

          error: (err) => {
            console.error('auth error', err);
          },
        });

      this.userStatusLastUpdate = new Date().getTime();
    }
  }
}
