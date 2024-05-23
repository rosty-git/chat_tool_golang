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

const getLastCreatedAt = (posts: PostItem[]): string => {
  if (posts.length === 0) {
    return '';
  }

  return posts.reduce(
    (latest, post) => (post.created_at > latest ? post.created_at : latest),
    posts[0].created_at,
  );
};

@Injectable({
  providedIn: 'root',
})
export class DataService {
  private posts = new BehaviorSubject<PostItem[]>([]);

  private lastCreatedAt = new BehaviorSubject<string>('');

  posts$ = this.posts.asObservable();

  lastCreatedAt$ = this.lastCreatedAt.asObservable();

  setPosts(newPosts: PostItem[]) {
    this.posts.next(newPosts);

    this.lastCreatedAt.next(getLastCreatedAt(newPosts));
  }

  addPost(newPost: PostItem) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, newPost];
    this.posts.next(updatedPosts);

    this.lastCreatedAt.next(newPost.created_at);
  }

  addPosts(newPosts: PostItem[]) {
    const currentPosts = this.posts.getValue();
    const updatedPosts = [...currentPosts, ...newPosts];
    this.posts.next(updatedPosts);

    this.lastCreatedAt.next(getLastCreatedAt(newPosts));
  }
}
