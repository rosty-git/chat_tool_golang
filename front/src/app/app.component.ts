import { HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';

import { MessengerComponent } from './messenger/messenger.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    HttpClientModule,
    MessengerComponent,
    InfiniteScrollModule,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'chat_tool_front';
}
