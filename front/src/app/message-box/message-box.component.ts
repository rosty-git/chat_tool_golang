import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';

import { GlobalVariable } from '../../global';
import { DataService, type PostItem } from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';

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

  postsLoading$ = false;

  posts: PostItem[] = [];

  throttle = 50;

  scrollDistance = 2;

  scrollUpDistance = 2;

  firstTime = true;

  activeChannel$: string = '';

  constructor(private dataService: DataService) {
    this.scrollFrame = new ElementRef('');
    this.scrollContainer = this.scrollFrame as unknown as HTMLElement;

    this.dataService.posts$.subscribe((value) => {
      this.posts = value;
    });

    this.dataService.postsLoading$.subscribe((value) => {
      this.postsLoading$ = value;
    });

    this.dataService.channelsActive$.subscribe((value) => {
      this.activeChannel$ = value.active;
    });
  }

  ngAfterViewInit() {
    this.scrollContainer = this.scrollFrame.nativeElement;
    this.dataService.posts$.subscribe(() => this.onItemElementsChanged());
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
    this.dataService.getPostsBefore({
      channelId: this.activeChannel$,
      limit: GlobalVariable.POSTS_PAGE_SIZE,
    });
  }

  // onScrollDown() {
  //   console.log('scrolled down!');
  // }
}
