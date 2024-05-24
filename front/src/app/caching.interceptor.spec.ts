import { HttpInterceptorFn } from '@angular/common/http';
import { TestBed } from '@angular/core/testing';

import { cachingInterceptor } from './caching.interceptor';

describe('cachingInterceptor', () => {
  // eslint-disable-next-line max-len
  const interceptor: HttpInterceptorFn = (req, next) => TestBed.runInInjectionContext(() => cachingInterceptor(req, next));

  beforeEach(() => {
    TestBed.configureTestingModule({});
  });

  it('should be created', () => {
    expect(interceptor).toBeTruthy();
  });
});
