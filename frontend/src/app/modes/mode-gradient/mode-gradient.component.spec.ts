import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeGradientComponent } from './mode-gradient.component';

describe('ModeGradientComponent', () => {
  let component: ModeGradientComponent;
  let fixture: ComponentFixture<ModeGradientComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeGradientComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeGradientComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
