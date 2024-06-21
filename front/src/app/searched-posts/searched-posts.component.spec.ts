import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchedPostsComponent } from './searched-posts.component';

describe('SearchedPostsComponent', () => {
  let component: SearchedPostsComponent;
  let fixture: ComponentFixture<SearchedPostsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SearchedPostsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SearchedPostsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
