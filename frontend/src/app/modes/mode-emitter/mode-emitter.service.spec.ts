import { TestBed } from '@angular/core/testing';

import { ModeEmitterService } from './mode-emitter.service';

describe('ModeEmitterService', () => {
  let service: ModeEmitterService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ModeEmitterService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
