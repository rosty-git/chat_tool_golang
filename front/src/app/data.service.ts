import { HttpErrorResponse, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { environment } from '../environments/environment';
import { ApiService } from './api.service';

export type PostItem = {
  id: string;
  message: string;
  created_at: string;
  channel_id: string;
  user: {
    name: string;
  };
  offline?: boolean;
  frontId?: string;
  files: FrontFile[];
};

export type GetPostsResp = {
  posts: PostItem[];
};

export type ChannelsState = {
  isOpenActive: boolean;
  isDirectActive: boolean;
  active: string;
  channels?: {
    [key: string]: Channel;
  };
};

export type Channel = {
  id: string;
  name: string;
  membersIds: string[];
  unread: number;
  type?: 'O' | 'D';
  lastCreatedAt?: string;
  firstCreatedAt?: string;
  posts?: PostItem[];
};

export type ChannelsResp = {
  channels: Channel[];
};

export type FrontFile = {
  id?: string;
  name: string;
  size: number;
  type: string;
  ext: string;
  blobUrl?: string;
  s3_key?: string;
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
  constructor(private api: ApiService) {}

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
    channels: {},
  });

  channelsActive$ = this.channelsActive.asObservable();

  private userId = new BehaviorSubject<string>('');

  userId$ = this.userId.asObservable();

  private userName = new BehaviorSubject<string>('');

  userName$ = this.userName.asObservable();

  offlineMessages: PostItem[] = [];

  offlineMessagesSending = false;

  private openChannelCollapse = new BehaviorSubject<boolean>(true);

  openChannelCollapse$ = this.openChannelCollapse.asObservable();

  private directChannelCollapse = new BehaviorSubject<boolean>(true);

  directChannelCollapse$ = this.directChannelCollapse.asObservable();

  private searchedPosts = new BehaviorSubject<PostItem[]>([]);

  searchedPosts$ = this.searchedPosts.asObservable();

  private showSearchedPosts = new BehaviorSubject<boolean>(false);

  showSearchedPosts$ = this.showSearchedPosts.asObservable();

  private showChannelSearchModal = new BehaviorSubject<boolean>(false);

  showChannelSearchModal$ = this.showChannelSearchModal.asObservable();

  private setPosts(channelId: string, newPosts: PostItem[]) {
    const currentChannelsState = this.channelsActive.getValue();

    const { last, first } = getFirstAndLastCreatedAt(newPosts);

    if (!currentChannelsState.channels) {
      currentChannelsState.channels = {};
    }

    currentChannelsState.channels[channelId] = {
      ...currentChannelsState.channels[channelId],
      posts: newPosts,
      firstCreatedAt: first,
      lastCreatedAt: last,
    };

    this.channelsActive.next(currentChannelsState);
  }

  private addPostAfter(options: { channelId: string; post: PostItem }) {
    const currentChannelsState = this.channelsActive.getValue();

    if (currentChannelsState.channels?.[options.channelId]) {
      const existingPosts =
        currentChannelsState.channels[options.channelId].posts;

      const updatedPosts = [...existingPosts!, options.post];

      const updatedChannel = {
        ...currentChannelsState.channels[options.channelId],
        posts: updatedPosts,
        lastCreatedAt: options.post.created_at,
      };

      currentChannelsState.channels[options.channelId] = updatedChannel;

      this.channelsActive.next(currentChannelsState);
    }
  }

  private addPostsAfter(options: { channelId: string; posts: PostItem[] }) {
    console.log('addPostsAfter', options);

    const { last } = getFirstAndLastCreatedAt(options.posts);

    const currentChannelsState = this.channelsActive.getValue();

    if (
      options.posts.length &&
      currentChannelsState.channels?.[options.channelId]
    ) {
      const existingPosts =
        currentChannelsState.channels[options.channelId].posts;

      const updatedPosts = [...existingPosts!, ...options.posts];

      const updatedChannel = {
        ...currentChannelsState.channels[options.channelId],
        posts: updatedPosts,
        lastCreatedAt: last,
      };

      currentChannelsState.channels[options.channelId] = updatedChannel;

      this.channelsActive.next(currentChannelsState);
    }
  }

  private addPostsBefore(channelId: string, newPosts: PostItem[]): void {
    console.log('addPostsBefore', { channelId, newPosts });

    if (newPosts.length) {
      const currentChannelsState = this.channelsActive.getValue();

      const { first } = getFirstAndLastCreatedAt(newPosts);

      if (currentChannelsState.channels?.[channelId]) {
        const existingPosts = currentChannelsState.channels[channelId].posts;

        const updatedPosts = [...newPosts, ...existingPosts!];

        const updatedChannel = {
          ...currentChannelsState.channels[channelId],
          posts: updatedPosts,
          firstCreatedAt: first,
        };

        currentChannelsState.channels[channelId] = updatedChannel;

        this.channelsActive.next(currentChannelsState);
      }
    }
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

    this.setPosts(options.channelId, posts);

    this.postsLoading.next(false);
  }

  getPostsAfter(options: { channelId: string; limit: number }) {
    console.log('getPostsAfter', { options });

    return new Promise((resolve, reject) => {
      const channelsState = this.channelsActive.getValue();

      if (channelsState.channels?.[options.channelId].lastCreatedAt) {
        const last = channelsState.channels[options.channelId].lastCreatedAt;

        let params: HttpParams;

        if (last && last !== '') {
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
            this.addPostsAfter({ channelId: options.channelId, posts });
            resolve(posts);
          })
          .catch((err) => {
            console.error('error', err);
            reject(err);
          });
      } else {
        reject(new Error('getPostsAfter'));
      }
    });
  }

  getPostsBefore(options: { channelId: string; limit: number }) {
    console.log('Get posts before', options);

    return new Promise((resolve, reject) => {
      const channelsState = this.channelsActive.getValue();

      if (
        channelsState.channels &&
        channelsState.channels[options.channelId]?.firstCreatedAt
      ) {
        const first = channelsState.channels[options.channelId].firstCreatedAt;

        let params: HttpParams;

        if (first && first !== '') {
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

            this.addPostsBefore(options.channelId, posts);

            resolve(posts);
          })
          .catch((err) => {
            console.error(err);
            reject(err);
          });
      } else {
        reject(
          new Error(
            'channelsState.channels?.[options.channelId].firstCreatedAt is undefined',
          ),
        );
      }
    });
  }

  sendPost(options: {
    message: string;
    channelId: string;
    frontId: string;
    withoutAdd?: boolean;
    files?: FrontFile[];
  }) {
    return new Promise((resolve, reject) => {
      if (!options.withoutAdd) {
        this.addPostAfter({
          channelId: options.channelId,
          post: {
            id: options.frontId,
            frontId: options.frontId,
            message: options.message,
            created_at: new Date().toISOString(),
            channel_id: options.channelId,
            user: {
              name: this.userName.getValue(),
            },
            files: options.files?.length ? options.files : [],
          },
        });
      }

      this.api
        .post('/v1/api/posts', {
          message: options.message,
          channelId: options.channelId,
          frontId: options.frontId,
          files: options.files?.map((file) => file.id),
        })
        .then((resp) => {
          const { post } = resp as { post: PostItem };

          resolve(post);
        })
        .catch((err: HttpErrorResponse) => {
          reject(err);
        });
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
    const currentChannelsState = this.channelsActive.getValue();

    localStorage.setItem('activeChannel', `open_${channelId}`);

    const newState = {
      ...currentChannelsState,
      active: channelId,
      isOpenActive: true,
      isDirectActive: false,
    };
    this.channelsActive.next(newState);
  }

  setDirectActive(channelId: string): void {
    const currentChannelsState = this.channelsActive.getValue();

    localStorage.setItem('activeChannel', `direct_${channelId}`);

    const newState = {
      ...currentChannelsState,
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

    const currentChannelsState = this.channelsActive.getValue();

    // eslint-disable-next-line no-restricted-syntax
    for (const channel of openChannels) {
      // eslint-disable-next-line no-await-in-loop
      const { members } = (await this.api.get(
        `/v1/api/channels/${channel.id}/members`,
      )) as { members: string[] };

      channel.membersIds = members;

      // eslint-disable-next-line no-await-in-loop
      const { unread } = (await this.api.get(
        `/v1/api/channels/${channel.id}/unread`,
      )) as { unread: number };

      channel.unread = unread;

      if (!currentChannelsState.channels) {
        currentChannelsState.channels = {};
      }

      currentChannelsState.channels[channel.id] = channel;
    }

    this.channelsActive.next(currentChannelsState);

    this.getStatusesForChannelMembers(openChannels);
  }

  async getDirectChannels() {
    const openParams = new HttpParams().append('channelType', 'D');

    const resp = await this.api.get('/v1/api/channels', openParams);

    const directChannels = (resp as ChannelsResp).channels;

    const currentChannelsState = this.channelsActive.getValue();

    // eslint-disable-next-line no-restricted-syntax
    for (const channel of directChannels) {
      // eslint-disable-next-line no-await-in-loop
      const { members } = (await this.api.get(
        `/v1/api/channels/${channel.id}/members`,
      )) as { members: string[] };

      channel.membersIds = members;

      // eslint-disable-next-line no-await-in-loop
      const { unread } = (await this.api.get(
        `/v1/api/channels/${channel.id}/unread`,
      )) as { unread: number };

      channel.unread = unread;

      if (!currentChannelsState.channels) {
        currentChannelsState.channels = {};
      }

      currentChannelsState.channels[channel.id] = channel;
    }

    this.channelsActive.next(currentChannelsState);

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

  markChannelAsRead(channelId: string) {
    this.api
      .put(`/v1/api/channels/${channelId}/markasread`, {})
      .then(() => {
        const currentChannelsState = this.channelsActive.getValue();

        if (currentChannelsState.channels?.[channelId]) {
          const updatedChannel = {
            ...currentChannelsState.channels[channelId],
            unread: 0,
          };

          currentChannelsState.channels[channelId] = updatedChannel;

          this.channelsActive.next(currentChannelsState);
        }
      })
      .catch((err) => console.error(err));
  }

  incUnread(channelId: string) {
    const currentChannelsState = this.channelsActive.getValue();

    if (currentChannelsState.channels) {
      const updatedChannels = {
        ...currentChannelsState.channels,
        [channelId]: {
          ...currentChannelsState.channels[channelId],
          unread: currentChannelsState.channels[channelId].unread + 1,
        },
      };

      const updatedChannelsState = {
        ...currentChannelsState,
        channels: updatedChannels,
      };

      this.channelsActive.next(updatedChannelsState);
    }
  }

  sendOfflineMessages(
    channelsIds: Set<string> = new Set<string>([]),
  ): Promise<string[]> {
    return new Promise((resolve) => {
      const post = this.offlineMessages.shift();

      if (post) {
        let setIntervalId = 0;

        setIntervalId = setInterval(() => {
          this.sendPost({
            channelId: post.channel_id,
            message: post.message,
            withoutAdd: true,
            frontId: post.frontId ? post.frontId : '',
          }).then(() => {
            clearInterval(setIntervalId);

            channelsIds.add(post.channel_id);

            this.sendOfflineMessages(channelsIds).then(resolve);
          });
        }, 1000) as unknown as number;
      } else {
        resolve(Array.from(channelsIds));
      }
    });
  }

  addOfflineMessage(options: {
    channelId: string;
    message: string;
    frontId: string;
  }) {
    const currentChannelsState = this.channelsActive.getValue();

    const userName = this.userName.getValue();

    if (currentChannelsState.channels?.[options.channelId]) {
      currentChannelsState.channels[options.channelId].posts?.push({
        id: crypto.randomUUID().toString(),
        message: options.message,
        channel_id: options.channelId,
        created_at: new Date().toUTCString(),
        user: { name: userName },
        offline: true,
        frontId: options.frontId,
        files: [],
      });

      this.channelsActive.next(currentChannelsState);

      this.offlineMessages.push({
        id: crypto.randomUUID().toString(),
        message: options.message,
        channel_id: options.channelId,
        created_at: new Date().toUTCString(),
        user: { name: userName },
        frontId: options.frontId,
        files: [],
      });

      if (!this.offlineMessagesSending) {
        this.offlineMessagesSending = true;

        this.sendOfflineMessages().then(async (channelsIds) => {
          this.offlineMessagesSending = false;

          // eslint-disable-next-line no-restricted-syntax
          for (const channelId of channelsIds) {
            // eslint-disable-next-line no-await-in-loop
            await this.getPosts({
              channelId,
              limit: environment.POSTS_PAGE_SIZE,
            });
          }
        });
      }
    }
  }

  searchPosts(text: string) {
    this.api
      .get(`/v1/api/posts/search/${text}`)
      .then((resp) => {
        console.log(resp);

        this.searchedPosts.next((resp as { posts: PostItem[] }).posts);
        this.showSearchedPosts.next(true);
      })
      .catch((err) => console.error(err));
  }

  clearSearchedPosts() {
    this.searchedPosts.next([]);
    this.showSearchedPosts.next(false);
  }

  changeDirectChannelCollapse() {
    const value = this.directChannelCollapse.getValue();

    this.directChannelCollapse.next(!value);
  }

  changeOpenChannelCollapse() {
    const value = this.openChannelCollapse.getValue();

    this.openChannelCollapse.next(!value);
  }

  openChannelSearchModal() {
    this.showChannelSearchModal.next(true);
  }

  closeChannelSearchModal() {
    this.showChannelSearchModal.next(false);
  }
}
