import { Component, OnDestroy, OnInit } from '@angular/core';
import { retry } from 'rxjs';

import { GlobalVariable } from '../../global';
import { HeaderComponent } from '../header/header.component';
import { LoginComponent } from '../login/login.component';
import { MessageListComponent } from '../message-list/message-list.component';
import { SidebarComponent } from '../sidebar/sidebar.component';
import { WebSocketService } from '../web-socket.service';

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
  constructor(private webSocketService: WebSocketService) {}

  ngOnInit(): void {
    this.webSocketService.connect(`${GlobalVariable.BASE_WS_URL}/v1/ws/`);

    this.webSocketService
      .getMessages()
      .pipe(retry({ delay: 5000 }))
      .subscribe((msg) => {
        console.log(msg);
      });
  }

  ngOnDestroy(): void {
    this.webSocketService.close();
  }
}
