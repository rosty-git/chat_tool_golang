/* eslint-disable no-await-in-loop */

import { NgClass } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { PickerComponent } from '@ctrl/ngx-emoji-mart';
import axios from 'axios';

import { ApiService } from '../api.service';
import { ChannelsState, DataService, FrontFile } from '../data.service';

@Component({
  selector: 'app-message-input',
  standalone: true,
  templateUrl: './message-input.component.html',
  styleUrl: './message-input.component.scss',
  imports: [ReactiveFormsModule, NgClass, PickerComponent],
})
export class MessageInputComponent {
  channelsState$: ChannelsState = {
    isOpenActive: false,
    isDirectActive: false,
    active: '',
  };

  files: FrontFile[] = [];

  filesUploading = false;

  constructor(
    private dataService: DataService,
    private api: ApiService,
  ) {
    this.dataService.channelsActive$.subscribe((value) => {
      this.channelsState$ = value;
    });
  }

  messageForm = new FormGroup({
    message: new FormControl(''),
  });

  isOpen = false;

  sendMessage() {
    if (this.messageForm.value.message) {
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

  async onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;

    this.filesUploading = true;

    if (input.files?.length) {
      await this.processFiles(input.files)
    }

    this.filesUploading = false;
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

  toggleDropdown() {
    this.isOpen = !this.isOpen;
  }

  addEmoji(event: { emoji: { native: string } }) {
    const message = this.messageForm.value.message!;

    this.messageForm.setValue({
      message: `${message}${event.emoji.native}`,
    });

    this.isOpen = !this.isOpen;
  }

  isFileOver: boolean = false;

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

  // eslint-disable-next-line class-methods-use-this
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
}
