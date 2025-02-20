import { TestBed } from '@angular/core/testing';

import { ModeRunningLedService } from './mode-running-led.service';

describe('ModeRunningLedService', () => {
  let service: ModeRunningLedService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeRunningLedService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
