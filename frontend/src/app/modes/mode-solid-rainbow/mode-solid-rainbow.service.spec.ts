import { TestBed } from '@angular/core/testing';

import { ModeSolidRainbowService } from './mode-solid-rainbow.service';

describe('ModeSolidRainbowService', () => {
  let service: ModeSolidRainbowService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeSolidRainbowService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
