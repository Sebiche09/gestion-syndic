import { TestBed } from '@angular/core/testing';

import { DisplayCondominiumService } from './display-condominium.service';

describe('DisplayCondominiumService', () => {
  let service: DisplayCondominiumService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(DisplayCondominiumService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
