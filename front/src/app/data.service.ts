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

const getFirstAndLastCreatedAt = (
  posts: PostItem[],
): { last: string; first: string } => {
  if (posts.length === 0) {
    return { last: '', first: '' };
  }

  const { last, first } = posts.reduce(
    (acc, post) => ({
      last: post.created_at > acc.last ? post.created_at : acc.last,
      first: post.created_at < acc.first ? post.created_at : acc.first,
    }),
    { last: posts[0].created_at, first: posts[0].created_at },
  );

  return { last, first };
};

@Injectable({
  providedIn: 'root',
})
export class DataService {
  constructor(private api: ApiService) {}

  private posts = new BehaviorSubject<PostItem[]>([]);

  posts$ = this.posts.asObservable();

  private lastCreatedAt = new BehaviorSubject<string>('');

  lastCreatedAt$ = this.lastCreatedAt.asObservable();

  private firstCreatedAt = new BehaviorSubject<string>('');

  firstCreatedAt$ = this.firstCreatedAt.asObservable();

  userStatus = 'online';

  userStatusLastUpdate = new Date().getTime() - 100_000;

  private postsLoading = new BehaviorSubject<boolean>(false);

  postsLoading$ = this.postsLoading.asObservable();

  private setPosts(newPosts: PostItem[]) {
    this.posts.next(newPosts);

    const { last, first } = getFirstAndLastCreatedAt(newPosts);

    this.lastCreatedAt.next(last);
    this.firstCreatedAt.next(first);
  }

  private addPostAfter(newPost: PostItem) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, newPost];
    this.posts.next(updatedPosts);

    this.lastCreatedAt.next(newPost.created_at);
  }

  private addPostsAfter(newPosts: PostItem[]) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, ...newPosts];
    this.posts.next(updatedPosts);

    const { last } = getFirstAndLastCreatedAt(newPosts);

    this.lastCreatedAt.next(last);
  }

  private addPostsBefore(newPosts: PostItem[]) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...newPosts, ...currentPosts];
    this.posts.next(updatedPosts);

    const { first } = getFirstAndLastCreatedAt(newPosts);

    this.firstCreatedAt.next(first);
  }

  getPosts(options: { channelId: string; limit: number }) {
    this.postsLoading.next(true);

    const params = new HttpParams().append('limit', options.limit);

    this.api.get(`/v1/api/posts/${options.channelId}`, params).subscribe({
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

  getPostsAfter(options: { channelId: string; limit: number }) {
    return new Promise((resolve, reject) => {
      const last = this.lastCreatedAt.getValue();

      let params: HttpParams;

      if (last !== '') {
        params = new HttpParams()
          .append('limit', options.limit)
          .append('after', last);
      } else {
        params = new HttpParams().append('limit', options.limit);
      }

      this.api.get(`/v1/api/posts/${options.channelId}`, params).subscribe({
        next: (response) => {
          const posts = (response as GetPostsResp).posts.sort(
            (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          );

          this.addPostsAfter(posts);

          resolve(posts);
        },
        error: (err: unknown) => {
          console.error('error', err);

          reject(err);
        },
      });
    });
  }

  getPostsBefore(options: { channelId: string; limit: number }) {
    return new Promise((resolve, reject) => {
      const first = this.firstCreatedAt.getValue();

      let params: HttpParams;

      if (first !== '') {
        params = new HttpParams()
          .append('limit', options.limit)
          .append('before', first);
      } else {
        params = new HttpParams().append('limit', options.limit);
      }

      this.api.get(`/v1/api/posts/${options.channelId}`, params).subscribe({
        next: (response) => {
          const posts = (response as GetPostsResp).posts.sort(
            (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          );

          this.addPostsBefore(posts);

          resolve(posts);
        },
        error: (err: unknown) => {
          console.error('error', err);

          reject(err);
        },
      });
    });
  }

  sendPost(options: { message: string; channelId: string }) {
    this.api
      .post('/v1/api/posts', {
        message: options.message,
        channelId: options.channelId,
      })
      .subscribe({
        next: (response: unknown) => {
          console.log('post created', response);

          this.addPostAfter((response as { post: PostItem }).post);
        },
        error: (err) => {
          console.error('auth error', err);
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
            console.log('status updated', response);
          },

          error: (err) => {
            console.error('status updated error', err);
          },
        });

      this.userStatusLastUpdate = new Date().getTime();
    }
  }
}
