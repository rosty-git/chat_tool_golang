import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { retry } from 'rxjs';

import { GlobalVariable } from '../../global';
import { DataService } from '../data.service';
import { HeaderComponent } from '../header/header.component';
import { LoginComponent } from '../login/login.component';
import { MessageListComponent } from '../message-list/message-list.component';
import { SidebarComponent } from '../sidebar/sidebar.component';
import { ChannelsStore } from '../store/channels.store';
import { WebSocketService } from '../web-socket.service';

type WebSocketMsg = {
  ToUsersIDs: string[];
  Action: 'new-post' | 'status-updated';
  Payload: NewPostPayload;
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
  readonly channelsStore = inject(ChannelsStore);

  constructor(
    private webSocketService: WebSocketService,
    private dataService: DataService,
  ) {}

  ngOnInit(): void {
    this.webSocketService.connect(`${GlobalVariable.BASE_WS_URL}/v1/ws/`);

    this.webSocketService
      .getMessages()
      .pipe(retry({ delay: 5000 }))
      .subscribe(async (msg) => {
        console.log(msg);

        console.log(this.channelsStore.active());

        const webSocketMsg = msg as WebSocketMsg;

        if (webSocketMsg.Action === 'new-post') {
          const audio = new Audio('assets/new-message-notification.wav');

          if (this.channelsStore.active() === webSocketMsg.Payload.channel_id) {
            this.dataService
              .getPostsAfter({
                channelId: webSocketMsg.Payload.channel_id,
                limit: GlobalVariable.POSTS_PAGE_SIZE,
              })
              .then(() => {
                audio.play();
              });
          } else {
            audio.play();
          }
        }
      });
  }

  ngOnDestroy(): void {
    this.webSocketService.close();
  }

  mouseMove(): void {
    this.dataService.updateOnlineStatus();
  }
}
