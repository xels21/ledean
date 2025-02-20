import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeSpectrumComponent } from './mode-spectrum.component';

describe('ModeSpectrumComponent', () => {
  let component: ModeSpectrumComponent;
  let fixture: ComponentFixture<ModeSpectrumComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ModeSpectrumComponent]
    });
    fixture = TestBed.createComponent(ModeSpectrumComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
