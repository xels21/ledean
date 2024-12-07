import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ModePictureComponent } from './mode-picture.component';

describe('ModePictureComponent', () => {
  let component: ModePictureComponent;
  let fixture: ComponentFixture<ModePictureComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ModePictureComponent]
    });
    fixture = TestBed.createComponent(ModePictureComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
