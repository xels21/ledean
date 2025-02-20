import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeRunningLedComponent } from './mode-running-led.component';

describe('ModeRunningLedComponent', () => {
  let component: ModeRunningLedComponent;
  let fixture: ComponentFixture<ModeRunningLedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeRunningLedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeRunningLedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
