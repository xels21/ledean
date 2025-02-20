import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeSolidRainbowComponent } from './mode-solid-rainbow.component';

describe('ModeSolidRainbowComponent', () => {
  let component: ModeSolidRainbowComponent;
  let fixture: ComponentFixture<ModeSolidRainbowComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeSolidRainbowComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeSolidRainbowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
