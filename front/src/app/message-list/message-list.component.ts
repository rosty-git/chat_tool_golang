/* eslint-disable no-await-in-loop */
import { CommonModule } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { PickerComponent } from '@ctrl/ngx-emoji-mart';
import axios from 'axios';
import { InfiniteScrollModule } from 'ngx-infinite-scroll';
import { BehaviorSubject, distinctUntilKeyChanged } from 'rxjs';

import { environment } from '../../environments/environment';
import { ApiService } from '../api.service';
import {
  ChannelsState,
  DataService,
  FrontFile,
  PostItem,
} from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';

@Component({
  selector: 'app-message-list',
  standalone: true,
  templateUrl: './message-list.component.html',
  styleUrl: './message-list.component.scss',
  imports: [
    ReactiveFormsModule,
    InfiniteScrollModule,
    MessageItemComponent,
    CommonModule,
    PickerComponent,
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

  isFileOver: boolean = false;

  filesUploading = false;

  files: FrontFile[] = [];

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  isOpen = false;

  constructor(
    private dataService: DataService,
    private api: ApiService,
  ) {
    this.scrollFrameDiv = new ElementRef('');
    this.scrollContainer = this.scrollFrameDiv as unknown as HTMLElement;

    this.dataService.channelsState$.subscribe((value) => {
      const posts = value.channels?.[value.active]?.posts ?? [];
      if (posts.length) {
        this.posts.next(posts);
      } else {
        this.posts.next([]);
      }
    });

    this.dataService.channelsState$
      .pipe(distinctUntilKeyChanged('active'))
      .subscribe((value) => {
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

        if (value.active) {
          this.dataService.markChannelAsRead(value.active);
        }
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

  onDragOver(event: DragEvent) {
    event.preventDefault();
    this.isFileOver = true;
  }

  onDrop(event: DragEvent) {
    event.preventDefault();
    this.isFileOver = false;

    if (event.dataTransfer && event.dataTransfer.files.length > 0) {
      const { files } = event.dataTransfer;
      this.uploadFiles(files);
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onDragLeave(_event: DragEvent) {
    this.isFileOver = false;
  }

  async uploadFiles(files: FileList) {
    this.filesUploading = true;

    if (files?.length) {
      await this.processFiles(files);
    }

    this.filesUploading = false;
  }

  async processFiles(files: FileList) {
    for (let i = 0; i < files.length; i += 1) {
      const file = files[i];

      this.files.push({
        name: file.name,
        size: file.size,
        type: file.type,
        ext: file.name.split('.').pop()!,
        progress: 1,
      });
    }

    for (let i = 0; i < files.length; i += 1) {
      const file = files[i];

      const createdFile = await this.api.post('/v1/api/files', {
        name: file.name,
        size: file.size,
        type: file.type,
      });

      const objIndex = this.files.findIndex(
        (fileItem) =>
          fileItem.name === file.name &&
          fileItem.size === file.size &&
          fileItem.type === file.type,
      );

      const createdFileId = (createdFile as { id: string }).id;

      this.files[objIndex].id = createdFileId;

      const buffer = await file.arrayBuffer();

      const blob = new Blob([buffer], { type: file.type });

      const url = URL.createObjectURL(blob);

      this.files[objIndex].blobUrl = url;

      const fileExt = file.name.split('.').pop();

      const s3Key = `${createdFileId}.${fileExt}`;

      const presignedUrl = await this.api.post(
        `/v1/api/files/get-presigned-url/${s3Key}`,
        {},
      );

      await axios.put((presignedUrl as { url: string }).url, buffer, {
        onUploadProgress: (progressEvent) => {
          this.files[objIndex].progress = Math.round(
            progressEvent.progress! * 100,
          );
        },
      });
      this.files[objIndex].progress = 100;

      const updatedFile = await this.api.put(
        `/v1/api/files/${createdFileId}/${s3Key}`,
        {},
      );

      console.log({ updatedFile });
    }
  }

  sendMessage() {
    if (this.messageForm.value.message) {
      const channelsOrderStringify = localStorage.getItem('channelsOrder');
      let channelsOrder = channelsOrderStringify
        ? JSON.parse(channelsOrderStringify)
        : [];

      channelsOrder = channelsOrder.filter(
        (i: string) => i !== this.channelsState$.active,
      );

      channelsOrder.unshift(this.channelsState$.active);

      localStorage.setItem('channelsOrder', JSON.stringify(channelsOrder));

      const message = this.messageForm.value.message as string;

      const frontId = crypto.randomUUID().toString();

      this.dataService
        .sendPost({
          message,
          channelId: this.channelsState$.active,
          frontId,
          files: this.files,
        })
        .catch((err: HttpErrorResponse) => {
          if (err.status === 0) {
            this.dataService.addOfflineMessage({
              message,
              channelId: this.channelsState$.active,
              frontId,
            });
          }
        });

      this.messageForm.reset();
      this.files = [];
    }
  }

  // eslint-disable-next-line class-methods-use-this
  handleEnterKey(event: Event) {
    const keyboardEvent = event as KeyboardEvent;

    if (keyboardEvent.key === 'Enter' && !keyboardEvent.shiftKey) {
      event.preventDefault();
    }
  }

  async deleteFile(name: string, size: number, type: string) {
    const fileForDelete = this.files.find(
      (file) => file.name === name && file.size === size && file.type === type,
    );

    this.files = this.files.filter(
      (file) =>
        !(file.name === name && file.size === size && file.type === type),
    );

    if (fileForDelete?.id) {
      await this.api.delete(`/v1/api/files/${fileForDelete.id}`);
    }
  }

  async onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;

    this.filesUploading = true;

    if (input.files?.length) {
      await this.processFiles(input.files);
    }

    this.filesUploading = false;
  }

  toggleDropdown() {
    this.isOpen = !this.isOpen;
  }

  addEmoji(event: { emoji: { native: string } }) {
    const message = this.messageForm.value.message!;

    this.messageForm.setValue({
      message: `${message || ''}${event.emoji.native}`,
    });

    this.isOpen = !this.isOpen;
  }
}
