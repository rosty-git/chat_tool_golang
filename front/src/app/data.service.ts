import { HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { ApiService } from './api.service';

export type PostItem = {
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

  private postsLoading = new BehaviorSubject<boolean>(false);

  postsLoading$ = this.postsLoading.asObservable();

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

  getPosts(channelId: string, params: HttpParams) {
    this.postsLoading.next(true);

    this.api.get(`/v1/api/posts/${channelId}`, params).subscribe({
      next: (response) => {
        const posts = (response as GetPostsResp).posts.sort(
          (a, b) =>
            new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
        );

        this.setPosts(posts);

        this.postsLoading.next(false);
      },
      error: (err: unknown) => {
        console.error('error', err);

        this.postsLoading.next(false);
      },
    });
  }

  updateOnlineStatus() {
    if (
      new Date().getTime() - this.userStatusLastUpdate >
      USER_UPDATE_STATUS_INTERVAL
    ) {
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
