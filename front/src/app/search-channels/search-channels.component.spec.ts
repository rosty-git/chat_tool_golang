import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchChannelsComponent } from './search-channels.component';

describe('SearchChannelsComponent', () => {
  let component: SearchChannelsComponent;
  let fixture: ComponentFixture<SearchChannelsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SearchChannelsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SearchChannelsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
