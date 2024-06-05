import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';
import { BehaviorSubject } from 'rxjs';

import { GlobalVariable } from '../../global';
import { DataService, type PostItem } from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';

@Component({
  selector: 'app-message-box',
  standalone: true,
  templateUrl: './message-box.component.html',
  styleUrl: './message-box.component.scss',
  imports: [MessageItemComponent, InfiniteScrollModule, CommonModule],
})
export class MessageBoxComponent implements AfterViewInit {
  @ViewChild('scrollFrame', { static: false }) scrollFrame: ElementRef;

  private scrollContainer: HTMLElement;

  private isNearBottom = true;

  postsLoading$ = false;

  private posts = new BehaviorSubject<PostItem[]>([]);

  posts$ = this.posts.asObservable();

  throttle = 50;

  scrollDistance = 2;

  scrollUpDistance = 2;

  private activeChannel = new BehaviorSubject<string>('');

  activeChannel$ = this.activeChannel.asObservable();

  constructor(private dataService: DataService) {
    this.scrollFrame = new ElementRef('');
    this.scrollContainer = this.scrollFrame as unknown as HTMLElement;

    this.dataService.channelsActive$.subscribe((value) => {
      this.activeChannel.next(value.active);

      if (!value.channels?.[value.active]) {
        this.dataService.getPosts({
          channelId: value.active,
          limit: GlobalVariable.POSTS_PAGE_SIZE,
        });
      }

      if (value?.channels?.[value.active]?.posts) {
        this.posts.next(value.channels[value.active].posts);
      }
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
      });
    }, 1);
  }

  scrolled(event: Event): void {
    this.isNearBottom = this.isUserNearBottom();
  }

  onUp() {
    this.dataService.getPostsBefore({
      channelId: this.activeChannel.getValue(),
      limit: GlobalVariable.POSTS_PAGE_SIZE,
    });
  }
}
