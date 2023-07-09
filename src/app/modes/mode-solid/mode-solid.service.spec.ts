import { TestBed } from '@angular/core/testing';

import { ModeSolidService } from './mode-solid.service';

describe('ModeSolidService', () => {
  let service: ModeSolidService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeSolidService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
