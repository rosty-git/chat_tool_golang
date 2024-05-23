import { HttpParams } from '@angular/common/http';
import { Component, inject, OnDestroy, OnInit } from '@angular/core';
import { firstValueFrom, retry } from 'rxjs';

import { GlobalVariable } from '../../global';
import { ApiService } from '../api.service';
import { DataService } from '../data.service';
import { HeaderComponent } from '../header/header.component';
import { LoginComponent } from '../login/login.component';
import { GetPostsResp } from '../message-box/message-box.component';
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
    private api: ApiService,
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
          if (this.channelsStore.active() === webSocketMsg.Payload.channel_id) {
            const lastCreatedAt = await firstValueFrom(
              this.dataService.lastCreatedAt$,
            );
            console.log('lastCreatedAt', lastCreatedAt);

            let params: HttpParams;

            if (lastCreatedAt !== '') {
              console.log('append');

              params = new HttpParams()
                .append('limit', 20)
                .append('afterCreatedAt', lastCreatedAt);
            } else {
              params = new HttpParams().append('limit', 20);
            }

            this.api
              .get(`/v1/api/posts/${webSocketMsg.Payload.channel_id}`, params)
              .subscribe({
                next: (response) => {
                  console.log('Get Channels', response);

                  this.dataService.addPosts((response as GetPostsResp).posts);
                },
                error: (err: unknown) => {
                  console.error('error', err);
                },
              });
          }
        }
      });
  }

  ngOnDestroy(): void {
    this.webSocketService.close();
  }
}
