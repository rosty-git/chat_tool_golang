import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { interval, retry, tap } from 'rxjs';

import { environment } from '../../environments/environment';
import { ChannelsState, DataService } from '../data.service';
import { GalleryComponent } from '../gallery/gallery.component';
import { HeaderComponent } from '../header/header.component';
import { LoginComponent } from '../login/login.component';
import { MessageListComponent } from '../message-list/message-list.component';
import { SearchChannelsComponent } from '../search-channels/search-channels.component';
import { SearchedPostsComponent } from '../searched-posts/searched-posts.component';
import { SidebarComponent } from '../sidebar/sidebar.component';
import { WebSocketService } from '../web-socket.service';

type WebSocketMsg = {
  ToUsersIDs: string[];
  Action: 'new-post' | 'status-updated' | 'new-own-post';
  Payload: NewPostPayload | NewStatusPayload | NewOwnPostPayload;
};

type NewPostPayload = {
  id: string;
  created_at: string;
  updated_at: string;
  deleted_at: string;
  user_id: string;
  channel_id: string;
  message: string;
};

type NewStatusPayload = {
  userId: string;
  status: 'online' | 'away' | 'dnd' | 'offline';
};

type NewOwnPostPayload = {
  createdPost: NewPostPayload;
  frontId: string;
};

const USER_UPDATE_ONLINE_STATUS_INTERVAL = 90_000;
const USER_UPDATE_AWAY_STATUS_INTERVAL = 300_000;

@Component({
  selector: 'app-messenger',
  standalone: true,
  templateUrl: './messenger.component.html',
  styleUrl: './messenger.component.scss',
  imports: [
    LoginComponent,
    HeaderComponent,
    MessageListComponent,
    SidebarComponent,
    CommonModule,
    SearchedPostsComponent,
    SearchChannelsComponent,
    GalleryComponent
  ]
})
export class MessengerComponent implements OnInit, OnDestroy {
  private awayTimeoutId: ReturnType<typeof setTimeout>;

  private userOnlineStatusLastUpdate = new Date().getTime() - 100_000;

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  showSearchedPosts = false;

  constructor(
    private webSocketService: WebSocketService,
    private dataService: DataService,
  ) {
    this.awayTimeoutId = setInterval(() => {
      this.dataService.setAwayStatus();
    }, USER_UPDATE_AWAY_STATUS_INTERVAL);

    this.dataService.channelsState$.subscribe((value) => {
      this.channelsState$ = value;
    });

    this.dataService.showSearchedPosts$.subscribe((value) => {
      this.showSearchedPosts = value;
    });
  }

  ngOnInit(): void {
    this.webSocketService.connect(`${environment.BASE_WS_URL}/v1/ws/`);

    const ping$ = interval(30000).pipe(
      tap(() => {
        this.webSocketService.sendMessage({ type: 'ping' });
      }),
    );

    ping$.subscribe();

    this.webSocketService
      .getMessages()
      .pipe(retry({ delay: 2000 }))
      .subscribe({
        next: async (msg) => {
          console.log(msg);

          const webSocketMsg = msg as WebSocketMsg;

          if (webSocketMsg.Action === 'new-post') {
            const audio = new Audio('assets/new-message-notification.wav');

            webSocketMsg.Payload = webSocketMsg.Payload as NewPostPayload;

            const newChannelsOrder = this.updateChannelsOrder(
              webSocketMsg.Payload.channel_id,
            );

            this.dataService.updateChannelsIndexes(newChannelsOrder);

            if (
              this.channelsState$.active === webSocketMsg.Payload.channel_id
            ) {
              this.dataService
                .getPostsAfter({
                  channelId: webSocketMsg.Payload.channel_id,
                  limit: environment.POSTS_PAGE_SIZE,
                })
                .then(() => {
                  audio.play();
                });
            } else if (
              this.channelsState$.active !== webSocketMsg.Payload.channel_id
            ) {
              this.dataService.incUnread(webSocketMsg.Payload.channel_id);

              audio.play();
            } else {
              audio.play();
            }
          } else if (webSocketMsg.Action === 'status-updated') {
            webSocketMsg.Payload = webSocketMsg.Payload as NewStatusPayload;

            this.dataService.setStatus(
              webSocketMsg.Payload.userId,
              webSocketMsg.Payload.status,
            );
          } else if (webSocketMsg.Action === 'new-own-post') {
            const payload = webSocketMsg.Payload as NewOwnPostPayload;

            if (
              this.channelsState$.channels?.[
                payload.createdPost.channel_id
              ].posts?.some((post) => post.frontId === payload.frontId)
            ) {
              // console.log();
            } else {
              this.dataService.getPostsAfter({
                channelId: payload.createdPost.channel_id,
                limit: environment.POSTS_PAGE_SIZE,
              });
            }
          }
        },
        error: (e) => console.error('WebSocket error', e),
        complete: () => console.log('WebSocket connection closed'),
      });

    this.dataService.getUser();

    Promise.all([
      this.dataService.getOpenChannels(),
      this.dataService.getDirectChannels(),
    ]).then(() => {
      const activeChannel = localStorage.getItem('activeChannel');

      if (activeChannel) {
        const [channelType, channelId] = activeChannel.split('_');

        if (channelType === 'open') {
          this.dataService.setOpenActive(channelId);
        } else if (channelType === 'direct') {
          this.dataService.setDirectActive(channelId);
        }

        this.dataService.setOpenChannelCollapse(false);
        this.dataService.setDirectChannelCollapse(false);
      }
    });
  }

  ngOnDestroy(): void {
    this.webSocketService.close();
  }

  mouseMove(): void {
    if (
      new Date().getTime() - this.userOnlineStatusLastUpdate >
      USER_UPDATE_ONLINE_STATUS_INTERVAL
    ) {
      this.userOnlineStatusLastUpdate = new Date().getTime();
      this.dataService.updateOnlineStatus().finally(() => {
        clearTimeout(this.awayTimeoutId);

        this.awayTimeoutId = setInterval(() => {
          this.dataService.setAwayStatus();
        }, USER_UPDATE_AWAY_STATUS_INTERVAL);
      });
    }
  }

  // eslint-disable-next-line class-methods-use-this
  updateChannelsOrder(channelId: string): string[] {
    const channelsOrderStringify = localStorage.getItem('channelsOrder');
    let channelsOrder = channelsOrderStringify
      ? JSON.parse(channelsOrderStringify)
      : [];

    channelsOrder = channelsOrder.filter((i: string) => i !== channelId);

    channelsOrder.unshift(channelId);

    localStorage.setItem('channelsOrder', JSON.stringify(channelsOrder));

    return channelsOrder;
  }
}
