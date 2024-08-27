import { TestBed } from '@angular/core/testing';

import { UniqueCheckService } from './unique-check.service';

describe('UniqueCheckService', () => {
  let service: UniqueCheckService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(UniqueCheckService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
