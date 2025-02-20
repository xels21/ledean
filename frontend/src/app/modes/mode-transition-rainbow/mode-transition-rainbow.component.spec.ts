import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeTransitionRainbowComponent } from './mode-transition-rainbow.component';

describe('ModeTransitionRainbowComponent', () => {
  let component: ModeTransitionRainbowComponent;
  let fixture: ComponentFixture<ModeTransitionRainbowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeTransitionRainbowComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeTransitionRainbowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
