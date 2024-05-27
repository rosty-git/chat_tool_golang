import { HttpClientModule } from '@angular/common/http';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { provideIcons } from '@ng-icons/core';
import { cssMenuGridR } from '@ng-icons/css.gg';
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
  viewProviders: [provideIcons({ cssMenuGridR })],
})
export class AppComponent {
  title = 'chat_tool_front';
}
