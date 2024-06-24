import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OpenChannelComponent } from './open-channel.component';

describe('OpenChannelComponent', () => {
  let component: OpenChannelComponent;
  let fixture: ComponentFixture<OpenChannelComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [OpenChannelComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(OpenChannelComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
