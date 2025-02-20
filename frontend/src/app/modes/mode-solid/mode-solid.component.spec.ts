import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeSolidComponent } from './mode-solid.component';

describe('ModeSolidComponent', () => {
  let component: ModeSolidComponent;
  let fixture: ComponentFixture<ModeSolidComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeSolidComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeSolidComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
