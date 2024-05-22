import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpenChannelsComponent } from './openChannels.component';

describe('OpenChannelsComponent', () => {
  let component: OpenChannelsComponent;
  let fixture: ComponentFixture<OpenChannelsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpenChannelsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(OpenChannelsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
