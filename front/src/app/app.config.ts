import { HTTP_INTERCEPTORS, provideHttpClient, withInterceptors } from '@angular/common/http';
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';

// import { SocketIoConfig, SocketIoModule } from 'ngx-socket-io';
// import { GlobalVariable } from '../global';
import { routes } from './app.routes';
import { unauthorizedInterceptor } from './unauthorized.interceptor';

// const config: SocketIoConfig = {
//   url: `${GlobalVariable.BASE_API_URL}/v1/ws`, // socket server url;
//   options: {},
// };

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([unauthorizedInterceptor])),
    // importProvidersFrom(SocketIoModule),
    {
      provide: HTTP_INTERCEPTORS,
      useFactory: unauthorizedInterceptor,
      multi: true,
    },
  ],
};
