import { HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { GlobalVariable } from '../global';
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

export type ChannelsState = {
  isOpenActive: boolean;
  isDirectActive: boolean;
  active: string;
};

export type Channel = {
  id: string;
  name: string;
  membersIds: string[];
};

export type ChannelsResp = {
  channels: Channel[];
};

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
  constructor(private api: ApiService) {
    this.channelsActive$.subscribe((value) => {
      this.getPosts({
        channelId: value.active,
        limit: GlobalVariable.POSTS_PAGE_SIZE,
      });
    });
  }

  private posts = new BehaviorSubject<PostItem[]>([]);

  posts$ = this.posts.asObservable();

  private lastCreatedAt = new BehaviorSubject<string>('');

  lastCreatedAt$ = this.lastCreatedAt.asObservable();

  private firstCreatedAt = new BehaviorSubject<string>('');

  firstCreatedAt$ = this.firstCreatedAt.asObservable();

  private userStatus = new BehaviorSubject<string>('online');

  userStatus$ = this.userStatus.asObservable();

  userStatusLastUpdate = new Date().getTime() - 100_000;

  private postsLoading = new BehaviorSubject<boolean>(false);

  postsLoading$ = this.postsLoading.asObservable();

  private statuses = new BehaviorSubject<Record<string, string>>({});

  statuses$ = this.statuses.asObservable();

  private channelsActive = new BehaviorSubject<ChannelsState>({
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  });

  channelsActive$ = this.channelsActive.asObservable();

  private openChannels = new BehaviorSubject<Channel[]>([]);

  openChannels$ = this.openChannels.asObservable();

  private directChannels = new BehaviorSubject<Channel[]>([]);

  directChannels$ = this.directChannels.asObservable();

  private userId = new BehaviorSubject<string>('');

  userId$ = this.userId.asObservable();

  private userName = new BehaviorSubject<string>('');

  userName$ = this.userName.asObservable();

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

  async getPosts(options: { channelId: string; limit: number }) {
    this.postsLoading.next(true);

    const params = new HttpParams().append('limit', options.limit);

    const resp = await this.api.get(
      `/v1/api/posts/${options.channelId}`,
      params,
    );

    const posts = (resp as GetPostsResp).posts.sort(
      (a, b) =>
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
    );

    this.setPosts(posts);

    this.postsLoading.next(false);
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

      this.api
        .get(`/v1/api/posts/${options.channelId}`, params)
        .then((resp) => {
          const posts = (resp as GetPostsResp).posts.sort(
            (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          );
          this.addPostsAfter(posts);
          resolve(posts);
        })
        .catch((err) => {
          console.error('error', err);
          reject(err);
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

      this.api
        .get(`/v1/api/posts/${options.channelId}`, params)
        .then((resp) => {
          const posts = (resp as GetPostsResp).posts.sort(
            (a, b) =>
              new Date(a.created_at).getTime() -
              new Date(b.created_at).getTime(),
          );
          this.addPostsBefore(posts);
          resolve(posts);
        })
        .catch((err) => {
          console.error(err);
          reject(err);
        });
    });
  }

  sendPost(options: { message: string; channelId: string }) {
    this.api
      .post('/v1/api/posts', {
        message: options.message,
        channelId: options.channelId,
      })
      .then((resp) => {
        console.log('post created', resp);

        this.addPostAfter((resp as { post: PostItem }).post);
      })
      .catch((err) => {
        console.error(err);
      });
  }

  updateStatus(options: {
    status: string;
    manual: boolean;
    dndEndTime?: string;
  }) {
    return new Promise((resolve, reject) => {
      const payload: {
        status: string;
        manual: boolean;
        dnd_end_time?: string;
      } = { status: options.status, manual: options.manual };
      if (options.dndEndTime) {
        payload.dnd_end_time = options.dndEndTime;
      }

      this.api
        .put<{ status: { status: string } }>('/v1/api/statuses', payload)
        .then((resp) => {
          console.log(payload.status)
          console.log(resp.status.status);

          this.userStatus.next(resp.status.status);

          resolve(resp);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }

  updateOnlineStatus() {
    console.log('updateOnlineStatus');

    return this.updateStatus({
      status: 'online',
      manual: false,
    });
  }

  setAwayStatus() {
    console.log('set away');

    return this.updateStatus({
      status: 'away',
      manual: false,
    });
  }

  setStatus(userId: string, status: string) {
    const currentStatuses = this.statuses.value;
    currentStatuses[userId] = status;
    this.statuses.next(currentStatuses);
  }

  setOpenActive(channelId: string): void {
    const newState = {
      active: channelId,
      isOpenActive: true,
      isDirectActive: false,
    };
    this.channelsActive.next(newState);
  }

  setDirectActive(channelId: string): void {
    const newState = {
      active: channelId,
      isDirectActive: true,
      isOpenActive: false,
    };
    this.channelsActive.next(newState);
  }

  async getStatusesForChannelMembers(channels: Channel[]) {
    // eslint-disable-next-line no-restricted-syntax
    for (const channel of channels) {
      // eslint-disable-next-line no-await-in-loop, no-restricted-syntax
      for (const userId of channel.membersIds) {
        // eslint-disable-next-line no-await-in-loop
        await this.getAndSetStatus(userId);
      }
    }
  }

  async getAndSetStatus(userId: string) {
    const { status } = (await this.api.get(`/v1/api/statuses/${userId}`)) as {
      status: { user_id: string; status: string };
    };

    this.setStatus(status.user_id, status.status);
  }

  async getOpenChannels() {
    const openParams = new HttpParams().append('channelType', 'O');

    const resp = await this.api.get('/v1/api/channels', openParams);

    const openChannels = (resp as ChannelsResp).channels;

    // eslint-disable-next-line no-restricted-syntax
    for (const channel of openChannels) {
      // eslint-disable-next-line no-await-in-loop
      const { members } = (await this.api.get(
        `/v1/api/channels/${channel.id}/members`,
      )) as { members: string[] };

      channel.membersIds = members;
    }

    this.openChannels.next(openChannels);

    this.getStatusesForChannelMembers(openChannels);
  }

  async getDirectChannels() {
    const openParams = new HttpParams().append('channelType', 'D');

    const resp = await this.api.get('/v1/api/channels', openParams);

    const directChannels = (resp as ChannelsResp).channels;

    // eslint-disable-next-line no-restricted-syntax
    for (const channel of directChannels) {
      // eslint-disable-next-line no-await-in-loop
      const { members } = (await this.api.get(
        `/v1/api/channels/${channel.id}/members`,
      )) as { members: string[] };

      channel.membersIds = members;
    }

    this.directChannels.next(directChannels);

    this.getStatusesForChannelMembers(directChannels);
  }

  getUser() {
    this.api
      .get('/v1/api/users/iam')
      .then((resp) => {
        const { user } = resp as { user: { id: string; name: string } };

        this.userId.next(user.id);

        this.userName.next(user.name);

        this.getAndSetStatus(user.id);
      })
      .catch((err) => console.error(err));
  }
}
