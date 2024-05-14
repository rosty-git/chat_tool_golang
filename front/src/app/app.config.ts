import { provideHttpClient } from '@angular/common/http';
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';

// import { SocketIoConfig, SocketIoModule } from 'ngx-socket-io';
// import { GlobalVariable } from '../global';
import { routes } from './app.routes';

// const config: SocketIoConfig = {
//   url: `${GlobalVariable.BASE_API_URL}/v1/ws`, // socket server url;
//   options: {},
// };

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(),
    // importProvidersFrom(SocketIoModule),
  ],
};
