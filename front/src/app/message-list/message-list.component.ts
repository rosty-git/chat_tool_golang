import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';
import { BehaviorSubject, distinctUntilKeyChanged } from 'rxjs';

import { environment } from '../../environments/environment';
import { ChannelsState, DataService, PostItem } from '../data.service';
import { MessageInputComponent } from '../message-input/message-input.component';
import { MessageItemComponent } from '../message-item/message-item.component';

@Component({
  selector: 'app-message-list',
  standalone: true,
  templateUrl: './message-list.component.html',
  styleUrl: './message-list.component.scss',
  imports: [
    MessageInputComponent,
    InfiniteScrollModule,
    MessageItemComponent,
    CommonModule,
  ],
})
export class MessageListComponent implements AfterViewInit {
  @ViewChild('scrollFrame')
  private scrollFrameDiv: ElementRef;

  private scrollContainer: HTMLElement;

  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  private isNearBottom = true;

  throttle = 50;

  scrollDistance = 2;

  scrollUpDistance = 2;

  postsLoading$ = false;

  private posts = new BehaviorSubject<PostItem[]>([]);

  posts$ = this.posts.asObservable();

  searchedPosts: PostItem[] = [];

  constructor(private dataService: DataService) {
    this.scrollFrameDiv = new ElementRef('');
    this.scrollContainer = this.scrollFrameDiv as unknown as HTMLElement;

    this.dataService.channelsActive$.subscribe((value) => {
      console.log('this.dataService.channelsActive$.subscribe');
      const posts = value.channels?.[value.active]?.posts ?? [];
      if (posts.length) {
        this.posts.next(posts);
      } else {
        this.posts.next([]);
      }
    });

    this.dataService.channelsActive$
      .pipe(distinctUntilKeyChanged('active'))
      .subscribe((value) => {
        console.log('active');
        this.channelsState$ = value;

        if (
          value.channels &&
          value.channels[value.active] &&
          (value.channels?.[value.active].posts?.length === 0 ||
            !value.channels?.[value.active].posts)
        ) {
          this.dataService.getPosts({
            channelId: value.active,
            limit: environment.POSTS_PAGE_SIZE,
          });
        }

        this.dataService.markChannelAsRead(value.active);
      });

    this.dataService.postsLoading$.subscribe((value) => {
      this.postsLoading$ = value;
    });

    this.dataService.searchedPosts$.subscribe((value) => {
      this.searchedPosts = value;
    });
  }

  ngAfterViewInit() {
    this.scrollContainer = this.scrollFrameDiv.nativeElement;
    this.posts$.subscribe(() => this.onItemElementsChanged());
  }

  private onItemElementsChanged() {
    if (this.isNearBottom) {
      this.scrollToBottom();
    }
  }

  private scrollToBottom(): void {
    setTimeout(() => {
      this.scrollContainer.scroll({
        top: this.scrollContainer.scrollHeight,
        left: 0,
      });
    }, 1);
  }

  private isUserNearBottom(): boolean {
    const threshold = 150;
    const position =
      this.scrollContainer.scrollTop + this.scrollContainer.offsetHeight;
    const height = this.scrollContainer.scrollHeight;
    return position > height - threshold;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  scrolled(_event: Event): void {
    this.isNearBottom = this.isUserNearBottom();
  }

  onUp() {
    this.dataService.getPostsBefore({
      channelId: this.channelsState$.active,
      limit: environment.POSTS_PAGE_SIZE,
    });
  }
}
