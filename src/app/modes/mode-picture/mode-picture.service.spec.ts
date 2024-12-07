import { TestBed } from '@angular/core/testing';

import { ModePictureService } from './mode-picture.service';

describe('ModePictureService', () => {
  let service: ModePictureService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModePictureService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
