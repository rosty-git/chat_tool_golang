import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DirectChannelsComponent } from './directChannels.component';

describe('ContactsComponent', () => {
  let component: DirectChannelsComponent;
  let fixture: ComponentFixture<DirectChannelsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DirectChannelsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(DirectChannelsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
