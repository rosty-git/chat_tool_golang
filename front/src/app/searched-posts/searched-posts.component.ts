import { Component } from '@angular/core';

import { DataService, PostItem } from '../data.service';
import { MessageItemComponent } from '../message-item/message-item.component';

@Component({
  selector: 'app-searched-posts',
  standalone: true,
  imports: [MessageItemComponent],
  templateUrl: './searched-posts.component.html',
  styleUrl: './searched-posts.component.scss',
})
export class SearchedPostsComponent {
  posts: PostItem[] = [];

  constructor(private dataService: DataService) {
    this.dataService.searchedPosts$.subscribe((value) => {
      this.posts = value;
    });
  }

  close() {
    this.dataService.clearSearchedPosts()
  }
}
