import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModeEmitterComponent } from './mode-emitter.component';

describe('ModeEmitterComponent', () => {
  let component: ModeEmitterComponent;
  let fixture: ComponentFixture<ModeEmitterComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ModeEmitterComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ModeEmitterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
