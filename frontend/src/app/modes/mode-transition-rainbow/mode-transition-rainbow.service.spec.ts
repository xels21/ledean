import { TestBed } from '@angular/core/testing';

import { ModeTransitionRainbowService } from './mode-transition-rainbow.service';

describe('ModeTransitionRainbowService', () => {
  let service: ModeTransitionRainbowService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeTransitionRainbowService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
