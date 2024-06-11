import { Component, OnDestroy, OnInit } from '@angular/core';
import { retry } from 'rxjs';

import { GlobalVariable } from '../../global';
import { ChannelsState, DataService } from '../data.service';
import { HeaderComponent } from '../header/header.component';
import { LoginComponent } from '../login/login.component';
import { MessageListComponent } from '../message-list/message-list.component';
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
  imports: [
    LoginComponent,
    HeaderComponent,
    MessageListComponent,
    SidebarComponent,
  ],
  templateUrl: './messenger.component.html',
  styleUrl: './messenger.component.scss',
})
export class MessengerComponent implements OnInit, OnDestroy {
  private awayTimeoutId: ReturnType<typeof setTimeout>;

  private userOnlineStatusLastUpdate = new Date().getTime() - 100_000;

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  constructor(
    private webSocketService: WebSocketService,
    private dataService: DataService,
  ) {
    this.awayTimeoutId = setInterval(() => {
      this.dataService.setAwayStatus();
    }, USER_UPDATE_AWAY_STATUS_INTERVAL);

    this.dataService.channelsActive$.subscribe((value) => {
      this.channelsState$ = value;
    });
  }

  ngOnInit(): void {
    this.webSocketService.connect(`${GlobalVariable.BASE_WS_URL}/v1/ws/`);

    this.webSocketService
      .getMessages()
      .pipe(retry({ delay: 5000 }))
      .subscribe(async (msg) => {
        console.log(msg);

        const webSocketMsg = msg as WebSocketMsg;

        if (webSocketMsg.Action === 'new-post') {
          const audio = new Audio('assets/new-message-notification.wav');

          webSocketMsg.Payload = webSocketMsg.Payload as NewPostPayload;

          if (this.channelsState$.active === webSocketMsg.Payload.channel_id) {
            this.dataService
              .getPostsAfter({
                channelId: webSocketMsg.Payload.channel_id,
                limit: GlobalVariable.POSTS_PAGE_SIZE,
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
              limit: GlobalVariable.POSTS_PAGE_SIZE,
            });
          }
        }
      });

    this.dataService.getUser();

    this.dataService.getOpenChannels();
    this.dataService.getDirectChannels();
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
}
