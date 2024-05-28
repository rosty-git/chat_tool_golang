import { HttpParams } from '@angular/common/http';
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

  constructor(private dataService: DataService) {
    this.scrollFrame = new ElementRef('');
    this.scrollContainer = this.scrollFrame as unknown as HTMLElement;

    effect(() => {
      const channelsState = getState(this.channelsStore);

      const params = new HttpParams().append(
        'limit',
        GlobalVariable.POSTS_PAGE_SIZE,
      );

      this.dataService.getPosts(channelsState.active, params);

      this.posts$.subscribe((posts) => {
        this.postItems = posts;

        setTimeout(() => {
          this.scrollContainer.scroll({
            top: this.scrollContainer.scrollHeight,
            left: 0,
            // behavior: 'smooth',
          });
        }, 1);
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
    this.scrollContainer.scroll({
      top: this.scrollContainer.scrollHeight,
      left: 0,
      behavior: 'smooth',
    });
  }

  scrolled(event: unknown): void {
    this.isNearBottom = this.isUserNearBottom();
  }

  onUp() {
    console.log('scrolled up!');
  }

  onScrollDown() {
    console.log('scrolled down!');
  }
}
