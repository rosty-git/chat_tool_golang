import {
  HTTP_INTERCEPTORS,
  provideHttpClient,
  withInterceptors,
} from '@angular/common/http';
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
// import { cachingInterceptor } from './caching.interceptor';
import { unauthorizedInterceptor } from './unauthorized.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([unauthorizedInterceptor])),
    {
      provide: HTTP_INTERCEPTORS,
      useFactory: unauthorizedInterceptor,
      multi: true,
    },
  ],
};
