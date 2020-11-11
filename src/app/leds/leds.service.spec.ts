import { TestBed } from '@angular/core/testing';

import { LedsService } from './leds.service';

describe('LedsService', () => {
  let service: LedsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LedsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
