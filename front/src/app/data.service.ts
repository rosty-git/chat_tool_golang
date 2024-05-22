import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

type PostItem = {
  id: string;
  message: string;
  created_at: string;
  user: {
    name: string;
  };
};

@Injectable({
  providedIn: 'root',
})
export class DataService {
  private posts = new BehaviorSubject<PostItem[]>([]);

  posts$ = this.posts.asObservable();

  setPosts(newPosts: PostItem[]) {
    this.posts.next(newPosts);
  }

  addPost(newPost: PostItem) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, newPost];
    this.posts.next(updatedPosts);
  }
}
