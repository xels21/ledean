import { TestBed } from '@angular/core/testing';

import { ModeGradientService } from './mode-gradient.service';

describe('ModeGradientService', () => {
  let service: ModeGradientService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeGradientService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
