import { TestBed } from '@angular/core/testing';

import { ModeSpectrumService } from './mode-spectrum.service';

describe('ModeSpectrumService', () => {
  let service: ModeSpectrumService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeSpectrumService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
