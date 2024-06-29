import { Component } from '@angular/core';

import { environment } from '../../environments/environment';
import { DataService, GalleryState } from '../data.service';

@Component({
  selector: 'app-gallery',
  standalone: true,
  imports: [],
  templateUrl: './gallery.component.html',
  styleUrl: './gallery.component.scss',
})
export class GalleryComponent {
  state: GalleryState = {
    isOpen: false,
  };

  s3Prefix = environment.S3_PREFIX;

  constructor(private dataService: DataService) {
    this.dataService.galleryState$.subscribe((value) => {
      this.state = value;
    });
  }

  close() {
    this.dataService.closeGallery();
  }

  increaseIndex() {
    this.dataService.increaseGalleryIndex();
  }

  decreaseIndex() {
    this.dataService.decreaseGalleryIndex();
  }
}
