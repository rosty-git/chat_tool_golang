import {
  AfterViewInit,
  Component,
  effect,
  ElementRef,
  inject,
  ViewChild,
} from '@angular/core';
import { getState } from '@ngrx/signals';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';

import { GlobalVariable } from '../../global';
import { DataService, type PostItem } from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';
import { ChannelsStore } from '../store/channels.store';

@Component({
  selector: 'app-message-box',
  standalone: true,
  templateUrl: './message-box.component.html',
  styleUrl: './message-box.component.scss',
  imports: [MessageItemComponent, InfiniteScrollModule],
})
export class MessageBoxComponent implements AfterViewInit {
  @ViewChild('scrollFrame', { static: false }) scrollFrame: ElementRef;

  private scrollContainer: HTMLElement;

  private isNearBottom = true;

  readonly channelsStore = inject(ChannelsStore);

  posts$ = this.dataService.posts$;

  postsLoading$ = false;

  postItems: PostItem[] = [];

  throttle = 50;

  scrollDistance = 2;

  scrollUpDistance = 2;

  firstTime = true;

  constructor(private dataService: DataService) {
    this.scrollFrame = new ElementRef('');
    this.scrollContainer = this.scrollFrame as unknown as HTMLElement;

    effect(() => {
      const channelsState = getState(this.channelsStore);

      this.dataService.getPosts({
        channelId: channelsState.active,
        limit: GlobalVariable.POSTS_PAGE_SIZE,
      });

      this.posts$.subscribe((posts) => {
        this.postItems = posts;
      });
    });

    this.dataService.postsLoading$.subscribe((value) => {
      this.postsLoading$ = value;
    });
  }

  ngAfterViewInit() {
    this.scrollContainer = this.scrollFrame.nativeElement;
    this.posts$.subscribe(() => this.onItemElementsChanged());
  }

  private onItemElementsChanged() {
    if (this.isNearBottom) {
      this.scrollToBottom();
    }
  }

  private isUserNearBottom(): boolean {
    const threshold = 150;
    const position =
      this.scrollContainer.scrollTop + this.scrollContainer.offsetHeight;
    const height = this.scrollContainer.scrollHeight;
    return position > height - threshold;
  }

  private scrollToBottom(): void {
    setTimeout(() => {
      this.scrollContainer.scroll({
        top: this.scrollContainer.scrollHeight,
        left: 0,
        // behavior: this.firstTime ? 'smooth',
        // behavior: this.firstTime ? 'instant' : 'smooth',
      });
      // this.firstTime = false;
    }, 1);
  }

  scrolled(event: Event): void {
    this.isNearBottom = this.isUserNearBottom();
  }

  onUp() {
    const channelsState = getState(this.channelsStore);

    this.dataService.getPostsBefore({
      channelId: channelsState.active,
      limit: GlobalVariable.POSTS_PAGE_SIZE,
    });
  }

  // onScrollDown() {
  //   console.log('scrolled down!');
  // }
}
